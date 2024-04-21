package dk

import (
	pb_d "Distributed_file_system/internals/pb/data_node"
	pb_m "Distributed_file_system/internals/pb/master_node"
	"Distributed_file_system/internals/utils"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"

	"context"

	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type File struct {
	FileName    string
	FilePath    string
	FileContent []byte
}

const (
	defaultAdress = "localhost:8080"
)

// DataKeeperNode represents a node responsible for keeping track of files.
type DataKeeperNode struct {
	pb_d.UnimplementedDataNodeServer
	ID    int      // Unique identifier for the node
	IP    string   // IP address of the node
	Port  string   // Port number the node is listening on
	Files []string // List of filenames stored on the node
	currBusyPorts int // a number indicating the current busy ports 
}

// NewDataKeeperNode creates a new DataKeeperNode instance with the provided parameters.
func NewDataKeeperNode(id int, ip string, port string, files []string) *DataKeeperNode {
	return &DataKeeperNode{
		ID:    id,
		IP:    ip,
		Port:  port,
		Files: files,
	}
}

// AddFile adds a new file to the list of files stored on the node.
func (n *DataKeeperNode) AddFile(file string) {
	n.Files = append(n.Files, file)
}

func (n *DataKeeperNode) ConnectToServer(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return conn, err
}

// RemoveFile removes a file from the list of files stored on the node.
func (n *DataKeeperNode) RemoveFile(filename string) {
	for i, file := range n.Files {
		if file == filename {
			n.Files = append(n.Files[:i], n.Files[i+1:]...)
			break
		}
	}
}

func (n *DataKeeperNode) GetFileByName(filename string) string {
	for _, file := range n.Files {
		if file == filename {
			return file

		}
	}
	return ""
}

func (n *DataKeeperNode) GetCurrentEmptyPort() string {
	portAsNum, err := utils.ConvertStrIntoInt(n.Port)
	if err != nil {
		fmt.Print("Error converting the port into an integer")
		return ""
	}
	portAsNum += n.currBusyPorts
	return fmt.Sprint(portAsNum)
}


// grpc function to send the current free port
func (n* DataKeeperNode) AskForFreePort(ctx context.Context, req* pb_d.AskForFreePortRequest) (*pb_d.AskForFreePortResponse, error) {
	n.currBusyPorts += 1
	return &pb_d.AskForFreePortResponse{Port: n.GetCurrentEmptyPort()}, nil
}

// grpc function to receive the uploaded file from the client
func (n *DataKeeperNode) UploadFile(ctx context.Context, req *pb_d.UploadFileRequest) (*pb_d.UploadFileResponse, error) {


	fileName := req.FileName
	listeningIp := n.IP
	listeningPort := req.Port

	newRecord := &pb_m.Record{FileName: fileName, FilePath: fileName, DataKeeperNodeID: int32(n.ID), Alive: true}

	fmt.Print("Uploading ....\n")
	go func() error {

		n.currBusyPorts -= 1
		conn, err := utils.ReceiveTCP(listeningIp, listeningPort)

		if err != nil {

			fmt.Printf("error in establishing tcp connection: %v\n", err)

			return err
		}

		
		// Read the data
		data, readContentErr := ioutil.ReadAll(conn)

		if readContentErr != nil {
			fmt.Printf("error reading the file content: %v\n", readContentErr)
			return readContentErr
		}

		deserializeError := utils.Deserialize(data, true) 

		
		if deserializeError != nil {
			
			fmt.Printf("error in deserializing the file content: %v\n", deserializeError)

			return deserializeError
		}

		connMaster, errConn := n.ConnectToServer(defaultAdress)

		masterClient := pb_m.NewMasterNodeClient(connMaster)

		if errConn != nil {
			fmt.Printf("error connecting to the master: %v\n", errConn)
			return errConn
		}

		defer connMaster.Close()

		_, notificationErr := masterClient.UploadNotification(context.Background(), &pb_m.UploadNotificationRequest{NewRecord: newRecord})

		if notificationErr != nil {
			fmt.Printf("Failed to notify the master: %v\n", notificationErr)
			return  notificationErr
		}

		n.AddFile(fileName)

		return  nil
	}()


	return &pb_d.UploadFileResponse{Success: true}, nil
}


// grpc function to download the file from the data node
func (n *DataKeeperNode) DownloadFile(ctx context.Context, req *pb_d.DownloadFileRequest) (*pb_d.DownloadFileResponse, error) {

	fileName := req.FileName

	// open the file
	file, errOpenFile := os.Open(fileName)

	if errOpenFile != nil {
		fmt.Printf("Failed to open the file: %v\n", errOpenFile)
		return  &pb_d.DownloadFileResponse{Success: false} , errOpenFile
	}

	defer file.Close()

	// get the file size
	fileInfo, errFileInfo := file.Stat()

	if errFileInfo != nil {
		fmt.Printf("Failed to get the file info: %v\n", errFileInfo)
		return &pb_d.DownloadFileResponse{Success: false}, errFileInfo
	}

	fileSize := fileInfo.Size()

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Printf("Failed to read the file: %v", err)
		return &pb_d.DownloadFileResponse{Success: false}, err
	}

	go func() error {	

		tcpConn, tcpConnErr := utils.ReceiveTCP(n.IP, n.GetCurrentEmptyPort())

		if tcpConnErr != nil {
			fmt.Printf("Failed to establish tcp connection: %v", tcpConnErr)
			return tcpConnErr
		}

		// close the connection
		defer tcpConn.Close()

		fmt.Print("Establish connection\n")

		response := &pb_d.FileResponse{FileName: fileName, FileContent: data, FileSize: fileSize}

		//serialize the request
		errSerialize := utils.Serialize(response, tcpConn)
		fmt.Print("Establish Serialization\n")

		if errSerialize != nil {
			return errSerialize
		}

		

		fmt.Print("The file was sent to the client, successfully :) \n")
		return nil
	}()

	return &pb_d.DownloadFileResponse{Success: true}, nil
}

//TODO: Implement the following functions
//////  Replication Handling //////

// grpc function to receive the file to another data node for replication
func (n *DataKeeperNode) ReceiveFileForReplica(ctx context.Context, req *pb_d.ReceiveFileForReplicaRequest) (*pb_d.ReceiveFileForReplicaRespone, error) {

	//get the destination machine ip and port
	srcMachineIP := req.Ip
	srcMachinePort := req.Port
	rand_seed := req.PortRandomSeed

	//change port
	var portNum int
	var portStr string
	_, err := fmt.Sscan(srcMachinePort, &portNum) // Handle potential parsing errors
	if err != nil {
		log.Printf("Failed to parse port number: %v", err)
	}
	portNum += int(rand_seed)

	// Convert port number back to string
	portStr = fmt.Sprint(portNum)

	log.Println("From replicate file: port number: ", portNum)

	listener, err := net.Listen("tcp", srcMachineIP+":"+portStr)
	if err != nil {
		fmt.Printf("error creating listener: %v\n", err)
		return &pb_d.ReceiveFileForReplicaRespone{Success: false}, err
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("error accepting connection: %v\n", err)
		return &pb_d.ReceiveFileForReplicaRespone{Success: false}, err

	}
	defer conn.Close()

	n.ReceiveFileForReplicaHandler(conn)
	// return n.ReceiveFileForReplicaHandler(conn)
	return &pb_d.ReceiveFileForReplicaRespone{Success: true}, nil
}

// func to handle the file received for replication
func (n *DataKeeperNode) ReceiveFileForReplicaHandler(conn net.Conn) {
	//Read the data
	data, err := ioutil.ReadAll(conn)

	if err != nil {
		fmt.Printf("error reading the file content: %v\n", err)
	}

	receivedFile := &pb_d.UploadFileRequest{}
	err = json.Unmarshal(data, receivedFile)
	if err != nil {
		fmt.Println("Error decoding file:", err)
	}
	print(data)

	// store the file content in the local file
	err = ioutil.WriteFile(receivedFile.FileName, receivedFile.FileContent, 0644)
	if err != nil {
		fmt.Printf("error saving the file: %v\n", err)
	}

}

// GRPC function to send the file to another data node for replication
func (n *DataKeeperNode) ReplicateFile(ctx context.Context, req *pb_d.ReplicaRequest) (*pb_d.NotifyReplicaResponse, error) {

	//get the destination machine ip and port
	sourceIP := req.Ip
	sourcePort := req.Port
	fileName := req.FileName
	rand_seed := req.PortRandomSeed

	//change port
	portNum, err := utils.ConvertStrIntoInt(sourcePort)
	
	if err != nil {
		return &pb_d.NotifyReplicaResponse{Success: false}, err
	}

	portNum += int(rand_seed)

	log.Println("From replicate file: port number: ", portNum)
	// Convert port number back to string
	portStr := fmt.Sprint(portNum)

	//connect to the destination machine tcp
	conn, err := net.Dial("tcp", sourceIP+":"+portStr)
	if err != nil {
		fmt.Printf("error connecting to the destination machine: %v\n", err)
		return &pb_d.NotifyReplicaResponse{Success: false}, err
	}
	defer conn.Close()

	//read the file content
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("error reading the file content: %v\n", err)
		return &pb_d.NotifyReplicaResponse{Success: false}, err
	}

	//convert the file content to json
	file := &pb_d.UploadFileRequest{FileName: fileName, FileContent: fileContent}

	serializedFile, err := json.Marshal(file)

	if err != nil {
		fmt.Printf("error serializing the file: %v\n", err)
		return &pb_d.NotifyReplicaResponse{Success: false}, err
	}

	//send the file content to the destination machine
	_, err = conn.Write(serializedFile)
	if err != nil {
		fmt.Printf("error sending the file content: %v\n", err)
		return &pb_d.NotifyReplicaResponse{Success: false}, err
	}

	return &pb_d.NotifyReplicaResponse{Success: true}, nil
}
