package http

import (
	"net/http"
	"os"

	"golang-im/pkg/logger"
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

// StartHttpServer 启动HTTP框架 监听端口
func StartHttpServer(ipAddr string) {
	logger.Logger.Info("http server start")

	err := http.ListenAndServe(ipAddr, &router{})
	if err != nil {
		panic(err)
	}
}

// StaticServer 静态文件处理
func StaticServer(w http.ResponseWriter, req *http.Request) {
	indexs := []string{"index.html", "index.htm"}
	basedir, _ := os.Getwd()                     // 获取当前目录路径 /webser/go_wepapp/golang-im
	filePath := basedir + "/dist" + req.URL.Path //注意 注意 注意:这里只能处理 dist 目录下的文件
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

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
