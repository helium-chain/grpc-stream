// 这里说明我们使用的是proto3语法

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: hello.proto

package service

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
	SayHello_SayHello_FullMethodName = "/SayHello/SayHello"
	SayHello_Channel_FullMethodName  = "/SayHello/Channel"
)

// SayHelloClient is the client API for SayHello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SayHelloClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	Channel(ctx context.Context, opts ...grpc.CallOption) (SayHello_ChannelClient, error)
}

type sayHelloClient struct {
	cc grpc.ClientConnInterface
}

func NewSayHelloClient(cc grpc.ClientConnInterface) SayHelloClient {
	return &sayHelloClient{cc}
}

func (c *sayHelloClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, SayHello_SayHello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sayHelloClient) Channel(ctx context.Context, opts ...grpc.CallOption) (SayHello_ChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &SayHello_ServiceDesc.Streams[0], SayHello_Channel_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &sayHelloChannelClient{stream}
	return x, nil
}

type SayHello_ChannelClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type sayHelloChannelClient struct {
	grpc.ClientStream
}

func (x *sayHelloChannelClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sayHelloChannelClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SayHelloServer is the server API for SayHello service.
// All implementations must embed UnimplementedSayHelloServer
// for forward compatibility
type SayHelloServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	Channel(SayHello_ChannelServer) error
	mustEmbedUnimplementedSayHelloServer()
}

// UnimplementedSayHelloServer must be embedded to have forward compatible implementations.
type UnimplementedSayHelloServer struct {
}

func (UnimplementedSayHelloServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedSayHelloServer) Channel(SayHello_ChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method Channel not implemented")
}
func (UnimplementedSayHelloServer) mustEmbedUnimplementedSayHelloServer() {}

// UnsafeSayHelloServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SayHelloServer will
// result in compilation errors.
type UnsafeSayHelloServer interface {
	mustEmbedUnimplementedSayHelloServer()
}

func RegisterSayHelloServer(s grpc.ServiceRegistrar, srv SayHelloServer) {
	s.RegisterService(&SayHello_ServiceDesc, srv)
}

func _SayHello_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SayHelloServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SayHello_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SayHelloServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SayHello_Channel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SayHelloServer).Channel(&sayHelloChannelServer{stream})
}

type SayHello_ChannelServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type sayHelloChannelServer struct {
	grpc.ServerStream
}

func (x *sayHelloChannelServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sayHelloChannelServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SayHello_ServiceDesc is the grpc.ServiceDesc for SayHello service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SayHello_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SayHello",
	HandlerType: (*SayHelloServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _SayHello_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Channel",
			Handler:       _SayHello_Channel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hello.proto",
}
