// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: data_node.proto

package data_node

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DataNode_UploadFile_FullMethodName            = "/data_node.DataNode/UploadFile"
	DataNode_GetFileSize_FullMethodName           = "/data_node.DataNode/GetFileSize"
	DataNode_DownloadFile_FullMethodName          = "/data_node.DataNode/DownloadFile"
	DataNode_ReceiveFileForReplica_FullMethodName = "/data_node.DataNode/ReceiveFileForReplica"
)

// DataNodeClient is the client API for DataNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataNodeClient interface {
	// Receive the file from the client
	UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error)
	// RPC for getting the file size
	GetFileSize(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error)
	// RPC for downloading a file
	DownloadFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (DataNode_DownloadFileClient, error)
	// RPC for replicating a file (sending a file to another data node)
	ReceiveFileForReplica(ctx context.Context, in *ReceiveFileForReplicaRequest, opts ...grpc.CallOption) (*ReceiveFileForReplicaRespone, error)
}

type dataNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewDataNodeClient(cc grpc.ClientConnInterface) DataNodeClient {
	return &dataNodeClient{cc}
}

func (c *dataNodeClient) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error) {
	out := new(UploadFileResponse)
	err := c.cc.Invoke(ctx, DataNode_UploadFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) GetFileSize(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, DataNode_GetFileSize_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) DownloadFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (DataNode_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &DataNode_ServiceDesc.Streams[0], DataNode_DownloadFile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dataNodeDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DataNode_DownloadFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type dataNodeDownloadFileClient struct {
	grpc.ClientStream
}

func (x *dataNodeDownloadFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dataNodeClient) ReceiveFileForReplica(ctx context.Context, in *ReceiveFileForReplicaRequest, opts ...grpc.CallOption) (*ReceiveFileForReplicaRespone, error) {
	out := new(ReceiveFileForReplicaRespone)
	err := c.cc.Invoke(ctx, DataNode_ReceiveFileForReplica_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataNodeServer is the server API for DataNode service.
// All implementations must embed UnimplementedDataNodeServer
// for forward compatibility
type DataNodeServer interface {
	// Receive the file from the client
	UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error)
	// RPC for getting the file size
	GetFileSize(context.Context, *FileRequest) (*FileResponse, error)
	// RPC for downloading a file
	DownloadFile(*FileRequest, DataNode_DownloadFileServer) error
	// RPC for replicating a file (sending a file to another data node)
	ReceiveFileForReplica(context.Context, *ReceiveFileForReplicaRequest) (*ReceiveFileForReplicaRespone, error)
	mustEmbedUnimplementedDataNodeServer()
}

// UnimplementedDataNodeServer must be embedded to have forward compatible implementations.
type UnimplementedDataNodeServer struct {
}

func (UnimplementedDataNodeServer) UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedDataNodeServer) GetFileSize(context.Context, *FileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileSize not implemented")
}
func (UnimplementedDataNodeServer) DownloadFile(*FileRequest, DataNode_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedDataNodeServer) ReceiveFileForReplica(context.Context, *ReceiveFileForReplicaRequest) (*ReceiveFileForReplicaRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveFileForReplica not implemented")
}
func (UnimplementedDataNodeServer) mustEmbedUnimplementedDataNodeServer() {}

// UnsafeDataNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataNodeServer will
// result in compilation errors.
type UnsafeDataNodeServer interface {
	mustEmbedUnimplementedDataNodeServer()
}

func RegisterDataNodeServer(s grpc.ServiceRegistrar, srv DataNodeServer) {
	s.RegisterService(&DataNode_ServiceDesc, srv)
}

func _DataNode_UploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).UploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataNode_UploadFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).UploadFile(ctx, req.(*UploadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_GetFileSize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).GetFileSize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataNode_GetFileSize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).GetFileSize(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataNodeServer).DownloadFile(m, &dataNodeDownloadFileServer{stream})
}

type DataNode_DownloadFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type dataNodeDownloadFileServer struct {
	grpc.ServerStream
}

func (x *dataNodeDownloadFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

func _DataNode_ReceiveFileForReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveFileForReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).ReceiveFileForReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataNode_ReceiveFileForReplica_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).ReceiveFileForReplica(ctx, req.(*ReceiveFileForReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataNode_ServiceDesc is the grpc.ServiceDesc for DataNode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataNode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data_node.DataNode",
	HandlerType: (*DataNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadFile",
			Handler:    _DataNode_UploadFile_Handler,
		},
		{
			MethodName: "GetFileSize",
			Handler:    _DataNode_GetFileSize_Handler,
		},
		{
			MethodName: "ReceiveFileForReplica",
			Handler:    _DataNode_ReceiveFileForReplica_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadFile",
			Handler:       _DataNode_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "data_node.proto",
}
