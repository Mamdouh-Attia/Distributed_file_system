package main

import (
	dk "Distributed_file_system/internals/data_keeper_node/packages"
	pb "Distributed_file_system/internals/pb/data_node"
	mt "Distributed_file_system/internals/pb/master_node"
	utils "Distributed_file_system/internals/utils"
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//generate random port number
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(10000) + 10000
	portStr := strconv.Itoa(port)
	// Create a new DataKeeperNode instance

	node := dk.NewDataKeeperNode(1, "localhost", portStr, []string{"2MB"})
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
	log.Printf("Datakeeper node ID: %v", regResult.NodeID)
	node.ID = int(regResult.NodeID)

	// 2. scan the current directory for files
	files, err := utils.FindMP4Files()
	if err != nil {
		log.Fatalf("Failed to find mp4 files: %v", err)
	}
	//print the files
	log.Printf("Files found: %v", files)
	// 3. send the file list to the server
	updateFilesListResult, err := client.ReceiveFileList(context.Background(), &mt.ReceiveFileListRequest{NodeID: int32(node.ID), Files: files})
	if err != nil {
		log.Fatalf("Failed to update file list: %v", err)
	}
	log.Printf("File list update response: %v", updateFilesListResult)

	//sepreate goroutine to send the heartbeat
	go func() {
		for {
			// send the heartbeat
			_, err := client.HeartbeatUpdate(context.Background(), &mt.HeartbeatUpdateRequest{NodeID: int32(node.ID)})
			if err != nil {
				log.Fatalf("Failed to send heartbeat: %v", err)
			}
			// sleep for 1 seconds
			time.Sleep(1 * time.Second)
		}
	}()

	//sepate goroutine to serve the data node
	go func() {
		// Set up the gRPC server to listen on its port
		lis, err := net.Listen("tcp", node.IP+":"+node.Port)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		// Create a new gRPC server
		server := grpc.NewServer()
		// Register the DataNode service with the server
		pb.RegisterDataNodeServer(server, node)
		// Serve the server
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
		// close the done channel
		done <- struct {
		}{}
	}()

	// keep the main function alive
	<-done //this statement will block the main function until the done channel is closed

}
