// Code generated by protoc-gen-go. DO NOT EDIT.
// source: connect.int.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("connect.int.proto", fileDescriptor_ea805b356a4eabbd) }

var fileDescriptor_ea805b356a4eabbd = []byte{
	// 114 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xce, 0xcf, 0xcb,
	0x4b, 0x4d, 0x2e, 0xd1, 0xcb, 0xcc, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a,
	0x48, 0x92, 0xe2, 0x2b, 0x28, 0x2d, 0xce, 0xd0, 0x4b, 0xad, 0x80, 0x8a, 0x19, 0x39, 0x70, 0x71,
	0x39, 0x43, 0x14, 0x7a, 0xe6, 0x95, 0x08, 0x19, 0x71, 0xf1, 0xb9, 0xa4, 0xe6, 0x64, 0x96, 0xa5,
	0x16, 0xf9, 0xa6, 0x16, 0x17, 0x27, 0xa6, 0xa7, 0x0a, 0xf1, 0xe9, 0x15, 0x24, 0xe9, 0x05, 0x94,
	0x16, 0x67, 0xf8, 0x16, 0xa7, 0x07, 0xa5, 0x16, 0x4a, 0x09, 0xa0, 0xf0, 0x0b, 0x72, 0x2a, 0x93,
	0xd8, 0xc0, 0x06, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x06, 0x2d, 0xbe, 0xd9, 0x71, 0x00,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConnectIntClient is the client API for ConnectInt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConnectIntClient interface {
	// 发送消息
	DeliverMessage(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error)
}

type connectIntClient struct {
	cc *grpc.ClientConn
}

func NewConnectIntClient(cc *grpc.ClientConn) ConnectIntClient {
	return &connectIntClient{cc}
}

func (c *connectIntClient) DeliverMessage(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error) {
	out := new(PushMsgReply)
	err := c.cc.Invoke(ctx, "/pb.ConnectInt/DeliverMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectIntServer is the server API for ConnectInt service.
type ConnectIntServer interface {
	// 发送消息
	DeliverMessage(context.Context, *PushMsgReq) (*PushMsgReply, error)
}

// UnimplementedConnectIntServer can be embedded to have forward compatible implementations.
type UnimplementedConnectIntServer struct {
}

func (*UnimplementedConnectIntServer) DeliverMessage(ctx context.Context, req *PushMsgReq) (*PushMsgReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeliverMessage not implemented")
}

func RegisterConnectIntServer(s *grpc.Server, srv ConnectIntServer) {
	s.RegisterService(&_ConnectInt_serviceDesc, srv)
}

func _ConnectInt_DeliverMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectIntServer).DeliverMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ConnectInt/DeliverMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectIntServer).DeliverMessage(ctx, req.(*PushMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConnectInt_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ConnectInt",
	HandlerType: (*ConnectIntServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeliverMessage",
			Handler:    _ConnectInt_DeliverMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "connect.int.proto",
}