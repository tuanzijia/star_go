package starGo

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		ErrorLog("webSocket获取连接时出错,错误信息:%v", err)
		return
	}

	// 新注册客户端
	client := newWebSocketClient(conn)
	client.start()
	RegisterWebSocketClient(client)

	DebugLog("收到客户端:%v连接请求", client.GetConn().RemoteAddr().String())
}

// 启动服务器
func StartWebSocketServer(addr string, url string, handler ClientCallBack, headerLen int32) error {
	DebugLog("开始监听WebSocket地址:%v,url:%v", addr, url)
	http.HandleFunc(url, handleConn)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		ErrorLog("监听webSocket地址:%v出错，错误信息:%v", addr, err)
		return err
	}

	// 注册回调方法
	wsHandlerReceiveFunc = handler

	// 记录头部数据长度
	wsReceiveDataHeaderLen = headerLen

	// 启动客户端处理协程
	clearExpireWebSocketClient()

	return nil
}
