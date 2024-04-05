package client

import (
	"io"
	"log"
	"os"

	"context"
	"fmt"

	pb_d "Distributed_file_system/internals/pb/data_node"
	pb_m "Distributed_file_system/internals/pb/master_node"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	
}

func NewClient() *Client{
	return &Client{};
}

func (c* Client) ConnectToServer(address string) (*grpc.ClientConn, error ) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return conn, err;
}

func (c* Client) AskForUpload(master pb_m.MasterNodeClient) (string, error) {
	
  uploadPort, err := master.AskForUpload(context.Background(), &pb_m.Empty{});
		fmt.Print(uploadPort)

	if err != nil {
		fmt.Printf("Error in getting the upload response %v\n", err)
		return "", err
	}

  return uploadPort.Port, err
}

func (c* Client) UploadFileToServer(master pb_m.MasterNodeClient, filename string) error {

	uploadPort, errGettingPort := c.AskForUpload(master);

	if errGettingPort != nil {
		fmt.Printf("Failed to get the upload port: %v\n", errGettingPort)
		return errGettingPort
	}


	// connect to the datakeeper node
	dataConn, errDataConn := c.ConnectToServer("localhost:"+uploadPort)

	if errDataConn != nil {
		fmt.Printf("Failed to connect to server: %v", errDataConn)
		return errDataConn
	}

	defer dataConn.Close()

	// create a client
	dataNodeClient := pb_d.NewDataNodeClient(dataConn)
	log.Printf("Connected to datakeeper node: %v\n", uploadPort)


	// open the file
	file, errOpenFile := os.Open(filename)

	if errOpenFile != nil {
		fmt.Printf("Failed to open the file: %v\n", errOpenFile)
		return errOpenFile
	}

	defer file.Close()

	// Buffer to read file contents
	buffer := make([]byte, 1024)


	// read the file contents
	for {
		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Failed to read the file: %v\n", err)
			return err
		}

		// send the file contents to the server
		_, errSendFile := dataNodeClient.UploadFile(context.Background(), &pb_d.UploadFileRequest{FileContent: buffer[:n], FileName: filename})

		if errSendFile != nil {
			fmt.Printf("Failed to send the file: %v\n", errSendFile)
			return errSendFile
		}

	}

	return nil;
}


func (c* Client) DownloadFile(masterClient pb_m.MasterNodeClient, filename string) error {
	
	//1. the client should send a request to the master to get the datakeeper node that has the file
	machines, errDownload := masterClient.AskForDownload(context.Background(), &pb_m.AskForDownloadRequest{FileName: filename})

	if errDownload != nil || len(machines.DataKeepers) == 0 {
		fmt.Printf("Failed to get the datakeeper node: %v\n", errDownload)
		return errDownload
	}

	// 2. the client now has list of datakeeper nodes IP and port to request the file from

	// 3. Client MUST request from every port uniformly. (Parallel download is considered a bonus)

	//uniformly means that the client should request from every port in the list
	
	//parallel download means that the client should request from every port in the list at the same time

	// so the client should create a go routine for each port in the list
	
	// and each go routine should request the file from the datakeeper node

	//create a channel to wait for the go routines to finish
	ch := make(chan int)
	for _, machine := range machines.DataKeepers {
		go func(machine *pb_m.DataKeeper) {
			
			// connect to the datakeeper node
			dataConn, errDataConn := c.ConnectToServer(machine.Ip+":"+machine.Port)

			if errDataConn != nil {
				fmt.Printf("Failed to connect to server: %v", errDataConn)
				return
			}

			defer dataConn.Close()

			// create a client
			dataNodeClient := pb_d.NewDataNodeClient(dataConn)
			log.Printf("Connected to datakeeper node: %v", machine)
			
			// request the file size
			fileSize, errGetFile := dataNodeClient.GetFileSize(context.Background(), &pb_d.FileRequest{FileName: filename})
			if errGetFile != nil {
				fmt.Printf("Failed to get the file size: %v", errGetFile)
				return
			}
			//print the file size
			fmt.Printf("File size: %v\n", fileSize.FileSize)

			//request the file
			stream, errDownloadFile := dataNodeClient.DownloadFile(context.Background(), &pb_d.FileRequest{FileName: filename})
			if errDownloadFile != nil {
				fmt.Printf("Failed to download the file: %v\n", errDownloadFile)
				return
			}
			
			//create a file to write the file contents
			os.Chdir("data")

			// TODO: Check while there exists another file with that file name append the word copy to the begining of the name

			file, errCreateFile := os.Create(filename)
			if errCreateFile != nil {
				fmt.Printf("Failed to create the file: %v", errCreateFile)
				return
			}
			defer file.Close()

			// read the file contents from the stream
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

	return nil;
}
