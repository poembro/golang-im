package service

import (
	"context"
	"golang-im/internal/logic/cache"
	"golang-im/pkg/logger"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 长连接登录
// 方案一: body 是一个jwt token 值 去其他服务拿到对应的 头像昵称等信息
// 方案二: demo 中 body 是一个json 已经包含了头像昵称等信息
func (*authService) SignIn(ctx context.Context, body []byte, connAddr string, clientAddr string) (string, int64, error) {
	var (
		deviceId string
		userId   int64
	)
	deviceId = "11011"
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
	   }
	*/
	//解析body  得到 deviceId, userId

	// 标记用户在设备上登录
	err := cache.Online.AddMapping(userId, "md5", connAddr, string(body))
	logger.Sugar.Infow("---->", "SignIn", 2, "desc", "标记用户在设备上登录")

	return deviceId, userId, err
}

func (s *authService) Offline(ctx context.Context, userId int64, deviceId string, clientAddr string) error {
	cache.Online.DelMapping(userId, deviceId)
	return nil
}
