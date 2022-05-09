package conf

import (
	"golang-im/pkg/logger"
	"os"

	"go.uber.org/zap"
)

// 真正的监听端口 不变，外部参数
//docker run -p 50000:50000 -p 50100:50100 -p 7923:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.84.168:50100 --env GRPC_CONNECT_ADDR=192.168.84.168:50000 --rm golang-im:1.0.17
func initLocalConf() {
	grpcConnectAddr := os.Getenv("GRPC_CONNECT_ADDR")
	grpcLogicAddr := os.Getenv("GRPC_LOGIC_ADDR")
	if grpcConnectAddr == "" {
		grpcConnectAddr = "192.168.84.168:50000"
	}
	if grpcLogicAddr == "" {
		grpcLogicAddr = "192.168.84.168:50100"
	}

	Conf = &Config{
		// 全局配置
		Global: &GlobalConf{
			ProjectName:   "golang-im 一个运行在[golang](#)上的实时通信软件。", //暂未使用
			GrpcSchema:    "goim",
			EtcdAddr:      "http://10.0.41.145:2379,http://10.0.41.145:2479,http://10.0.41.145:2579",
			RedisIP:       "10.0.41.145:6379", //在使用
			RedisPassword: "",
			PushAllTopic:  "push_all_topic", // 全服消息队列
		},
		// logic 服务相关配置
		Logic: &LogicConf{
			HttpListenAddr:  ":8090",                                                                 //外部HTTP 监听8888 在使用
			MySQL:           "root:root@tcp(192.168.82.36:3306)/default?charset=utf8&parseTime=true", //暂未使用
			RPCListenAddr:   ":50100",
			LogicIntSerName: "logicint_grpc_service", // logic 服务名用来标识 logic grpc服务                                                     //内部logic grpc服务监听50100 在使用
			LocalAddr:       grpcLogicAddr,
		},
		// connect 服务相关配置
		Connect: &ConnectConf{
			TCPListenAddr:     ":6923",                   //外部TCP 监听8080 在使用
			WSListenAddr:      ":7923",                   //外部websocket监听8081 在使用
			RPCListenAddr:     ":50000",                  //内部connect grpc服务监听50000 在使用
			ConnectIntSerName: "connectint_grpc_service", // connect 服务名用来标识 connect grpc服务
			LocalAddr:         grpcConnectAddr,           //connect服务本机局域网ip、端口,用来标识当前用户在哪个节点
			SubscribeNum:      100,                       //开启多少个groutine去redis取数据
		},
	}
	logger.Level = zap.DebugLevel
	logger.Target = logger.Console
}
