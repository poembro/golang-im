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



# docker run -v $(pwd)/:/app -p 8080:8080 -p 8081:8081 -p 50100:50100 alpine .//app/main