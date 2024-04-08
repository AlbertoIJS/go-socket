package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

		command := strings.Split(line, " ")[0]
		if command == "get" || command == "put" {
			filesMenu(conn, line)
		} else {
			res, err := consoleMenu(line)
			if err != nil {
				fmt.Println("Error sending to client: ", err)
				return
			}
			_, err = conn.Write(res)
		}

	}
}

func consoleMenu(line string) ([]byte, error) {
	arr := strings.Split(line, " ")
	command := arr[0]

	switch command {
	case "list":
		return exec.Command("ls", "-a").Output()
	case "mkdir":
		return exec.Command("mkdir", arr[1]).Output()
	case "rmdir":
		fileInfo, err := os.Stat(arr[1])
		if err != nil {
			return []byte("Error al obtener información sobre el archivo o carpeta."), err
		}
		if fileInfo.IsDir() {
			return exec.Command("rm", "-rf", arr[1]).Output()
		} else {
			return exec.Command("rm", arr[1]).Output()
		}

	default:
		return []byte("Comando desconocido."), nil
	}
}

func filesMenu(conn net.Conn, line string) {
	arr := strings.Split(line, " ")
	command := arr[0]

	switch command {
	case "get":
		file, err := os.Open(arr[1])
		if err != nil {
			fmt.Println("Error al abrir el archivo: ", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(conn, file)
		if err != nil {
			fmt.Println("Error al copiar el archivo: ", err)
			return
		}
	case "put":
		file, err := os.Create(arr[1])

		if err != nil {
			fmt.Println("Error al abrir el archivo: ", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, conn)
		if err != nil {
			fmt.Println("Error al copiar el archivo: ", err)
			return
		}
	}
}
