
 # demo 

## 特点
- 前端样式采用mui框架
- js脚本 采用jq 
- 采用golang 内置net/http做静态文件服务 和 api 接口服务
- 目前支持基础的聊天客服功能 (图文聊天)


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
 
``` 

## 描述
- 框架层面参考了 gin 和 goim 项目，但是市面上没有发现完整的开源demo 于是作为练手实践写了该项目 


