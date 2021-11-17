package connect

import (
	"golang-im/pkg/protocol"
	"sync"
	"sync/atomic"
	"time"
)

var ConnsManager = sync.Map{}

// SetConn 存储
func SetConn(deviceId int64, conn *Conn) {
	ConnsManager.Store(deviceId, conn)
}

// GetConn 获取
func GetConn(deviceId int64) *Conn {
	value, ok := ConnsManager.Load(deviceId)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除
func DeleteConn(deviceId int64) {
	ConnsManager.Delete(deviceId)
}

// PushAll 给所有人推送消息
func PushAll(speed int32, p *protocol.Proto) {
	var (
		ops int32
	)

	ConnsManager.Range(func(key, value interface{}) bool {
		if ops%1024 == speed {
			time.Sleep(time.Duration(speed) * time.Second)
		}

		atomic.AddInt32(&ops, 1)

		conn := value.(*Conn)
		conn.Send(p)
		return true
	})
}
