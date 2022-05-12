package main

import (
	"context"
	"golang-im/conf"
	"golang-im/internal/connect"
	"golang-im/pkg/interceptor"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"golang-im/pkg/rpc"
	"golang-im/pkg/rpc/etcdv3"
	"golang-im/pkg/urlwhitelist"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Init()
	connect.New(conf.Conf)

	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(time.Second * 60), //60s 连接最大闲置时间
		MaxConnectionAgeGrace: time.Duration(time.Second * 20), //20s 连接最大闲置时间
		Time:                  time.Duration(time.Second * 60), //60s
		Timeout:               time.Duration(time.Second * 20), //20s
		MaxConnectionAge:      time.Duration(time.Hour * 2),    //2h  小时
	})
	//UnaryInterceptor 返回一个为服务器设置 UnaryServerInterceptor 的 ServerOption 。只能安装一个一元拦截器。多个拦截器(如链接)的构造可以在调用者处实现。
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("connect_int_interceptor", urlwhitelist.Connect)), keepParams)

	pb.RegisterConnectIntServer(server, connect.NewConnIntServer())
	listener, err := net.Listen("tcp", conf.Conf.Connect.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	// 初始化RpcClient
	closeEtcd, err := Register(conf.Conf)
	if err != nil {
		panic(err)
	}

	// 监听服务关闭信号，grpc 服务优雅的关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		s := <-c
		logger.Logger.Info("connect server stop ", zap.Any("signal", s))
		_, err := rpc.LogicInt(conf.Conf).ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: conf.Conf.Connect.LocalAddr})
		if err != nil {
			panic(err)
		}

		server.GracefulStop()

		closeEtcd()
		logger.Logger.Info("server stop end")
	}()

	logger.Logger.Info("rpc服务已经开启",
		zap.String("EtcdAddr", conf.Conf.Global.EtcdAddr),
		zap.String("connect_rpc_server_ip_port",
			conf.InternalIP()+conf.Conf.Connect.RPCListenAddr))
	err = server.Serve(listener)
	if err != nil {
		logger.Logger.Error("Serve", zap.Error(err))
	}
}

// 服务注册
func Register(c *conf.Config) (func(), error) {
	schema := c.Global.GrpcSchema
	etcdAddr := c.Global.EtcdAddr
	srvIpPort := c.Connect.LocalAddr
	srvName := c.Connect.ConnectIntSerName
	// 服务注册至ETCD
	return etcdv3.RegisterEtcd(schema, etcdAddr, srvIpPort, srvName, 5)
}
