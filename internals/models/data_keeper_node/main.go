package main

import (
	pb "Distributed_file_system/internals/proto/dfs"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

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

func main() {
	// Create a new DataKeeperNode instance
	node := NewDataKeeperNode(1, "Node1", "localhost", 8080, []string{"2MB"})
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect to server: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewHeartbeatUpdateClient(conn)
	log.Printf("Connected to server: %v", "localhost:8080")
	done := make(chan struct{})
	defer close(done) // close the done channel when the main function returns
	// thread to send heartbeat to server
	go func() {
		// Send a heartbeat to the server every 5 seconds
		for {
			res, err := client.HeartbeatUpdate(context.Background(), &pb.Heartbeat{Id: int32(node.ID), Ip: node.IP, Port: int32(node.Port)})
			if err != nil {
				fmt.Printf("Failed to send heartbeat to server: %v", err)
			}
			fmt.Printf("Heartbeat response: %v \n", res)

			// Sleep for 5 seconds
			time.Sleep(1 * time.Second)
		}
		// Close the done channel when the thread exits
		done <- struct{}{}
	}()
	// Wait for the thread to exit
	<-done

}
