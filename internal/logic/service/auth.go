package service

import (
    "context"
    "golang-im/pkg/logger"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 长连接登录
func (*authService) SignIn(ctx context.Context, body []byte, connAddr string, clientAddr string) (int64, int64, error) {
    var (
        deviceId int64
        userId   int64
    )
    deviceId = int64(14)
    userId = int64(7)
    /*
        p.Body = {
            room_id:'live://1000', //将消息发送到指定房间
            accepts:'[1000,1001,1002]',//接收系统全局房间的消息
            key:'1xxxxxDeviceIdxxx假定为设备idxx',

            user_id: '663291537152950273',
            user_name:'随机用户001',
            user_face:'/static/wap/img/portrait.jpg',

            shop_id:0,
            shop_name:'杂货铺老板',
            shop_face:'/static/wap/img/portrait.jpg',

            platform:'web',

            suburl : "ws://192.168.3.222:9999/sub",
            pushurl:"http://192.168.3.222:9999/open/push",
        }*/
    //解析body  得到 deviceId, userId

    logger.Sugar.Infow("---->", "SignIn", 2, "desc", "标记用户在设备上登录")
    // 标记用户在设备上登录
    //err := DeviceService.Online(ctx, deviceId, userId, connAddr, clientAddr)
    //logger.Sugar.Infow("---->", "SignIn", 11, "desc", 9999)
    //if err != nil {
    //    return deviceId, userId, err
    //}
    return deviceId, userId, nil
}
