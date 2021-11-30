package config

import (
	"golang-im/pkg/logger"

	"go.uber.org/zap"
)

func initProdConf(ip string) {
	// 全局配置
	Global = GlobalConf{
		ProjectName: "golang-im 一个运行在[golang](#)上的实时通信软件。", //暂未使用
		GrpcSchema:  "goim",
	}

	// connect 服务相关配置
	Connect = ConnectConf{
		TCPListenAddr: ":8090",       //外部TCP 监听8080 在使用
		WSListenAddr:  ":8081",       //外部websocket监听8081 在使用
		RPCListenAddr: ":50000",      //内部connect grpc服务监听50000 在使用
		LocalAddr:     ip + ":50000", //connect服务本机局域网ip、端口,用来标识当前用户在哪个节点
		RedisIP:       ip + ":6379",  //在使用
		RedisPassword: "",
		SubscribeNum:  100, //开启多少个groutine去redis取数据
	}

	// logic 服务相关配置
	Logic = LogicConf{
		MySQL:         "root:root@tcp(" + ip + ":3306)/default?charset=utf8&parseTime=true",
		RedisIP:       ip + ":6379", //在使用
		RedisPassword: "",
		RPCListenAddr: ":50100", //内部logic grpc服务监听50100 在使用
		LocalAddr:     ip + ":50100",
		EtcdIPs:       "http://127.0.0.1:2379,http://127.0.0.1:2479,http://127.0.0.1:2579",
	}

	logger.Level = zap.DebugLevel
	logger.Target = logger.Console
}
