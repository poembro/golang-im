package rpc

import (
	"context"
	"fmt"
	"golang-im/conf"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"golang-im/pkg/rpc/etcdv3"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcKeepAliveTime         = time.Duration(10) * time.Second
	grpcKeepAliveTimeout      = time.Duration(3) * time.Second
	grpcBackoffMaxDelay       = time.Duration(3) * time.Second
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
)

// 连接句柄
var logicIntClient pb.LogicIntClient = nil

/////////////////////////////对外////////////////////////////

// LogicInt grpc server 服务名方式,grpc自动轮询
func LogicInt(c *conf.Config) pb.LogicIntClient {
	if logicIntClient == nil {
		newServiceNameGrpc(c.Global.GrpcSchema, c.Global.EtcdAddr, c.Logic.LogicIntSerName)
	}
	return logicIntClient
}

func newServiceNameGrpc(schema, etcdAddr, serviceName string) {
	rr := etcdv3.NewDiscovery(schema, etcdAddr, serviceName)
	resolver.Register(rr) //向resolver/resolver.go 中m变量追加参数和值 m[b.Scheme()] = b

	conn, err := grpc.DialContext(
		context.TODO(),
		etcdv3.GetPrefix4Unique(schema, serviceName),
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
			grpc.WithUnaryInterceptor(interceptor), // 拦截器，适用于普通rpc连接
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		}...)

	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	logicIntClient = pb.NewLogicIntClient(conn)
}

//////////////////////////////////////////////////////////////////////

var connClients map[string]*grpc.ClientConn

// ConnectInt grpc server 访问指定connect服务节点
func ConnectInt(addr string) pb.ConnectIntClient {
	var (
		ok bool
		c  *grpc.ClientConn
	)

	c, ok = connClients[addr]
	if ok {
		return pb.NewConnectIntClient(c)
	} else {
		c = newGrpc(addr)
		if c == nil {
			logger.Sugar.Error(fmt.Errorf("connect 节点%s不在线", addr))
			return nil
		}
		connClients[addr] = c
		return pb.NewConnectIntClient(c)
	}
}

// ConnectIntSrv()  服务名方式,grpc自动轮询 TODO
func newGrpc(addr string) *grpc.ClientConn {
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
			grpc.WithKeepaliveParams(keepalive.ClientParameters{ //自动重连，所以某节点挂了后，再恢复并不会影响
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
			grpc.WithUnaryInterceptor(interceptor),
		}...)

	if err != nil {
		logger.Sugar.Error(err)
		return nil
	}

	return conn
}
