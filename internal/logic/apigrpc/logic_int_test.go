package apigrpc

import (
	"context"
	"golang-im/pkg/grpclib"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"

	//"golang-im/pkg/util"
	"testing"

	//"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	logger.Init()
}

func getCtx() context.Context {
	token := "0"
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", "2",
		"device_id", "1",
		"token", token,
		"request_id", "11111"))
}

func getLogicIntClient() pb.LogicIntClient {
	conn, err := grpc.Dial("192.168.84.168:50100", grpc.WithInsecure())
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}
	return pb.NewLogicIntClient(conn)
}

func TestLogicIntServer_SignIn(t *testing.T) {
	token := `{"created_at":"2022-04-29 15:54:24","device_id":"7241d0deb3fcf45dda85901acb59b1f1","face":"http://img.touxiangwu.com/2020/3/uq6Bja.jpg","nickname":"user193610","platform":"web","pushurl":"http://localhost:8090/open/push?&platform=web","referer":"http://192.168.84.168:8083/im.html?shop_id=13200000000","remote_addr":"192.168.84.168","room_id":"d339a209ccbaca713fa5407a79a3c17d","shop_face":"https://img.wxcha.com/m00/86/59/7c6242363084072b82b6957cacc335c7.jpg","shop_id":"13200000000","shop_name":"shop13200000000","suburl":"ws://localhost:7923/ws","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNDA1NDg3Mjk5OTU5MTkzNjEwIiwiZGV2aWNlX2lkIjoiNzI0MWQwZGViM2ZjZjQ1ZGRhODU5MDFhY2I1OWIxZjEiLCJuaWNrbmFtZSI6InVzZXIxOTM2MTAiLCJleHAiOjE2ODI3NTQ4NjQsImlzcyI6ImdvbGFuZ3Byb2plY3QifQ.IBBpWzjMBTjhskA5G1BpXv5hOux4WIsXicMqOtlgmYI","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:96.0) Gecko/20100101 Firefox/96.0","user_id":"405487299959193610"}`
	resp, err := getLogicIntClient().ConnSignIn(grpclib.ContextWithRequstId(context.TODO(), int64(23333)), &pb.ConnSignInReq{
		Body:       []byte(token),
		ConnAddr:   "127.0.0.1:5000",
		ClientAddr: "127.0.0.1",
	})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Info("结果:", resp)
}

func TestLogicIntServer_PushRoom(t *testing.T) {
	// 加上grpc头防止api授权拦截

	buf := &pb.PushMsg{
		Type:      pb.PushMsg_ROOM,
		Operation: 5,
		Speed:     2,
		Server:    "127.0.0.1:5000",
		RoomId:    "live://800020210817001",
		Msg:       []byte("hello world"),
		DeviceId:  []string{"1"},
	}

	_, err := getLogicIntClient().SendMessage(getCtx(), &pb.PushMsgReq{Message: buf})
	if err != nil {
		t.Errorf("结果:%s \n", err.Error())
		return
	}
	t.Errorf("结果:- \n")
}
