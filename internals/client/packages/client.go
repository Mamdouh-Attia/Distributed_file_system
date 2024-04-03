package client

import (
	mt "Distributed_file_system/internals/master_tracker/packages"
	pb_d "Distributed_file_system/internals/pb/master_node"
	"Distributed_file_system/pb"
	"context"

	// "log"
	// "time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	
}

func NewClient() *Client{
	return &Client{};
}

func (c* Client) ConnectToServer(server string) (*grpc.ClientConn, error ) {
	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return conn, err;
}

func (c* Client) AskForUpload(master* mt.Master) (string, error) {
  uploadPort, err := master.AskForUploadRequest(context.Background(), &pb.Empty{});

  return uploadPort.Port, err
}

func (c* Client) UploadFileToServer(port string) error {

	
	conn, err := c.ConnectToServer("" + port)
	
	dataNodeClient := pb_d.NewMasterNodeClient(conn)

	return nil
}

