package config

import (
	"net"
	"os"
	"strings"
)

var (
	Global  GlobalConf
	Logic   LogicConf
	Connect ConnectConf
)

// GlobalConf RPC配置
type GlobalConf struct {
	ProjectName string
	GrpcSchema  string
}

// ConnectConf Connect配置
type ConnectConf struct {
	TCPListenAddr string
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	RedisIP       string
	RedisPassword string
	SubscribeNum  int
}

// LogicConf logic配置
type LogicConf struct {
	EtcdIPs       string
	MySQL         string
	NSQIP         string
	RedisIP       string
	RedisPassword string
	RPCListenAddr string
	LocalAddr     string
}

func init() {
	ip := InternalIP()
	env := os.Getenv("gim_env")
	switch env {
	case "dev":
		initDevConf(ip)
	case "prod":
		initProdConf(ip)
	default:
		initLocalConf(ip)
	}
}

// InternalIP return internal ip.
func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}
