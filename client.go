package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	defer conn.Close()

	// Get input from console
	reader := bufio.NewReader(os.Stdin)
	// Read line each time the users presses enter
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch line {
		case "quit":
			break
		case "list local":
			output, err := exec.Command("tree", ".").Output()
			if err != nil {
				fmt.Println("Error listing local files: ", err)
				return
			}
			fmt.Println(string(output))
		default:
			// Send input to the server
			_, err = conn.Write([]byte(line))
			if err != nil {
				fmt.Println("Error ", err)
				return
			}

			// Read response from the server
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading data from the server: ", err)
				return
			}

			fmt.Println(string(buffer[:n]))
		}
	}
}
