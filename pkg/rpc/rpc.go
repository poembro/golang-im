package rpc

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang-im/pkg/grpclib/etcdv3"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"

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

	Client *client
)

// ConnectInt grpc server 服务需要知道是哪个节点
func ConnectInt(addr string) pb.ConnectIntClient {
	conn := Client.GetConn(Client.Schema, Client.EtcdAddr, addr)
	if conn == nil {
		err := fmt.Errorf("grpc client failed %s 不在线", addr)
		logger.Logger.Error("LogicInt", zap.Any("conn", err))
		return nil
	}
	return pb.NewConnectIntClient(conn)
}

// LogicInt grpc server 服务不用知道哪个节点
func LogicInt() pb.LogicIntClient {
	conn := Client.GetConn(Client.Schema, Client.EtcdAddr, LogicIntSerName)
	if conn == nil {
		err := fmt.Errorf("grpc client failed %s 不在线", LogicIntSerName)
		logger.Logger.Error("LogicInt", zap.Any("conn", err))
		return nil
	}
	return pb.NewLogicIntClient(conn)
}

func newServiceNameGrpc(schema, etcdAddr, serviceName string) (*grpc.ClientConn, error) {
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

	return conn, err
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

type client struct {
	ClientConn map[string]*grpc.ClientConn
	Lock       sync.RWMutex
	AllService map[string]string
	Schema     string
	EtcdAddr   string
}

func NewClient(schema, etcdAddr, serviceName string) {
	Client = &client{
		ClientConn: make(map[string]*grpc.ClientConn),
		AllService: make(map[string]string),
		Schema:     schema,
		EtcdAddr:   etcdAddr,
	}

	// 去etcd 检查某服务所有在线节点.  注意:只针对非服务发现的
	if serviceName == ConnectIntSerName {
		Client.checkNode(schema, etcdAddr, ConnectIntSerName)
	}
}

func (c *client) GetConn(schema, etcdAddr, serviceName string) *grpc.ClientConn {
	var (
		err error
		ok  bool
		key string
		r   *grpc.ClientConn
	)
	key = etcdv3.GetPrefix4Unique(schema, serviceName)
	//先判断是否已经有了该节点
	c.Lock.RLock()
	if r, ok = c.ClientConn[key]; ok {
		c.Lock.RUnlock()
		return r
	}
	c.Lock.RUnlock()

	c.Lock.Lock()

	u, err := url.Parse(key) // 判断当前是否 ip:port
	if u.Port() != "" {
		// 判断节点是否在线
		if _, ok := c.AllService[serviceName]; !ok {
			err = fmt.Errorf("将要访问的节点 %s 不在线", serviceName)
			logger.Logger.Error("GetConn", zap.Any("conn", err))
			return nil
		}
		r, err = newGrpc(serviceName)
	} else {
		r, err = newServiceNameGrpc(schema, etcdAddr, serviceName)
	}

	if err != nil {
		c.Lock.Unlock()
		return nil
	}

	c.ClientConn[key] = r
	c.Lock.Unlock()
	return r
}

func (c *client) delConn(schema, serviceName string) {
	c.Lock.Lock()
	key := etcdv3.GetPrefix4Unique(schema, serviceName)
	delete(c.ClientConn, key)
	c.Lock.Unlock()
}

func (c *client) checkNode(schema, etcdAddr, serviceName string) {
	timeoutTicker := 1 * time.Minute // 每隔1分钟检查一次在线节点
	ns := etcdv3.GetAllService(schema, etcdAddr, serviceName)
	for k, v := range ns {
		c.AllService[k] = v
	}

	// checkTimeout 定时检查在线的connect节点
	go func() {
		ticker := time.NewTicker(timeoutTicker)
		for {
			select {
			case <-ticker.C:
				ns = etcdv3.GetAllService(schema, etcdAddr, serviceName)
				c.handleService(schema, ns)
			}
		}
	}()
}

// 处理服务节点
func (c *client) handleService(schema string, ns map[string]string) {
	// 剔除掉旧的 grpc 连接
	for k, _ := range c.AllService {
		// 判断旧的节点 是否在最新的etcd 服务节点里
		if _, ok := ns[k]; !ok {
			delete(c.AllService, k)
			c.delConn(schema, k)
		}
	}

	// 写入新的节点进来
	for k, v := range ns {
		c.AllService[k] = v
	}
}
