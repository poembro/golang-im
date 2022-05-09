package connect

import (
	"golang-im/pkg/protocol"
	"sync"
	"sync/atomic"
	"time"
)

var ConnsManager = sync.Map{}

// SetConn 存储
func SetConn(deviceId string, conn *Conn) {
	ConnsManager.Store(deviceId, conn)
}

// GetConn 获取
func GetConn(deviceId string) *Conn {
	value, ok := ConnsManager.Load(deviceId)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// DeleteConn 删除
func DeleteConn(deviceId string) {
	ConnsManager.Delete(deviceId)
}

// 下发频率限制
var Ops = uint64(0)

// PushAll 给所有人推送消息
func PushAll(speed int32, p *protocol.Proto) {
	ConnsManager.Range(func(key, value interface{}) bool {
		opsFinal := atomic.AddUint64(&Ops, 1)
		if opsFinal%1024 == 0 {
			time.Sleep(time.Duration(speed) * time.Second)
		}

		conn := value.(*Conn)
		conn.Send(p)
		return true
	})
}
