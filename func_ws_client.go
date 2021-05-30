package starGo

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn       *websocket.Conn
	stop       bool
	activeTime int64
	sendCh     chan []byte
	mutex      sync.RWMutex
}

func newWebSocketClient(conn *websocket.Conn) *WebSocketClient {
	return &WebSocketClient{
		conn:       conn,
		stop:       false,
		activeTime: time.Now().Unix(),
		sendCh:     make(chan []byte, 1024),
	}
}

func (c *WebSocketClient) GetConn() *websocket.Conn {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.conn
}

func (c *WebSocketClient) GetStop() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stop
}

func (c *WebSocketClient) SetStop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.stop = true
}

func (c *WebSocketClient) GetActiveTime() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.activeTime
}

func (c *WebSocketClient) SetActiveTime() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.activeTime = time.Now().Unix()
}

func (c *WebSocketClient) GetReceiveData(headerLen int32, data []byte) (message []byte, exists bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(data) < int(headerLen) {
		return
	}

	// 获取头部信息
	header := data[:headerLen]

	// 将头部数据转换为内容的长度
	contentLength := BytesToInt32(header, true)

	// 判断长度是否满足
	if len(data) < int(headerLen+contentLength) {
		return
	}

	// 提取消息内容
	message = data[headerLen : headerLen+contentLength]

	// 存在合理的数据
	exists = true

	return
}

func (c *WebSocketClient) AppendSendQueue(message []byte) {
	c.sendCh <- message
}

func (c *WebSocketClient) start() {
	Go(func(Stop chan struct{}) {
		defer func() {
			_ = c.GetConn().Close()
			c.SetStop()
		}()
		for !c.GetStop() {
			_, data, err := c.GetConn().ReadMessage()
			if err != nil {
				ErrorLog("读取消息错误:%v", err)
				break
			}

			Go2(func() {
				message, exists := c.GetReceiveData(wsReceiveDataHeaderLen, data)
				if exists && wsHandlerReceiveFunc != nil {
					wsHandlerReceiveFunc(message, c.GetConn().RemoteAddr().String())
				}
			})
		}
	})
	Go(func(Stop chan struct{}) {
		for !c.GetStop() {
			select {
			case message := <-c.sendCh:
				Go2(func() {
					err := c.GetConn().WriteMessage(websocket.BinaryMessage, message)
					if err != nil {
						ErrorLog("向客户端:%v发送数据出错,错误信息:%v", c.GetConn().RemoteAddr().String(), err)
						_ = c.GetConn().Close()
						c.SetStop()
					}
				})
			case <-Stop:
				return
			}
		}
	})
}

func registerWebSocketClient(c *WebSocketClient) {
	wsClientMap.Store(c.GetConn().RemoteAddr().String(), c)
}

func clearExpireWebSocketClient() {
	Go(func(Stop chan struct{}) {
		for allForStopSignal == 0 {
			t := time.NewTicker(5 * time.Second)
			<-t.C
			removeClient := make([]string, 0)
			wsClientMap.Range(func(key, value interface{}) bool {
				client := value.(*WebSocketClient)
				if client.GetActiveTime()+clientExpireTime <= time.Now().Unix() {
					removeClient = append(removeClient, key.(string))
				}

				return true
			})

			// 移除过期的客户端
			callBackList := make([]string, 0)
			for _, key := range removeClient {
				value, exists := wsClientMap.Load(key)
				if !exists {
					continue
				}

				// 再次判断是否过期，防止将要移除时有发生通信的事件
				client := value.(*WebSocketClient)
				if client.GetActiveTime()+clientExpireTime > time.Now().Unix() {
					continue
				}

				// 移除过期客户端
				client.SetStop()
				_ = client.GetConn().Close()
				wsClientMap.Delete(key)
				callBackList = append(callBackList, key)
			}

			if len(callBackList) > 0 && wsClientExpireHandleFunc != nil {
				wsClientExpireHandleFunc(callBackList)
			}
		}
	})
}
