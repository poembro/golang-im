package cache

import (
    "fmt"
    "github.com/go-redis/redis"
    "golang-im/pkg/db"
    "golang-im/pkg/gerrors"
    "time"
)

type online struct{}

var (
    Online = new(online)
)

const (
    _prefixMidServer    = "userId_%d" // mid -> DeviceId:server
    _prefixKeyServer    = "deviceId_%s" // deviceId -> server
    _prefixServerOnline = "ol_%s"  // server -> online

    _prefixMessageAck = "user_msg_ack_%d"  // user -> ack
    Expire = 75 * time.Second
)

func KeyUserIdServer(userId int64) string {
    return fmt.Sprintf(_prefixMidServer, userId)
}

func KeyDeviceIdServer(deviceId string) string {
    return fmt.Sprintf(_prefixKeyServer, deviceId)
}

func keyServerOnline(deviceId string) string {
    return fmt.Sprintf(_prefixServerOnline, deviceId)
}

func keyMessageAck(userId int64) string {
    return fmt.Sprintf(_prefixMessageAck, userId)
}

// KeysByUserIds get a deviceId server by userId.  HSET userId_123 2000aa78df60000 192.168.3.222
// HGETALL userId_123
func (c *online) KeysByUserIds(userIds []int64) ([]string, error) {
    dst := make([]string, 0)
    for _, userId := range userIds {
        data, err := db.RedisCli.HGetAll(KeyUserIdServer(userId)).Result()
        if err != nil {
            continue
        }

        for k, v := range data {
            if v != "" {
                dst = append(dst, data[k])
            }
        }
    }
    return dst, nil
}

// AddMapping add a mapping.
//    HSET userId_123 2000aa78df60000 {id:1,nickname:张三,face:p.png,}
//    SET  deviceId_2000aa78df60000  192.168.3.222
func (c *online) AddMapping(userId int64, deviceId, server, userinfo string) error {
    // 一个用户有N个设备 全部在hset上面
    _, err := db.RedisCli.HSet(KeyUserIdServer(userId), deviceId, userinfo).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }

    _, err = db.RedisCli.Expire(KeyUserIdServer(userId), Expire).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }

    _, err = db.RedisCli.Set(KeyDeviceIdServer(deviceId), server, Expire).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }

    return nil
}

// ExpireMapping expire a mapping.
//EXPIRE userId_123 2000aa78df60000 1000
//EXPIRE deviceId_2000aa78df60000 1000
func (c *online) ExpireMapping(userId int64, deviceId string) error {
    _, err := db.RedisCli.Expire(KeyUserIdServer(userId), Expire).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }

    _, err = db.RedisCli.Expire(KeyDeviceIdServer(deviceId), Expire).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }
    return nil
}

// DelMapping del a mapping.
// HDEL userId_123 2000aa78df60000
// DEL  deviceId_2000aa78df60000
func (c *online) DelMapping(userId int64, deviceId string) error {
    _, err := db.RedisCli.HDel(KeyUserIdServer(userId), deviceId).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }
    _, err = db.RedisCli.Del(KeyDeviceIdServer(deviceId)).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }

    return nil
}



// AddShopList 将用户添加到商户列表
// zadd  shop_id  time() user_id
func (c *online) AddShopList(shopId string, userId string) error {
    score := time.Now().Unix()
    err := db.RedisCli.ZAdd("userList:" + shopId, redis.Z{
        Score: float64(score),
        Member:userId,
    }).Err()

    if err != nil {
        return gerrors.WrapError(err)
    }

    return nil
}


// AddMsgAckMapping add a msg ack mapping. 记录用户已读偏移
//    HSET userId_123 8000 100000000
func (c *online) AddMsgAckMapping(userId int64, roomId string, seq int64) error {
    // 一个用户有N个房间 每个房间都有个已读偏移位置
    _, err := db.RedisCli.HSet(keyMessageAck(userId), roomId, seq).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }

    _, err = db.RedisCli.Expire(keyMessageAck(userId), Expire * 10000).Result()
    if err != nil {
        return gerrors.WrapError(err)
    }

    return nil
}
// GetMsgAckMapping 读取某个用户的已读偏移
func (c *online) GetMsgAckMapping(userId int64, roomId string) (string, error) {
    // 一个用户有N个房间 每个房间都有个已读偏移位置
    dst, err := db.RedisCli.HGet(keyMessageAck(userId), roomId).Result()
    if err != nil {
        return dst, gerrors.WrapError(err)
    }

    return dst, err
}


// AddMessageList 将消息添加到对应房间 roomId
// zadd  roomId  time() msg
func (c *online) AddMessageList(roomId string, id int64, msg string) error {
    err := db.RedisCli.ZAddNX("messagelist:" + roomId, redis.Z{
        Score: float64(id),
        Member: msg,
    }).Err()

    if err != nil {
        return gerrors.WrapError(err)
    }

    return nil
}

// GetMessageCount 统计未读
func (c *online) GetMessageCount(roomId , start, stop string) (int64, error) {
    dst, err := db.RedisCli.ZCount("messagelist:" + roomId, start, stop).Result()

    if err != nil {
        return dst, gerrors.WrapError(err)
    }

    return dst, nil
}

// GetMessageList 取回消息
func (c *online) GetMessageList(roomId string, start, stop int64) ([]string, error) {
    dst, err := db.RedisCli.ZRevRange("messagelist:" + roomId, start, stop).Result()

    if err != nil {
        return dst, gerrors.WrapError(err)
    }

    return dst, nil
}

