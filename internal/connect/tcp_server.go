package connect

import (
	"golang-im/config"
	"golang-im/pkg/logger"
	"time"

	"golang-im/pkg/gn"

	"go.uber.org/zap"
)

var encoder = gn.NewHeaderLenEncoder(2, 1024)

var server *gn.Server

// StartTCPServer 启动gn TCP框架 监听端口
func StartTCPServer() {
	gn.SetLogger(logger.Sugar)

	var err error
	server, err = gn.NewServer(config.Connect.TCPListenAddr, &handler{},
		gn.NewHeaderLenDecoder(4),
		//限制了客户端发送数据的最大长度, 好处是采用sync.pool 内存复用
		//比如申请1024个字节长度 第一次使用了169字节，第二次使用16个字节,则第二次的16字节覆盖第一次169字节的前面 0-16字节 我们利用偏移只取前16字节即可
		//参考 https://mp.weixin.qq.com/s/6Nx7IGFU_FbM5AOdUzmvcw
		gn.WithReadBufferLen(65535),
		gn.WithTimeout(5*time.Minute, 11*time.Minute),
		gn.WithAcceptGNum(10),
		gn.WithIOGNum(100))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	server.Run()
}

type handler struct{}

var Handler = new(handler)

func (*handler) OnConnect(c *gn.Conn) {
	// 初始化连接数据
	conn := &Conn{
		CoonType:          CoonTypeTCP,
		TCP:               c,
		LastHeartbeatTime: time.Now(),
	}
	c.SetData(conn)
	logger.Logger.Debug("OnConnect", zap.Int32("fd", c.GetFd()), zap.String("addr", c.GetAddr()))
}

func (h *handler) OnMessage(c *gn.Conn, bytes []byte) {
	conn := c.GetData().(*Conn)
	conn.HandleMessage(bytes)
}

func (*handler) OnClose(c *gn.Conn, err error) {
	conn := c.GetData().(*Conn)
	logger.Logger.Debug("OnClose", zap.String("addr", c.GetAddr()), zap.Int64("user_id", conn.UserId),
		zap.String("device_id", conn.DeviceId), zap.Error(err))

	conn.Close()
}
