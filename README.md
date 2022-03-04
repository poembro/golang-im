# golang-im  
一个运行在[golang](#)上的实时通信软件。

---

## 特点
- 简洁 
- 高性能
- 支持心跳来维持在线
- 采用 redis 发布订阅做广播推送 (目前demo代码为了阅读简洁暂采用redis，生产环境可更换为 kafka 等)
- 采用 TLV 协议格式，保持与[goim](https://github.com/Terry-Mao/goim)一致
- 采用框架设计参考[gim](https://github.com/alberliu/gim),代码观赏性极好....
- 采用 etcd 作为服务发现,grpc客户端负载均衡 实现分布式高可用


---

## 描述
- 在学习goim 与 gim 等项目代码后，做的二者结合，用最精简的方式达到练手实践效果。 


---

### 设计方案
**所有的聊天会话都是向某一个room_id中写消息**：
- 新成员进入聊天窗口，都必先有个唯一device_id设备标识,两个成员直接相互聊天，本质上就是向两个device_id拼接后的字符串(room_id)写消息；如果有多人群聊即多对多，需额外处理，并不影响该设计


**golang负责接受并管理连接**：
- connect 接入层 如  建立连接后 要多一次订阅房间操作,  将以房间号为key  Room结构为值,存放至全局sync.Map ; 所有连接句柄根据房间号,放入对应Room结构的Conns字段(链表)内  有消息过来 根据房间号 遍历链表的句柄 达到推送效果
- logic   逻辑处理层，如 在线状态存入redis，房间成员信息 ，已读 未读标识。


---

## 项目目录简介
``` 
cmd                        golang 服务启动入口
    |___connect            接入层 提供对外长连接端口 如 websocket tcp ，提供对内grpc服务端口
    |___logic              逻辑处理层 提供对内grpc服务端口 
config                     配置 开发环境 本地环境 生产环境
dist                       静态文件 提供一套完整的 用户端聊天界面 与 客服人员回答界面
    |___ admin           客服人员聊天静态页面  (一个客服可以跟多个用户聊天)
    |___ upload          聊天的图片，上传目录
    |___ im.html         用户端网页咨询窗口的静态页面
internal             
    |___connect             长连接协议处理 ,mq 订阅推送处理)
    |___logic               内部鉴权,消息逻辑 处理
       |___ api            grpc 服务方法
       |___ cache          消息缓存
       |___ model          用户模型
       |___ service        服务层为grpc 提供服务逻辑处理
pkg 
    |___ db                外部数据源
    |___ gerrors           grpc 错误处理
    |___ gn                epoll tcp服务框架 注:这里解释下,由于改了协议所以没有直接引用 [gn](https://github.com/alberliu/gn)
    |___ grpclib           采用etcd 做grpc 服务注册、服务发现
    |___ interceptor       grpc 服务拦截 
    |___ logger            采用zap 做日志库
    |___ pb                proto 生成后的文件
    |___ proto             proto定义grpc 方法和消息结构
    |___ protocol          TLV消息头 同 [goim](https://github.com/Terry-Mao/goim)  
    |___ rpc               构建grpc客户端 及 处理服务发现节点
    |___ urlwhitelist      grpc 服务白名单，如有grpc服务方法需要授权访问 防止外部人员向任意房间发消息
    |___ util              工具
 
run.sh                     普通方式构建 并且 运行 
Dockerfile                 用来构建docker镜像
docker.run.build.sh        用来构建为可执行文件，方便拷贝到docker镜像里面去
docker-start.sh            docker启动时脚本   
tcp_client_testing.go      TCP客户端测试脚本  go run tcp_client_testing.go  222 1 192.168.83.165:6923 
``` 

---

## 安装
``` 
1. 安装docker redis  省略
2. 启动etcd 
docker run --name etcd1 -d -p 2379:2379 -p 2380:2380 -v /Users/luoyuxiang/.laradock/data/etcd:/var/etcd -v /etc/localtime:/etc/localtime registry.cn-hangzhou.aliyuncs.com/google_containers/etcd:3.2.24 etcd --name etcd-s1 --auto-compaction-retention=1 --data-dir=/var/etcd/etcd-data  --listen-client-urls http://0.0.0.0:2379  --listen-peer-urls http://0.0.0.0:2380  --initial-advertise-peer-urls http://192.168.83.165:2380  --advertise-client-urls http://192.168.83.165:2379,http://192.168.83.165:2380 -initial-cluster-token etcd-cluster  -initial-cluster "etcd-s1=http://192.168.83.165:2380,etcd-s2=http://192.168.83.165:2480,etcd-s3=http://192.168.83.165:2580" -initial-cluster-state new

docker run --name etcd2  -d -p 2479:2379 -p 2480:2380 -v /Users/luoyuxiang/.laradock/data/etcd2:/var/etcd -v /etc/localtime:/etc/localtime  registry.cn-hangzhou.aliyuncs.com/google_containers/etcd:3.2.24 etcd --name etcd-s2 --auto-compaction-retention=1 --data-dir=/var/etcd/etcd-data  --listen-client-urls http://0.0.0.0:2479  --listen-peer-urls http://0.0.0.0:2480  --initial-advertise-peer-urls http://192.168.83.165:2480  --advertise-client-urls http://192.168.83.165:2479,http://192.168.83.165:2480 -initial-cluster-token etcd-cluster  -initial-cluster "etcd-s1=http://192.168.83.165:2380,etcd-s2=http://192.168.83.165:2480,etcd-s3=http://192.168.83.165:2580" -initial-cluster-state new


########进入镜像 查看集群状态########
/ # etcdctl --endpoints=http://192.168.83.165:2379,http://192.168.83.165:2479 member list 
df339c03e281023: name=etcd-s1 peerURLs=http://192.168.83.165:2380 clientURLs=http://192.168.83.165:2379,http://192.168.83.165:2380 isLeader=true
4087e512fb03a648: name=etcd-s3 peerURLs=http://192.168.83.165:2580 clientURLs= isLeader=false
68cbc22003188bb9: name=etcd-s2 peerURLs=http://192.168.83.165:2480 clientURLs=http://192.168.83.165:2479,http://192.168.83.165:2480 isLeader=false


########待golang-im节点启动后,执行key前缀匹配,查看etcd中已经存在的节点信息######
# export ETCDCTL_API=3; etcdctl-3.2.24 --endpoints=http://192.168.83.165:2379,http://192.168.83.165:2479 --write-out="simple"  get g --prefix --keys-only
goim:///connectint_grpc_service/192.168.83.165:50000
goim:///connectint_grpc_service/192.168.83.165:50002
goim:///logicint_grpc_service/192.168.83.165:50100
goim:///logicint_grpc_service/192.168.83.165:50102



3. 安装golang-im 服务
[root@iZ~]#cd /data/web
[root@iZ~]#git clone git@github.com:poembro/golang-im.git 
[root@iZ~]#cd /data/web/golang-im
[root@iZ~]#sh run.sh

4.运行测试网页  
 4.1 浏览器打开 http://192.168.83.165:8090/admin/login.html
 4.2 注册 输入手机号 密码 --> 登录  输入手机号 密码
 4.3 点击 底部导航 ”发现“页面  --> 点击浮动头像  即:表示 打开用户端咨询窗口并发送1条消息、
 4.4 点击 底部导航 ”消息“ 页面  --> 可以看到 用户列表  即:表示 当前所有找我咨询的用户  --> 点击对应用户头像 即:回复咨询页面 



5. 多节点集群运行
- golang 编译构建可执行程序
[root@iZ~]#sh ./docker.run.build.sh

- 将可执行程序拷贝到 docker 镜像
[root@iZ~]#docker image build -t golang-im:1.0.18 .

- 用构建好的镜像 启动实例
第一台机器192.168.83.165,启动第一个实例
[root@iZ~]#docker run --name golang-im01 -d -p 50000:50000 -p 50100:50100 -p 7923:7923 -p 6923:6923 -p 8090:8090 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.83.165:50100 --env GRPC_CONNECT_ADDR=192.168.83.165:50000 --rm golang-im:1.0.18

第一台机器192.168.83.165,启动第二个实例
[root@iZ~]#docker run --name golang-im02 -d -p 50002:50000 -p 50102:50100 -p 7924:7923 -p 6923:6923 -p 8090:8090 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.83.165:50102 --env GRPC_CONNECT_ADDR=192.168.83.165:50002 --rm golang-im:1.0.18

第一台机器192.168.83.165,启动第三个实例
[root@iZ~]#docker run  --name golang-im03 -d -p 50003:50000 -p 50103:50100 -p 7925:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.83.165:50103 --env GRPC_CONNECT_ADDR=192.168.83.165:50003 --rm golang-im:1.0.18

第二台机器192.168.82.220,启动第四个实例 (局域网其他机器 互通)
[root@iZ~]#docker run  --name golang-im04 -d -p 50000:50000 -p 50100:50100 -p 7923:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.82.220:50100 --env GRPC_CONNECT_ADDR=192.168.82.220:50000 --rm golang-im:1.0.18


``` 





## 协议格式  
#### 二进制，请求和返回协议一致 
| 参数名     | 必选  | 类型 | 说明       |
| :-----     | :---  | :--- | :---       |
| package length        | true  | int32 bigendian | 包长度 |
| header Length         | true  | int16 bigendian    | 包头长度 |
| ver        | true  | int16 bigendian    | 协议版本 |
| operation          | true | int32 bigendian | 协议指令 |
| seq         | true | int32 bigendian | 序列号 |
| body         | false | binary | $(package lenth) - $(header length) |


#### 协议指令
| 指令     | 说明  | 
| :-----     | :---  |
| 2 | 客户端请求心跳 |
| 3 | 服务端心跳答复 |
| 5 | 下行消息 |
| 7 | auth认证 |
| 8 | auth认证返回 |


---
 
## 感谢

#### 感谢 gim, goim 等开源项目,有冒犯到原作者的地方请及时指正