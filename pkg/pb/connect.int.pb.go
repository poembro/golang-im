// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.13.0
// source: connect.int.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_connect_int_proto protoreflect.FileDescriptor

var file_connect_int_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2e, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0e, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x65, 0x78,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x40, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x49, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_connect_int_proto_goTypes = []interface{}{
	(*PushMsgReq)(nil),   // 0: pb.PushMsgReq
	(*PushMsgReply)(nil), // 1: pb.PushMsgReply
}
var file_connect_int_proto_depIdxs = []int32{
	0, // 0: pb.ConnectInt.DeliverMessage:input_type -> pb.PushMsgReq
	1, // 1: pb.ConnectInt.DeliverMessage:output_type -> pb.PushMsgReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_connect_int_proto_init() }
func file_connect_int_proto_init() {
	if File_connect_int_proto != nil {
		return
	}
	file_push_ext_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_connect_int_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_connect_int_proto_goTypes,
		DependencyIndexes: file_connect_int_proto_depIdxs,
	}.Build()
	File_connect_int_proto = out.File
	file_connect_int_proto_rawDesc = nil
	file_connect_int_proto_goTypes = nil
	file_connect_int_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ConnectIntClient is the client API for ConnectInt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConnectIntClient interface {
	// ????????????
	DeliverMessage(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error)
}

type connectIntClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectIntClient(cc grpc.ClientConnInterface) ConnectIntClient {
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
	// ????????????
	DeliverMessage(context.Context, *PushMsgReq) (*PushMsgReply, error)
}

// UnimplementedConnectIntServer can be embedded to have forward compatible implementations.
type UnimplementedConnectIntServer struct {
}

func (*UnimplementedConnectIntServer) DeliverMessage(context.Context, *PushMsgReq) (*PushMsgReply, error) {
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
