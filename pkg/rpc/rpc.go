package rpc

import (
	"context"
	"fmt"
	"golang-im/config"
	"golang-im/pkg/grpclib/etcdv3"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc/resolver"
)

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
)

var (
	// grpc options
	grpcKeepAliveTime    = time.Duration(10) * time.Second
	grpcKeepAliveTimeout = time.Duration(3) * time.Second
	grpcBackoffMaxDelay  = time.Duration(3) * time.Second
	grpcMaxSendMsgSize   = 1 << 24
	grpcMaxCallMsgSize   = 1 << 24

	// 全局对象 用来访问对应grpc方法
	LogicIntClient   pb.LogicIntClient
	ConnectIntClient pb.ConnectIntClient

	// grpc 服务名称
	LogicIntSerName   = "logicint_grpc_service"
	ConnectIntSerName = "connectint_grpc_service"
)

// InitLogicIntClient connect访问logic不需要知道具体访问哪个节点
func InitLogicIntClient() {
	r := etcdv3.NewDiscovery(config.Logic.EtcdIPs)
	resolver.Register(r) //向resolver/resolver.go 中m变量追加参数和值 m[b.Scheme()] = b

	conn, err := grpc.DialContext(
		context.TODO(),
		fmt.Sprintf("%s:///%s", r.Scheme(), LogicIntSerName),
		[]grpc.DialOption{
			grpc.WithInsecure(), //禁用传输认证，没有这个选项必须设置一种认证方式
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
			grpc.WithUnaryInterceptor(interceptor), // 一元拦截器，适用于普通rpc连接
			//当前不走配置 grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, db.EtcdCli.Scheme())),
			grpc.WithBalancerName("weight"),
		}...)

	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	LogicIntClient = pb.NewLogicIntClient(conn)
}

// InitConnectIntClient logic 访问 connect 服务最好知道是哪个节点
func InitConnectIntClient(addr string) (pb.ConnectIntClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second))
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
			grpc.WithUnaryInterceptor(interceptor),
		}...)

	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	return pb.NewConnectIntClient(conn), err
}
