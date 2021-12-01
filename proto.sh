cd pkg/proto/
protoc -I=./  --go_out=plugins=grpc:../pb/ *.proto
