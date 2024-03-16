package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	// Read filename from user input
	var filename string
	fmt.Print("Enter the name of the MP4 file: ")
	fmt.Scanln(&filename)

	// Send filename to server
	_, err = conn.Write([]byte(filename))
	if err != nil {
		fmt.Println("Error sending filename:", err.Error())
		return
	}

	// Read file size from server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving file size:", err.Error())
		return
	}

	// Convert received file size to integer
	fileSizeStr := strings.TrimSpace(string(buffer[:n]))
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		fmt.Println("Error parsing file size:", err.Error())
		return
	}

	// Create a buffer to store file contents
	fileBuffer := make([]byte, fileSize)

	// Read file contents from server
	totalRead := 0
	for totalRead < int(fileSize) {
		n, err := conn.Read(fileBuffer[totalRead:])
		if err != nil {
			fmt.Println("Error receiving file contents:", err.Error())
			return
		}
		totalRead += n
	}

	// Write file contents to a new file
	filename += ".mp4"
	newFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
		return
	}
	defer newFile.Close()

	_, err = newFile.Write(fileBuffer)
	if err != nil {
		fmt.Println("Error writing to file:", err.Error())
		return
	}

	fmt.Println("File received successfully.")
}
