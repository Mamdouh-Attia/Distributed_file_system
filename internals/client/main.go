// imports
package main

import (
	client "Distributed_file_system/internals/client/packages"
	mt "Distributed_file_system/internals/pb/master_node"
	"fmt"
)

const (
	address         = "localhost:8080"
	defaultFilename = "2MB.mp4"
)

func main() {

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

	fmt.Print("Test 0\n")
	uploadErr := Client.UploadFileToServer(masterClient, "2MB.mp4")

	if uploadErr != nil {
		fmt.Printf("Error while uploading the file: %v\n", uploadErr)
	}
}
