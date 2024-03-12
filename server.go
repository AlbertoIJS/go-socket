package main

import (
	"fmt"
	"net"
	"os/exec"
)

func main() {
	// Listen for incoming connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer listener.Close()

	// Accept incoming connections and handle them
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
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
		if err != nil {
			fmt.Println("Error reading data: ", err)
		}
		if line == "list remote" {
			output, err := exec.Command("tree", "remote").Output()
			if err != nil {
				fmt.Println("Error executing command: ", err)
				return
			}

			// Send result to client
			_, err = conn.Write(output)
			if err != nil {
				fmt.Println("Error sending to client: ", err)
				return
			}
		}
	}
}
