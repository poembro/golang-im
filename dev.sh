GO111MODULE="on"
gim_env=dev
cd cmd/logic
rm -f devlogic
go build -ldflags "-w -s" -v -o devlogic main.go
echo "打包devlogic成功"
pkill devlogic
echo "停止devlogic服务"
export gim_env=dev && nohup ./devlogic &
echo "启动devlogic服务"

