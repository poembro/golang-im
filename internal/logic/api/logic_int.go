package api

import (
    "context"
    "golang-im/internal/logic/service"
    "golang-im/pkg/grpclib"

    //"golang-im/pkg/gerrors"
    "golang-im/pkg/logger"
    "golang-im/pkg/pb"

    "go.uber.org/zap"
)

type LogicIntServer struct{}

// ConnSignIn 设备登录
func (*LogicIntServer) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) (*pb.ConnSignInResp, error) {
    logger.Sugar.Infow("ConnSignIn", "SignIn", "设备登录", "desc_requeset_id", grpclib.GetCtxRequestId(ctx))
    deviceId, userId, err := service.AuthService.SignIn(ctx, req.Body, req.ConnAddr, req.ClientAddr)
    return &pb.ConnSignInResp{UserId: userId, DeviceId: deviceId}, err
}

// SendMessage 发送消息
func (*LogicIntServer) SendMessage(ctx context.Context, req *pb.PushMsgReq) (*pb.PushMsgReply, error) {
    var err error
    /*buf := &pb.PushMsg{
          Type :pb.PushMsg_ROOM,
          Operation:2,
          Speed:2,
          Server: "127.0.0.1:5000",
          RoomId:"live://8000-20210817001",
          Msg: []byte("hello world 123s"),
      }
       &pb.PushMsgReq{
          Message: buf,
      }
    */
    logger.Logger.Debug("logic服务 SendMessage方法收到消息", zap.Any("msg", req.Message))

    switch req.Message.Type {
    case pb.PushMsg_PUSH:
        err = service.MessageService.SendOne(ctx, req)
    case pb.PushMsg_ROOM:
        err = service.MessageService.SendRoom(ctx, req.Message)
    default:
        // TODO
    }
    if err != nil {
        return nil, err
    }
    return &pb.PushMsgReply{}, nil
}

// MessageACK 设备收到消息回执
func (s *LogicIntServer) MessageACK(ctx context.Context, req *pb.MessageACKReq) (*pb.MessageACKResp, error) {
    err := service.MessageService.MessageACK(ctx, req.DeviceId, req.RoomId, req.UserId, req.DeviceAck, req.ReceiveTime)
    return &pb.MessageACKResp{}, err
}

// Sync 设备同步消息
func (*LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
    return service.MessageService.Sync(ctx, req.UserId, req.Seq)
}

// Heartbeat 心跳包
func (s *LogicIntServer) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatResp, error) {
    service.AuthService.Heartbeat(ctx, req.UserId, req.DeviceId, req.ConnAddr)
    return &pb.HeartbeatResp{}, nil
}

// Offline 设备离线
func (s *LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*pb.OfflineResp, error) {
    logger.Sugar.Infow("Offline", "Offline", "设备离线", "desc_requeset_id", grpclib.GetCtxRequestId(ctx))

    service.AuthService.Offline(ctx, req.UserId, req.DeviceId, req.ClientAddr)
    return &pb.OfflineResp{}, nil
}

// ServerStop 服务停止
func (s *LogicIntServer) ServerStop(ctx context.Context, in *pb.ServerStopReq) (*pb.ServerStopResp, error) {
    return &pb.ServerStopResp{}, nil
}
