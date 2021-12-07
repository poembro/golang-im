package service

import (
    "context"
    "encoding/json"
    "fmt"
    "golang-im/internal/logic/cache"
    "golang-im/internal/logic/model"
    "golang-im/pkg/logger"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 长连接登录
// 方案一: body 是一个jwt token 值 去其他服务拿到对应的 头像昵称等信息
// 方案二: demo 中 body 是一个json 已经包含了头像昵称等信息
func (*authService) SignIn(ctx context.Context, body []byte, connAddr string, clientAddr string) (string, int64, error) {
    var (
        user model.User
        deviceId = "11011"
        userId   = int64(7)
        err error
    )
    //解析body  得到 deviceId, userId
    json.Unmarshal(body, &user)

    userId = int64(user.UserId)
    deviceId = user.DeviceId
    // 标记用户在设备上登录
    err = cache.Online.AddMapping(userId, deviceId, connAddr, string(body))
    logger.Sugar.Infow("SignIn", "user", user)

    // 写入商户列表
    err = cache.Online.AddShopList(user.ShopId, fmt.Sprintf("%d", userId))

    return deviceId, userId, err
}

func (s *authService) Heartbeat(ctx context.Context, userId int64, deviceId, connAddr string) error {
    cache.Online.ExpireMapping(userId, deviceId)
    return nil
}

func (s *authService) Offline(ctx context.Context, userId int64, deviceId, connAddr string) error {
    cache.Online.DelMapping(userId, deviceId)
    return nil
}
