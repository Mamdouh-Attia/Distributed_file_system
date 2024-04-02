package dk

import (
	pb "Distributed_file_system/internals/pb/data_node"
	"context"

	// "errors"
	"fmt"
	"io"
	"log"
	"os"
)

// DataKeeperNode represents a node responsible for keeping track of files.
type DataKeeperNode struct {
	pb.UnimplementedDataNodeServer
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
func (n *DataKeeperNode) AddFile(filename string) {
	n.Files = append(n.Files, filename)
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

// // gRPC function to get file size
// func (n *DataKeeperNode) GetFileSize(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
// 	log.Printf("Received request to get file size for file: %s", req.FileName)

// 	// Open the requested file
// 	file, err := os.Open(".\\" + req.FileName)
// 	print("file", file)
// 	if err != nil {
// 		return nil, errors.New("error opening file")
// 	}
// 	print("file opened")
// 	defer file.Close()

// 	// Read file contents
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		return nil, errors.New("error reading file")
// 	}

// 	// Send file size to client
// 	return &pb.FileResponse{FileName: req.FileName, FileSize: int64(fileInfo.Size())}, nil
// }

// // grpc function to download the file from the datakeeper node to the client
// func (n *DataKeeperNode) DownloadFile(req *pb.FileRequest, stream pb.DataNode_DownloadFileServer) error {
// 	// Open the requested file
// 	file, err := os.Open("data/" + req.FileName)
// 	if err != nil {
// 		return errors.New("error opening file")
// 	}
// 	defer file.Close()

// 	// Buffer to read file contents
// 	buffer := make([]byte, 1024)

// 	// Read file contents and send to client
// 	for {
// 		bytesRead, err := file.Read(buffer)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			return errors.New("error reading file")
// 		}
// 		// Send file chunk to client
// 		if err := stream.Send(&pb.FileChunk{Data: buffer[:bytesRead]}); err != nil {
// 			return errors.New("error sending file chunk")
// 		}
// 	}

// 	return nil
// }

func (n *DataKeeperNode) GetFileSize(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	log.Printf("Received request to get file size for file: %s", req.FileName)

	// Open the requested file
	file, err := os.Open(".\\" + req.FileName)
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
	return &pb.FileResponse{FileName: req.FileName, FileSize: fileInfo.Size()}, nil
}

func (n *DataKeeperNode) DownloadFile(req *pb.FileRequest, stream pb.DataNode_DownloadFileServer) error {
	// Open the requested file
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
		if err := stream.Send(&pb.FileChunk{Data: buffer[:bytesRead]}); err != nil {
			return fmt.Errorf("error sending file chunk: %v", err)
		}
	}

	return nil
}
