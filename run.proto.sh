cd pkg/proto/
protoc --go_out=plugins=grpc:../pb/ *.proto
