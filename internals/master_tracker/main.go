package main

import (
	mt "Distributed_file_system/internals/master_tracker/packages"
	pb "Distributed_file_system/internals/pb/master_node"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// Set up the master
	master := mt.NewMaster()
	// print the master
	log.Printf("Master: %v", master)
	// Set up the gRPC server
	lis, err := net.Listen("tcp", "localhost:8080")
	log.Printf("Master is listening on localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//register methods
	//1.RegisterDataNode

	s := grpc.NewServer()
	pb.RegisterMasterNodeServer(s, master)
	log.Printf("Master is serving")
	// pb.RegisterHeartbeatUpdateServer(s, master)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
