package weight

import (
	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

// Name is the name of weight balancer.
const Name = "weight"

var (
	minWeight = 1
	maxWeight = 5
)

// attributeKey is the type used as the key to store AddrInfo in the Attributes
// field of resolver.Address.
type attributeKey struct{}

// AddrInfo will be stored inside Address metadata in order to use weighted balancer.
type AddrInfo struct {
	Weight int
}

// SetAddrInfo returns a copy of addr in which the Attributes field is updated
// with addrInfo.
func SetAddrInfo(addr resolver.Address, addrInfo AddrInfo) resolver.Address {
	obj := attributes.New()
	addr.Attributes = obj.WithValues(attributeKey{}, addrInfo)
	return addr
}

// GetAddrInfo returns the AddrInfo stored in the Attributes fields of addr.
func GetAddrInfo(addr resolver.Address) AddrInfo {
	v := addr.Attributes.Value(attributeKey{})
	node, _ := v.(AddrInfo)
	return node
}

// NewBuilder creates a new weight balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(Name, &rrPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.V2Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	scs := make([]balancer.SubConn, 0, 5)
	for subConn, addr := range info.ReadySCs {
		node := GetAddrInfo(addr.Address)
		if node.Weight <= 0 {
			node.Weight = minWeight
		}
		if node.Weight > 5 {
			node.Weight = maxWeight
		}
		log.Printf("-服务发现--ip:%v---weight: %d :", addr.Address.Addr, node.Weight)
		// 将不同的ip 复制1-5份 到切片  实现不够美 待优化
		for i := 0; i < node.Weight; i++ {
			scs = append(scs, subConn)
		}
	}

	return &rrPicker{
		subConns: scs,
	}
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn
	mu       sync.Mutex
}

func (p *rrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	index := rand.Intn(len(p.subConns))
	log.Printf("共有%d个节点---任选1个节点--Pick-> %v :", len(p.subConns), index)
	sc := p.subConns[index]
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}
