package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang-im/config"
	"golang-im/pkg/db"
	"golang-im/pkg/logger"

	"go.uber.org/zap"

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
	case "/open/finduserlist": 	// 在线列表
		// TODO
		apiFindUserList(w, r)
	case "/upload/file": 	//文件上传接口 
		apiUpload(w, r)
	default:
		StaticServer(w, r)
		//notFound(w)
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
	filePath := "./dist" + req.URL.Path
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

func apiUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
        if part.FileName() == "" {  // this is FormData
            data, _ := ioutil.ReadAll(part)
            fmt.Printf("FormData=[%s]\n", string(data))
        } else {    // This is FileData
            dst, _ := os.Create("./" + part.FileName() + ".upload")
            defer dst.Close()
            io.Copy(dst, part)
        }
    }
}

type resp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg,omitempty"`
	Success bool   `json:"success,omitempty"`

	Mobile string `json:"mobile,omitempty"`
	Password string `json:"password,omitempty"`
}

func apiRegister(w http.ResponseWriter, r *http.Request) {
	var (
		mobile string
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
		Msg: "",
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
		mobile string
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
		Msg: "",
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

func apiFindUserList(w http.ResponseWriter, r *http.Request) {
	var (
		roomId string
	)

	if r.Method == "POST" {
		mobile = r.FormValue("mobile")
		password = r.FormValue("password")
	}

	oldPwd, err := db.RedisCli.HGet("user_table", mobile).Result()

	res := resp{
		Code:    0,
		Success: false,
		Msg: "",
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

func apiPush(w http.ResponseWriter, r *http.Request) {
	var ( 
		userId string
		roomId string
		face string
		nickname string 
		typ string
		msg string
		mobile string

		userIdInt64 int64
		PushAllTopic  = "push_all_topic"  // 全服消息队列
	)

	if r.Method == "POST" {
		roomId = r.FormValue("room_id")
		typ = r.FormValue("type")
		mobile = r.FormValue("mobile") 
		msg = r.FormValue("msg") 
	} else {
		userId = r.Form["mid"][0]
		userIdInt64, _ := strconv.ParseInt(userId, 10, 64)

		face := r.Form["face"][0]
		nickname := r.Form["nickname"][0] 
	}

	body := fmt.Sprintf(`{
        "me" : { "mid" : %d, "nickname" : "%s", "mobile" : "%s", "face" : "%s"}, --记录发送人
        "type" : "%s",
        "msg" : "%s",
        "room_id" : "%s", 
        "dateline" : %d,
        "id" : "%s",
    }`, userIdInt64, nickname, mobile, face, typ, msg, roomId, time.Now().Unix(), time.Now().UnixNano())

	buf := &pb.PushMsg{
        Type:      pb.PushMsg_ROOM,
        Operation: protocol.OpSendMsgReply,
        Speed:     2,
        Server:    config.Connect.LocalAddr,
        RoomId:    roomId,
        Msg:       body,
	}
    bytes, err := proto.Marshal(buf)
    if err == nil {
		err = cache.Queue.Publish(PushAllTopic, bytes)
	}

	res := resp{
		Code:    0,
		Success: false,
		Msg: "",
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
