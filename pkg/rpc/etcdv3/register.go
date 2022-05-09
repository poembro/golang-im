package etcdv3

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func newClient(etcdAddr string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:          strings.Split(etcdAddr, ","),
		DialTimeout:        time.Second * time.Duration(5),
		MaxCallSendMsgSize: 2 * 1024 * 1024,
	})
}

// GetPrefix 格式化 "%s:///%s/"
func GetPrefix(schema, svcname string) string {
	return fmt.Sprintf("%s:///%s/", schema, svcname)
}

// GetPrefix4Unique 格式化 "%s:///%s"
func GetPrefix4Unique(schema, svcname string) string {
	return fmt.Sprintf("%s:///%s", schema, svcname)
}

// RegisterEtcd4Unique "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(schema, etcdAddr, myHost string, myPort string, svcname string, ttl int) (func(), error) {
	val := net.JoinHostPort(myHost, myPort)
	svcname = svcname + ":" + val
	return RegisterEtcd(schema, etcdAddr, val, svcname, ttl)
}

// RegisterEtcd  注册服务
func RegisterEtcd(schema, etcdAddr, value string, svcname string, ttl int) (func(), error) {
	cli, err := newClient(etcdAddr)

	if err != nil {
		return nil, err
	}

	//lease
	ctx, cancel := context.WithCancel(context.Background())
	lease, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		cancel()
		return nil, err
	}

	key := GetPrefix(schema, svcname) + value

	if _, err := cli.Put(ctx, key, value, clientv3.WithLease(lease.ID)); err != nil {
		cancel()
		return nil, err
	}
	keepAlive, err := cli.KeepAlive(ctx, lease.ID)
	if err != nil {
		cancel()
		return nil, err
	}

	go func() {
		for ka := range keepAlive { // keepAlive是1个channel
			if ka == nil {
				break
			}
			//fmt.Println("-->续约成功", ka)
		}
		fmt.Printf("2022-02-25 09:44:38.690	DEBUG	etcdv3/register.go:73	%s \r\n", "关闭续租")
	}()

	closeEtcd := func() {
		_, _ = cli.Revoke(ctx, lease.ID)
		cancel()
		fmt.Printf("2022-02-25 09:44:38.690	DEBUG	etcdv3/register.go:79	%s \r\n", "关闭etcd连接")
	}

	return closeEtcd, nil
}
