package mt

import (
	dk "Distributed_file_system/internals/data_keeper_node/packages"
	pb_d "Distributed_file_system/internals/pb/data_node"
	pb_m "Distributed_file_system/internals/pb/master_node"
	utils "Distributed_file_system/internals/utils"
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
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
	mu         sync.Mutex
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
func (mt *Master) ReplicateFile(ctx context.Context, req *pb_d.ReplicaRequest) (*pb_d.NotifyReplicaResponse, error) {
	// return
	return &pb_d.NotifyReplicaResponse{Success: true}, nil
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

// remove dead records from the master
func (m *Master) RemoveDeadRecords() {
	for i, record := range m.Records {
		if !record.alive {
			m.Records = append(m.Records[:i], m.Records[i+1:]...)
		}
	}
}

// grpc function to handle the request from data node to register itself in the master
func (m *Master) RegisterDataNode(ctx context.Context, req *pb_m.RegisterDataNodeRequest) (*pb_m.RegisterDataNodeResponse, error) {
	//generate a new data node id
	id := utils.GenerateID()
	//TODO : if duplicate id, generate a new one

	//mutex lock
	m.mu.Lock()
	defer m.mu.Unlock()

	// add the data node to the master
	m.DataKeeperNodes = append(m.DataKeeperNodes, *dk.NewDataKeeperNode(int(id), req.DataKeeper.Ip, req.DataKeeper.Port, []string{}))
	//print the data node
	log.Printf("DataKeeperNode: %v", m.DataKeeperNodes)
	// register file records in the master

	//start the heartbeat timer for this data node
	m.Heartbeats[int32(id)] = time.AfterFunc(5*time.Second, func() {
		m.KillDataNode(int32(id))
	})
	log.Printf("DataKeeperNode: %v is registered", id)
	return &pb_m.RegisterDataNodeResponse{Success: true, NodeID: int32(id)}, nil
}

// grpc function to update the master with the new status of the data node
func (m *Master) HeartbeatUpdate(ctx context.Context, dk *pb_m.HeartbeatUpdateRequest) (*pb_m.HeartbeatUpdateResponse, error) {

	log.Printf("HeartbeatUpdateRequest: ")
	// get all records with the same datakeeper node id
	records := m.GetRecordsByDataKeeperNode(int(dk.NodeID))

	// if the data node is not in the master, return an error
	if len(records) == 0 {
		return &pb_m.HeartbeatUpdateResponse{Success: false}, nil
	}

	//check if the data node exists in the master
	if _, ok := m.Heartbeats[dk.NodeID]; !ok {
		//start the heartbeat timer for this data node
		m.Heartbeats[dk.NodeID] = time.AfterFunc(5*time.Second, func() {
			m.KillDataNode(dk.NodeID)
		})
	}
	//reset the heartbeat timer
	m.Heartbeats[dk.NodeID].Reset(5 * time.Second)
	// if the datakeeper node is in the master, update its status
	for _, record := range records {
		record.alive = true
	}
	//TODO::log the data node is alive
	// log.Printf("DataKeeperNode: %v is alive", dk.NodeID)

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
func (m *Master) ReplicateFiles() {

	log.Printf("Replicating files")
	// define set of distinct files -- extracted from Records
	distinctFiles := m.getDistinctFiles()
	if len(distinctFiles) == 0 {
		log.Printf("No files found for replication")
		return
	}

	for file := range distinctFiles {

		fileRecords := m.GetRecordsByFilename(file)

		//get the source data node machine
		if len(fileRecords) == 0 {
			log.Printf("No records found for file: %v", file)
			continue
		}
		srcRecord := fileRecords[0]
		sourceDataNodeID := srcRecord.DataKeeperNodeID
		sourceDataNode := m.GetDataKeeperNodeById(sourceDataNodeID)

		//get the number of data nodes that have this file
		numDataNodes := len(fileRecords)

		//while the number of data nodes is less than 3, replicate the file
		for numDataNodes < 3 {

			//get the destination data node machine
			//returns a valid IP and a valid port of a machine to copy a file instance to.
			destinationMachine := m.selectDestinationMachine(file)
			if destinationMachine == nil {
				log.Printf("No destination machine found for file: %v", file)
				break
			}
			log.Printf("Source: %v will replicate file: %v to destination: %v\n", sourceDataNode, file, destinationMachine)

			//connect to the source data node
			conn, err := grpc.Dial(fmt.Sprintf("%s:%s", sourceDataNode.IP, sourceDataNode.Port), grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to connect to the source data node: %v", err)
				break
			}
			defer conn.Close()
			//create a new client
			// client := pb_d.NewDataNodeClient(conn)
			//1- notify source data node to replicate the file
			//2- get the destination data node machine

			//client instance
			client := pb_d.NewDataNodeClient(conn)
			//1- notify source data node to replicate the file
			_, err = client.ReceiveFileForReplica(context.Background(), &pb_d.ReceiveFileForReplicaRequest{Ip: destinationMachine.IP, Port: destinationMachine.Port})
			if err != nil {
				log.Printf("Failed to notify the source data node: %v", err)
				break
			}
			//print the destination machine
			log.Printf("Destination Machine: %v", destinationMachine.IP)
			log.Printf("Destination Machine: %v", destinationMachine.Port)

			//connect to the destination data node
			conn2, err := grpc.Dial(fmt.Sprintf("%s:%s", destinationMachine.IP, destinationMachine.Port), grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to connect to the destination data node: %v", err)
				break
			}
			defer conn2.Close()
			//create a new client
			client2 := pb_d.NewDataNodeClient(conn2)
			//2- get the destination data node machine
			_, err = client2.ReplicateFile(context.Background(), &pb_d.ReplicaRequest{Ip: sourceDataNode.IP, Port: sourceDataNode.Port, FileName: file})
			if err != nil {
				log.Printf("Failed to notify the destination data node: %v", err)
				break
			}

			// update the records of the master
			m.AddRecord(Record{FileName: file, FilePath: file, alive: true, DataKeeperNodeID: destinationMachine.ID})

			//update the number of data nodes that have this file
			numDataNodes++
		}

	}
}

// function to select the destination machine to replicate the file
// returns a valid IP and a valid port of a machine to copy a file instance to.
func (m *Master) selectDestinationMachine(file string) *dk.DataKeeperNode {

	//get the list of all data nodes that have this file
	fileRecords := m.GetRecordsByFilename(file)
	if len(fileRecords) == 0 {
		log.Printf("No records found for file: %v", file)
		return nil
	}

	//get the list of all data nodes that do not have this file
	destinationDataNodes := make([]dk.DataKeeperNode, 0)
	for _, record := range m.DataKeeperNodes {
		if !contains(fileRecords, record.ID) {
			destinationDataNodes = append(destinationDataNodes, record)
		}
	}

	//if there is no destination data node, return nil
	if len(destinationDataNodes) == 0 {
		log.Printf("No destination data nodes found for file: %v", file)
		return nil
	}

	log.Printf("Debug|| Destination Data Nodes: %v", destinationDataNodes)
	//get the destination data node machine
	destinationMachine := destinationDataNodes[rand.Intn(len(destinationDataNodes))]
	return &destinationMachine
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

// function to get destinct files
func (m *Master) getDistinctFiles() map[string]int {

	distinctFiles := make(map[string]int)
	for _, record := range m.Records {
		distinctFiles[record.FileName] = 0
	}

	return distinctFiles
}
