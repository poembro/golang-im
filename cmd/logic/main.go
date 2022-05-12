package main

import (
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"google.golang.org/grpc/keepalive"

	"golang-im/conf"
	"golang-im/internal/logic/apigrpc"
	"golang-im/internal/logic/apihttp"
	"golang-im/pkg/interceptor"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"golang-im/pkg/rpc/etcdv3"
	"golang-im/pkg/urlwhitelist"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Init()

	// 启动HTTP服务器
	go func() {
		apihttp.StartHttpServer(conf.Conf)
	}()

	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(time.Second * 60), //60s 连接最大闲置时间
		MaxConnectionAgeGrace: time.Duration(time.Second * 20), //20s 连接最大闲置时间
		Time:                  time.Duration(time.Second * 60), //60s
		Timeout:               time.Duration(time.Second * 20), //20s
		MaxConnectionAge:      time.Duration(time.Hour * 2),    //2h  小时
	})
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("logic_int_interceptor", urlwhitelist.Logic)), keepParams)

	pb.RegisterLogicIntServer(server, apigrpc.NewLogicIntServer(conf.Conf))
	listen, err := net.Listen("tcp", conf.Conf.Logic.RPCListenAddr)
	if err != nil {
		panic(err)
	}

	// 初始化RpcClient
	closeEtcd, err := Register(conf.Conf)
	if err != nil {
		panic(err)
	}

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
		closeEtcd()
	}()

	logger.Logger.Info("rpc服务已经开启",
		zap.String("EtcdAddr", conf.Conf.Global.EtcdAddr),
		zap.String("logic_rpc_server_ip_port",
			conf.InternalIP()+conf.Conf.Logic.RPCListenAddr))
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("Serve error", zap.Error(err))
	}
}

// 服务注册
func Register(c *conf.Config) (func(), error) {
	schema := c.Global.GrpcSchema
	etcdAddr := c.Global.EtcdAddr
	srvIpPort := c.Logic.LocalAddr
	srvName := c.Logic.LogicIntSerName

	// 服务注册至ETCD
	return etcdv3.RegisterEtcd(schema, etcdAddr, srvIpPort, srvName, 5)
}
