package etcdv3

import (
	"context"
	"fmt"
	"log"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"

	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/resolver"
)

type Discovery struct {
	cli                *clientv3.Client
	cc                 resolver.ClientConn
	srvName            string
	schema             string
	watchStartRevision int64
	SrvList            sync.Map //服务列表
}

// NewDiscovery  新建服务发现
func NewDiscovery(schema, etcdAddr, srvName string) resolver.Builder {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Discovery{
		cli:     cli,
		schema:  schema,
		srvName: srvName,
	}
}

func (r *Discovery) Scheme() string {
	return r.schema
}

func (r *Discovery) ResolveNow(rn resolver.ResolveNowOptions) {
	//fmt.Println(rn)
}

func (r *Discovery) Close() {
}

func (r *Discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	//     "%s:///%s/"
	prefix := GetPrefix(r.schema, r.srvName)

	// get key first
	resp, err := r.cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("2022-02-25 09:44:38.690	DEBUG	etcdv3/descovery.go:64 Build ", err.Error())
		return nil, err
	}

	for i := range resp.Kvs {
		k := string(resp.Kvs[i].Key)
		v := string(resp.Kvs[i].Value)
		fmt.Printf("2022-02-25 09:44:38.690	DEBUG	etcdv3/descovery.go:71	Build	{\"desc\": \"已存在节点 k:%s v:%s\"} \r\n", k, v)
		r.Set(k, v)
	}
	r.cc.UpdateState(resolver.State{Addresses: r.Gets()})
	r.watchStartRevision = resp.Header.Revision + 1

	go r.watch(prefix)
	return r, nil
}

func (r *Discovery) watch(prefix string) {
	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrefix())
	for resp := range rch {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				r.Set(string(ev.Kv.Key), string(ev.Kv.Value))
				r.cc.UpdateState(resolver.State{Addresses: r.Gets()})
			case mvccpb.DELETE:
				r.Del(string(ev.Kv.Key))
				r.cc.UpdateState(resolver.State{Addresses: r.Gets()})
			default:
				// TODO
			}
		}
	}
}

// Set 设置服务地址
func (r *Discovery) Set(key, val string) {
	fmt.Printf("2022-02-25 09:44:38.690	DEBUG	etcdv3/descovery.go:101	Set	{\"desc\": \"有新节点加入 %s\"} \r\n", val)
	r.SrvList.Store(key, val)
}

// Del 删除服务地址
func (r *Discovery) Del(key string) {
	r.SrvList.Delete(key)
}

// Gets 获取服务地址
func (r *Discovery) Gets() []resolver.Address {
	addrs := make([]resolver.Address, 0)
	r.SrvList.Range(func(k, v interface{}) bool {
		addrs = append(addrs, resolver.Address{Addr: v.(string)})
		return true
	})
	//fmt.Printf("etcd --Gets()---当前-> %+v \r\n", addrs)
	return addrs
}

////////////////////////////Check检查节点离线逻辑////////////////////////////

type Check struct {
	cli     *clientv3.Client
	srvName string
	schema  string
}

// NewNodeCheck  检查节点
func NewHealthCheck(schema, etcdAddr, srvName string) *Check {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Check{
		cli:     cli,
		schema:  schema,
		srvName: srvName,
	}
}

func (c *Check) GetEtcdConns() map[string]string {
	conns := make(map[string]string)

	// "%s:///%s"
	prefix := GetPrefix4Unique(c.schema, c.srvName)

	resp, err := c.cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	//  "%s:///%s/ip:port"   -> %s:ip:port
	if err != nil {
		fmt.Println("2022-02-25 09:44:38.690	DEBUG	etcdv3/descovery.go:155	GetEtcdConns", err.Error())
		return conns
	}

	for i := range resp.Kvs {
		k := string(resp.Kvs[i].Key)
		v := string(resp.Kvs[i].Value)
		conns[k] = v
	}
	return conns
}
