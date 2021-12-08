package main

import (
	"encoding/json"
	"fmt"
	"golang-im/pkg/util"

	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang-im/config"
	"golang-im/internal/logic/cache"
	"golang-im/internal/logic/model"
	"golang-im/pkg/db"
	"golang-im/pkg/gerrors"
	"golang-im/pkg/logger"
	"golang-im/pkg/pb"
	"golang-im/pkg/protocol"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

type router struct {
}

func (ro *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth/logout": // 退出
		http.Redirect(w, r, "/admin/login.html", http.StatusFound)
	case "/auth/register": // 注册
		apiRegister(w, r)
	case "/auth/login": // 登录
		apiLogin(w, r)
	case "/open/url": // 获取授权参数 后才允许连接推送服务
		apiUrl(w, r)
	case "/open/push": // 接收消息写入redis or mq
		apiPush(w, r)
	case "/open/finduserlist": // 在线列表
		apiFindUserList(w, r)
	case "/upload/file": //文件上传接口
		apiUpload(w, r)
	default:
		StaticServer(w, r)
	}
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

// StaticServer 静态文件处理
func StaticServer(w http.ResponseWriter, req *http.Request) {
	indexs := []string{"index.html", "index.htm"}
	filePath := "./dist" + req.URL.Path //注意 注意 注意:这里只能处理 dist 目录下的文件
	fi, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.NotFound(w, req)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if fi.IsDir() {
		if req.URL.Path[len(req.URL.Path)-1] != '/' {
			http.Redirect(w, req, req.URL.Path+"/", 301)
			return
		}
		for _, index := range indexs {
			fi, err = os.Stat(filePath + index)
			if err != nil {
				continue
			}
			http.ServeFile(w, req, filePath+index)
			return
		}
		http.NotFound(w, req)
		return
	}
	http.ServeFile(w, req, filePath)
}

// api 接口区

// api响应结构体
type resp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`

	Mobile   string       `json:"mobile,omitempty"`
	Password string       `json:"password,omitempty"`
	FileName string       `json:"filename,omitempty"`
	UserList []model.User `json:"user_list,omitempty"`
	Data     model.User   `json:"data,omitempty"`
}

func serveJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Server", "poembro")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	body, err := json.Marshal(data)
	if err == nil {
		w.Write(body)
	} else {
		w.Write([]byte(`{"code":0, "error":"解析JSON出错"}`))
	}
}

// 客服聊天场景 外链url获取
func apiUrl(w http.ResponseWriter, r *http.Request) {
	var (
		mobile string
	)

	if r.Method == "POST" {
		mobile = r.FormValue("mobile")
	}

	if mobile == "" {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     "手机号不能为空",
		})
		return
	}

	sID, err := util.SFlake.GetID()
	if err != nil {
		fmt.Printf("雪花算法 error: %s \r\n", err.Error())
	}
	platform := "web"
	// 客服聊天场景
	dst := model.User{
		UserId:   model.Int64(sID),                                                       //   用户大多是临时过来咨询,所以这里采用随机唯一
		Nickname: fmt.Sprintf("用户%d", sID),                                               // 随机昵称
		Face:     "http://img.touxiangwu.com/2020/3/uq6Bja.jpg",                          // 随机头像
		DeviceId: fmt.Sprintf("%s_%d", platform, sID),                                    // 多个平台达到的效果不一样
		RoomId:   fmt.Sprintf("%d", sID),                                                 //房间号唯一否则消息串房间
		ShopId:   mobile,                                                                 // 登录该后台的手机号
		ShopName: fmt.Sprintf("客服%s", mobile),                                            // 客服昵称
		ShopFace: "https://img.wxcha.com/m00/86/59/7c6242363084072b82b6957cacc335c7.jpg", // 客服头像
		Platform: platform,
		Suburl:   "ws://localhost:7923/ws",                       // 订阅地址
		Pushurl:  "http://localhost:8888/open/push&platform=web", // 发布地址
	}

	serveJSON(w, resp{
		Code:    1,
		Success: true,
		Msg:     "success",
		Data:    dst,
	})
}

func apiUpload(w http.ResponseWriter, r *http.Request) {
	var (
		newPath string // 暂时只处理1个文件上传
	)
	if r.Method != "POST" {
		return
	}
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		fmt.Printf("FileName=[%s], FormName=[%s]\n", part.FileName(), part.FormName())
		if part.FileName() == "" { // this is FormData
			data, _ := ioutil.ReadAll(part)
			fmt.Printf("FormData=[%s]\n", string(data))
		} else { // This is FileData
			newPath = fmt.Sprintf("/upload/%d_%s", time.Now().Unix(), part.FileName())
			dst, _ := os.Create("./dist" + newPath) // 写入时需要dist 访问路径上不能带有 /dist
			defer dst.Close()
			io.Copy(dst, part)
		}
	}

	if newPath != "" {
		serveJSON(w, resp{
			Code:     1,
			Success:  true,
			FileName: newPath,
		})
	} else {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     "上传失败",
		})
	}
}

// apiRegister 用户注册接口  为了演示,临时采用redis存储
func apiRegister(w http.ResponseWriter, r *http.Request) {
	var (
		mobile   string
		password string
	)

	if r.Method == "POST" {
		mobile = r.FormValue("mobile")
		password = r.FormValue("password")
	}
	if mobile == "" || password == "" {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     "参数mobile or password不能为空",
		})
		return
	}
	_, err := db.RedisCli.HSet("user_table", mobile, password).Result()
	serveJSON(w, resp{
		Code:    1,
		Success: true,
		Msg:     err.Error(),
	})
}

func apiLogin(w http.ResponseWriter, r *http.Request) {
	var (
		mobile   string
		password string
	)

	if r.Method == "POST" {
		mobile = r.FormValue("mobile")
		password = r.FormValue("password")
	}

	oldPwd, err := db.RedisCli.HGet("user_table", mobile).Result()

	if oldPwd == password {
		serveJSON(w, resp{
			Code:    1,
			Success: true,
		})
	} else {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     "参数" + mobile + "不能为空" + err.Error(),
		})
	}
}

// 查看所有与自己聊天的用户
func apiFindUserList(w http.ResponseWriter, r *http.Request) {
	var (
		shopId string
		u      model.User
		dst    []model.User
	)

	if r.Method == "POST" {
		shopId = r.FormValue("shop_id")
	}
	if shopId == "" {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     "参数shop_id不能为空",
		})
		return
	}
	ids := make([]int64, 0)
	// 查询在线人数
	idsTmp, err := GetShopList(shopId, 0, 50)
	for _, v := range idsTmp {
		id, _ := strconv.ParseInt(v, 10, 64)
		ids = append(ids, id)
	}
	userIds, err := cache.Online.KeysByUserIds(ids)

	// 查询已读/未读
	for _, v := range userIds {
		if v == "" {
			continue
		}
		logger.Logger.Debug("apiFindUserList", zap.Any("userJson", v))
		json.Unmarshal([]byte(v), &u)
		fmt.Printf("解析用户  : %+v \r\n", u)
		if u.DeviceId == "" {
			continue
		}
		index, err := cache.Online.GetMessageAckMapping(u.DeviceId, u.RoomId) // 获取消息已读偏移
		if err != nil {
			fmt.Printf("获取消息已读偏移 error : %+v  \n", err)
			continue
		}
		count, err := GetMessageCount(u.RoomId, index, "+inf") // 拿到偏移去统计未读
		if err != nil {
			fmt.Printf("拿到偏移去统计未读 error : %+v  \n", err)
			continue
		}

		lastMessage, err := GetMessageList(u.RoomId, 0, 0) // 取回消息
		if err != nil {
			fmt.Printf("取回消息 error : %+v  \n", err)
			continue
		}

		u.Unread = model.Int64(count)
		u.LastMessage = lastMessage

		dst = append(dst, u)
	}

	if err == nil {
		serveJSON(w, resp{
			Code:     1,
			Success:  true,
			UserList: dst,
		})
	} else {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     err.Error(),
		})
	}
}

func apiPush(w http.ResponseWriter, r *http.Request) {
	var (
		userId string
		roomId string
		shopId string
		typ    string
		msg    string
		msgId  int64
	)

	if r.Method == "POST" {
		roomId = r.FormValue("room_id")
		typ = r.FormValue("type")
		msg = r.FormValue("msg")
		userId = r.FormValue("user_id")
		shopId = r.FormValue("shop_id")
	}
	if roomId == "" || typ == "" || msg == "" || userId == "" || shopId == "" {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     "参数room_id type msg user_id shop_id不能为空",
		})
		return
	}
	msgId = time.Now().UnixNano() // 消息唯一id 为了方便临时demo采用该方案， 后期线上可以用雪花算法
	body := fmt.Sprintf(`{"user_id":%s, "shop_id":%s, "type":"%s", "msg":"%s", "room_id":"%s", "dateline":%d, "id":"%d"}`,
		userId, shopId, typ, msg, roomId, time.Now().Unix(), msgId)

	buf := &pb.PushMsg{
		Type:      pb.PushMsg_ROOM,
		Operation: protocol.OpSendMsgReply,
		Speed:     2,
		Server:    config.Connect.LocalAddr,
		RoomId:    roomId,
		Msg:       []byte(body),
	}
	bytes, err := proto.Marshal(buf)
	if err == nil {
		// 写入商户列表
		err = AddShopList(shopId, userId)

		// 推送 或者 写入kafka 队列等
		err = cache.Queue.Publish(config.Global.PushAllTopic, bytes)
		if err == nil {
			// 消息持久化
			err = AddMessageList(roomId, msgId, body)
		}
	}
	if err == nil {
		serveJSON(w, resp{
			Code:    1,
			Success: true,
		})
	} else {
		serveJSON(w, resp{
			Code:    0,
			Success: false,
			Msg:     err.Error(),
		})
	}
}

// redis 操作区

// AddShopList 将用户添加到商户列表
// zadd  shop_id  time() user_id
func AddShopList(shopId string, userId string) error {
	score := time.Now().Unix()
	key := fmt.Sprintf("userList:%s", shopId)
	err := db.RedisCli.ZAdd(key, redis.Z{
		Score:  float64(score),
		Member: userId,
	}).Err()

	if err != nil {
		return gerrors.WrapError(err)
	}

	return nil
}

// GetShopList 查询在线人数
// zrevrange  shop_id  0, 50
func GetShopList(shopId string, start, stop int64) ([]string, error) {
	key := fmt.Sprintf("userList:%s", shopId)
	ids, err := db.RedisCli.ZRevRange(key, start, stop).Result()

	if err != nil {
		return ids, gerrors.WrapError(err)
	}

	return ids, nil
}

// AddMessageList 将消息添加到对应房间 roomId
// zadd  roomId  time() msg
func AddMessageList(roomId string, id int64, msg string) error {
	key := fmt.Sprintf("messagelist:%s", roomId)
	err := db.RedisCli.ZAddNX(key, redis.Z{
		Score:  float64(id),
		Member: msg,
	}).Err()

	if err != nil {
		return gerrors.WrapError(err)
	}

	return nil
}

// GetMessageCount 统计未读
func GetMessageCount(roomId, start, stop string) (int64, error) {
	key := fmt.Sprintf("messagelist:%s", roomId)
	dst, err := db.RedisCli.ZCount(key, start, stop).Result()

	if err != nil {
		return dst, gerrors.WrapError(err)
	}

	return dst, nil
}

// GetMessageList 取回消息
func GetMessageList(roomId string, start, stop int64) ([]string, error) {
	key := fmt.Sprintf("messagelist:%s", roomId)
	dst, err := db.RedisCli.ZRevRange(key, start, stop).Result()

	if err != nil {
		return dst, gerrors.WrapError(err)
	}

	return dst, nil
}

func main() {
	logger.Init()

	db.InitRedis(config.Global.RedisIP, config.Global.RedisPassword)

	ipAddr := ":8888"

	t := time.Now().Local().Format("2006-01-02 15:04:05 -0700")

	logger.Logger.Info("http demo 服务已经开启", zap.String("demo_http_server_ip_port", ipAddr+"  "+t))

	err := http.ListenAndServe(ipAddr, &router{})
	if err != nil {
		fmt.Println(err)
	}
}
