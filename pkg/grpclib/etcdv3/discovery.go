package etcdv3

import (
	"context"
	"fmt"
	"log"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"

	//"google.golang.org/genproto/googleapis/ads/googleads/v1/services"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/resolver"
)

type Discovery struct {
	cli                *clientv3.Client
	cc                 resolver.ClientConn
	serviceName        string
	schema             string
	watchStartRevision int64
	SrvList            sync.Map //服务列表
}

// NewDiscovery  新建服务发现
func NewDiscovery(schema, etcdAddr, serviceName string) resolver.Builder {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Discovery{
		cli:         cli,
		schema:      schema,
		serviceName: serviceName,
	}
}

func (r *Discovery) Scheme() string {
	return r.schema
}

func (r *Discovery) ResolveNow(rn resolver.ResolveNowOptions) {
	fmt.Println(rn)
}

func (r *Discovery) Close() {
}

func (r *Discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if r.cli == nil {
		return nil, fmt.Errorf("etcd clientv3 client failed, etcd:%s", target)
	}
	r.cc = cc

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s/"
	prefix := GetPrefix(r.schema, r.serviceName)

	// get key first
	resp, err := r.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Errorf(err.Error())
		return nil, err
	}

	for i := range resp.Kvs {
		k := string(resp.Kvs[i].Key)
		v := string(resp.Kvs[i].Value)
		//fmt.Println("-k->" , k, " ---v->",v )
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
	fmt.Println("etcd SrvList set addr : ", val)
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
	fmt.Printf("etcd ---> %+v \r\n", addrs)
	return addrs
}

func GetAllService(schema, etcdaddr, servicename string) map[string]string {
	allService := make(map[string]string)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdaddr, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("etcd err", err.Error())
		return allService
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s"
	prefix := GetPrefix4Unique(schema, servicename)

	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	//  "%s:///%s/ip:port"   -> %s:ip:port
	if err != nil {
		fmt.Println("etcd err", err.Error())
	}

	for i := range resp.Kvs {
		key := string(resp.Kvs[i].Value)
		allService[key] = servicename
	}
	cli.Close()
	return allService
}
