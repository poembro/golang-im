package connect

import (
	"context"
	"golang-im/pkg/grpclib"
	"golang-im/pkg/logger"

	//"fmt"
	//"golang-im/pkg/grpclib"
	//"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	//"go.uber.org/zap"
	//"golang-im/pkg/protocol"
)

type ConnIntServer struct{}

// DeliverMessage 投递消息
func (s *ConnIntServer) DeliverMessage(ctx context.Context, pushMsg *pb.PushMsgReq) (*pb.PushMsgReply, error) {
	logger.Sugar.Infow("DeliverMessage", "method", "投递消息", "requeset_id", grpclib.GetCtxRequestId(ctx))
	resp := &pb.PushMsgReply{}

	Dispatch(pushMsg.Message)

	return resp, nil
}
