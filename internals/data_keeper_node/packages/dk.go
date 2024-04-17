package dk

import (
	pb_d "Distributed_file_system/internals/pb/data_node"
	pb_m "Distributed_file_system/internals/pb/master_node"
	"Distributed_file_system/internals/utils"
	"encoding/json"
	"io/ioutil"
	"net"

	"context"
	"os"

	"fmt"
	"io"
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

// grpc function to receive the uploaded file from the client
func (n *DataKeeperNode) UploadFile(ctx context.Context, req *pb_d.UploadFileRequest) (*pb_d.UploadFileResponse, error) {

	fileName := req.FileName
	clintIp := req.Ip
	clintPort := req.Port


	go func() {

		data, err := utils.ReceiveTCP(clintIp, clintPort)

		if err != nil {
			return
		}


		receivedFile := &pb_d.UploadFileRequest{}
		
		err = utils.Deserialize(data, receivedFile, receivedFile.FileName, receivedFile.FileContent) 

		if err != nil {
			return
		}

	}()

	n.AddFile(fileName)
	connMaster, errConn := n.ConnectToServer(defaultAdress)

	if errConn != nil {
		fmt.Printf("error connecting to the master: %v\n", errConn)
		return &pb_d.UploadFileResponse{Success: false}, errConn
	}

	defer connMaster.Close()

	masterClient := pb_m.NewMasterNodeClient(connMaster)

	_, notificationErr := masterClient.UploadNotification(context.Background(), &pb_m.UploadNotificationRequest{NewRecord: &pb_m.Record{FileName: fileName, FilePath: fileName, DataKeeperNodeID: int32(n.ID), Alive: true}})

	if notificationErr != nil {
		fmt.Printf("Failed to notify the master: %v\n", notificationErr)
		return &pb_d.UploadFileResponse{Success: false}, notificationErr
	}

	return &pb_d.UploadFileResponse{Success: true}, nil
}

// GetFiles returns the list of files stored on the node.
func (n *DataKeeperNode) GetFileSize(ctx context.Context, req *pb_d.FileRequest) (*pb_d.FileResponse, error) {
	log.Printf("Received request to get file size for file: %s", req.FileName)

	// Open the file
	file, err := os.Open(req.FileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %v", err)
	}

	// Send file size to client
	return &pb_d.FileResponse{FileName: req.FileName, FileSize: fileInfo.Size()}, nil
}

func (n *DataKeeperNode) DownloadFile(req *pb_d.FileRequest, stream pb_d.DataNode_DownloadFileServer) error {

	file, err := os.Open(req.FileName)

	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Buffer to read file contents
	buffer := make([]byte, 1024)

	// Read file contents and send to client
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading file: %v", err)
		}
		// Send file chunk to client
		if err := stream.Send(&pb_d.FileChunk{Data: buffer[:bytesRead]}); err != nil {
			return fmt.Errorf("error sending file chunk: %v", err)
		}
	}

	return nil
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
