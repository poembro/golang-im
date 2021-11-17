package api

import (
	"context"
	"fmt"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"strconv"

	//"golang-im/pkg/util"
	"testing"
	"time"

	//"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	logger.Init()
	fmt.Println("init logger")
}

func getCtx() context.Context {
	token := "0"
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", "2",
		"device_id", "1",
		"token", token,
		"request_id", strconv.FormatInt(time.Now().UnixNano(), 10)))
}

func getLogicIntClient() pb.LogicIntClient {
	conn, err := grpc.Dial("192.168.3.222:50100", grpc.WithInsecure())
	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}
	return pb.NewLogicIntClient(conn)
}

func TestLogicIntServer_SignIn(t *testing.T) {
	token := `{
		room_id:'live://1000', //将消息发送到指定房间
		accepts:'[1000,1001,1002]',//接收指定房间的消息
		key:'1xxxxxDeviceIdxxx假定为设备idxx', 

		user_id: '663291537152950273',
		user_name:'随机用户001',
		user_face:'/static/wap/img/portrait.jpg',

		shop_id:0,
		shop_name:'杂货铺老板', 
		shop_face:'/static/wap/img/portrait.jpg',

		platform:'web',
		
		suburl : "ws://192.168.3.222:9999/sub", 
		pushurl:"http://192.168.3.222:9999/open/push",
	}`

	resp, err := getLogicIntClient().ConnSignIn(context.TODO(),
		&pb.ConnSignInReq{
			Body:     []byte(token),
			ConnAddr: "127.0.0.1:5000",
		})
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	logger.Sugar.Info(resp)
}

func TestLogicIntServer_PushRoom(t *testing.T) {
	buf := &pb.PushMsg{
		Type:      pb.PushMsg_ROOM,
		Operation: 5,
		Speed:     2,
		Server:    "127.0.0.1:5000",
		RoomId:    "live://800020210817001",
		Msg:       []byte("hello world"),
	}

	resp, err := getLogicIntClient().SendMessage(getCtx(),
		&pb.PushMsgReq{
			Message: buf,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
	t.Errorf("%+v \n", resp)
}
