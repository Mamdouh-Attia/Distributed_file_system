package dk

import (
	pb_d "Distributed_file_system/internals/pb/data_node"
	"context"
	"os"

	"fmt"
	"io"
	"log"
)


type File struct {
	FileName string
	FilePath string
	FileContent []byte
}


const (
	dataDir = "./data/"
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
// func (n* DataKeeperNode)  ReceiveUploadedFile(req* pb_d.ReceiveUploadedFileRequest) (*pb_d.ReceiveUploadedFileResponse, error) {

// 	fileName := req.FileName
// 	fileContent := req.FileContent

// 	err := utils.SaveFile("data", fileName, []byte(fileContent))

// 	if err != nil {
// 		n.AddFile(fileName)
// 		return &pb_d.ReceiveUploadedFileResponse{Success: true}, nil
// 	}

// 	return &pb_d.ReceiveUploadedFileResponse{Success: false}, err
// }

// GetFiles returns the list of files stored on the node.
func (n *DataKeeperNode) GetFileSize(ctx context.Context, req *pb_d.FileRequest) (*pb_d.FileResponse, error) {
	log.Printf("Received request to get file size for file: %s", req.FileName)

	// Open the file
	file, err := os.Open(dataDir+req.FileName)
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

	file, err := os.Open(dataDir+req.FileName)

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
