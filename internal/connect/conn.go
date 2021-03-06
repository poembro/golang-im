package connect

import (
	"container/list"
	"context"
	"fmt"
	"golang-im/conf"
	"golang-im/pkg/gn"
	"golang-im/pkg/grpclib"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"golang-im/pkg/protocol"
	"golang-im/pkg/rpc"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	CoonTypeTCP int8 = 1 // tcp连接
	ConnTypeWS  int8 = 2 // websocket连接
)

type Conn struct {
	CoonType          int8 // 连接类型
	Config            *conf.Config
	TCP               *gn.Conn        // tcp连接
	FdMutex           sync.Mutex      // 写锁
	WS                *websocket.Conn // websocket连接
	UserId            int64           // 用户ID
	DeviceId          string          // 设备ID
	RoomId            string          // 订阅的房间ID
	Element           *list.Element   // 链表节点
	LastHeartbeatTime time.Time       // 最后一次读取数据的时间
}

// Close 关闭
func (c *Conn) Close() error {
	// 取消订阅，需要异步出去，防止重复加锁造成死锁
	go func() {
		logger.Logger.Debug("Close", zap.String("DeviceId", c.DeviceId))
		SubscribedRoom(c, "")
	}()

	// 取消设备和连接的对应关系
	if c.DeviceId != "" {
		DeleteConn(c.DeviceId)

		// 通知Logic服务谁已经下线
		_, _ = rpc.LogicInt(c.Config).Offline(context.TODO(), &pb.OfflineReq{
			UserId:     c.UserId,
			DeviceId:   c.DeviceId,
			ClientAddr: c.GetAddr(),
		})
	}

	if c.CoonType == CoonTypeTCP {
		return c.TCP.Close()
	} else if c.CoonType == ConnTypeWS {
		return c.WS.Close()
	}
	return nil
}

func (c *Conn) GetAddr() string {
	if c.CoonType == CoonTypeTCP {
		return c.TCP.GetAddr()
	} else if c.CoonType == ConnTypeWS {
		return c.WS.RemoteAddr().String()
	}
	return ""
}

// Write 写入数据
func (c *Conn) Write(bytes []byte) error {
	c.FdMutex.Lock()
	defer c.FdMutex.Unlock()

	if c.CoonType == CoonTypeTCP {
		return c.TCP.WriteWithEncoder(bytes)
	} else if c.CoonType == ConnTypeWS {
		err := c.WS.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
		if err != nil {
			return err
		}
		return c.WS.WriteMessage(websocket.BinaryMessage, bytes)
	}
	logger.Logger.Error("unknown conn type", zap.Any("conn", c))
	return nil
}

// Send 下发消息
func (c *Conn) Send(output *protocol.Proto) {
	outputBytes, err := output.Encode()
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	msg := ""
	switch output.Op {
	case protocol.OpAuthReply:
		msg = "回登录"
	case protocol.OpHeartbeatReply:
		msg = "回心跳"
	case protocol.OpMessageAckReply:
		msg = "回收到偏移上报"
	case protocol.OpSyncReply:
		msg = "回同步历史消息"
	case protocol.OpSendMsgReply:
		msg = "回收到客户端发来的聊天消息"
	case protocol.OpSubReply:
		msg = "回收到订阅房间"
	case protocol.OpUnsubReply:
		msg = "回收到取消订阅房间"
	case protocol.OpChangeRoomReply:
		msg = "回修改房间"
	default:
	}
	logger.Logger.Debug("Send", zap.String("desc", fmt.Sprintf("op:%d msg:%s", output.Op, msg)), zap.String("body", string(output.Body)))

	err = c.Write(outputBytes)
	if err != nil {
		logger.Sugar.Error(err)
		c.Close()
		return
	}
}

func (c *Conn) HandleMessage(bytes []byte) {
	var (
		input *protocol.Proto
	)
	input = new(protocol.Proto)
	input.Decode(bytes)

	logger.Logger.Info("HandleMessage", zap.String("desc", fmt.Sprintf("op:%d msg:%s", input.Op, "收到")), zap.Any("input", string(input.Body)))

	// 对未登录的用户进行拦截
	if input.Op != protocol.OpAuth && c.UserId == 0 {
		// TODO 应该告诉用户没有登录
		return
	}

	switch input.Op {
	case protocol.OpAuth:
		c.SignIn(input)
	case protocol.OpSync:
		c.Sync(input)
	case protocol.OpHeartbeat:
		c.Heartbeat(input)
	case protocol.OpSendMsg:
		c.OpSendMsg(input) // OpSendMsgReply
	case protocol.OpMessageAck:
		c.MessageACK(input)
	case protocol.OpSub:
		c.OpSub(input) // OpSubReply
	case protocol.OpUnsub:
		c.OpUnsub(input) // OpUnsubReply
	case protocol.OpChangeRoom:
		c.OpChangeRoom(input) // OpChangeRoomReply
	default:
		logger.Logger.Error("handler switch other")
	}

	return
}

//MessageACK 消息偏移上报
func (c *Conn) MessageACK(p *protocol.Proto) {
	tmp := string(p.Body)
	index, _ := strconv.Atoi(tmp)
	_, err := rpc.LogicInt(c.Config).MessageACK(grpclib.ContextWithRequstId(context.TODO(), int64(p.Seq)), &pb.MessageACKReq{
		UserId:      c.UserId,
		DeviceId:    c.DeviceId,
		RoomId:      c.RoomId,
		DeviceAck:   int64(index), //这里需要 body转int64 标识已经读到哪里了
		ReceiveTime: time.Now().UnixNano(),
	})
	if err != nil {
		return
	}
	p.Op = protocol.OpMessageAckReply
	p.Body = nil
	c.Send(p)
}

// Sync 同步历史聊天记录
func (c *Conn) Sync(p *protocol.Proto) {
	resp, err := rpc.LogicInt(c.Config).Sync(grpclib.ContextWithRequstId(context.TODO(), int64(p.Seq)), &pb.SyncReq{
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		RoomId:   c.RoomId,
		Seq:      2,
	})
	if err != nil {
		return
	}
	p.Op = protocol.OpSyncReply
	p.Body = resp.Messages
	c.Send(p)
}

// SignIn 登录
func (c *Conn) SignIn(p *protocol.Proto) {
	resp, err := rpc.LogicInt(c.Config).ConnSignIn(grpclib.ContextWithRequstId(context.TODO(), int64(p.Seq)), &pb.ConnSignInReq{
		Body:       p.Body,
		ConnAddr:   c.Config.Connect.LocalAddr,
		ClientAddr: c.GetAddr(),
	})
	if err != nil {
		logger.Logger.Debug("SignIn", zap.String("error", err.Error()))
		return
	}
	p.Op = protocol.OpAuthReply
	p.Body = []byte("ok")
	c.Send(p)

	c.UserId = resp.UserId
	c.DeviceId = resp.DeviceId
	SetConn(resp.DeviceId, c)
}

// OpSendMsg 接收客户端发来的消息
func (c *Conn) OpSendMsg(p *protocol.Proto) {
	if c.RoomId == "" {
		p.Op = protocol.OpSendMsgReply
		p.Body = []byte("Not subscribing to a room")
		c.Send(p)
		return
	}
	buf := &pb.PushMsg{
		Type:      pb.PushMsg_ROOM, //  测试时改 pb.PushMsg_PUSH 并给 DeviceId 属性赋值
		Operation: protocol.OpSendMsgReply,
		Speed:     2,
		Server:    c.Config.Connect.LocalAddr,
		RoomId:    c.RoomId,
		Msg:       p.Body,
		DeviceId:  []string{c.DeviceId},
	}
	// 加上grpc头防止api授权拦截
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", strconv.FormatInt(c.UserId, 10),
		"device_id", c.DeviceId,
		"token", "TODO token verify",
		"request_id", fmt.Sprintf("%d", p.Seq)))

	_, err := rpc.LogicInt(c.Config).SendMessage(ctx, &pb.PushMsgReq{Message: buf})
	if err != nil {
		return
	}
}

// Heartbeat 心跳
func (c *Conn) Heartbeat(p *protocol.Proto) {
	c.LastHeartbeatTime = time.Now()
	_, err := rpc.LogicInt(c.Config).Heartbeat(grpclib.ContextWithRequstId(context.TODO(), int64(p.Seq)), &pb.HeartbeatReq{
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		ConnAddr: c.Config.Connect.LocalAddr,
	})
	if err != nil {
		logger.Logger.Debug("Heartbeat", zap.Error(err))
		return
	}

	p.Op = protocol.OpHeartbeatReply
	p.Body = nil
	c.Send(p)
}

// OpSub 订阅房间
func (c *Conn) OpSub(p *protocol.Proto) {
	logger.Logger.Info("OpSub", zap.String("desc", "订阅房间"), zap.String("input", string(p.Body)))
	SubscribedRoom(c, string(p.Body))
	p.Op = protocol.OpSubReply
	p.Body = []byte("subscribed room ok")
	c.Send(p)
	/*
	   TODO
	   1.接收房间号
	   2.丢给logic 服务 验证房间号 防止串台
	*/
}

// OpUnsub 取消订阅房间
func (c *Conn) OpUnsub(p *protocol.Proto) {
	SubscribedRoom(c, "")
	p.Op = protocol.OpUnsubReply
	p.Body = []byte("unsubscribed room ok")
	c.Send(p)
}

// OpChangeRoom 修改房间
func (c *Conn) OpChangeRoom(p *protocol.Proto) {
	/*
	   TODO
	   1.接收房间号
	   2.丢给logic 服务 验证房间号 防止串台
	*/
	logger.Logger.Info("OpChangeRoom", zap.String("desc", "订阅房间"), zap.String("input", string(p.Body)))

	SubscribedRoom(c, "")
	SubscribedRoom(c, string(p.Body))

	p.Op = protocol.OpChangeRoomReply
	p.Body = []byte("subscribed room ok")
	c.Send(p)
}
