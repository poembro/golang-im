package service

import (
	"context"
	"golang-im/pkg/grpclib"
	"golang-im/pkg/rpc"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"

	"golang-im/pkg/logger"
	"golang-im/pkg/pb"

	"golang-im/pkg/gerrors"
)

// SendOne 一对一消息发送
func (s *Service) SendOne(ctx context.Context, msg *pb.PushMsgReq) error {
	requestId := grpclib.GetCtxRequestIdStr(ctx)
	userId, deviceId, err := grpclib.GetCtxDataStr(ctx)
	if err != nil {
		logger.Sugar.Infow("logic 服务 SendMessage 头信息error")
		return err
	}
	// 加上grpc头防止api授权拦截
	metaCtx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", userId,
		"device_id", deviceId,
		"token", "TODO token verify",
		"request_id", requestId))

	rpc.ConnectInt(msg.Message.Server).DeliverMessage(metaCtx, msg)

	return nil
}

// SendRoom 群组消息发送
func (s *Service) SendRoom(ctx context.Context, msg *pb.PushMsg) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	err = s.dao.Publish(s.c.Global.PushAllTopic, bytes)
	if err != nil {
		return err
	}

	return nil
}

// Sync 消息同步
func (*Service) Sync(ctx context.Context, userId, seq int64) (*pb.SyncResp, error) {
	msg := []byte(`{"name":"zhangsan","sex":2}`)
	resp := &pb.SyncResp{Messages: msg, HasMore: false}
	return resp, nil
}

// MessageACK 消息确认机制
func (s *Service) MessageACK(ctx context.Context, deviceId, roomId string, userId, deviceAck, receiveTime int64) error {
	s.dao.AddMessageACKMapping(deviceId, roomId, deviceAck)
	return nil
}

// GetMessageCount 统计未读
func (s *Service) GetMessageCount(roomId, start, stop string) (int64, error) {
	return s.dao.GetMessageCount(roomId, start, stop)
}

// GetMessageList 取回消息
func (s *Service) GetMessageList(roomId string, start, stop int64) ([]string, error) {
	return s.dao.GetMessageList(roomId, start, stop)
}

// AddMessageList 将消息添加到对应房间 roomId.
func (s *Service) AddMessageList(roomId string, id int64, msg string) error {
	return s.dao.AddMessageList(roomId, id, msg)
}
