// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	auth.proto

It has these top-level messages:
	GetAccountRequest
	GetAccountResponse
	ResetPasswordRequest
	ResetPasswordResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type GetAccountRequest struct {
	EmailAddress string `protobuf:"bytes,1,opt,name=emailAddress" json:"emailAddress,omitempty"`
}

func (m *GetAccountRequest) Reset()                    { *m = GetAccountRequest{} }
func (m *GetAccountRequest) String() string            { return proto1.CompactTextString(m) }
func (*GetAccountRequest) ProtoMessage()               {}
func (*GetAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetAccountRequest) GetEmailAddress() string {
	if m != nil {
		return m.EmailAddress
	}
	return ""
}

type GetAccountResponse struct {
	AccountId int64 `protobuf:"varint,1,opt,name=accountId" json:"accountId,omitempty"`
}

func (m *GetAccountResponse) Reset()                    { *m = GetAccountResponse{} }
func (m *GetAccountResponse) String() string            { return proto1.CompactTextString(m) }
func (*GetAccountResponse) ProtoMessage()               {}
func (*GetAccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetAccountResponse) GetAccountId() int64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

type ResetPasswordRequest struct {
	AccountId   int64  `protobuf:"varint,1,opt,name=accountId" json:"accountId,omitempty"`
	NewPassword string `protobuf:"bytes,2,opt,name=newPassword" json:"newPassword,omitempty"`
}

func (m *ResetPasswordRequest) Reset()                    { *m = ResetPasswordRequest{} }
func (m *ResetPasswordRequest) String() string            { return proto1.CompactTextString(m) }
func (*ResetPasswordRequest) ProtoMessage()               {}
func (*ResetPasswordRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ResetPasswordRequest) GetAccountId() int64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *ResetPasswordRequest) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

type ResetPasswordResponse struct {
	AccountId int64 `protobuf:"varint,1,opt,name=accountId" json:"accountId,omitempty"`
}

func (m *ResetPasswordResponse) Reset()                    { *m = ResetPasswordResponse{} }
func (m *ResetPasswordResponse) String() string            { return proto1.CompactTextString(m) }
func (*ResetPasswordResponse) ProtoMessage()               {}
func (*ResetPasswordResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ResetPasswordResponse) GetAccountId() int64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func init() {
	proto1.RegisterType((*GetAccountRequest)(nil), "proto.GetAccountRequest")
	proto1.RegisterType((*GetAccountResponse)(nil), "proto.GetAccountResponse")
	proto1.RegisterType((*ResetPasswordRequest)(nil), "proto.ResetPasswordRequest")
	proto1.RegisterType((*ResetPasswordResponse)(nil), "proto.ResetPasswordResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AccountService service

type AccountServiceClient interface {
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
}

type accountServiceClient struct {
	cc *grpc.ClientConn
}

func NewAccountServiceClient(cc *grpc.ClientConn) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	out := new(GetAccountResponse)
	err := grpc.Invoke(ctx, "/proto.AccountService/GetAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AccountService service

type AccountServiceServer interface {
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
}

func RegisterAccountServiceServer(s *grpc.Server, srv AccountServiceServer) {
	s.RegisterService(&_AccountService_serviceDesc, srv)
}

func _AccountService_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AccountService/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccountService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAccount",
			Handler:    _AccountService_GetAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

// Client API for ResetPasswordService service

type ResetPasswordServiceClient interface {
	ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*ResetPasswordResponse, error)
}

type resetPasswordServiceClient struct {
	cc *grpc.ClientConn
}

func NewResetPasswordServiceClient(cc *grpc.ClientConn) ResetPasswordServiceClient {
	return &resetPasswordServiceClient{cc}
}

func (c *resetPasswordServiceClient) ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*ResetPasswordResponse, error) {
	out := new(ResetPasswordResponse)
	err := grpc.Invoke(ctx, "/proto.ResetPasswordService/ResetPassword", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ResetPasswordService service

type ResetPasswordServiceServer interface {
	ResetPassword(context.Context, *ResetPasswordRequest) (*ResetPasswordResponse, error)
}

func RegisterResetPasswordServiceServer(s *grpc.Server, srv ResetPasswordServiceServer) {
	s.RegisterService(&_ResetPasswordService_serviceDesc, srv)
}

func _ResetPasswordService_ResetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResetPasswordServiceServer).ResetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ResetPasswordService/ResetPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResetPasswordServiceServer).ResetPassword(ctx, req.(*ResetPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ResetPasswordService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ResetPasswordService",
	HandlerType: (*ResetPasswordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResetPassword",
			Handler:    _ResetPasswordService_ResetPassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

func init() { proto1.RegisterFile("auth.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x2c, 0x2d, 0xc9,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0xe6, 0x5c, 0x82, 0xee, 0xa9,
	0x25, 0x8e, 0xc9, 0xc9, 0xf9, 0xa5, 0x79, 0x25, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42,
	0x4a, 0x5c, 0x3c, 0xa9, 0xb9, 0x89, 0x99, 0x39, 0x8e, 0x29, 0x29, 0x45, 0xa9, 0xc5, 0xc5, 0x12,
	0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x28, 0x62, 0x4a, 0x46, 0x5c, 0x42, 0xc8, 0x1a, 0x8b, 0x0b,
	0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x64, 0xb8, 0x38, 0x13, 0x21, 0x42, 0x9e, 0x29, 0x60, 0x6d, 0xcc,
	0x41, 0x08, 0x01, 0xa5, 0x30, 0x2e, 0x91, 0xa0, 0xd4, 0xe2, 0xd4, 0x92, 0x80, 0xc4, 0xe2, 0xe2,
	0xf2, 0xfc, 0xa2, 0x14, 0x98, 0x7d, 0x78, 0x75, 0x09, 0x29, 0x70, 0x71, 0xe7, 0xa5, 0x96, 0xc3,
	0xf4, 0x48, 0x30, 0x81, 0x1d, 0x83, 0x2c, 0xa4, 0x64, 0xca, 0x25, 0x8a, 0x66, 0x2e, 0x31, 0xce,
	0x31, 0x0a, 0xe5, 0xe2, 0x83, 0xba, 0x3f, 0x38, 0xb5, 0xa8, 0x2c, 0x33, 0x39, 0x55, 0xc8, 0x99,
	0x8b, 0x0b, 0xe1, 0x29, 0x21, 0x09, 0x48, 0x50, 0xe9, 0x61, 0x04, 0x90, 0x94, 0x24, 0x16, 0x19,
	0x88, 0x95, 0x4a, 0x0c, 0x46, 0x29, 0x68, 0xbe, 0x84, 0x19, 0xee, 0xc3, 0xc5, 0x8b, 0x22, 0x2e,
	0x24, 0x0d, 0x35, 0x05, 0x5b, 0x98, 0x48, 0xc9, 0x60, 0x97, 0x84, 0xd9, 0x92, 0xc4, 0x06, 0x96,
	0x36, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x03, 0x9a, 0x06, 0x01, 0xd4, 0x01, 0x00, 0x00,
}
