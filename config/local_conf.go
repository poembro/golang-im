package config

import (
	"golang-im/pkg/logger"
	"os"

	"go.uber.org/zap"
)

// 真正的监听端口 不变，外部参数
//docker run -p 50000:50000 -p 50100:50100 -p 7923:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.82.110:50100 --env GRPC_CONNECT_ADDR=192.168.82.110:50000 --rm golang-im:1.0.17
func initLocalConf() {
	grpcConnectAddr := os.Getenv("GRPC_CONNECT_ADDR")
	grpcLogicAddr := os.Getenv("GRPC_LOGIC_ADDR")
	if grpcConnectAddr == "" {
		grpcConnectAddr = ":50000"
	}
	if grpcLogicAddr == "" {
		grpcLogicAddr = ":50100"
	}

	// 全局配置
	Global = GlobalConf{
		ProjectName:   "golang-im 一个运行在[golang](#)上的实时通信软件。", //暂未使用
		GrpcSchema:    "goim",
		EtcdAddr:      "http://192.168.82.110:2379,http://192.168.82.110:2479,http://192.168.82.110:2579",
		RedisIP:       "10.0.41.145:6379", //在使用
		RedisPassword: "",
		PushAllTopic:  "push_all_topic", // 全服消息队列
	}

	// connect 服务相关配置
	Connect = ConnectConf{
		TCPListenAddr: ":8090",         //外部TCP 监听8080 在使用
		WSListenAddr:  ":7923",         //外部websocket监听8081 在使用
		RPCListenAddr: ":50000",        //内部connect grpc服务监听50000 在使用
		LocalAddr:     grpcConnectAddr, //connect服务本机局域网ip、端口,用来标识当前用户在哪个节点
		SubscribeNum:  100,             //开启多少个groutine去redis取数据
	}

	// logic 服务相关配置
	Logic = LogicConf{
		MySQL:         "root:root@tcp(192.168.82.36:3306)/default?charset=utf8&parseTime=true", //暂未使用
		RPCListenAddr: ":50100",                                                                //内部logic grpc服务监听50100 在使用
		LocalAddr:     grpcLogicAddr,
	}

	logger.Level = zap.DebugLevel
	logger.Target = logger.Console
}
