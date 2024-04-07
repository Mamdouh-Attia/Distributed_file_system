package main

import (
	mt "Distributed_file_system/internals/master_tracker/packages"
	pb "Distributed_file_system/internals/pb/master_node"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func main() {

	//Declarations

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

	// separate goroutine to handle the heartbeat updates
	go func() {
		for {
			//change all the records to false
			// for i := range master.Records {
			// master.DataKeeperNodes[master.Records[i].DataKeeperNodeID].
			// }
			// sleep for 2 second
			time.Sleep(2 * time.Second)
			// check if the records are alive
			for i := range master.Records {
				if master.DataKeeperNodes[master.Records[i].DataKeeperNodeID].Files == nil {
					// remove the record
					master.RemoveRecord(master.Records[i].FileName, master.Records[i].DataKeeperNodeID)
				}
			}

		}
	}()

	// Make a separate thread to handle the replication process every 10 seconds
	go func() {
		for {
			// sleep for 10 seconds
			time.Sleep(10 * time.Second)
			// replicate the files
			master.ReplicateFiles()
		}
	}()

	// pb.RegisterHeartbeatUpdateServer(s, master)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
