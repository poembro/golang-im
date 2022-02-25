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
func GetPrefix(schema, srvName string) string {
	return fmt.Sprintf("%s:///%s/", schema, srvName)
}

// "%s:///%s"
func GetPrefix4Unique(schema, srvName string) string {
	return fmt.Sprintf("%s:///%s", schema, srvName)
}

// "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(schema, etcdAddr, myHost string, myPort string, srvName string, ttl int) error {
	srvName = srvName + ":" + net.JoinHostPort(myHost, myPort)
	return RegisterEtcd(schema, etcdAddr, myHost, myPort, srvName, ttl)
}

//etcdAddr separated by commas  注册服务
func RegisterEtcd(schema, etcdAddr, myHost string, myPort string, srvName string, ttl int) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","),
	})

	if err != nil {
		return fmt.Errorf("create etcd clientv3 client failed, errmsg:%v, etcd addr:%s", err, etcdAddr)
	}

	//lease
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		cancel()
		return fmt.Errorf("grant failed")
	}

	//在connect那边注册 ---k-> goim:///logicint_grpc_service/192.168.83.165:50100  ---v-> 192.168.83.165:50100
	//在logic那边注册 ---k-> goim:///connectint_grpc_service/192.168.83.165:50002  ---v-> 192.168.83.165:50002
	serviceValue := net.JoinHostPort(myHost, myPort)
	serviceKey := GetPrefix(schema, srvName) + serviceValue

	fmt.Printf("2022-02-25 09:44:38.690	DEBUG	etcdv3/register.go:60	RegisterEtcd %s \r\n", serviceKey)

	//set key->value
	if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		cancel()
		return fmt.Errorf("put failed, errmsg:%v， key:%s, value:%s", err, serviceKey, serviceValue)
	}

	//keepalive
	kresp, err := cli.KeepAlive(ctx, resp.ID)
	if err != nil {
		cancel()
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

	rEtcd = &RegEtcd{
		ctx:    ctx,
		cli:    cli,
		cancel: cancel,
		key:    serviceKey,
	}

	return nil
}

func UnRegisterEtcd() {
	//delete
	rEtcd.cancel()
	rEtcd.cli.Delete(rEtcd.ctx, rEtcd.key)
}
