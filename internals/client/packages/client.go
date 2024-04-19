package client

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"context"
	"fmt"

	pb_d "Distributed_file_system/internals/pb/data_node"
	pb_m "Distributed_file_system/internals/pb/master_node"
	"Distributed_file_system/internals/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ConnectToServer(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return conn, err
}

func (c *Client) AskForUpload(master pb_m.MasterNodeClient) (string, string, error) {

	response, err := master.AskForUpload(context.Background(), &pb_m.Empty{})
	fmt.Print(response)

	if err != nil {
		fmt.Printf("Error in getting the upload response %v\n", err)
		return "", "", err
	}

	return response.Port, response.Ip, nil
}
func (m *Client) ReplicateFile(ctx context.Context, req *pb_d.ReplicaRequest) (*pb_d.NotifyReplicaResponse, error) {
	log.Printf("Replicating file: %v", req.FileName)
	return &pb_d.NotifyReplicaResponse{Success: true}, nil

}
func (c *Client) UploadFileToServer(master pb_m.MasterNodeClient, fileName string) error {

	nodePort, nodeIP, errGettingPort := c.AskForUpload(master)

	if errGettingPort != nil {
		fmt.Printf("Failed to get the upload port: %v\n", errGettingPort)
		return errGettingPort
	}

	// connect to the datakeeper node
	dataConn, errDataConn := c.ConnectToServer(nodeIP + ":" + nodePort)

	if errDataConn != nil {
		fmt.Printf("Failed to connect to server: %v", errDataConn)
		return errDataConn
	}

	defer dataConn.Close()

	// create a client
	dataNodeClient := pb_d.NewDataNodeClient(dataConn)
	log.Printf("Connected to datakeeper node: %v\n", nodePort)

	// open the file
	file, errOpenFile := os.Open(fileName)

	if errOpenFile != nil {
		fmt.Printf("Failed to open the file: %v\n", errOpenFile)
		return errOpenFile
	}

	defer file.Close()

	// get the file size
	fileInfo, errFileInfo := file.Stat()

	if errFileInfo != nil {
		fmt.Printf("Failed to get the file info: %v\n", errFileInfo)
		return errFileInfo
	}

	fileSize := fileInfo.Size()

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Printf("Failed to read the file: %v", err)
		return err
	}

	uploadPortResponse, errPortRequest := dataNodeClient.AskForFreePort(context.Background(), &pb_d.AskForFreePortRequest{})

	if errPortRequest != nil {
		fmt.Print("Error getting the free port for uploading")
		return errPortRequest
	}

	uploadPort := uploadPortResponse.Port
	// grpc call the client
	_, errUpload := dataNodeClient.UploadFile(context.Background(), &pb_d.UploadFileRequest{FileName: fileName, FileSize: fileSize, Port: uploadPort, FileContent: data})

	if errUpload != nil {
		fmt.Printf("Failed to upload the file: %v", errUpload)
		return errUpload
	}

	// connect to the datakeeper node as TCP
	tcpConn, errTcpConn := utils.SendTCP(nodeIP, uploadPort)

	if errTcpConn != nil {
		return errTcpConn
	}
	defer tcpConn.Close()

	//convert file to bytes
	//instantiate request
	request := &pb_d.UploadFileRequest{FileName: fileName, FileSize: fileSize, Port: uploadPort, FileContent: data}

	//serialize the request
	errSerialize := utils.Serialize(request, tcpConn)


	return errSerialize
}

func (c *Client) DownloadFile(masterClient pb_m.MasterNodeClient, fileName string) error {

	//1. the client should send a request to the master to get the datakeeper node that has the file
	machines, errAskForDownload := masterClient.AskForDownload(context.Background(), &pb_m.AskForDownloadRequest{FileName: fileName})

	if errAskForDownload != nil || len(machines.DataKeepers) == 0 {
		fmt.Printf("Failed to get the datakeeper node: %v\n", errAskForDownload)
		return errAskForDownload
	}

	// 2. the client now has list of datakeeper nodes IP and port to request the file from
	
	// For now, we will download from one machine
	// We will implement the download from multiple machines later

	// connect to the datakeeper node
	chosenDataNode := machines.DataKeepers[rand.Intn(len(machines.DataKeepers))]

	dataConn, errDataConn := c.ConnectToServer(chosenDataNode.Ip + ":" + chosenDataNode.Port)

	if errDataConn != nil {
		fmt.Printf("Failed to connect to server: %v", errDataConn)
		return errDataConn
	}

	defer dataConn.Close()

	// create a client
	dataNodeClient := pb_d.NewDataNodeClient(dataConn)
	log.Printf("Connected to datakeeper node: %v\n", chosenDataNode.Port)

	// ask the data keeper node for empty port for tcp connection
	emptyPortResponse, errEmptyPort := dataNodeClient.AskForFreePort(context.Background(), &pb_d.AskForFreePortRequest{})

	if errEmptyPort != nil {
		fmt.Print("Error getting the empty port")
		return errEmptyPort
	}

	_, errDownload := dataNodeClient.DownloadFile(context.Background(), &pb_d.DownloadFileRequest{FileName: fileName})

	if errDownload != nil {
		fmt.Printf("Error downloading the file %v, from the dataNodeKeeper: %v\n", fileName, chosenDataNode.Id)

		return errDownload
	}


	emptyPort := emptyPortResponse.Port

	// connect to the datakeeper node as TCP
	tcpConn, errTcpConn := utils.SendTCP(chosenDataNode.Ip, emptyPort)

	if errTcpConn != nil {
		return errTcpConn
	}
	defer tcpConn.Close()


	
	// Read the data
	data, readContentErr := ioutil.ReadAll(tcpConn)
	
	if readContentErr != nil {
		fmt.Printf("error reading the file content: %v\n", readContentErr)
		return readContentErr
	}

	deserializeError := utils.Deserialize(data, false) 

	
	if deserializeError != nil {
		
		fmt.Printf("error in deserializing the file content: %v\n", deserializeError)

		return deserializeError
	}
	
	fmt.Print("The file has been downloaded successfully :)\n")

	return nil
	

}
