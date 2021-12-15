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
	ProjectName   string
	GrpcSchema    string
	EtcdAddr      string
	RedisIP       string
	RedisPassword string
	PushAllTopic  string
}

// ConnectConf Connect配置
type ConnectConf struct {
	TCPListenAddr string
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	SubscribeNum  int
}

// LogicConf logic配置
type LogicConf struct {
	MySQL         string
	RPCListenAddr string
	LocalAddr     string
}

func init() {
	env := os.Getenv("APP_ENV")
	switch env {
	case "dev":
		initDevConf()
	case "prod":
		initProdConf()
	default:
		initLocalConf()
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
