package apihttp

import (
	"golang-im/conf"
	"golang-im/internal/logic/service"
	"golang-im/pkg/logger"
	"golang-im/pkg/util"
	"net/http"
)

func cosMiddleware(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return true
	}

	return true
}

func verifyMiddleware(w http.ResponseWriter, r *http.Request, next func(w http.ResponseWriter, r *http.Request)) {
	if r.Method == "OPTIONS" {
		return
	}
	// 解析token
	var token string
	token = r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("token")
	}
	if token == "" {
		OutJson(w, OutData{Code: -1, Success: false, Msg: "参数token不能为空"})
		return
	}
	tokenInfo, err := util.DecryptToken(token)
	if tokenInfo == nil || err != nil {
		OutJson(w, OutData{Code: -1, Success: false, Msg: "参数token认证错误"})
		return
	}
	//logger.Logger.Debug("auth", zap.Any("token", tokenInfo))
	// 去执行后续handler逻辑
	next(w, r)
}

type Router struct{}

func (h *Router) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(`pong`))
}
func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//defer util.RecoverPanic()
	if !cosMiddleware(w, r) {
		return
	}

	switch r.URL.Path {
	case "/ping": // ping
		h.Ping(w, r)
	case "/auth/logout": // 退出
		http.Redirect(w, r, "/admin/login.html", http.StatusFound)
	case "/open/adduser": // 用户创建
		h.apiAddUser(w, r)
	case "/open/push": // 接收消息写入mq
		verifyMiddleware(w, r, h.apiPush)
	case "/auth/register": // 注册
		h.apiRegister(w, r)
	case "/auth/login": // 登录
		h.apiLogin(w, r)
	case "/open/finduserinfo": // 用户详细
		h.apiFindUserInfo(w, r)
	case "/open/finduserlist": // 在线列表
		verifyMiddleware(w, r, h.apiFindUserList)
	case "/upload/file": // 文件上传接口
		verifyMiddleware(w, r, h.apiUpload)
	case "/clearData": // 清理数据
		h.apiClearData(w, r)
	case "/open/listIpblack": // 黑名单列表
		verifyMiddleware(w, r, h.listIpblack)
	case "/open/addIpblack": // ip添加至黑名单
		verifyMiddleware(w, r, h.addIpblack)
	case "/open/delIpblack": // ip从黑名单删除
		verifyMiddleware(w, r, h.delIpblack)
	default:
		h.StaticServer(w, r)
	}
}

var (
	svc    *service.Service
	config *conf.Config
)

// StartHttpServer 启动HTTP框架 监听端口
func StartHttpServer(c *conf.Config) {
	logger.Logger.Info("http server start")
	config = c
	svc = service.New(c)

	err := http.ListenAndServe(c.Logic.HttpListenAddr, &Router{})
	if err != nil {
		svc.Close()
		panic(err)
	}
	svc.Close()
}
