package etcdv3

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"

	"sync"

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

// NewDiscovery
func NewDiscovery(schema, etcdAddr, srvName string) resolver.Builder {
	cli, err := newClient(etcdAddr)

	if err != nil {
		fmt.Printf("etcd 连接错误 %s \r\n", err.Error())
	}

	r := &Discovery{
		cli:     cli,
		schema:  schema,
		srvName: srvName,
	}
	// 第一次先将所有的节点拿出来
	prefix := GetPrefix(r.schema, r.srvName)
	dst, err := r.cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if dst != nil && err == nil {
		for i := range dst.Kvs {
			k := string(dst.Kvs[i].Key)
			v := string(dst.Kvs[i].Value)
			r.set(k, v)
		}
		r.watchStartRevision = dst.Header.Revision + 1
	}
	return r
}

func (r *Discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	r.cc.UpdateState(resolver.State{Addresses: r.getAll()})
	prefix := GetPrefix(r.schema, r.srvName)
	go r.Watch(prefix)
	return r, nil
}

func (r *Discovery) Scheme() string {
	return r.schema
}

////////////////////////////////////////////////////////////

func (r *Discovery) ResolveNow(rn resolver.ResolveNowOptions) {
	// todo
}

func (r *Discovery) Close() {
	// todo
}

func (r *Discovery) Watch(prefix string) {
	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrefix())
	for resp := range rch {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				r.set(string(ev.Kv.Key), string(ev.Kv.Value))
				r.cc.UpdateState(resolver.State{Addresses: r.getAll()})
			case mvccpb.DELETE:
				r.delete(string(ev.Kv.Key))
				r.cc.UpdateState(resolver.State{Addresses: r.getAll()})
			default:
				// TODO
			}
		}
	}
}

// set 设置服务地址
func (r *Discovery) set(key, val string) {
	fmt.Printf("2022-02-25 09:44:38.690	DEBUG	etcdv3/descovery.go:93	set	{\"desc\": \"有新节点加入 %s\"} \r\n", val)
	r.SrvList.Store(key, val)
}

// delete 删除服务地址
func (r *Discovery) delete(key string) {
	r.SrvList.Delete(key)
}

// getAll 获取服务地址
func (r *Discovery) getAll() []resolver.Address {
	dst := make([]resolver.Address, 0)
	r.SrvList.Range(func(k, v interface{}) bool {
		dst = append(dst, resolver.Address{Addr: v.(string)})
		return true
	})
	//fmt.Printf("etcd --getAll()---当前-> %+v \r\n", dst)
	return dst
}
