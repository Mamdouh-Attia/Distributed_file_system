package mt

import (
	dk "Distributed_file_system/internals/data_keeper_node/packages"
	mt "Distributed_file_system/internals/pb/master_node"
	utils "Distributed_file_system/internals/utils"
	"context"
	"log"
	"math/rand"
	"time"
)

// Record represents a single record in the master
type Record struct {
	FileName         string
	FilePath         string
	alive            bool
	DataKeeperNodeID int
}

// Master represents the master data structure containing records
type Master struct {
	mt.UnimplementedMasterNodeServer
	Records         []Record // List of records in the master
	DataKeeperNodes []dk.DataKeeperNode
	//map of datakeeper node id to timer to keep track of the heartbeat
	Heartbeats map[int32]*time.Timer
}

// NewMaster creates a new Master instance'
func NewMaster() *Master {
	return &Master{
		Records:         []Record{}, // Initialize the list of records to be empty
		DataKeeperNodes: []dk.DataKeeperNode{},
		Heartbeats:      make(map[int32]*time.Timer),
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

// function to kill the data node, will be associated with the heartbeat timer
func (m *Master) KillDataNode(dataKeeperNodeId int32) {
	// get all records with the same datakeeper node id
	records := m.GetRecordsByDataKeeperNode(int(dataKeeperNodeId))

	// if the datakeeper node is in the master, update its status
	for _, record := range records {
		record.alive = false
	}
	//log the data node is dead
	log.Printf("DataKeeperNode: %v is dead", dataKeeperNodeId)
}

// grpc function to handle the request from data node to register itself in the master
func (m *Master) RegisterDataNode(ctx context.Context, req *mt.RegisterDataNodeRequest) (*mt.RegisterDataNodeResponse, error) {
	//generate a new data node id
	id := utils.GenerateID()
	//TODO : if duplicate id, generate a new one

	// add the data node to the master
	m.DataKeeperNodes = append(m.DataKeeperNodes, *dk.NewDataKeeperNode(int(id), req.DataKeeper.Ip, req.DataKeeper.Port, []string{}))
	//print the data node
	log.Printf("DataKeeperNode: %v", m.DataKeeperNodes)
	//start the heartbeat timer for this data node
	m.Heartbeats[int32(id)] = time.AfterFunc(5*time.Second, func() {
		m.KillDataNode(int32(id))
	})
	return &mt.RegisterDataNodeResponse{Success: true, NodeID: int32(id)}, nil
}

// grpc function to update the master with the new status of the data node
func (m *Master) HeartbeatUpdate(ctx context.Context, dk *mt.HeartbeatUpdateRequest) (*mt.HeartbeatUpdateResponse, error) {

	log.Printf("HeartbeatUpdateRequest: ")
	// get all records with the same datakeeper node id
	records := m.GetRecordsByDataKeeperNode(int(dk.NodeID))

	//reset the heartbeat timer
	m.Heartbeats[dk.NodeID].Reset(5 * time.Second)
	// if the datakeeper node is in the master, update its status
	for _, record := range records {
		record.alive = true
	}
	//log the data node is alive
	log.Printf("DataKeeperNode: %v is alive", dk.NodeID)

	return &mt.HeartbeatUpdateResponse{Success: true}, nil
}


// grpc function to handle the request from client to download a file
func (m *Master) AskForDownload(ctx context.Context, file *mt.AskForDownloadRequest) (*mt.AskForDownloadResponse, error) {

	// get all records with the same filename
	records := m.GetRecordsByFilename(file.FileName)

	// if the file is not in the master, return an error
	if len(records) == 0 {

		// return error, that the file does not exist
		return &mt.AskForDownloadResponse{}, nil
	}

	// get all Data nodes that has this file and return the port and Ip of each one in a map
	FileLocations := make(map[string]string)
	for _, record := range records {
		dataKeeperNode := m.GetDataKeeperNodeById(record.DataKeeperNodeID)

		FileLocations[dataKeeperNode.Port] = dataKeeperNode.IP
	}

	// return the datakeeper nodes to the client
	return &mt.AskForDownloadResponse{FileLocations: FileLocations}, nil
}

// grpc function to recieve Files list from data node
func (m *Master) ReceiveFileList(ctx context.Context, filesRequest *mt.ReceiveFileListRequest) (*mt.ReceiveFileListResponse, error) {

	// get the data node keeper by this id
	dataKeeperNode := m.GetDataKeeperNodeById(int(filesRequest.NodeID))

	for _, file := range filesRequest.Files {
		dataKeeperNode.AddFile(file)
	}

	// update the master with the new files list
	for i, node := range m.DataKeeperNodes {
		if node.ID == int(filesRequest.NodeID) {
			m.DataKeeperNodes[i] = dataKeeperNode
		}
	}

	// print the data node
	log.Printf("DataKeeperNode: %v files: %v", dataKeeperNode, dataKeeperNode.Files)

	// return success
	return &mt.ReceiveFileListResponse{Success: true}, nil
}

// grpc function to handle the request from client to upload a file
func (m* Master) AskForUpload(ctx context.Context, req *mt.Empty) (*mt.AskForUploadResponse, error) {
	// Choose a random datakeeper node to store the file
	dataKeeperNodeID := m.Records[rand.Intn(len(m.Records))].DataKeeperNodeID

	// get the data node keeper by this id
	dataKeeperNode := m.GetDataKeeperNodeById(dataKeeperNodeID)

	// return the datakeeper node to the client
	return &mt.AskForUploadResponse{Port: dataKeeperNode.Port}, nil
}