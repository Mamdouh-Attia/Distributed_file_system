// imports
package main

import (
	dk "Distributed_file_system/internals/pb/data_node"
	mt "Distributed_file_system/internals/pb/master_node"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type Client struct{}

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

const (
	address         = "localhost:8080"
	defaultFilename = "data\\2MB.mp4"
)

func main() {

	// the client should connect to the master
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		fmt.Printf("Failed to connect to server: %v", err)
		return
	}

	defer conn.Close()

	fmt.Printf("Connected to server: %v\n", "localhost:8080")
	//keep the connection open
	//channl

	// ch := make(chan int)

	//downlaod steps:

	//1. the client should send a request to the master to get the datakeeper node that has the file
	masterClient := mt.NewMasterNodeClient(conn)
	machines, err := masterClient.DownloadRequest(context.Background(), &mt.DownloadFileRequest{FileName: defaultFilename})

	if err != nil {
		fmt.Printf("Failed to get the datakeeper node: %v", err)
		return
	}
	// 2. the client now has list of datakeeper nodes IP and port to request the file from
	//print the datakeeper node
	fmt.Printf("Datakeeper nodes for file %v: %v", defaultFilename, machines.DataKeepers)

	//3. 3.	Client MUST request from every port uniformly. (Parallel download is considered a bonus)

	//uniformly means that the client should request from every port in the list
	//parallel download means that the client should request from every port in the list at the same time
	// so the client should create a go routine for each port in the list
	// and each go routine should request the file from the datakeeper node

	//create a channel to wait for the go routines to finish
	ch := make(chan int)
	for _, machine := range machines.DataKeepers {
		go func(machine *mt.DataKeeper) {
			//connect to the datakeeper node
			conn, err := grpc.Dial(machine.Ip+":"+machine.Port, grpc.WithInsecure())
			if err != nil {
				fmt.Printf("Failed to connect to server: %v", err)
				return
			}
			defer conn.Close()
			//create a client
			dataNodeClient := dk.NewDataNodeClient(conn)
			log.Printf("Connected to datakeeper node: %v", machine)
			//request the file size
			fileSize, err := dataNodeClient.GetFileSize(context.Background(), &dk.FileRequest{FileName: defaultFilename})
			if err != nil {
				fmt.Printf("Failed to get the file size: %v", err)
				return
			}
			//print the file size
			fmt.Printf("File size: %v", fileSize.FileSize)
			//request the file
			stream, err := dataNodeClient.DownloadFile(context.Background(), &dk.FileRequest{FileName: defaultFilename})
			if err != nil {
				fmt.Printf("Failed to download the file: %v", err)
				return
			}
			//create a file to write the file contents
			splitted := strings.Split(defaultFilename, "\\")
			fileName := splitted[len(splitted)-1]
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Failed to create the file: %v", err)
				return
			}
			defer file.Close()

			//read the file contents from the stream
			for {
				chunk, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("Error while receiving chunk: %v", err)
					return
				}
				//write the chunk to the file
				if _, err := file.Write(chunk.Data); err != nil {
					fmt.Printf("Failed to write to the file: %v", err)
					return
				}
			}
			//close the channel
			ch <- 1
		}(machine)
	}
	//wait for the go routines to finish
	for i := 0; i < len(machines.DataKeepers); i++ {
		<-ch
	}
	//print that the download is finished
	fmt.Println("Download finished")

}
