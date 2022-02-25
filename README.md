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

## 项目目录简介
``` 
cmd                        golang 服务启动入口
    |___connect            接入层 提供对外长连接端口 如 websocket tcp ，提供对内grpc服务端口
    |___logic              逻辑处理层 提供对内grpc服务端口 
config                     配置 开发环境 本地环境 生产环境
demo                       提供一套完整的 用户端聊天界面 与 客服人员回答界面
    |___dist                 静态文件
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

run.data.race.sh           竞争检测方式构建
run.sh                     普通方式构建
proto.sh                   编译protoc文件为pb文件
Dockerfile                 用来构建docker镜像
docker.run.build.sh        用来构建为可执行文件，方便拷贝到docker镜像里面去
docker-start.sh            docker启动时脚本   
 
``` 

---

### 设计方案
**所有的聊天会话都是向某一个room_id中写消息**：
- 新成员进入聊天窗口，都必先有个唯一device_id设备标识,两个成员直接相互聊天，本质上就是向两个device_id拼接后的字符串(room_id)写消息；如果有多人群聊即多对多，需额外处理，并不影响该设计


**golang负责接收并管理连接**：
- connect 接入层 如  建立连接后 要多一次订阅房间操作,  将以房间号为key  Room结构为值,存放至全局sync.Map ; 所有连接句柄根据房间号,放入对应Room结构的Conns字段(链表)内  有消息过来 根据房间号 遍历链表的句柄 达到推送效果
- logic   逻辑处理层，如 在线状态存入redis，房间成员信息 ，已读 未读标识。

---

## 安装
``` 
1. 安装redis  省略
2. 安装etcd 
2.1 安装docker-compose 
# yum -y install libcurl libcurl-devel
# curl -L https://github.com/docker/compose/releases/download/1.21.2/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
# chmod +x /usr/local/bin/docker-compose
# docker-compose --version

2.2 [root@iZ~]#mkdir /data/web/etcd/ && cd /data/web/etcd
2.3 [root@iZ~]#vi docker-compose.yml
version: "3.5"
services:
  etcd:
    hostname: etcd
    image: bitnami/etcd:3
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
    # ports:
    #   - "2379:2379"
    #   - "2380:2380"
    #   - "4001:4001"
    #   - "7001:7001"
    privileged: true
    volumes:
      - "~/.laradock/data/etcd/data:/opt/bitnami/etcd/data"  ##注意这里目录映射
    environment:
      - "ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379"
      - "ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379"
      - "ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380"
      - "ETCD_INITIAL_ADVERTISE_PEER_URLS=http://0.0.0.0:2380"
      - "ALLOW_NONE_AUTHENTICATION=yes"
      - "ETCD_INITIAL_CLUSTER=node1=http://0.0.0.0:2380"
      - "ETCD_NAME=node1"
      - "ETCD_DATA_DIR=/opt/bitnami/etcd/data"
    ports:
      - 2379:2379
      - 2380:2380
    networks:
      - etcdnet

networks:
  etcdnet:
    name: etcdnet




2.4 [root@iZ~]#docker-compose up -d 

3. 安装golang-im 服务
[root@iZ~]#cd /data/web
[root@iZ~]#git clone git@github.com:poembro/golang-im.git 
[root@iZ~]#cd /data/web/golang-im
[root@iZ~]#sh run.sh

4.运行测试网页  
 4.1 浏览器打开 http://localhost:8888/admin/login.html
 4.2 注册 输入手机号 密码 --> 登录  输入手机号 密码
 4.3 点击 底部导航 ”发现“页面  --> 点击浮动头像  即:表示 打开用户端咨询窗口并发送1条消息、
 4.4 点击 底部导航 ”消息“ 页面  --> 可以看到 用户列表  即:表示 当前所有找我咨询的用户  --> 点击对应用户头像 即:回复咨询页面




5. 多节点集群运行
[root@iZ~]#sh ./docker.run.build.sh

[root@iZ~]#docker image build -t golang-im:1.0.18 .

第一台机器192.168.83.165,启动第一个实例
[root@iZ~]#docker run -p 50000:50000 -p 50100:50100 -p 7923:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.83.165:50100 --env GRPC_CONNECT_ADDR=192.168.83.165:50000 --rm golang-im:1.0.18

第一台机器192.168.83.165,启动第二个实例
[root@iZ~]#docker run -p 50002:50000 -p 50102:50100 -p 7924:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.83.165:50102 --env GRPC_CONNECT_ADDR=192.168.83.165:50002 --rm golang-im:1.0.18

第一台机器192.168.83.165,启动第三个实例
[root@iZ~]#docker run -p 50003:50000 -p 50103:50100 -p 7925:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.83.165:50103 --env GRPC_CONNECT_ADDR=192.168.83.165:50003 --rm golang-im:1.0.18

第二台机器192.168.82.220,启动第四个实例 (局域网其他机器 互通)
[root@iZ~]#docker run -p 50000:50000 -p 50100:50100 -p 7923:7923 --env APP_ENV=local --env GRPC_LOGIC_ADDR=192.168.82.220:50100 --env GRPC_CONNECT_ADDR=192.168.82.220:50000 --rm golang-im:1.0.18

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