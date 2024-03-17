package main

import (
	pb "Distributed_file_system/internals/proto/dfs"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Record represents a single record in the master
type Record struct {
	FileName         string
	DataKeeperNodeID int
	FilePath         string
	IsDataNodeAlive  bool
}

// Master represents the master data structure containing records
type Master struct {
	Records []Record // List of records in the master
	pb.UnimplementedHeartbeatUpdateServer
}

// NewMaster creates a new Master instance'
func NewMaster() *Master {
	return &Master{
		Records: []Record{}, // Initialize the list of records to be empty
	}
}

// AddRecord adds a new record to the master
func (m *Master) AddRecord(record Record) {
	m.Records = append(m.Records, record)
}

// RemoveRecord removes a record from the master
// optional parameters to remove a record either by filename or by datakeeper node
func (m *Master) RemoveRecord(filename string, dataKeeperNode int) {
	for i, record := range m.Records {
		if record.FileName == filename || record.DataKeeperNodeID == dataKeeperNode {
			m.Records = append(m.Records[:i], m.Records[i+1:]...)
			break
		}
	}
}

// GetRecordByFilename returns a record from the master by filename
func (m *Master) GetRecordByFilename(filename string) *Record {
	for _, record := range m.Records {
		if record.FileName == filename {
			return &record
		}
	}
	return nil
}

// GetRecordByDataKeeperNode returns a record from the master by datakeeper node
func (m *Master) GetRecordByDataKeeperNode(dataKeeperNode int) *Record {
	for _, record := range m.Records {
		if record.DataKeeperNodeID == dataKeeperNode {
			return &record
		}
	}
	return nil
}

///Mamdouh : heartbeat

// grpc function to update the master with the new status of the data node
func (s *Master) HeartbeatUpdate(ctx context.Context, in *pb.Heartbeat) (*pb.HeartbeatResponse, error) {
	//check if the data node is registered in the master
	record := s.GetRecordByDataKeeperNode(int(in.Id))
	//if the data node is not registered in the master
	if record == nil {
		//add the data node to the master
		log.Printf("Data node with id %v is not registered in the master", in.Id)
		s.AddRecord(Record{FileName: "", DataKeeperNodeID: int(in.Id), FilePath: "", IsDataNodeAlive: true})
	} else {
		//update the status of the data node
		log.Printf("Data node with id %v is registered in the master", in.Id)
		record.IsDataNodeAlive = true
	}
	//return the response with id , ip, port and status of the data node
	return &pb.HeartbeatResponse{Id: in.Id, Ip: in.Ip, Port: int32(in.Port), Status: "OK"}, nil
}

func main() {
	// Set up the master
	master := NewMaster()
	//print the master
	log.Printf("Master: %v", master)
	// Set up the gRPC server
	lis, err := net.Listen("tcp", "localhost:8080")
	log.Printf("Master is listening on localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("Master is serving")
	pb.RegisterHeartbeatUpdateServer(s, master)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//

}
