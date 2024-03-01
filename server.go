package main

import (
	"fmt"
	"net"
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
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	//	Read incoming data
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err)
		return
	}

	//	Print the incoming data
	fmt.Printf("Received data: %s\n", string(buffer[:n]))
}
