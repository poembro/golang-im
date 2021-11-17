package etcdv3

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang-im/pkg/grpclib/weight"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

const schema = "poembro"

type Discovery struct {
	cli        *clientv3.Client //etcd client
	cc         resolver.ClientConn
	serverList sync.Map //服务列表
	prefix     string   //监视的前缀
}

// NewDiscovery  新建服务发现
func NewDiscovery(endpoints []string) resolver.Builder {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Discovery{
		cli: cli,
	}
}

/*
grpc服务发现使用  https://blog.csdn.net/weixin_39838758/article/details/111103014

1.在grpc.Dial() 之前 调用 resolver.Register(db.EtcdCli) 表示向resolver/resolver.go 中m变量追加参数和值 m[b.Scheme()] = b
2.创建resolver 用来解析服务端的地址，过程中 newCCResolverWrapper  方法里调用的 Discovery.Build(x,x)
3.将第一次从ETCD GET方式获取的所有节点拿出来作为切片 丢给cc.UpdateState方法
4.并且每个节点 都将 构建/new 1个 resolver.Address 结构 其中该结构有个参数Attributes 存放着权重值

5.由于调用了 NewDiscovery方法 以至于 weight.go文件中 balancer.Register(newBuilder()) 被执行
grpc客户端参数中再指定使用 grpc.WithBalancerName("weight") 作为路由选择器
6.于是 rrPickerBuilder.Build(info) 被隐式调用了 其中info 里面就是 resolver.Address 结构 这里就可以取出 权重值
*/

// Build 初始构建 实现接口中固定的方法
func (s *Discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	s.cc = cc
	s.prefix = "/" + target.Scheme + "/" + target.Endpoint + "/"
	resp, err := s.cli.Get(context.Background(), s.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, ev := range resp.Kvs {
		s.Set(string(ev.Key), string(ev.Value))
	}
	// 如果有固定地址的resolver则直接写死, 没有则从etcd取
	// cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint}}})
	cc.UpdateState(resolver.State{Addresses: s.Gets()})

	go s.watcher() // 监视前缀，修改变更的server 地址
	return s, nil
}

// ResolveNow 监视目标更新
func (s *Discovery) ResolveNow(rn resolver.ResolveNowOptions) {
	//log.Println("监视目标更新 ResolveNow")
}

// Scheme 实现接口中固定的方法, 如果不实现这个方法 默认名字叫做 passthrough
func (s *Discovery) Scheme() string {
	return schema
}

// Close 关闭
func (s *Discovery) Close() {
	s.cli.Close()
}

// watcher 监听前缀
func (s *Discovery) watcher() {
	ch := s.cli.Watch(context.Background(), s.prefix, clientv3.WithPrefix())
	//log.Printf("watching prefix:%s now...", s.prefix)  //watching prefix:/grpclb/connectint_grpc_service/ now...
	for resp := range ch {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				s.Set(string(ev.Kv.Key), string(ev.Kv.Value))
				s.cc.UpdateState(resolver.State{Addresses: s.Gets()})
			case mvccpb.DELETE:
				s.Del(string(ev.Kv.Key))
				s.cc.UpdateState(resolver.State{Addresses: s.Gets()})
			default:
				// TODO
			}
		}
	}
}

// Set 设置服务地址
func (s *Discovery) Set(key, val string) {
	// 获取服务地址 去除前缀 /grpclb/connectint_grpc_service/
	node := resolver.Address{Addr: strings.TrimPrefix(key, s.prefix)}
	nodeWeight, err := strconv.Atoi(val) //获取服务地址权重
	if err != nil {
		nodeWeight = 1 // 非数字字符默认权重为1
	}
	// 把服务地址权重数值作为参数 追加到resolver.Address结构体的元数据中
	node = weight.SetAddrInfo(node, weight.AddrInfo{Weight: nodeWeight})
	s.serverList.Store(key, node)
}

// Del 删除服务地址
func (s *Discovery) Del(key string) {
	s.serverList.Delete(key)
}

// Gets 获取服务地址
func (s *Discovery) Gets() []resolver.Address {
	addrs := make([]resolver.Address, 0, 10)
	s.serverList.Range(func(k, v interface{}) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}
