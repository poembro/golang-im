package service

import (
	"context"
	"golang-im/internal/logic/cache"
	"golang-im/pkg/grpclib"
	"golang-im/pkg/rpc"
	"strconv"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"

	"golang-im/pkg/logger"
	"golang-im/pkg/pb"

	"golang-im/config"
	"golang-im/pkg/gerrors"
)

type messageService struct{}

var MessageService = new(messageService)

// SendOne 一对一消息发送
func (s *messageService) SendOne(ctx context.Context, msg *pb.PushMsgReq) error {
	requestId := grpclib.GetCtxRequestId(ctx)
	userId, deviceId, err := grpclib.GetCtxData(ctx)
	if err != nil {
		logger.Sugar.Infow("logic 服务 SendMessage 头信息error")
		return err
	}
	// 加上grpc头防止api授权拦截
	metaCtx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", strconv.FormatInt(userId, 10),
		"device_id", deviceId,
		"token", "md5/jwt/xxx",
		"request_id", strconv.FormatInt(requestId, 10)))

	rpc.ConnectInt(msg.Message.Server).DeliverMessage(metaCtx, msg)

	return nil
}

// SendRoom 群组消息发送
func (s *messageService) SendRoom(ctx context.Context, msg *pb.PushMsg) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	err = cache.Queue.Publish(config.Global.PushAllTopic, bytes)
	if err != nil {
		return err
	}

	return nil
}

// Sync 消息同步
func (*messageService) Sync(ctx context.Context, userId, seq int64) (*pb.SyncResp, error) {
	msg := []byte(`{"name":"zhangsan","sex":2}`)
	resp := &pb.SyncResp{Messages: msg, HasMore: false}
	return resp, nil
}

// MessageACK 消息确认机制
func (s *messageService) MessageACK(ctx context.Context, deviceId, roomId string, userId, deviceAck, receiveTime int64) error {
	cache.Online.AddMessageACKMapping(deviceId, roomId, deviceAck)
	return nil
}
