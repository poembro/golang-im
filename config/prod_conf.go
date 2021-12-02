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
        EtcdAddr:    "http://" + ip + ":2379,http://" + ip + ":2479,http://" + ip + ":2579",
        RedisIP:       ip + ":6379",  //在使用
        RedisPassword: "",
    }

    // connect 服务相关配置
    Connect = ConnectConf{
        TCPListenAddr: ":8090",       //外部TCP 监听8080 在使用
        WSListenAddr:  ":7923",       //外部websocket监听8081 在使用
        RPCListenAddr: ":50000",      //内部connect grpc服务监听50000 在使用
        LocalAddr:     ip + ":50000", //connect服务本机局域网ip、端口,用来标识当前用户在哪个节点
        SubscribeNum:  100, //开启多少个groutine去redis取数据
    }

    // logic 服务相关配置
    Logic = LogicConf{
        MySQL:         "root:root@tcp(" + ip + ":3306)/default?charset=utf8&parseTime=true",
        RPCListenAddr: ":50100", //内部logic grpc服务监听50100 在使用
        LocalAddr:     ip + ":50100",
    }

    logger.Level = zap.DebugLevel
    logger.Target = logger.Console
}
