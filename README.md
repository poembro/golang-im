# golang-im  
一个运行在[golang](#)上的实时通信软件。

---

## 特点
- 简洁
- 高性能
- 支持心跳来维持在线
- 使用 redis 发布订阅做推送
- 采用 TLV 协议格式，保持与[goim](https://github.com/Terry-Mao/goim)一致
- 采用框架结构分层设计参考[gim](https://github.com/alberliu/gim)
- 采用 etcd 作为服务发现,grpc客户端负载均衡 实现分布式高可用
- c10K以内的并发连接完全够用


---

## 描述
- 在学习goim 与 gim 等项目代码后，做的二者结合，用最精简的方式达到练手实践效果。 



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

4.运行测试网页  点击启动按钮 建立websocket连接，点击发送按钮  发送消息
[root@iZ~]#cd golang-im/test
[root@iZ~]# go run main.go
http://localhost:1999/


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