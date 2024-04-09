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

	for {
		// Accept incoming connections and handle them
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			return
		}
		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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

		res, err := Menu(conn, line)
		if err != nil {
			fmt.Println("Error executing the command: ", err)
			return
		}
		_, err = conn.Write(res)
		if err != nil {
			fmt.Println("Error writing to client: ", err)
			return
		}
	}

}

func Menu(conn net.Conn, line string) ([]byte, error) {
	arr := strings.Split(line, " ")
	command := arr[0]

	switch command {
	case "quit":
		conn.Close()
	case "list":
		return exec.Command("ls", "-a").Output()
	case "mkdir":
		return []byte("Directorio creado correctamente"), exec.Command("mkdir", arr[1]).Run()
	case "rmdir":
		fileInfo, err := os.Stat(arr[1])
		if err != nil {
			return []byte("Error al obtener información sobre el archivo o carpeta."), err
		}
		if fileInfo.IsDir() {
			return []byte("Directorio eliminado correctamente."), exec.Command("rm", "-rf", arr[1]).Run()
		} else {
			return []byte("Archivo eliminado correctamente."), exec.Command("rm", arr[1]).Run()
		}
	case "get":
		file, err := os.Open(arr[1])
		if err != nil {
			return []byte("Error al abrir el archivo: "), err

		}
		defer file.Close()

		_, err = io.Copy(conn, file)
		if err != nil {
			return []byte("Error al copiar el archivo: "), err
		}
		return []byte("Archivo enviado correctamente"), nil
	case "put":
		file, err := os.Create(arr[1])

		if err != nil {
			return []byte("Error al abrir el archivo: "), err

		}
		defer file.Close()

		_, err = io.Copy(file, conn)
		if err != nil {
			return []byte("Error al copiar el archivo: "), err
		}
		return []byte("Archivo guardado correctamente"), nil
	}
	return []byte(""), nil
}
