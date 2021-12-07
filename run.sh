GO111MODULE="on"
cd cmd/logic
rm -f logic
go build -ldflags "-w -s" -v -o logic main.go
echo "打包logic成功"
pkill logic
echo "停止logic服务"
nohup ./logic &
echo "启动logic服务"

cd ../connect
rm -f connect
go build -ldflags "-w -s" -v -o connect main.go
echo "打包connect成功"
pkill connect
echo "停止connect服务"
sleep 2
nohup ./connect &
echo "启动connect服务"

cd ../../demo
rm -f demo
go build -ldflags "-w -s" -v -o demo main.go
echo "打包demo成功"
pkill demo
echo "停止demo服务"
sleep 2
nohup ./demo &
echo "启动demo服务"
