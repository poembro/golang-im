package main

import (
	"context"
	"golang-im/config"
	"golang-im/internal/connect"
	"golang-im/pkg/db"
	"golang-im/pkg/interceptor"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"golang-im/pkg/rpc"
	"golang-im/pkg/urlwhitelist"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	logger.Init()

	db.InitRedis(config.Global.RedisIP, config.Global.RedisPassword)

	// 启动TCP长链接服务器
	go func() {
		//connect.StartTCPServer()
	}()

	// 启动WebSocket长链接服务器
	go func() {
		connect.StartWSServer(config.Connect.WSListenAddr)
	}()

	// 启动服务订阅
	connect.StartSubscribe()

	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(time.Second * 60), //60s 连接最大闲置时间
		MaxConnectionAgeGrace: time.Duration(time.Second * 20), //20s 连接最大闲置时间
		Time:                  time.Duration(time.Second * 60), //60s
		Timeout:               time.Duration(time.Second * 20), //20s
		MaxConnectionAge:      time.Duration(time.Hour * 2),    //2h  小时
	})
	//UnaryInterceptor 返回一个为服务器设置 UnaryServerInterceptor 的 ServerOption 。只能安装一个一元拦截器。多个拦截器(如链接)的构造可以在调用者处实现。
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("connect_int_interceptor", urlwhitelist.Connect)), keepParams)

	// 监听服务关闭信号，grpc服务优雅的关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("server stop start", zap.Any("signal", s))
		_, err := rpc.LogicInt().ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.Connect.LocalAddr})
		if err != nil {
			panic(err)
		}
		logger.Logger.Info("server stop end")

		server.GracefulStop()
	}()

	pb.RegisterConnectIntServer(server, &connect.ConnIntServer{})
	listener, err := net.Listen("tcp", config.Connect.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	// 初始化RpcClient
	rpc.Init(config.Global.GrpcSchema, config.Global.EtcdAddr, rpc.ConnectIntSerName, config.Connect.LocalAddr)

	logger.Logger.Info("rpc服务已经开启", zap.String("EtcdAddr", config.Global.EtcdAddr), zap.String("connect_rpc_server_ip_port", config.InternalIP()+config.Connect.RPCListenAddr))
	err = server.Serve(listener)
	if err != nil {
		logger.Logger.Error("Serve", zap.Error(err))
	}
}
