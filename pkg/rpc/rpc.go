package rpc

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"
	
	"golang-im/pkg/grpclib/etcdv3"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/balancer/roundrobin"
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

// 不需要知道具体访问哪个节点
func newServicenameGrpc(schema, etcdAddr, servicename string) (*grpc.ClientConn, error) {
	rr := etcdv3.NewDiscovery(schema, etcdAddr, servicename)
	resolver.Register(rr) //向resolver/resolver.go 中m变量追加参数和值 m[b.Scheme()] = b

	conn, err := grpc.DialContext(
		context.TODO(),
		etcdv3.GetPrefix4Unique(schema, servicename),
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

	return conn, err
}

// ConnectInt grpc server 服务需要知道是哪个节点
func ConnectInt(addr string) pb.ConnectIntClient {
	conn := Client.GetConn(Client.Schema, Client.EtcdAddr, addr)
	if conn == nil {
		fmt.Errorf("grpc client failed")
		return nil
	}
	return pb.NewConnectIntClient(conn)
}

// LogicInt grpc server 服务不用知道哪个节点
func LogicInt() pb.LogicIntClient {
	conn := Client.GetConn(Client.Schema, Client.EtcdAddr, LogicIntSerName)
	if conn == nil {
		fmt.Errorf("grpc client failed")
		return nil
	}
	return pb.NewLogicIntClient(conn)
}


type client struct {
	ClientConn map[string]*grpc.ClientConn
	Lock       sync.RWMutex
	AllService map[string]string
	Schema     string
	EtcdAddr   string
}

func NewClient(schema, etcdAddr string) {
	Client = &client{
		ClientConn: make(map[string]*grpc.ClientConn),
		AllService: make(map[string]string),
		Schema:     schema,
		EtcdAddr:   etcdAddr,
	}

	// 去etcd 检查某服务所有在线节点  更多节点TODO  注意自定义grpc路由不用检查
	Client.checkNode(schema, etcdAddr, ConnectIntSerName)
}

func (c *client) GetConn(schema, etcdAddr, servicename string) *grpc.ClientConn {
	var (
		r   *grpc.ClientConn
		err error
		ok  bool
		key string
	)
	key = etcdv3.GetPrefix4Unique(schema, servicename)
	//先判断是否已经有了该节点
	c.Lock.RLock()
	if r, ok = c.ClientConn[key]; ok {
		c.Lock.RUnlock()
		return r
	}
	c.Lock.RUnlock()

	c.Lock.Lock()
	// 判断当前是否 ip:port
	u, err := url.Parse(key)
	if u.Port() != "" {
		// 判断节点是否在线 TODO
		if _, ok := c.AllService[servicename]; !ok {
			fmt.Errorf("将要访问的节点  %s 不在线", servicename)
			return nil
		}
		r, err = newGrpc(servicename)
	} else {
		r, err = newServicenameGrpc(schema, etcdAddr, servicename)
	}

	if err != nil {
		c.Lock.Unlock()
		return nil
	}

	c.ClientConn[key] = r
	c.Lock.Unlock()
	return r
}

func (c *client) delConn(schema, servicename string) {
	c.Lock.Lock()
	key := etcdv3.GetPrefix4Unique(schema, servicename)
	delete(c.ClientConn, key)
	c.Lock.Unlock()
}

func (c *client) checkNode(schema, etcdAddr, servicename string) {
	timeoutTicker := 1 * time.Minute
	ns := etcdv3.GetAllService(schema, etcdAddr, servicename)
	for k, v := range ns {
		c.AllService[k] = v
	}

	// checkTimeout 定时检查在线的connect节点
	go func() {
		ticker := time.NewTicker(timeoutTicker)
		for {
			select {
			case <-ticker.C:
				ns = etcdv3.GetAllService(schema, etcdAddr, servicename)
				c.handleService(schema, ns)
			}
		}
	}()
}

// 处理服务节点
func (c *client) handleService(schema string, ns map[string]string) {
	// 剔除掉旧的 grpc 连接
	for k, v := range c.AllService {
		// 判断旧的节点 在不在最新的etcd 服务节点里
		if newval, ok := ns[k]; newval == v && !ok {
			delete(c.AllService, k)
			c.delConn(schema, k)
		}
	}

	// 写入新的节点进来
	for k, v := range ns {
		c.AllService[k] = v
	}
}
