package connect

import (
	"fmt"
	"golang-im/config"
	"golang-im/pkg/db"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"golang-im/pkg/protocol"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

const (
	PushRoomTopic = "push_room_topic" // 房间消息队列
	PushAllTopic  = "push_all_topic"  // 全服消息队列
)

// StartSubscribe 将redis的数据 推送到全局Map
func StartSubscribe() {
	channel := db.RedisCli.Subscribe(PushRoomTopic, PushAllTopic).Channel()
	for i := 0; i < config.Connect.SubscribeNum; i++ {
		go handleRedisMsg(channel)
	}
}

func handleRedisMsg(channel <-chan *redis.Message) {
	for msg := range channel {
		if msg.Channel == PushAllTopic {
			pushMsg := new(pb.PushMsg)
			err := proto.Unmarshal([]byte(msg.Payload), pushMsg)
			if err != nil {
				logger.Logger.Debug("StartSubscribe", zap.Error(err))
				continue
			}
			logger.Logger.Debug("RedisCli_Subscribe_msg", zap.Any("body", pushMsg))
			Dispatch(pushMsg)
		}
	}
}

// Dispatch 下发消息
func Dispatch(m *pb.PushMsg) {
	switch m.Type {
	case pb.PushMsg_PUSH:
		_pushKeys(m.Operation, m.Server, m.DeviceId, m.Msg)
	case pb.PushMsg_ROOM:
		_pushRoom(m.Operation, m.RoomId, m.Msg)
	case pb.PushMsg_BROADCAST:
		_pushAll(m.Operation, m.Msg, m.Speed)
	default:
		err := fmt.Errorf("no match push type: %s", m.Type)
		logger.Logger.Debug("handlePushAll", zap.Any("err", err))
	}
}

func _pushKeys(op int32, serverID string, DeviceId []string, body []byte) (err error) {
	// TODO 如果当前节点 与 serverID 不相等直接return

	for _, key := range DeviceId {
		// 获取设备对应的TCP连接
		conn := GetConn(key)
		if conn == nil {
			logger.Logger.Warn("GetConn warn", zap.String("device_id", key))
			return nil
		}

		if conn.DeviceId != key {
			logger.Logger.Warn("GetConn warn", zap.String("device_id", key))
			return nil
		}

		p := new(protocol.Proto)
		p.Op = op
		p.Body = body
		conn.Send(p)
	}

	return
}

func _pushRoom(op int32, roomid string, body []byte) {
	PushRoom(roomid, &protocol.Proto{
		Op:   op,
		Body: body,
	})
}

func _pushAll(op int32, body []byte, speed int32) {
	//TODO 如果单个节点连接数太多 需要用time.Sleep(Speed/节点数)间隔一下
	PushAll(speed, &protocol.Proto{
		Op:   op,
		Body: body,
	})
}
