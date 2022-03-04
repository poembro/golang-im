GO111MODULE="on"

cd pkg/proto/
#protoc -I=./  --go_out=plugins=grpc:../pb/ *.proto


cd ../../cmd/logic
rm -f logic

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -v -o logic main.go
echo "打包logic成功"


cd ../connect
rm -f connect
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -ldflags "-w -s" -v -o connect main.go
echo "打包connect成功" 
