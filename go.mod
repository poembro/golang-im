module golang-im

go 1.17

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v0.0.0-20170926233335-4201258b820c
	github.com/json-iterator/go v1.1.7
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/smartystreets/goconvey v1.7.2
	github.com/sony/sonyflake v1.0.0
	go.etcd.io/etcd v0.0.0-20200402134248-51bdeb39e698
	go.uber.org/zap v1.14.1
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/coreos/go-semver v0.2.0 // indirect
	github.com/coreos/go-systemd/v22 v22.0.0 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/google/uuid v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace golang-im/pkg/gn v1.0.0 => ./pkg/gn

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1 // indirect
