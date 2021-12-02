package service

import (
    "context"
    "golang-im/internal/logic/cache"
    "golang-im/pkg/grpclib"
    "golang-im/pkg/rpc"
    "strconv"

    "google.golang.org/grpc/metadata"

        "github.com/golang/protobuf/proto"
    "go.uber.org/zap"

    //"golang-im/pkg/grpclib"
    "golang-im/pkg/logger"
    "golang-im/pkg/pb"

    //"golang-im/pkg/rpc"
    "golang-im/pkg/gerrors"
)

const (
    PushRoomTopic = "push_room_topic" // 房间消息队列
    PushAllTopic  = "push_all_topic"  // 全服消息队列
)

type messageService struct{}

var MessageService = new(messageService)

// SendOne 一对一消息发送
func (s *messageService) SendOne(ctx context.Context, msg *pb.PushMsgReq) error {
    // 1.拿到房间号去redis 取该房间下的所有人
    // 2.拿到每个人所在节点 new 1个grpc客户端
    // 3.grpc请求到每个人对应节点去

    requestId := grpclib.GetCtxRequestId(ctx)
    UserId, DeviceId, err := grpclib.GetCtxData(ctx)
    if err != nil {
        logger.Sugar.Infow("logic 服务 SendMessage 头信息error")
        return err
    }
    // 加上grpc头防止api授权拦截
    toConnectCtx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
        "user_id", strconv.FormatInt(UserId, 10),
        "device_id", strconv.FormatInt(DeviceId, 10),
        "token", "md5/jwt/xxx",
        "request_id", strconv.FormatInt(requestId, 10)))

    client, err := rpc.InitConnectIntClient(msg.Message.Server)
    if err != nil {
        return err
    }
    client.DeliverMessage(toConnectCtx, msg)

    return nil
}

// SendRoom 群组消息发送
func (s *messageService) SendRoom(ctx context.Context, msg *pb.PushMsg) error {
    logger.Logger.Debug("Push", zap.Any("msg", msg))
    bytes, err := proto.Marshal(msg)
    if err != nil {
        return gerrors.WrapError(err)
    }
    err = cache.Queue.Publish(PushAllTopic, bytes)
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
