package connect

import (
	"golang-im/conf"
)

var (
	config *conf.Config
)

func New(c *conf.Config) {
	config = c
	// 启动服务订阅
	StartSubscribe(c)

	// 启动TCP长链接服务器
	go func() {
		//StartTCPServer(c.Connect.TCPListenAddr)
	}()

	// 启动WebSocket长链接服务器
	go func() {
		StartWSServer(c.Connect.WSListenAddr)
	}()
}
