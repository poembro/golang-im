// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.7.0
// source: logic.int.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConnSignInReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body       []byte `protobuf:"bytes,1,opt,name=Body,proto3" json:"Body,omitempty"`                               // body 是1个json字符串 包含了设备id 双方头像昵称等信息
	ConnAddr   string `protobuf:"bytes,2,opt,name=conn_addr,json=connAddr,proto3" json:"conn_addr,omitempty"`       // 服务器地址
	ClientAddr string `protobuf:"bytes,3,opt,name=client_addr,json=clientAddr,proto3" json:"client_addr,omitempty"` // 客户端地址
}

func (x *ConnSignInReq) Reset() {
	*x = ConnSignInReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnSignInReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnSignInReq) ProtoMessage() {}

func (x *ConnSignInReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnSignInReq.ProtoReflect.Descriptor instead.
func (*ConnSignInReq) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{0}
}

func (x *ConnSignInReq) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *ConnSignInReq) GetConnAddr() string {
	if x != nil {
		return x.ConnAddr
	}
	return ""
}

func (x *ConnSignInReq) GetClientAddr() string {
	if x != nil {
		return x.ClientAddr
	}
	return ""
}

type ConnSignInResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId int64 `protobuf:"varint,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"` // 设备id  用来区分一个用户多个设备 之间消息同步问题
	UserId   int64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`       // 用户id
}

func (x *ConnSignInResp) Reset() {
	*x = ConnSignInResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnSignInResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnSignInResp) ProtoMessage() {}

func (x *ConnSignInResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnSignInResp.ProtoReflect.Descriptor instead.
func (*ConnSignInResp) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{1}
}

func (x *ConnSignInResp) GetDeviceId() int64 {
	if x != nil {
		return x.DeviceId
	}
	return 0
}

func (x *ConnSignInResp) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type MessageACKReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                // 用户id
	DeviceId    int64 `protobuf:"varint,2,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`          // 设备id
	DeviceAck   int64 `protobuf:"varint,3,opt,name=device_ack,json=deviceAck,proto3" json:"device_ack,omitempty"`       // 设备收到消息的确认号
	ReceiveTime int64 `protobuf:"varint,4,opt,name=receive_time,json=receiveTime,proto3" json:"receive_time,omitempty"` // 消息接收时间戳，精确到毫秒
}

func (x *MessageACKReq) Reset() {
	*x = MessageACKReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageACKReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageACKReq) ProtoMessage() {}

func (x *MessageACKReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageACKReq.ProtoReflect.Descriptor instead.
func (*MessageACKReq) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{2}
}

func (x *MessageACKReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *MessageACKReq) GetDeviceId() int64 {
	if x != nil {
		return x.DeviceId
	}
	return 0
}

func (x *MessageACKReq) GetDeviceAck() int64 {
	if x != nil {
		return x.DeviceAck
	}
	return 0
}

func (x *MessageACKReq) GetReceiveTime() int64 {
	if x != nil {
		return x.ReceiveTime
	}
	return 0
}

type MessageACKResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MessageACKResp) Reset() {
	*x = MessageACKResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageACKResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageACKResp) ProtoMessage() {}

func (x *MessageACKResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageACKResp.ProtoReflect.Descriptor instead.
func (*MessageACKResp) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{3}
}

type SyncReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`       // 用户id
	DeviceId int64 `protobuf:"varint,2,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"` // 设备id
	Seq      int64 `protobuf:"varint,3,opt,name=seq,proto3" json:"seq,omitempty"`                           // 客户端已经同步的序列号
}

func (x *SyncReq) Reset() {
	*x = SyncReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncReq) ProtoMessage() {}

func (x *SyncReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncReq.ProtoReflect.Descriptor instead.
func (*SyncReq) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{4}
}

func (x *SyncReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *SyncReq) GetDeviceId() int64 {
	if x != nil {
		return x.DeviceId
	}
	return 0
}

func (x *SyncReq) GetSeq() int64 {
	if x != nil {
		return x.Seq
	}
	return 0
}

type SyncResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []byte `protobuf:"bytes,1,opt,name=messages,proto3" json:"messages,omitempty"`               // 消息列表
	HasMore  bool   `protobuf:"varint,2,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"` // 是否有更多数据
}

func (x *SyncResp) Reset() {
	*x = SyncResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncResp) ProtoMessage() {}

func (x *SyncResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncResp.ProtoReflect.Descriptor instead.
func (*SyncResp) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{5}
}

func (x *SyncResp) GetMessages() []byte {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *SyncResp) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

type OfflineReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`            // 用户id
	DeviceId   int64  `protobuf:"varint,2,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`      // 设备id
	ClientAddr string `protobuf:"bytes,3,opt,name=client_addr,json=clientAddr,proto3" json:"client_addr,omitempty"` // 客户端地址
}

func (x *OfflineReq) Reset() {
	*x = OfflineReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OfflineReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OfflineReq) ProtoMessage() {}

func (x *OfflineReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OfflineReq.ProtoReflect.Descriptor instead.
func (*OfflineReq) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{6}
}

func (x *OfflineReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *OfflineReq) GetDeviceId() int64 {
	if x != nil {
		return x.DeviceId
	}
	return 0
}

func (x *OfflineReq) GetClientAddr() string {
	if x != nil {
		return x.ClientAddr
	}
	return ""
}

type OfflineResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OfflineResp) Reset() {
	*x = OfflineResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OfflineResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OfflineResp) ProtoMessage() {}

func (x *OfflineResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OfflineResp.ProtoReflect.Descriptor instead.
func (*OfflineResp) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{7}
}

type ServerStopReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnAddr string `protobuf:"bytes,1,opt,name=conn_addr,json=connAddr,proto3" json:"conn_addr,omitempty"`
}

func (x *ServerStopReq) Reset() {
	*x = ServerStopReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerStopReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerStopReq) ProtoMessage() {}

func (x *ServerStopReq) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerStopReq.ProtoReflect.Descriptor instead.
func (*ServerStopReq) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{8}
}

func (x *ServerStopReq) GetConnAddr() string {
	if x != nil {
		return x.ConnAddr
	}
	return ""
}

type ServerStopResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServerStopResp) Reset() {
	*x = ServerStopResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logic_int_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerStopResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerStopResp) ProtoMessage() {}

func (x *ServerStopResp) ProtoReflect() protoreflect.Message {
	mi := &file_logic_int_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerStopResp.ProtoReflect.Descriptor instead.
func (*ServerStopResp) Descriptor() ([]byte, []int) {
	return file_logic_int_proto_rawDescGZIP(), []int{9}
}

var File_logic_int_proto protoreflect.FileDescriptor

var file_logic_int_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0e, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x65, 0x78, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x61, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6e, 0x53, 0x69, 0x67,
	0x6e, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6f,
	0x6e, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63,
	0x6f, 0x6e, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x22, 0x46, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e,
	0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x87, 0x01, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x52,
	0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x5f, 0x61, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x41, 0x63, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x52, 0x65, 0x73, 0x70, 0x22, 0x51, 0x0a, 0x07,
	0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x73, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x73, 0x65, 0x71, 0x22,
	0x41, 0x0a, 0x08, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x5f, 0x6d,
	0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x61, 0x73, 0x4d, 0x6f,
	0x72, 0x65, 0x22, 0x63, 0x0a, 0x0a, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x22, 0x0d, 0x0a, 0x0b, 0x4f, 0x66, 0x66, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x2c, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x53, 0x74, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x6e, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x6e,
	0x41, 0x64, 0x64, 0x72, 0x22, 0x10, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74,
	0x6f, 0x70, 0x52, 0x65, 0x73, 0x70, 0x32, 0xa9, 0x02, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x69, 0x63,
	0x49, 0x6e, 0x74, 0x12, 0x33, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x49,
	0x6e, 0x12, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x49,
	0x6e, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x53, 0x69,
	0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2f, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x33, 0x0a, 0x0a, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x12, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x43, 0x4b, 0x52, 0x65, 0x73, 0x70, 0x12, 0x21,
	0x0a, 0x04, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x79, 0x6e, 0x63,
	0x52, 0x65, 0x71, 0x1a, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x2a, 0x0a, 0x07, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x0e, 0x2e, 0x70,
	0x62, 0x2e, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70,
	0x62, 0x2e, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x33, 0x0a,
	0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x11, 0x2e, 0x70, 0x62,
	0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x1a, 0x12,
	0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_logic_int_proto_rawDescOnce sync.Once
	file_logic_int_proto_rawDescData = file_logic_int_proto_rawDesc
)

func file_logic_int_proto_rawDescGZIP() []byte {
	file_logic_int_proto_rawDescOnce.Do(func() {
		file_logic_int_proto_rawDescData = protoimpl.X.CompressGZIP(file_logic_int_proto_rawDescData)
	})
	return file_logic_int_proto_rawDescData
}

var file_logic_int_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_logic_int_proto_goTypes = []interface{}{
	(*ConnSignInReq)(nil),  // 0: pb.ConnSignInReq
	(*ConnSignInResp)(nil), // 1: pb.ConnSignInResp
	(*MessageACKReq)(nil),  // 2: pb.MessageACKReq
	(*MessageACKResp)(nil), // 3: pb.MessageACKResp
	(*SyncReq)(nil),        // 4: pb.SyncReq
	(*SyncResp)(nil),       // 5: pb.SyncResp
	(*OfflineReq)(nil),     // 6: pb.OfflineReq
	(*OfflineResp)(nil),    // 7: pb.OfflineResp
	(*ServerStopReq)(nil),  // 8: pb.ServerStopReq
	(*ServerStopResp)(nil), // 9: pb.ServerStopResp
	(*PushMsgReq)(nil),     // 10: pb.PushMsgReq
	(*PushMsgReply)(nil),   // 11: pb.PushMsgReply
}
var file_logic_int_proto_depIdxs = []int32{
	0,  // 0: pb.LogicInt.ConnSignIn:input_type -> pb.ConnSignInReq
	10, // 1: pb.LogicInt.SendMessage:input_type -> pb.PushMsgReq
	2,  // 2: pb.LogicInt.MessageACK:input_type -> pb.MessageACKReq
	4,  // 3: pb.LogicInt.Sync:input_type -> pb.SyncReq
	6,  // 4: pb.LogicInt.Offline:input_type -> pb.OfflineReq
	8,  // 5: pb.LogicInt.ServerStop:input_type -> pb.ServerStopReq
	1,  // 6: pb.LogicInt.ConnSignIn:output_type -> pb.ConnSignInResp
	11, // 7: pb.LogicInt.SendMessage:output_type -> pb.PushMsgReply
	3,  // 8: pb.LogicInt.MessageACK:output_type -> pb.MessageACKResp
	5,  // 9: pb.LogicInt.Sync:output_type -> pb.SyncResp
	7,  // 10: pb.LogicInt.Offline:output_type -> pb.OfflineResp
	9,  // 11: pb.LogicInt.ServerStop:output_type -> pb.ServerStopResp
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_logic_int_proto_init() }
func file_logic_int_proto_init() {
	if File_logic_int_proto != nil {
		return
	}
	file_push_ext_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_logic_int_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnSignInReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnSignInResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageACKReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageACKResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OfflineReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OfflineResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerStopReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logic_int_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerStopResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_logic_int_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_logic_int_proto_goTypes,
		DependencyIndexes: file_logic_int_proto_depIdxs,
		MessageInfos:      file_logic_int_proto_msgTypes,
	}.Build()
	File_logic_int_proto = out.File
	file_logic_int_proto_rawDesc = nil
	file_logic_int_proto_goTypes = nil
	file_logic_int_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LogicIntClient is the client API for LogicInt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LogicIntClient interface {
	// 登录
	ConnSignIn(ctx context.Context, in *ConnSignInReq, opts ...grpc.CallOption) (*ConnSignInResp, error)
	// 发送消息
	SendMessage(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error)
	// 设备收到消息回执
	MessageACK(ctx context.Context, in *MessageACKReq, opts ...grpc.CallOption) (*MessageACKResp, error)
	// 同步历史聊天记录
	Sync(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (*SyncResp, error)
	// 设备离线
	Offline(ctx context.Context, in *OfflineReq, opts ...grpc.CallOption) (*OfflineResp, error)
	// 服务停止
	ServerStop(ctx context.Context, in *ServerStopReq, opts ...grpc.CallOption) (*ServerStopResp, error)
}

type logicIntClient struct {
	cc grpc.ClientConnInterface
}

func NewLogicIntClient(cc grpc.ClientConnInterface) LogicIntClient {
	return &logicIntClient{cc}
}

func (c *logicIntClient) ConnSignIn(ctx context.Context, in *ConnSignInReq, opts ...grpc.CallOption) (*ConnSignInResp, error) {
	out := new(ConnSignInResp)
	err := c.cc.Invoke(ctx, "/pb.LogicInt/ConnSignIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicIntClient) SendMessage(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error) {
	out := new(PushMsgReply)
	err := c.cc.Invoke(ctx, "/pb.LogicInt/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicIntClient) MessageACK(ctx context.Context, in *MessageACKReq, opts ...grpc.CallOption) (*MessageACKResp, error) {
	out := new(MessageACKResp)
	err := c.cc.Invoke(ctx, "/pb.LogicInt/MessageACK", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicIntClient) Sync(ctx context.Context, in *SyncReq, opts ...grpc.CallOption) (*SyncResp, error) {
	out := new(SyncResp)
	err := c.cc.Invoke(ctx, "/pb.LogicInt/Sync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicIntClient) Offline(ctx context.Context, in *OfflineReq, opts ...grpc.CallOption) (*OfflineResp, error) {
	out := new(OfflineResp)
	err := c.cc.Invoke(ctx, "/pb.LogicInt/Offline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicIntClient) ServerStop(ctx context.Context, in *ServerStopReq, opts ...grpc.CallOption) (*ServerStopResp, error) {
	out := new(ServerStopResp)
	err := c.cc.Invoke(ctx, "/pb.LogicInt/ServerStop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogicIntServer is the server API for LogicInt service.
type LogicIntServer interface {
	// 登录
	ConnSignIn(context.Context, *ConnSignInReq) (*ConnSignInResp, error)
	// 发送消息
	SendMessage(context.Context, *PushMsgReq) (*PushMsgReply, error)
	// 设备收到消息回执
	MessageACK(context.Context, *MessageACKReq) (*MessageACKResp, error)
	// 同步历史聊天记录
	Sync(context.Context, *SyncReq) (*SyncResp, error)
	// 设备离线
	Offline(context.Context, *OfflineReq) (*OfflineResp, error)
	// 服务停止
	ServerStop(context.Context, *ServerStopReq) (*ServerStopResp, error)
}

// UnimplementedLogicIntServer can be embedded to have forward compatible implementations.
type UnimplementedLogicIntServer struct {
}

func (*UnimplementedLogicIntServer) ConnSignIn(context.Context, *ConnSignInReq) (*ConnSignInResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnSignIn not implemented")
}
func (*UnimplementedLogicIntServer) SendMessage(context.Context, *PushMsgReq) (*PushMsgReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (*UnimplementedLogicIntServer) MessageACK(context.Context, *MessageACKReq) (*MessageACKResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageACK not implemented")
}
func (*UnimplementedLogicIntServer) Sync(context.Context, *SyncReq) (*SyncResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (*UnimplementedLogicIntServer) Offline(context.Context, *OfflineReq) (*OfflineResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Offline not implemented")
}
func (*UnimplementedLogicIntServer) ServerStop(context.Context, *ServerStopReq) (*ServerStopResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServerStop not implemented")
}

func RegisterLogicIntServer(s *grpc.Server, srv LogicIntServer) {
	s.RegisterService(&_LogicInt_serviceDesc, srv)
}

func _LogicInt_ConnSignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnSignInReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicIntServer).ConnSignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LogicInt/ConnSignIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicIntServer).ConnSignIn(ctx, req.(*ConnSignInReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInt_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicIntServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LogicInt/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicIntServer).SendMessage(ctx, req.(*PushMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInt_MessageACK_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageACKReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicIntServer).MessageACK(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LogicInt/MessageACK",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicIntServer).MessageACK(ctx, req.(*MessageACKReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInt_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicIntServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LogicInt/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicIntServer).Sync(ctx, req.(*SyncReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInt_Offline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OfflineReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicIntServer).Offline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LogicInt/Offline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicIntServer).Offline(ctx, req.(*OfflineReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogicInt_ServerStop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerStopReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicIntServer).ServerStop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LogicInt/ServerStop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicIntServer).ServerStop(ctx, req.(*ServerStopReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _LogicInt_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.LogicInt",
	HandlerType: (*LogicIntServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnSignIn",
			Handler:    _LogicInt_ConnSignIn_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _LogicInt_SendMessage_Handler,
		},
		{
			MethodName: "MessageACK",
			Handler:    _LogicInt_MessageACK_Handler,
		},
		{
			MethodName: "Sync",
			Handler:    _LogicInt_Sync_Handler,
		},
		{
			MethodName: "Offline",
			Handler:    _LogicInt_Offline_Handler,
		},
		{
			MethodName: "ServerStop",
			Handler:    _LogicInt_ServerStop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logic.int.proto",
}
