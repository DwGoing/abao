// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: abao.proto

package abao

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

// AbaoClient is the client API for Abao service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AbaoClient interface {
	GetAddress(ctx context.Context, in *GetAddressRequest, opts ...grpc.CallOption) (*GetAddressResponse, error)
	GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error)
}

type abaoClient struct {
	cc grpc.ClientConnInterface
}

func NewAbaoClient(cc grpc.ClientConnInterface) AbaoClient {
	return &abaoClient{cc}
}

func (c *abaoClient) GetAddress(ctx context.Context, in *GetAddressRequest, opts ...grpc.CallOption) (*GetAddressResponse, error) {
	out := new(GetAddressResponse)
	err := c.cc.Invoke(ctx, "/abao.Abao/GetAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *abaoClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, "/abao.Abao/GetBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *abaoClient) Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error) {
	out := new(TransferResponse)
	err := c.cc.Invoke(ctx, "/abao.Abao/Transfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AbaoServer is the server API for Abao service.
// All implementations must embed UnimplementedAbaoServer
// for forward compatibility
type AbaoServer interface {
	GetAddress(context.Context, *GetAddressRequest) (*GetAddressResponse, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	Transfer(context.Context, *TransferRequest) (*TransferResponse, error)
	mustEmbedUnimplementedAbaoServer()
}

// UnimplementedAbaoServer must be embedded to have forward compatible implementations.
type UnimplementedAbaoServer struct {
}

func (UnimplementedAbaoServer) GetAddress(context.Context, *GetAddressRequest) (*GetAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddress not implemented")
}
func (UnimplementedAbaoServer) GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedAbaoServer) Transfer(context.Context, *TransferRequest) (*TransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
}
func (UnimplementedAbaoServer) mustEmbedUnimplementedAbaoServer() {}

// UnsafeAbaoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AbaoServer will
// result in compilation errors.
type UnsafeAbaoServer interface {
	mustEmbedUnimplementedAbaoServer()
}

func RegisterAbaoServer(s grpc.ServiceRegistrar, srv AbaoServer) {
	s.RegisterService(&Abao_ServiceDesc, srv)
}

func _Abao_GetAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AbaoServer).GetAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abao.Abao/GetAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AbaoServer).GetAddress(ctx, req.(*GetAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Abao_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AbaoServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abao.Abao/GetBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AbaoServer).GetBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Abao_Transfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AbaoServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abao.Abao/Transfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AbaoServer).Transfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Abao_ServiceDesc is the grpc.ServiceDesc for Abao service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Abao_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "abao.Abao",
	HandlerType: (*AbaoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAddress",
			Handler:    _Abao_GetAddress_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _Abao_GetBalance_Handler,
		},
		{
			MethodName: "Transfer",
			Handler:    _Abao_Transfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "abao.proto",
}
