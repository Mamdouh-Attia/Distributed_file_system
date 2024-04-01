package main

import (
	dk "Distributed_file_system/internals/data_keeper_node/packages"
	mt "Distributed_file_system/internals/pb/master_node"
	utils "Distributed_file_system/internals/utils"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//generat
	// Create a new DataKeeperNode instance

	node := dk.NewDataKeeperNode(1, "localhost", "8080", []string{"2MB"})
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server: %v", "localhost:8080")
	done := make(chan struct{}) // create a channel to keep the main function alive
	defer close(done)           // close the done channel when the main function returns
	// Create a client for the DistributedFileSystem service
	client := mt.NewMasterNodeClient(conn)
	//Flow :
	// 1. Send a datakeeper node registration request to the server

	regResult, err := client.RegisterDataNode(context.Background(), &mt.RegisterDataNodeRequest{DataKeeper: &mt.DataKeeper{Id: int32(node.ID), Ip: node.IP, Port: node.Port}})

	if err != nil {
		log.Fatalf("Failed to register datakeeper node: %v", err)
	}
	log.Printf("Datakeeper node registration response: %v", regResult)

	// scan the current directory for files
	files, err := utils.FindMP4Files("./data")
	if err != nil {
		log.Fatalf("Failed to find mp4 files: %v", err)
	}
	//print the files
	log.Printf("Files found: %v", files)
	// send the file list to the server
	updateFilesListResult, err := client.ReceiveFileList(context.Background(), &mt.ReceiveFileListRequest{NodeID: int32(node.ID), Files: files})
	if err != nil {
		log.Fatalf("Failed to update file list: %v", err)
	}
	log.Printf("File list update response: %v", updateFilesListResult)
}
