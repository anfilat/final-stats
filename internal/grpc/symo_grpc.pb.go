// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SymoClient is the client API for Symo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SymoClient interface {
	GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (Symo_GetStatsClient, error)
}

type symoClient struct {
	cc grpc.ClientConnInterface
}

func NewSymoClient(cc grpc.ClientConnInterface) SymoClient {
	return &symoClient{cc}
}

func (c *symoClient) GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (Symo_GetStatsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Symo_serviceDesc.Streams[0], "/stats.Symo/GetStats", opts...)
	if err != nil {
		return nil, err
	}
	x := &symoGetStatsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Symo_GetStatsClient interface {
	Recv() (*Stats, error)
	grpc.ClientStream
}

type symoGetStatsClient struct {
	grpc.ClientStream
}

func (x *symoGetStatsClient) Recv() (*Stats, error) {
	m := new(Stats)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SymoServer is the server API for Symo service.
// All implementations must embed UnimplementedSymoServer
// for forward compatibility
type SymoServer interface {
	GetStats(*StatsRequest, Symo_GetStatsServer) error
	mustEmbedUnimplementedSymoServer()
}

// UnimplementedSymoServer must be embedded to have forward compatible implementations.
type UnimplementedSymoServer struct {
}

func (UnimplementedSymoServer) GetStats(*StatsRequest, Symo_GetStatsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedSymoServer) mustEmbedUnimplementedSymoServer() {}

// UnsafeSymoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SymoServer will
// result in compilation errors.
type UnsafeSymoServer interface {
	mustEmbedUnimplementedSymoServer()
}

func RegisterSymoServer(s grpc.ServiceRegistrar, srv SymoServer) {
	s.RegisterService(&_Symo_serviceDesc, srv)
}

func _Symo_GetStats_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StatsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SymoServer).GetStats(m, &symoGetStatsServer{stream})
}

type Symo_GetStatsServer interface {
	Send(*Stats) error
	grpc.ServerStream
}

type symoGetStatsServer struct {
	grpc.ServerStream
}

func (x *symoGetStatsServer) Send(m *Stats) error {
	return x.ServerStream.SendMsg(m)
}

var _Symo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "stats.Symo",
	HandlerType: (*SymoServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetStats",
			Handler:       _Symo_GetStats_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "symo.proto",
}
