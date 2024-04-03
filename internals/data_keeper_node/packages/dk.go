package dk

import (
	pb_d "Distributed_file_system/internals/pb/data_node_keeper"
	utils "Distributed_file_system/internals/utils"
	// pb_m "Distributed_file_system/internals/pb/master_node"
)


type File struct {
	FileName string
	FilePath string
	FileContent []byte
}

// DataKeeperNode represents a node responsible for keeping track of files.
type DataKeeperNode struct {
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
func (n* DataKeeperNode)  ReceiveUploadedFile(req* pb_d.ReceiveUploadedFileRequest) (*pb_d.ReceiveUploadedFileResponse, error) {

	fileName := req.FileName
	fileContent := req.FileContent

	err := utils.SaveFile("data", fileName, []byte(fileContent))

	if err != nil {
		n.AddFile(fileName)
		return &pb_d.ReceiveUploadedFileResponse{Success: true}, nil
	}

	return &pb_d.ReceiveUploadedFileResponse{Success: false}, err

}
