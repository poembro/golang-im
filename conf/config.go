package conf

import (
	"net"
	"os"
	"strings"
)

var (
	Conf *Config
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
	TCPListenAddr     string
	WSListenAddr      string
	RPCListenAddr     string
	ConnectIntSerName string
	LocalAddr         string
	SubscribeNum      int
}

// LogicConf logic配置
type LogicConf struct {
	HttpListenAddr  string
	MySQL           string
	RPCListenAddr   string
	LogicIntSerName string
	LocalAddr       string
}

type Config struct {
	Global  *GlobalConf
	Logic   *LogicConf
	Connect *ConnectConf
}

func init() {
	env := os.Getenv("APP_ENV")
	switch env {
	case "dev":
		initLocalConf()
	case "prod":
		initLocalConf()
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
