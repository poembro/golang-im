package main

import (
	"encoding/json"
	"fmt"

	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang-im/config"
	"golang-im/pkg/db"
	"golang-im/pkg/logger"
	"golang-im/pkg/gerrors"
	"golang-im/internal/logic/cache"
	"golang-im/internal/logic/model"
	"golang-im/pkg/pb"
	"golang-im/pkg/protocol"
	
	"go.uber.org/zap"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

type router struct {
}

func (ro *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth/logout": // 退出
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	case "/auth/register": // 注册
		apiRegister(w, r)
	case "/auth/login": // 登录
		apiLogin(w, r)
	case "/open/push": // 接收消息 并发起grpc至logic服务的SendMessage方法
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

func serveJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Server", "poembro")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	content, err := json.Marshal(data)
	if err == nil {
		w.Write(content)
	} else {
		w.Write([]byte(`{"code":0, "error":"解析JSON出错"}`))
	}
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

type resp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg,omitempty"`
	Success bool   `json:"success,omitempty"`

	Mobile   string       `json:"mobile,omitempty"`
	Password string       `json:"password,omitempty"`
	FileName string       `json:"filename,omitempty"`
	UserList []model.User `json:"user_list,omitempty"`
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

	res := resp{
		Code:    0,
		Success: false,
		Msg:     "",
	}
	if newPath != "" {
		res.Code = 1
		res.Success = true
		res.FileName = newPath
	} else {
		res.Code = 0
		res.Msg = "上传失败"
	}
	serveJSON(w, res)
}

func apiRegister(w http.ResponseWriter, r *http.Request) {
	var (
		mobile   string
		password string
	)

	if r.Method == "POST" {
		mobile = r.FormValue("mobile")
		password = r.FormValue("password")
	}

	_, err := db.RedisCli.HSet("user_table", mobile, password).Result()

	res := resp{
		Code:    0,
		Success: false,
		Msg:     "",
	}
	if mobile != "" {
		res.Code = 1
		res.Success = true
	} else {
		res.Code = 0
		res.Msg = "参数" + mobile + "不能为空" + err.Error()
	}
	serveJSON(w, res)
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

	res := resp{
		Code:    0,
		Success: false,
		Msg:     "",
	}
	if oldPwd == password {
		res.Code = 1
		res.Success = true
	} else {
		res.Code = 0
		res.Msg = "参数" + mobile + "不能为空" + err.Error()
	}
	serveJSON(w, res)
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

	ids := make([]int64, 0)
	// 查询在线人数 
	idsTmp, err := GetShopList("userList:"+shopId, 0, 50)
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
		//fmt.Printf("解析用户  : %+v", u)
		if u.DeviceId == "" {
			continue
		}
		index, err := cache.Online.GeMessageAckMapping(int64(u.UserId), u.RoomId) // 获取消息已读偏移
		if err != nil {
			fmt.Printf("获取消息已读偏移 error : %+v", err)
			continue
		}
		count, err := GetMessageCount(u.RoomId, index, "+inf") // 拿到偏移去统计未读
		if err != nil {
			fmt.Printf("拿到偏移去统计未读 error : %+v", err)
			continue
		}

		lastMessage, err := GetMessageList(u.RoomId, 0, 0) // 取回消息
		if err != nil {
			fmt.Printf("取回消息 error : %+v", err)
			continue
		}

		u.Unread = model.Int64(count)
		u.LastMessage = lastMessage

		dst = append(dst, u)
	}

	res := resp{
		Code:    0,
		Success: false,
		Msg:     "",
	}
	if err == nil {
		res.Code = 1
		res.Success = true
		res.UserList = dst
	} else {
		res.Code = 0
		res.Msg = "参数" + shopId + "不能为空" + err.Error()
	}
	serveJSON(w, res)
}


func apiPush(w http.ResponseWriter, r *http.Request) {
	var (
		userId string
		roomId string
		shopId string
		typ    string
		msg    string 
		msgId        int64 
	)

	if r.Method == "POST" {
		roomId = r.FormValue("room_id")
		typ = r.FormValue("type")
		msg = r.FormValue("msg")
		userId = r.FormValue("user_id")
		shopId = r.FormValue("shop_id")
	}

	msgId = time.Now().UnixNano() // 消息唯一id 为了方便临时demo采用该方案， 后期线上可以用雪花算法

	body := fmt.Sprintf(`{"user_id" : %s,"shop_id":%s,"type" : "%s","msg" : "%s","room_id" : "%s","dateline" : %d,"id" : "%d"}`,
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
		// 推送 或者 写入kafka 队列等
		err = cache.Queue.Publish(config.Global.PushAllTopic, bytes)
		if err == nil {
			// 消息持久化
			err = AddMessageList(roomId, msgId, body) 
			// 写入商户列表
			AddShopList(shopId, userId)
		}
	}

	res := resp{
		Code:    0,
		Success: false,
		Msg:     "",
	}
	if err == nil {
		res.Code = 1
		res.Success = true
	} else {
		res.Code = 0
		res.Msg = "error" + err.Error()
	}
	serveJSON(w, res)
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

//下方 redis 操作区域

// AddShopList 将用户添加到商户列表
// zadd  shop_id  time() user_id
func AddShopList(shopId string, userId string) error {
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

// GetShopList 查询在线人数
// zrevrange  shop_id  0, 50
func GetShopList(shopId string, start, stop int64) ([]string, error)  {
    ids, err := db.RedisCli.ZRevRange("userList:" + shopId, start, stop).Result()

    if err != nil {
        return ids, gerrors.WrapError(err)
    }

    return ids, nil
}


// AddMessageList 将消息添加到对应房间 roomId
// zadd  roomId  time() msg
func AddMessageList(roomId string, id int64, msg string) error {
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
func GetMessageCount(roomId , start, stop string) (int64, error) {
    dst, err := db.RedisCli.ZCount("messagelist:" + roomId, start, stop).Result()

    if err != nil {
        return dst, gerrors.WrapError(err)
    }

    return dst, nil
}

// GetMessageList 取回消息
func GetMessageList(roomId string, start, stop int64) ([]string, error) {
    dst, err := db.RedisCli.ZRevRange("messagelist:" + roomId, start, stop).Result()

    if err != nil {
        return dst, gerrors.WrapError(err)
    }

    return dst, nil
}
