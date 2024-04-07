package dk

import (
	pb_d "Distributed_file_system/internals/pb/data_node"
	pb_m "Distributed_file_system/internals/pb/master_node"
	"encoding/json"
	"io/ioutil"
	"net"

	"Distributed_file_system/internals/utils"
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
	//convert the file size from string to int64
	// fileSize := req.FileSize
	//Listen to the client

	go func() {
		//create a listener
		listener, err := net.Listen("tcp", clintIp+":"+clintPort)
		if err != nil {
			fmt.Printf("error creating listener: %v\n", err)
			return
		}
		defer listener.Close()

		//accept the connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
			return
		}
		defer conn.Close()

		//Read the data
		data, err := ioutil.ReadAll(conn)

		if err != nil {
			fmt.Printf("error reading the file content: %v\n", err)
			return
		}

		receivedFile := &pb_d.UploadFileRequest{}
		err = json.Unmarshal(data, receivedFile)
		if err != nil {
			fmt.Println("Error decoding file:", err)
			return
		}
		print(data)

		// store the file content in the local file
		err = ioutil.WriteFile(receivedFile.FileName, receivedFile.FileContent, 0644)
		if err != nil {
			fmt.Printf("error saving the file: %v\n", err)
			return
		}

		//create a file buffer using the file size
		// fileBuffer := make([]byte, fileSize)

		// //read the file content from the client
		// totalRead := 0
		// for totalRead < int(fileSize) {
		// 	n, err := conn.Read(fileBuffer[totalRead:])
		// 	if err != nil {
		// 		if err == io.EOF {
		// 			break
		// 		}
		// 		fmt.Printf("error reading the file content: %v\n", err)
		// 		return
		// 	}
		// 	totalRead += n
		// }

		//write the file content to the local file
		// err = utils.SaveFile(fileName, fileBuffer)
		// if err != nil {
		// 	fmt.Printf("error saving the file: %v\n", err)
		// 	return
		// }

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

// grpc function to send the file to another data node for replication
func (n *DataKeeperNode) ReceiveFileForReplica(ctx context.Context, req *pb_d.ReceiveFileForReplicaRequest) (*pb_d.ReceiveFileForReplicaRespone, error) {

	fileName := req.FileName
	fileContent := req.FileContent
	fmt.Printf("Received request to replicate file: %s\n", fileName)

	err := utils.SaveFile(fileName, []byte(fileContent))

	if err != nil {
		fmt.Printf("error in saving the file locally, %v", err)
		return &pb_d.ReceiveFileForReplicaRespone{Success: false}, err
	}

	n.AddFile(fileName)

	return &pb_d.ReceiveFileForReplicaRespone{Success: true}, nil

}

// function to replicate the file to another data node in the network
func (n *DataKeeperNode) ReplicateFile(destinationMachineIP string, destinationMachinePort string, fileName string) error {

	conn, err := grpc.Dial(destinationMachineIP+":"+destinationMachinePort, grpc.WithInsecure())

	if err != nil {
		return fmt.Errorf("error connecting to the destination machine: %v", err)
	}

	defer conn.Close()

	// client := pb_d.NewDataNodeClient(conn)

	log.Printf("debug: filename: '%s'", fileName)
	// Open the file
	//TODO : fix the sys can't find the file issue
	//check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %v", err)
	}
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	defer file.Close()

	// Buffer to read file contents
	// buffer := make([]byte, 1024)

	// Read file contents and send to client
	// for {
	// 	bytesRead, err := file.Read(buffer)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		return fmt.Errorf("error reading file: %v", err)
	// 	}

	// Send file chunk to client
	//TODO: Send the file content to the destination machine
	// _, errSendFile := client.UploadFile(context.Background(), &pb_d.UploadFileRequest{FileContent: buffer[:bytesRead], FileName: fileName})

	// if errSendFile != nil {
	// 	return fmt.Errorf("error sending file chunk: %v", errSendFile)
	// }

	// }

	// print the success message
	fmt.Printf("File replicated successfully to %s:%s\n", destinationMachineIP, destinationMachinePort)

	return nil
}
