package main

import (
	dk "Distributed_file_system/internals/data_keeper_node/packages"
	mt "Distributed_file_system/internals/pb/master_node"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a new DataKeeperNode instance
	node := dk.NewDataKeeperNode(1, "localhost", "8080", []string{"2MB"})
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server: %v", "localhost:8080")
	done := make(chan struct{})
	defer close(done) // close the done channel when the main function returns
	// Create a client for the DistributedFileSystem service
	client := mt.NewMasterNodeClient(conn)
	//Flow :
	// 1. Send a datakeeper node registration request to the server
	
	regResult, err := client.RegisterDataNode(context.Background(), &mt.RegisterDataNodeRequest{DataKeeper: &mt.DataKeeper{Id: int32(node.ID), Ip: node.IP, Port: node.Port}})
	
	if err != nil {
		log.Fatalf("Failed to register datakeeper node: %v", err)
	}
	log.Printf("Datakeeper node registration response: %v", regResult)

	

}
