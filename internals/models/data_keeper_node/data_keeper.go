package main

// DataKeeperNode represents a node responsible for keeping track of files.
type DataKeeperNode struct {
	ID    int      // Unique identifier for the node
	Name  string   // Name of the node
	IP    string   // IP address of the node
	Port  int      // Port number the node is listening on
	Files []string // List of filenames stored on the node
}

// NewDataKeeperNode creates a new DataKeeperNode instance with the provided parameters.
func NewDataKeeperNode(id int, name, ip string, port int, files []string) *DataKeeperNode {
	return &DataKeeperNode{
		ID:    id,
		Name:  name,
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
