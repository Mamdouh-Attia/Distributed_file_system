package dk

import (
	pb_d "Distributed_file_system/internals/pb/data_node"
	pb_m "Distributed_file_system/internals/pb/master_node"

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
	FileName string
	FilePath string
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

func (n* DataKeeperNode) ConnectToServer(address string) (*grpc.ClientConn, error ) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return conn, err;
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

func (n* DataKeeperNode) GetFileByName(filename string) string {
	for _, file := range n.Files {
		if file == filename {
			return file
			
		}
	}
	return ""
}


// grpc function to receive the uploaded file from the client 
func (n* DataKeeperNode)  UploadFile (ctx context.Context, req *pb_d.UploadFileRequest) (*pb_d.UploadFileResponse, error)  {

	fileName := req.FileName
	fileContent := req.FileContent
	fmt.Print("Received request to upload file: %s\n", fileName)
	err := utils.SaveFile(fileName, []byte(fileContent))

	if err != nil {
		fmt.Print("error in saving the file locally, %v", err)
		return &pb_d.UploadFileResponse{Success: false}, err
	}
	
	n.AddFile(fileName)
	connMaster, errConn := n.ConnectToServer(defaultAdress)

	if errConn != nil {
		fmt.Print("error connecting to the master: %v\n", errConn)
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
