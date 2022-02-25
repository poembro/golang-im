package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang-im/pkg/grpclib/etcdv3"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc/balancer/roundrobin"
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

	// grpc 服务名称
	LogicIntSerName   = "logicint_grpc_service"
	ConnectIntSerName = "connectint_grpc_service"

	// 连接句柄
	logicIntClient pb.LogicIntClient

	cli *client
)

/////////////////////////////对外////////////////////////////

// LogicInt grpc server 服务名方式 grpc自动轮询
func LogicInt() pb.LogicIntClient {
	if logicIntClient == nil {
		newServiceNameGrpc(cli.Schema, cli.EtcdAddr, LogicIntSerName)
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
			grpc.WithUnaryInterceptor(interceptor), // 一元拦截器，适用于普通rpc连接
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		}...)

	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	logicIntClient = pb.NewLogicIntClient(conn)
}

// ConnectInt grpc server 访问指定节点
func ConnectInt(addr string) pb.ConnectIntClient {
	conn := cli.GetConn(ConnectIntSerName, addr)
	if conn == nil {
		err := fmt.Errorf("grpc client failed %s 不在线", addr)
		logger.Logger.Error("LogicInt", zap.Any("conn", err))
		return nil
	}
	return pb.NewConnectIntClient(conn)
}

func newGrpc(addr string) (*grpc.ClientConn, error) {
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

	return conn, err
}

// 初始化
func Init(schema, etcdAddr, srvName, srvIpPort string) {
	// 服务注册至ETCD
	host, port, err := net.SplitHostPort(srvIpPort)
	if err != nil {
		fmt.Println("port not int error ")
		return
	}
	etcdv3.RegisterEtcd(schema, etcdAddr, host, port, srvName, 5)

	// 初始化grpc客户端
	cli = &client{
		Conns:    make(map[string]*grpc.ClientConn),
		AllConns: make(map[string]string),
		Schema:   schema,
		EtcdAddr: etcdAddr,
	}

	// 去etcd 检查某服务所有在线节点.
	// 注意:在logic服务里调用connect服务时,检查 connect 服务节点
	if srvName == LogicIntSerName {
		cli.HealthCheck(ConnectIntSerName)
	}
}

/////////////////////////////对内////////////////////////////

type client struct {
	Conns    map[string]*grpc.ClientConn
	Lock     sync.RWMutex
	AllConns map[string]string
	Schema   string
	EtcdAddr string
}

func (c *client) GetConn(srvName, addr string) *grpc.ClientConn {
	var (
		ok  bool
		key string
		err error
		r   *grpc.ClientConn
	)
	key = etcdv3.GetPrefix(c.Schema, srvName) + addr

	//先判断是否已经有了该节点
	c.Lock.RLock()
	if r, ok = c.Conns[key]; ok {
		c.Lock.RUnlock()
		return r
	}
	c.Lock.RUnlock()
	logger.Logger.Debug("GetConn", zap.String("desc", "c.Conns中没有对应连接句柄"), zap.Any("key", key), zap.Any("c.Conns", c.Conns))
	c.Lock.Lock()

	// 判断节点是否在线
	if _, ok := c.AllConns[key]; !ok {
		c.Lock.Unlock()
		err = fmt.Errorf("将要访问的节点 %s 不在线", addr)
		logger.Logger.Error("GetConn", zap.Any("key", key), zap.Error(err))
		return nil
	}
	r, err = newGrpc(addr)

	if err != nil {
		c.Lock.Unlock()
		return nil
	}

	c.Conns[key] = r
	c.Lock.Unlock()
	return r
}

func (c *client) HealthCheck(srvName string) {
	ck := etcdv3.NewHealthCheck(c.Schema, c.EtcdAddr, srvName)
	ns := ck.GetEtcdConns()
	for k, v := range ns {
		c.AllConns[k] = v
	}

	// checkTimeout 定时检查在线的connect节点
	timeoutTicker := 1 * time.Minute // 每隔1分钟检查一次在线节点
	go func() {
		ticker := time.NewTicker(timeoutTicker)
		for {
			select {
			case <-ticker.C:
				ns = ck.GetEtcdConns()
				if len(ns) > 0 {
					c.handleConn(ns)
				}
			}
		}
	}()
}

// handleConn 处理服务节点
func (c *client) handleConn(ns map[string]string) {
	// etcd  GetAllService ---k-> goim:///connectint_grpc_service/192.168.83.165:50000  ---v-> 192.168.83.165:50000
	// 遍历旧节点,判断每个节点是否都在最新的etcd里面, 如果不在则 剔除掉旧的 grpc 连接
	for k, _ := range c.AllConns {
		// 判断旧的节点 是否在最新的etcd 服务节点里
		if _, ok := ns[k]; !ok {
			delete(c.AllConns, k)
			// 剔除连接句柄
			c.Lock.Lock()
			delete(c.Conns, k)
			c.Lock.Unlock()
		}
	}

	// 写入新的节点进来
	for k, v := range ns {
		c.AllConns[k] = v
	}
}
