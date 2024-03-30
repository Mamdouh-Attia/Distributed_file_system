// imports
package client

import (
	pb "Distributed_file_system/internals/pb"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type Client struct {}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ConnectToServer() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Client) CloseConnection(conn *grpc.ClientConn) {
	conn.Close()
}


func main () {

	// the client should connect to the master
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		fmt.Printf("Failed to connect to server: %v", err)
		return
	}

	defer conn.Close()

	fmt.Printf("Connected to server: %v", "localhost:8080")

	// the client should send a request to the master to upload a file
	dfs := pb.NewDistributedFileSystemClient(conn)
	uploadPort, _ := dfs.UploadRequest(context.Background(), &pb.Empty{})

	// the client should connect to the datakeeper node
	conn, err = grpc.Dial("localhost:"+uploadPort.Port, grpc.WithInsecure())

	if err != nil {
		fmt.Printf("Failed to connect to server: %v", err)
		return
	}

	defer conn.Close()

	fmt.Printf("Connected to server: %v", "localhost:"+uploadPort.Port)

	// the client should send the file to the datakeeper node

	// read the file from the client
	// open the file
	// send the file to the datakeeper node
	// close the file
	// close the connection to the datakeeper node

	




}
