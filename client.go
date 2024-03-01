package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error ", err)
		return
	}

	// Send some data to the server
	_, err = conn.Write([]byte("Hello, server!"))
	if err != nil {
		fmt.Println("Error ", err)
		return
	}

	conn.Close()
}
