package connect

import (
	"golang-im/pkg/logger"
	"golang-im/pkg/util"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWsConn(fd *websocket.Conn) *Conn {
	return &Conn{
		CoonType:          ConnTypeWS,
		Config:            config,
		WS:                fd,
		LastHeartbeatTime: time.Now(),
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	conn := NewWsConn(wsConn)
	DoConn(conn)
}

// DoConn 处理连接
func DoConn(conn *Conn) {
	defer util.RecoverPanic()

	for {
		err := conn.WS.SetReadDeadline(time.Now().Add(2 * time.Minute))
		_, data, err := conn.WS.ReadMessage()
		if err != nil {
			HandleReadErr(conn, err)
			return
		}

		conn.HandleMessage(data)
	}
}

// HandleReadErr 读取conn错误
func HandleReadErr(conn *Conn, err error) {
	logger.Logger.Debug("read tcp error：", zap.Int64("user_id", conn.UserId),
		zap.String("device_id", conn.DeviceId), zap.Error(err))
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	conn.Close()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

func StartWSServer(WSListenAddr string) {
	http.HandleFunc("/ws", wsHandler)
	logger.Logger.Info("websocket server start")
	err := http.ListenAndServe(WSListenAddr, nil)
	if err != nil {
		panic(err)
	}
}
