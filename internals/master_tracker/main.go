package main

import (
	"Distributed_file_system/internals/data_keeper_node"
	pb "Distributed_file_system/internals/pb"
	"context"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"
)

// Record represents a single record in the master
type Record struct {
	FileName         string
	FilePath         string
	alive			bool
	DataKeeperNodeID int
}

// Master represents the master data structure containing records
type Master struct {
	Records []Record // List of records in the master
	DataKeeperNodes []data_keeper_node.DataKeeperNode
}

// NewMaster creates a new Master instance'
func NewMaster() *Master {
	return &Master{
		Records: []Record{}, // Initialize the list of records to be empty
		DataKeeperNodes: []data_keeper_node.DataKeeperNode{},
	}
}

// AddRecord adds a new record to the master
func (m *Master) AddRecord(record Record) {
	m.Records = append(m.Records, record)
}

// RemoveRecord removes a record from the master
// optional parameters to remove a record either by filename or by datakeeper node
func (m *Master) RemoveRecord(filename string, dataKeeperNodeId int) {
	for i, record := range m.Records {
		if record.FileName == filename && record.DataKeeperNodeID == dataKeeperNodeId {
			m.Records = append(m.Records[:i], m.Records[i+1:]...)
			break
		}
	}
}

// GetRecordsByFilename returns all records from the master that has this filename
func (m *Master) GetRecordsByFilename(filename string) []Record {
	records := []Record{}
	for _, record := range m.Records {
		if record.FileName == filename {
			records = append(records, record)
		}
	}
	return records
}

// GetRecordsByDataKeeperNode returns all records from the master that have this datakeeper node
func (m *Master) GetRecordsByDataKeeperNode(dataKeeperNodeId int) []Record {
	records := []Record{}
	for _, record := range m.Records {
		if record.DataKeeperNodeID == dataKeeperNodeId {
			records = append(records, record)
		}
	}
	return records
}

func (m* Master) GetDataKeeperNodeById(dataKeeperNodeId int) data_keeper_node.DataKeeperNode {
	for _, dataKeeperNode := range m.DataKeeperNodes {
		if dataKeeperNode.ID == dataKeeperNodeId {
			return dataKeeperNode
		}
	}
	return data_keeper_node.DataKeeperNode{}
}


// grpc function to update the master with the new status of the data node
func (s *Master) HeartbeatUpdate(ctx context.Context, dk *pb.DataKeeper) (pb.Empty, error) {

	// get all records with the same datakeeper node id
	records := s.GetRecordsByDataKeeperNode(int(dk.Id))


	// if the datakeeper node is in the master, update its status
	for _, record := range records {
		record.alive = true
	}

	return pb.Empty{}, nil
}

// grpc function to handle the request from client to upload a file 
func (s *Master) UploadFile(ctx context.Context, file *pb.Empty) (pb.UploadResponse, error) {
	
	// Choose a random datakeeper node to store the file
	dataKeeperNodeID := s.Records[rand.Intn(len(s.Records))].DataKeeperNodeID

	// get the data node keeper by this id
	dataKeeperNode := s.GetDataKeeperNodeById(dataKeeperNodeID)

	// return the datakeeper node to the client
	return pb.UploadResponse{Port: dataKeeperNode.Port}, nil
}

// grpc function to handle the request from client to download a file
func (s *Master) AskForDownload(ctx context.Context, file *pb.AskForDownloadRequest) (pb.AskForDownloadResponse, error) {
	
	// get all records with the same filename
	records := s.GetRecordsByFilename(file.FileName)

	// if the file is not in the master, return an error
	if len(records) == 0 {
		
		// return error, that the file does not exist
		return pb.AskForDownloadResponse{}, nil
	}

	// get all Data nodes that has this file and return the port and Ip of each one in a map
	FileLocations := make(map[string]string)
	for _, record := range records {
		dataKeeperNode := s.GetDataKeeperNodeById(record.DataKeeperNodeID)
		
		FileLocations[dataKeeperNode.Port] = dataKeeperNode.IP
	}

	// return the datakeeper nodes to the client
	return pb.AskForDownloadResponse{FileLocations: FileLocations}, nil
}


func main() {
	// Set up the master
	master := NewMaster()
	// print the master
	log.Printf("Master: %v", master)
	// Set up the gRPC server
	lis, err := net.Listen("tcp", "localhost:8080")
	log.Printf("Master is listening on localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("Master is serving")
	// pb.RegisterHeartbeatUpdateServer(s, master)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//

}
