package main

import (
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

    "google.golang.org/grpc/keepalive"

    "golang-im/config"
    "golang-im/internal/logic/api"
    "golang-im/pkg/db"
    "golang-im/pkg/interceptor"
    "golang-im/pkg/logger"
    "golang-im/pkg/pb"
    "golang-im/pkg/rpc"
    "golang-im/pkg/urlwhitelist"

    "go.uber.org/zap"
    "google.golang.org/grpc"

    "golang-im/pkg/grpclib/etcdv3"
)

func main() {
    logger.Init()
    // db.InitEtcd(config.Global.EtcdAddr)
    // db.InitMysql(config.Logic.MySQL)
    db.InitRedis(config.Global.RedisIP, config.Global.RedisPassword)

    // 初始化RpcClient
	rpc.NewClient(config.Global.GrpcSchema, config.Global.EtcdAddr, rpc.ConnectIntSerName)

    keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
        MaxConnectionIdle:     time.Duration(time.Second * 60), //60s 连接最大闲置时间
        MaxConnectionAgeGrace: time.Duration(time.Second * 20), //20s 连接最大闲置时间
        Time:                  time.Duration(time.Second * 60), //60s
        Timeout:               time.Duration(time.Second * 20), //20s
        MaxConnectionAge:      time.Duration(time.Hour * 2),    //2h  小时
    })
    server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("logic_int_interceptor", urlwhitelist.Logic)), keepParams)

    // 监听服务关闭信号，服务平滑重启
    // 其实你是担心直接重启服务, 会有处理到一半的请求被中断了, 导致尴尬的局面.你要的并不是热重启, 而是优雅关闭.
    // grpc框架支持优雅关闭的.基本原理是, 你监听一个信号, 收到信号时调用grpc的GracefulStop接口, 这时grpc会首先关闭对外监听的fd, 这时就不会有新的请求进来.
    // 而已经在处理的请求则会继续处理完, 然后再关闭服务.在grpc关闭对外监听的fd后的那个瞬间, 你其实可以启动你的新程序了, 所以基本上中断时间很短, 而原来处理着的请求并不会有问题.
    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGTERM)
        s := <-c
        logger.Logger.Info("server stop", zap.Any("signal", s))
        server.GracefulStop()
    }()

    pb.RegisterLogicIntServer(server, &api.LogicIntServer{})
    listen, err := net.Listen("tcp", config.Logic.RPCListenAddr)
    if err != nil {
        panic(err)
    }

    err = etcdv3.Register(config.Global.GrpcSchema, config.Global.EtcdAddr, config.Logic.LocalAddr, rpc.LogicIntSerName, 5)
    if err != nil {
        logger.Logger.Error("register service err ", zap.Error(err))
    }

    logger.Logger.Info("rpc服务已经开启", zap.String("logic_rpc_server_ip_port", config.InternalIP()+config.Logic.RPCListenAddr))
    err = server.Serve(listen)
    if err != nil {
        logger.Logger.Error("Serve error", zap.Error(err))
    }
}
