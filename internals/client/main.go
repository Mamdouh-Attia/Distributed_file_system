// imports
package main

import (
	client "Distributed_file_system/internals/client/packages"
	mt "Distributed_file_system/internals/pb/master_node"
	"fmt"
	"os"
)

const (
	address         = "localhost:8080"
	defaultFilename = "2MB.mp4"
)

func main() {

	//take file name as argument
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %v <filename>\n", os.Args[0])
		return
	}

	//arg1 : "u" or "d"
	//arg2 : filename

	// the client should connect to the master
	Client := client.NewClient()

	conn, errConn := Client.ConnectToServer(address)

	if errConn != nil {
		fmt.Printf("Failed to connect to server: %v", errConn)
		return
	}

	defer conn.Close()

	fmt.Printf("Connected to server: %v\n", "localhost:8080")

	masterClient := mt.NewMasterNodeClient(conn)

	switch os.Args[1] {
	case "u":
		uploadErr := Client.UploadFileToServer(masterClient, os.Args[2])

		if uploadErr != nil {
			fmt.Printf("Error while uploading the file: %v\n", uploadErr)
		}
	case "d":
		Client.DownloadFile(masterClient, os.Args[2])
	default:
		fmt.Printf("Usage: %v <u/d> <filename>\n", os.Args[0])
	}
}

// fmt.Print("Test 0\n")
