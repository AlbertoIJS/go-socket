package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	// Listen for incoming connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening at port 8080: ", err)
		return
	}
	defer listener.Close()

	// Accept incoming connections and handle them
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Close the connection when we're done
	defer conn.Close()

	//	Handle connection logic
	fmt.Println("Accepted connection from:", conn.RemoteAddr())

	//	Read incoming data
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		line := string(buffer[:n])
		line = strings.TrimSpace(line)

		if err != nil {
			fmt.Println("Error reading data from client: ", err)
			return
		}

		res, err := menu(line)

		_, err = conn.Write(res)
		if err != nil {
			fmt.Println("Error sending to client: ", err)
			return
		}
	}
}

func menu(line string) ([]byte, error) {
	arr := strings.Split(line, " ")
	command := arr[0]

	switch command {
	case "list":
		return exec.Command("ls", "-a").Output()
	case "mkdir":
		return exec.Command("mkdir", arr[1]).Output()
	case "rmdir":
		return exec.Command("rm -rf", arr[2]).Output()
	default:
		return []byte("Unknown command"), nil
	}
}
