package starGo

import "github.com/gin-gonic/gin"

type webRequestType string

const (
	GET     webRequestType = "GET"
	POST    webRequestType = "POST"
	DELETE  webRequestType = "DELETE"
	PATCH   webRequestType = "PATCH"
	PUT     webRequestType = "PUT"
	OPTIONS webRequestType = "OPTIONS"
	HEAD    webRequestType = "HEAD"
)

type WebServer struct {
	router          *gin.Engine // 路由实例
	serverIpAddress string
	isUseMiddleware bool
	registerFuncMap map[webRequestType]map[string]bool
}

func NewWebServer(ip string, isUseMiddleware bool) *WebServer {
	return &WebServer{
		router:          gin.New(),
		serverIpAddress: ip,
		isUseMiddleware: isUseMiddleware,
		registerFuncMap: make(map[webRequestType]map[string]bool),
	}
}

func (w *WebServer) StartWebServer() {
	// 启用中间件
	if w.isUseMiddleware {
		w.router.Use(gin.Recovery())
		w.router.Use(gin.Logger())
	}

	if w.serverIpAddress == "" {
		err := w.router.Run()
		if err != nil {
			ErrorLog("启动webServer时出错,错误信息:%v", err)
		}
	} else {
		err := w.router.Run(w.serverIpAddress)
		if err != nil {
			ErrorLog("启动webServer时出错,错误信息:%v", err)
		}
	}
}

func (w *WebServer) RegisterRequestHandleFunc(requestType webRequestType, url string, handleFunc gin.HandlerFunc) {
	_, mapExists := w.registerFuncMap[requestType]
	if mapExists {
		_, urlExists := w.registerFuncMap[requestType][url]
		if urlExists {
			ErrorLog("已注册相同URL路径,requestType:%v, url:%v", requestType, url)
			return
		}
	} else {
		w.registerFuncMap[requestType] = make(map[string]bool)
	}

	// 注册方法
	w.registerFuncMap[requestType][url] = true

	switch requestType {
	case GET:
		w.router.GET(url, handleFunc)
	case POST:
		w.router.POST(url, handleFunc)
	case DELETE:
		w.router.DELETE(url, handleFunc)
	case PATCH:
		w.router.PATCH(url, handleFunc)
	case PUT:
		w.router.PUT(url, handleFunc)
	case OPTIONS:
		w.router.OPTIONS(url, handleFunc)
	case HEAD:
		w.router.HEAD(url, handleFunc)
	}
}
