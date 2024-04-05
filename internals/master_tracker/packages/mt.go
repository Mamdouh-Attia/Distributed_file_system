package mt

import (
	dk "Distributed_file_system/internals/data_keeper_node/packages"
	pb_m "Distributed_file_system/internals/pb/master_node"
	utils "Distributed_file_system/internals/utils"
	"context"
	"fmt"
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
	pb_m.UnimplementedMasterNodeServer
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
func (m *Master) RegisterDataNode(ctx context.Context, req *pb_m.RegisterDataNodeRequest) (*pb_m.RegisterDataNodeResponse, error) {
	//generate a new data node id
	id := utils.GenerateID()
	//TODO : if duplicate id, generate a new one

	// add the data node to the master
	m.DataKeeperNodes = append(m.DataKeeperNodes, *dk.NewDataKeeperNode(int(id), req.DataKeeper.Ip, req.DataKeeper.Port, []string{}))
	//print the data node
	log.Printf("DataKeeperNode: %v", m.DataKeeperNodes)
	// register file records in the master

	//start the heartbeat timer for this data node
	m.Heartbeats[int32(id)] = time.AfterFunc(5*time.Second, func() {
		m.KillDataNode(int32(id))
	})
	return &pb_m.RegisterDataNodeResponse{Success: true, NodeID: int32(id)}, nil
}

// grpc function to update the master with the new status of the data node
func (m *Master) HeartbeatUpdate(ctx context.Context, dk *pb_m.HeartbeatUpdateRequest) (*pb_m.HeartbeatUpdateResponse, error) {

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

	return &pb_m.HeartbeatUpdateResponse{Success: true}, nil
}

// grpc function to handle the request from client to download a file
func (m *Master) AskForDownload(ctx context.Context, file *pb_m.AskForDownloadRequest) (*pb_m.AskForDownloadResponse, error) {

	//print all records
	log.Printf("Master: %v", m.Records)
	log.Printf("Request for file: %v", file.FileName)

	// get all records with the same fileReqname
	records := m.GetRecordsByFilename(file.FileName)

	// if the file is not in the master, return an error
	if len(records) == 0 {

		// return error, that the file does not exist
		return &pb_m.AskForDownloadResponse{}, nil
	}

	// return list of DataKeeper struct used in response
	MachinesList := make([]*pb_m.DataKeeper, 0)
	for _, record := range records {
		dataKeeperNode := m.GetDataKeeperNodeById(record.DataKeeperNodeID)
		MachinesList = append(MachinesList, &pb_m.DataKeeper{Id: int32(dataKeeperNode.ID), Ip: dataKeeperNode.IP, Port: dataKeeperNode.Port})
	}

	log.Printf("Master: Found the file %v in the following datakeeper nodes: %v", file.FileName, MachinesList)

	// return the datakeeper nodes to the client
	return &pb_m.AskForDownloadResponse{DataKeepers: MachinesList}, nil
}

// grpc function to recieve Files list from data node
func (m *Master) ReceiveFileList(ctx context.Context, filesRequest *pb_m.ReceiveFileListRequest) (*pb_m.ReceiveFileListResponse, error) {

	// get the data node keeper by this id
	dataKeeperNode := m.GetDataKeeperNodeById(int(filesRequest.NodeID))

	for _, file := range filesRequest.Files {
		dataKeeperNode.AddFile(file)
	}

	// update the master with the new files list
	for i, node := range m.DataKeeperNodes {
		if node.ID == int(filesRequest.NodeID) {
			m.DataKeeperNodes[i] = dataKeeperNode
			// update the records
			for _, file := range filesRequest.Files {
				m.AddRecord(Record{FileName: file, FilePath: file, alive: true, DataKeeperNodeID: int(filesRequest.NodeID)})
			}
		}
	}

	// print the data node
	log.Printf("DataKeeperNode: %v files: %v", dataKeeperNode, dataKeeperNode.Files)

	// return success
	return &pb_m.ReceiveFileListResponse{Success: true}, nil
}

// grpc function to handle the request from client to upload a file
func (m *Master) AskForUpload(ctx context.Context, req *pb_m.Empty) (*pb_m.AskForUploadResponse, error) {
	// Choose a random datakeeper node to store the file
	dataKeeperNodeID := m.DataKeeperNodes[rand.Intn(len(m.DataKeeperNodes))]

	// return the datakeeper node to the client
	return &pb_m.AskForUploadResponse{Port: dataKeeperNodeID.Port}, nil
}

// grpc function to handle the notification from the dataNode when the uploading is done
func (m *Master) UploadNotification(ctx context.Context, notification *pb_m.UploadNotificationRequest) (*pb_m.UploadNotificationResponse, error) {

	newRecord := Record{FileName: notification.NewRecord.FileName, FilePath: notification.NewRecord.FilePath, alive: true, DataKeeperNodeID: int(notification.NewRecord.DataKeeperNodeID)}

	m.Records = append(m.Records, newRecord)

	// print the file name in the records of the master
	fmt.Printf("Master: File %v is uploaded to DataKeeperNode: %v\n", notification.NewRecord.FileName, notification.NewRecord.DataKeeperNodeID)

	return &pb_m.UploadNotificationResponse{Success: true}, nil
}

//// Replication functions //////

// replication thread function to replicate the files
// it checks every 10 seconds if there is a distinct file which is not replicated on at least 3 data nodes
func (m *Master) ReplicationThread() {
	//initialization
	//define set of distinct files -- extracted from Records
	distinctFiles := make(map[string]int)
	for _, record := range m.Records {
		distinctFiles[record.FileName] = 0
	}

	for {
		//loop over the distinct files
		for file := range distinctFiles {
			file_recs := m.GetRecordsByFilename(file)
			//get the source data node machine
			srcRecord := file_recs[0]
			sourceDataNodeID := srcRecord.DataKeeperNodeID
			//get the number of data nodes that have this file
			numDataNodes := len(file_recs)
			//while the number of data nodes is less than 3, replicate the file
			for numDataNodes < 3 {
				//get the destination data node machine
				//returns a valid IP and a valid port of a machine to copy a file instance to.
				destinationMachine := m.selectDestinationMachine(m.Records, file)
				if destinationMachine == nil {
					break
				}
				//notify both the source and the destination machines to replicate the file
				result := m.notifyMachineDataTransfer(sourceDataNodeID, destinationMachine, file)
				if result {
					log.Printf("Master: File %v is replicated from DataKeeperNode: %v to DataKeeperNode: %v\n", file, sourceDataNodeID, destinationMachine.ID)
				}
				//update the number of data nodes that have this file
				numDataNodes++
			}
			//sleep for 10 seconds
			time.Sleep(10 * time.Second)
		}
	}
}

// function to select the destination machine to replicate the file
// returns a valid IP and a valid port of a machine to copy a file instance to.
func (m *Master) selectDestinationMachine(records []Record, file string) *dk.DataKeeperNode {
	//get the list of all data nodes that have this file
	file_recs := m.GetRecordsByFilename(file)
	if len(file_recs) == 0 {
		return nil
	}
	//get the list of all data nodes that do not have this file
	destinationDataNodes := make([]dk.DataKeeperNode, 0)
	for _, record := range m.DataKeeperNodes {
		if !contains(file_recs, record.ID) {
			destinationDataNodes = append(destinationDataNodes, record)
		}
	}
	//if there is no destination data node, return nil
	if len(destinationDataNodes) == 0 {
		return nil
	}
	//get the destination data node machine
	destinationMachine := destinationDataNodes[rand.Intn(len(destinationDataNodes))]
	return &destinationMachine
}

// function to notify both the source and the destination machines to replicate the file
func (m *Master) notifyMachineDataTransfer(sourceDataNodeID int, destinationMachine *dk.DataKeeperNode, file string) bool {
	//get the source data node machine
	// sourceDataNode := m.GetDataKeeperNodeById(sourceDataNodeID)
	//TODO: notify the source machine to replicate the file
	return true
}

// function to check if a list contains a specific element
func contains(records []Record, id int) bool {
	for _, record := range records {
		if record.DataKeeperNodeID == id && record.alive {
			return true
		}
	}
	return false
}
