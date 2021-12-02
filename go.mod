module golang-im

go 1.13

replace golang-im/pkg/gn v1.0.0 => ./pkg/gn

require (
    github.com/go-redis/redis v6.15.9+incompatible
    github.com/go-sql-driver/mysql v1.6.0
    github.com/golang/protobuf v1.5.2
    github.com/gorilla/websocket v0.0.0-20170926233335-4201258b820c
    github.com/json-iterator/go v1.1.7
    github.com/onsi/ginkgo v1.16.5 // indirect
    github.com/onsi/gomega v1.17.0 // indirect
    go.etcd.io/etcd v0.0.0-20200402134248-51bdeb39e698
    go.uber.org/zap v1.14.1
    golang.org/x/sys v0.0.0-20210423082822-04245dca01da
    google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
    google.golang.org/grpc v1.29.1
    google.golang.org/protobuf v1.26.0
    gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1 // indirect
