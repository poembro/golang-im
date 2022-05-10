package apigrpc

import (
	"context"
	"golang-im/conf"
	"golang-im/internal/logic/service"

	//"golang-im/pkg/gerrors"

	"golang-im/pkg/pb"
)

var (
	svc *service.Service
)

type LogicIntServer struct{}

func NewLogicIntServer(c *conf.Config) *LogicIntServer {
	svc = service.New(c)
	return &LogicIntServer{}
}

// ConnSignIn 设备登录
func (*LogicIntServer) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) (*pb.ConnSignInResp, error) {
	//logger.Sugar.Infow("ConnSignIn", "SignIn", "设备登录", "desc_requeset_id", grpclib.GetCtxRequestId(ctx))
	deviceId, userId, err := svc.SignIn(ctx, req.Body, req.ConnAddr, req.ClientAddr)
	return &pb.ConnSignInResp{UserId: userId, DeviceId: deviceId}, err
}

// SendMessage 发送消息
func (*LogicIntServer) SendMessage(ctx context.Context, req *pb.PushMsgReq) (*pb.PushMsgReply, error) {
	var err error
	//logger.Logger.Debug("SendMessage", zap.String("desc", "grpc服务logic业务SendMessage方法 收到消息"), zap.Any("msg", req.Message))

	switch req.Message.Type {
	case pb.PushMsg_PUSH:
		err = svc.SendOne(ctx, req)
	case pb.PushMsg_ROOM:
		err = svc.SendRoom(ctx, req.Message)
	default:
		// TODO
	}
	if err != nil {
		return nil, err
	}
	return &pb.PushMsgReply{}, nil
}

// MessageACK 设备收到消息
func (*LogicIntServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*pb.MessageACKResp, error) {
	err := svc.MessageACK(ctx, req.DeviceId, req.RoomId, req.UserId, req.DeviceAck, req.ReceiveTime)
	return &pb.MessageACKResp{}, err
}

// Sync 设备同步消息
func (*LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	return svc.Sync(ctx, req.RoomId, req.Seq)
}

// Heartbeat 心跳包
func (*LogicIntServer) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatResp, error) {
	svc.Heartbeat(ctx, req.UserId, req.DeviceId, req.ConnAddr)
	return &pb.HeartbeatResp{}, nil
}

// Offline 设备离线
func (*LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*pb.OfflineResp, error) {
	svc.Offline(ctx, req.UserId, req.DeviceId, req.ClientAddr)
	return &pb.OfflineResp{}, nil
}

// ServerStop 服务停止
func (*LogicIntServer) ServerStop(ctx context.Context, in *pb.ServerStopReq) (*pb.ServerStopResp, error) {
	return &pb.ServerStopResp{}, nil
}
