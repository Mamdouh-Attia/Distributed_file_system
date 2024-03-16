package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer to read incoming data
	buffer := make([]byte, 1024)

	// Read filename from client
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// Convert received data to string
	filename := strings.TrimSpace(string(buffer[:n]))

	// Open the requested file (mp4 file)
	filename += ".mp4"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		return
	}
	defer file.Close()

	// Read file contents
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error reading file:", err.Error())
		return
	}

	// Send file size to client
	_, err = conn.Write([]byte(fmt.Sprintf("%d\n", fileInfo.Size())))
	if err != nil {
		fmt.Println("Error sending file size:", err.Error())
		return
	}

	// Send file contents to client
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file contents:", err.Error())
		return
	}
}

func main() {
	// Listen for incoming connections on port 8080
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Listening on port 8080...")

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle incoming connection in a separate goroutine
		go handleConnection(conn)
	}
}
