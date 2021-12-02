package etcdv3

import (
	"context"
	"fmt"
	"net"
	"strings"

	"go.etcd.io/etcd/clientv3"
)

type RegEtcd struct {
	cli    *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
	key    string
}

var rEtcd *RegEtcd

// "%s:///%s/"
func GetPrefix(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

// "%s:///%s"
func GetPrefix4Unique(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s", schema, serviceName)
}

// "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(schema, etcdAddr, myHost string, myPort string, serviceName string, ttl int) error {
	serviceName = serviceName + ":" + net.JoinHostPort(myHost, myPort)
	return RegisterEtcd(schema, etcdAddr, myHost, myPort, serviceName, ttl)
}

func Register(schema, etcdAddr, ipPort, serviceName string, ttl int) error {
	host, port, err := net.SplitHostPort(ipPort)
	if err != nil {
		return fmt.Errorf("port not int error ")
	}
	// TODO 可能更多 服务注册与发现 的持久化工具
	return RegisterEtcd(schema, etcdAddr, host, port, serviceName, ttl)
}

//etcdAddr separated by commas  注册服务
func RegisterEtcd(schema, etcdAddr, myHost string, myPort string, serviceName string, ttl int) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","),
	})
	fmt.Println("etcd Register")
	if err != nil {
		//        return fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
		return fmt.Errorf("create etcd clientv3 client failed, errmsg:%v, etcd addr:%s", err, etcdAddr)
	}

	//lease
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		return fmt.Errorf("grant failed")
	}

	//  schema:///serviceName/ip:port ->ip:port
	serviceValue := net.JoinHostPort(myHost, myPort)
	serviceKey := GetPrefix(schema, serviceName) + serviceValue

	//set key->value
	if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("put failed, errmsg:%v， key:%s, value:%s", err, serviceKey, serviceValue)
	}

	//keepalive
	kresp, err := cli.KeepAlive(ctx, resp.ID)
	if err != nil {
		return fmt.Errorf("keepalive faild, errmsg:%v, lease id:%d", err, resp.ID)
	}

	go func() {
	FLOOP:
		for {
			select {
			case _, ok := <-kresp:
				if ok == true {
				} else {
					break FLOOP
				}
			}
		}
	}()

	rEtcd = &RegEtcd{ctx: ctx,
		cli:    cli,
		cancel: cancel,
		key:    serviceKey}

	return nil
}

func UnRegisterEtcd() {
	//delete
	rEtcd.cancel()
	rEtcd.cli.Delete(rEtcd.ctx, rEtcd.key)
}
