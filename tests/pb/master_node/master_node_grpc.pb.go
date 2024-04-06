// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: master_node.proto

package master_node

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
	MasterNode_RegisterDataNode_FullMethodName   = "/master_node.MasterNode/RegisterDataNode"
	MasterNode_ReceiveFileList_FullMethodName    = "/master_node.MasterNode/ReceiveFileList"
	MasterNode_AskForDownload_FullMethodName     = "/master_node.MasterNode/AskForDownload"
	MasterNode_AskForUpload_FullMethodName       = "/master_node.MasterNode/AskForUpload"
	MasterNode_UploadNotification_FullMethodName = "/master_node.MasterNode/UploadNotification"
	MasterNode_HeartbeatUpdate_FullMethodName    = "/master_node.MasterNode/HeartbeatUpdate"
)

// MasterNodeClient is the client API for MasterNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MasterNodeClient interface {
	// Register a data node to the master
	RegisterDataNode(ctx context.Context, in *RegisterDataNodeRequest, opts ...grpc.CallOption) (*RegisterDataNodeResponse, error)
	// Receive the list of files from the Data node
	ReceiveFileList(ctx context.Context, in *ReceiveFileListRequest, opts ...grpc.CallOption) (*ReceiveFileListResponse, error)
	// Downloading sequence (from master to client)
	AskForDownload(ctx context.Context, in *AskForDownloadRequest, opts ...grpc.CallOption) (*AskForDownloadResponse, error)
	// The client should request from the master to upload a file
	AskForUpload(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AskForUploadResponse, error)
	UploadNotification(ctx context.Context, in *UploadNotificationRequest, opts ...grpc.CallOption) (*UploadNotificationResponse, error)
	// Check if the data node is alive
	HeartbeatUpdate(ctx context.Context, in *HeartbeatUpdateRequest, opts ...grpc.CallOption) (*HeartbeatUpdateResponse, error)
}

type masterNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterNodeClient(cc grpc.ClientConnInterface) MasterNodeClient {
	return &masterNodeClient{cc}
}

func (c *masterNodeClient) RegisterDataNode(ctx context.Context, in *RegisterDataNodeRequest, opts ...grpc.CallOption) (*RegisterDataNodeResponse, error) {
	out := new(RegisterDataNodeResponse)
	err := c.cc.Invoke(ctx, MasterNode_RegisterDataNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterNodeClient) ReceiveFileList(ctx context.Context, in *ReceiveFileListRequest, opts ...grpc.CallOption) (*ReceiveFileListResponse, error) {
	out := new(ReceiveFileListResponse)
	err := c.cc.Invoke(ctx, MasterNode_ReceiveFileList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterNodeClient) AskForDownload(ctx context.Context, in *AskForDownloadRequest, opts ...grpc.CallOption) (*AskForDownloadResponse, error) {
	out := new(AskForDownloadResponse)
	err := c.cc.Invoke(ctx, MasterNode_AskForDownload_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterNodeClient) AskForUpload(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AskForUploadResponse, error) {
	out := new(AskForUploadResponse)
	err := c.cc.Invoke(ctx, MasterNode_AskForUpload_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterNodeClient) UploadNotification(ctx context.Context, in *UploadNotificationRequest, opts ...grpc.CallOption) (*UploadNotificationResponse, error) {
	out := new(UploadNotificationResponse)
	err := c.cc.Invoke(ctx, MasterNode_UploadNotification_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterNodeClient) HeartbeatUpdate(ctx context.Context, in *HeartbeatUpdateRequest, opts ...grpc.CallOption) (*HeartbeatUpdateResponse, error) {
	out := new(HeartbeatUpdateResponse)
	err := c.cc.Invoke(ctx, MasterNode_HeartbeatUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterNodeServer is the server API for MasterNode service.
// All implementations must embed UnimplementedMasterNodeServer
// for forward compatibility
type MasterNodeServer interface {
	// Register a data node to the master
	RegisterDataNode(context.Context, *RegisterDataNodeRequest) (*RegisterDataNodeResponse, error)
	// Receive the list of files from the Data node
	ReceiveFileList(context.Context, *ReceiveFileListRequest) (*ReceiveFileListResponse, error)
	// Downloading sequence (from master to client)
	AskForDownload(context.Context, *AskForDownloadRequest) (*AskForDownloadResponse, error)
	// The client should request from the master to upload a file
	AskForUpload(context.Context, *Empty) (*AskForUploadResponse, error)
	UploadNotification(context.Context, *UploadNotificationRequest) (*UploadNotificationResponse, error)
	// Check if the data node is alive
	HeartbeatUpdate(context.Context, *HeartbeatUpdateRequest) (*HeartbeatUpdateResponse, error)
	mustEmbedUnimplementedMasterNodeServer()
}

// UnimplementedMasterNodeServer must be embedded to have forward compatible implementations.
type UnimplementedMasterNodeServer struct {
}

func (UnimplementedMasterNodeServer) RegisterDataNode(context.Context, *RegisterDataNodeRequest) (*RegisterDataNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDataNode not implemented")
}
func (UnimplementedMasterNodeServer) ReceiveFileList(context.Context, *ReceiveFileListRequest) (*ReceiveFileListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveFileList not implemented")
}
func (UnimplementedMasterNodeServer) AskForDownload(context.Context, *AskForDownloadRequest) (*AskForDownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AskForDownload not implemented")
}
func (UnimplementedMasterNodeServer) AskForUpload(context.Context, *Empty) (*AskForUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AskForUpload not implemented")
}
func (UnimplementedMasterNodeServer) UploadNotification(context.Context, *UploadNotificationRequest) (*UploadNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadNotification not implemented")
}
func (UnimplementedMasterNodeServer) HeartbeatUpdate(context.Context, *HeartbeatUpdateRequest) (*HeartbeatUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartbeatUpdate not implemented")
}
func (UnimplementedMasterNodeServer) mustEmbedUnimplementedMasterNodeServer() {}

// UnsafeMasterNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterNodeServer will
// result in compilation errors.
type UnsafeMasterNodeServer interface {
	mustEmbedUnimplementedMasterNodeServer()
}

func RegisterMasterNodeServer(s grpc.ServiceRegistrar, srv MasterNodeServer) {
	s.RegisterService(&MasterNode_ServiceDesc, srv)
}

func _MasterNode_RegisterDataNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterDataNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterNodeServer).RegisterDataNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterNode_RegisterDataNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterNodeServer).RegisterDataNode(ctx, req.(*RegisterDataNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterNode_ReceiveFileList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveFileListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterNodeServer).ReceiveFileList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterNode_ReceiveFileList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterNodeServer).ReceiveFileList(ctx, req.(*ReceiveFileListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterNode_AskForDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AskForDownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterNodeServer).AskForDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterNode_AskForDownload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterNodeServer).AskForDownload(ctx, req.(*AskForDownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterNode_AskForUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterNodeServer).AskForUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterNode_AskForUpload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterNodeServer).AskForUpload(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterNode_UploadNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterNodeServer).UploadNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterNode_UploadNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterNodeServer).UploadNotification(ctx, req.(*UploadNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterNode_HeartbeatUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterNodeServer).HeartbeatUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterNode_HeartbeatUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterNodeServer).HeartbeatUpdate(ctx, req.(*HeartbeatUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MasterNode_ServiceDesc is the grpc.ServiceDesc for MasterNode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MasterNode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "master_node.MasterNode",
	HandlerType: (*MasterNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterDataNode",
			Handler:    _MasterNode_RegisterDataNode_Handler,
		},
		{
			MethodName: "ReceiveFileList",
			Handler:    _MasterNode_ReceiveFileList_Handler,
		},
		{
			MethodName: "AskForDownload",
			Handler:    _MasterNode_AskForDownload_Handler,
		},
		{
			MethodName: "AskForUpload",
			Handler:    _MasterNode_AskForUpload_Handler,
		},
		{
			MethodName: "UploadNotification",
			Handler:    _MasterNode_UploadNotification_Handler,
		},
		{
			MethodName: "HeartbeatUpdate",
			Handler:    _MasterNode_HeartbeatUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "master_node.proto",
}
