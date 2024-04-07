package main

import (
	mt "Distributed_file_system/internals/master_tracker/packages"
	pb "Distributed_file_system/internals/pb/master_node"
	"fmt"
	"log"
	"net"
	"time"

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

	// separate goroutine to handle the heartbeat updates
	go func() {
		for {
			//change all the records to false
			// for i := range master.Records {
			// master.DataKeeperNodes[master.Records[i].DataKeeperNodeID].
			// }
			// sleep for 2 second
			time.Sleep(2 * time.Second)
			// Remove dead records
			master.RemoveDeadRecords()
		}
	}()

	// A separate thread to handle the replication process every 10 seconds
	go func() {
		fmt.Println("Replication process started")
		for {
			// replicate the files
			fmt.Println("Replicating files")
			master.ReplicateFiles()
			// sleep for 10 seconds
			time.Sleep(10 * time.Second)
		}
	}()
	fmt.Println("Master is serving")
	// pb.RegisterHeartbeatUpdateServer(s, master)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
