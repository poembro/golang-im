GO111MODULE="on"
APP_ENV=local


cd pkg/proto/
#protoc -I=./  --go_out=plugins=grpc:../pb/ *.proto


cd ../../cmd/logic
rm -f logic

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
GRPC_LOGIC_ADDR=192.168.84.168:50100 GRPC_CONNECT_ADDR=192.168.84.168:50000
go build -ldflags "-w -s" -v -o logic main.go
echo "打包logic成功"
pkill logic
echo "停止logic服务"
nohup ./logic &
echo "启动logic服务"

sleep 2

cd ../connect
rm -f connect
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
GRPC_LOGIC_ADDR=192.168.84.168:50100 GRPC_CONNECT_ADDR=192.168.84.168:50000
go build -ldflags "-w -s" -v -o connect main.go
echo "打包connect成功"
pkill connect
echo "停止connect服务"
sleep 2
nohup ./connect &
echo "启动connect服务"
