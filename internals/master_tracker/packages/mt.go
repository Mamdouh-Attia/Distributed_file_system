package mt

import (
	dk "Distributed_file_system/internals/data_keeper_node/packages"
	mt "Distributed_file_system/internals/pb/master_node"
	"context"
	"log"
	"math/rand"
)

// Record represents a single record in the master
type Record struct {
	FileName         string
	FilePath         string
	alive            bool
	DataKeeperNodeID int
	mt.UnimplementedMasterNodeServer
}

// Master represents the master data structure containing records
type Master struct {
	Records         []Record // List of records in the master
	DataKeeperNodes []dk.DataKeeperNode
	
}

// NewMaster creates a new Master instance'
func NewMaster() *Master {
	return &Master{
		Records:         []Record{}, // Initialize the list of records to be empty
		DataKeeperNodes: []dk.DataKeeperNode{},
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

func (m *Master) GetDataKeeperNodeById(dataKeeperNodeId int) dk.DataKeeperNode {
	for _, dataKeeperNode := range m.DataKeeperNodes {
		if dataKeeperNode.ID == dataKeeperNodeId {
			return dataKeeperNode
		}
	}
	return dk.DataKeeperNode{}
}

// grpc function to handle the request from data node to register itself in the master
func (s *Master) RegisterDataNode(ctx context.Context, req *mt.RegisterDataNodeRequest) (mt.RegisterDataNodeResponse, error) {
	// add the data node to the master
	s.DataKeeperNodes = append(s.DataKeeperNodes, *dk.NewDataKeeperNode(int(req.DataKeeper.Id), req.DataKeeper.Ip, req.DataKeeper.Port, []string{}))
	//print the data node
	log.Printf("DataKeeperNode: %v", s.DataKeeperNodes)
	return mt.RegisterDataNodeResponse{Success: true}, nil
}

// grpc function to update the master with the new status of the data node
func (s *Master) HeartbeatUpdate(ctx context.Context, dk *mt.DataKeeper) (mt.Empty, error) {

	// get all records with the same datakeeper node id
	records := s.GetRecordsByDataKeeperNode(int(dk.Id))

	// if the datakeeper node is in the master, update its status
	for _, record := range records {
		record.alive = true
	}

	return mt.Empty{}, nil
}

// grpc function to handle the request from client to upload a file
func (s *Master) UploadFile(ctx context.Context, file *mt.Empty) (mt.UploadResponse, error) {

	// Choose a random datakeeper node to store the file
	dataKeeperNodeID := s.Records[rand.Intn(len(s.Records))].DataKeeperNodeID

	// get the data node keeper by this id
	dataKeeperNode := s.GetDataKeeperNodeById(dataKeeperNodeID)

	// return the datakeeper node to the client
	return mt.UploadResponse{Port: dataKeeperNode.Port}, nil
}

// grpc function to handle the request from client to download a file
func (s *Master) AskForDownload(ctx context.Context, file *mt.AskForDownloadRequest) (mt.AskForDownloadResponse, error) {

	// get all records with the same filename
	records := s.GetRecordsByFilename(file.FileName)

	// if the file is not in the master, return an error
	if len(records) == 0 {

		// return error, that the file does not exist
		return mt.AskForDownloadResponse{}, nil
	}

	// get all Data nodes that has this file and return the port and Ip of each one in a map
	FileLocations := make(map[string]string)
	for _, record := range records {
		dataKeeperNode := s.GetDataKeeperNodeById(record.DataKeeperNodeID)

		FileLocations[dataKeeperNode.Port] = dataKeeperNode.IP
	}

	// return the datakeeper nodes to the client
	return mt.AskForDownloadResponse{FileLocations: FileLocations}, nil
}