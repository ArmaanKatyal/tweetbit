// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: proto/followers.proto

package fanoutService

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
	FollowerService_GetUserFollowers_FullMethodName = "/fanout.FollowerService/GetUserFollowers"
)

// FollowerServiceClient is the client API for FollowerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FollowerServiceClient interface {
	GetUserFollowers(ctx context.Context, in *GetFollowerRequest, opts ...grpc.CallOption) (*GetFollowerResponse, error)
}

type followerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFollowerServiceClient(cc grpc.ClientConnInterface) FollowerServiceClient {
	return &followerServiceClient{cc}
}

func (c *followerServiceClient) GetUserFollowers(ctx context.Context, in *GetFollowerRequest, opts ...grpc.CallOption) (*GetFollowerResponse, error) {
	out := new(GetFollowerResponse)
	err := c.cc.Invoke(ctx, FollowerService_GetUserFollowers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FollowerServiceServer is the server API for FollowerService service.
// All implementations must embed UnimplementedFollowerServiceServer
// for forward compatibility
type FollowerServiceServer interface {
	GetUserFollowers(context.Context, *GetFollowerRequest) (*GetFollowerResponse, error)
	mustEmbedUnimplementedFollowerServiceServer()
}

// UnimplementedFollowerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFollowerServiceServer struct {
}

func (UnimplementedFollowerServiceServer) GetUserFollowers(context.Context, *GetFollowerRequest) (*GetFollowerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFollowers not implemented")
}
func (UnimplementedFollowerServiceServer) mustEmbedUnimplementedFollowerServiceServer() {}

// UnsafeFollowerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FollowerServiceServer will
// result in compilation errors.
type UnsafeFollowerServiceServer interface {
	mustEmbedUnimplementedFollowerServiceServer()
}

func RegisterFollowerServiceServer(s grpc.ServiceRegistrar, srv FollowerServiceServer) {
	s.RegisterService(&FollowerService_ServiceDesc, srv)
}

func _FollowerService_GetUserFollowers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFollowerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServiceServer).GetUserFollowers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FollowerService_GetUserFollowers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServiceServer).GetUserFollowers(ctx, req.(*GetFollowerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FollowerService_ServiceDesc is the grpc.ServiceDesc for FollowerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FollowerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fanout.FollowerService",
	HandlerType: (*FollowerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserFollowers",
			Handler:    _FollowerService_GetUserFollowers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/followers.proto",
}
