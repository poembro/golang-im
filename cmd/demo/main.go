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

type resp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg,omitempty"`
	Success bool   `json:"success,omitempty"`

	Nickname string `json:"nickname,omitempty"`
	Password string `json:"password,omitempty"`
}

func apiRegister(w http.ResponseWriter, r *http.Request) {
	var (
		nickname string
		password string
	)

	if r.Method == "POST" {
		nickname = r.FormValue("mobile")
		password = r.FormValue("password")
	}

	_, err := db.RedisCli.HSet("user_table", nickname, password).Result()

	res := resp{
		Code:    0,
		Success: false,
		Msg: "",
	}
	if nickname != "" {
		res.Code = 1
		res.Success = true
	} else {
		res.Code = 0
		res.Msg = "参数" + nickname + "不能为空" + err.Error()
	}
	serveJSON(w, res)
}

func apiLogin(w http.ResponseWriter, r *http.Request) {
	var (
		nickname string
		password string
	)

	if r.Method == "POST" {
		nickname = r.FormValue("mobile")
		password = r.FormValue("password")
	}

	oldPwd, err := db.RedisCli.HGet("user_table", nickname).Result()

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
		res.Msg = "参数" + nickname + "不能为空" + err.Error()
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
