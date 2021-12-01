package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"

	//	"fmt"
	"golang-im/pkg/grpclib/etcdv3"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"time"

	"google.golang.org/grpc"
//	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
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
func InitLogicIntClient(schema, etcdaddr string) {
	rr := etcdv3.NewDiscovery(schema, etcdaddr, LogicIntSerName)
	resolver.Register(rr) //向resolver/resolver.go 中m变量追加参数和值 m[b.Scheme()] = b

	conn, err := grpc.DialContext(
		context.TODO(),
		etcdv3.GetPrefix(schema, LogicIntSerName),
		[]grpc.DialOption{
			grpc.WithInsecure(), //禁用传输认证，没有这个选项必须设置一种认证方式
			grpc.WithTimeout(time.Duration(5) * time.Second),
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
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		}...)

	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	LogicIntClient = pb.NewLogicIntClient(conn)
}

// InitConnectIntClient logic 访问 connect 服务最好知道是哪个节点
func InitConnectIntClient(addr string) (pb.ConnectIntClient, error) {
	// TODO 判断该addr 是否在服务发现列表里面
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
