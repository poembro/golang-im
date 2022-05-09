package connect

import (
	"context"

	"golang-im/pkg/pb"
)

type ConnIntServer struct{}

func NewConnIntServer() *ConnIntServer {
	return &ConnIntServer{}
}

// DeliverMessage 投递消息
func (s *ConnIntServer) DeliverMessage(ctx context.Context, pushMsg *pb.PushMsgReq) (*pb.PushMsgReply, error) {
	resp := &pb.PushMsgReply{}
	Dispatch(pushMsg.Message)
	return resp, nil
}
