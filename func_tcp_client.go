package starGo

import (
	"io"
	"net"
	"sync"
	"time"
)

type ClientCallBack func(message []byte, addr string)

type ClientExpireCallBack func(addr []string)

type Client struct {
	conn         net.Conn
	stop         bool
	activeTime   int64
	receiveQueue []byte
	sendCh       chan []byte
	mutex        sync.RWMutex
}

func newTcpClient(conn net.Conn) *Client {
	return &Client{
		conn:         conn,
		stop:         false,
		activeTime:   time.Now().Unix(),
		receiveQueue: make([]byte, 0),
		sendCh:       make(chan []byte, 1024),
	}
}

func (c *Client) GetConn() net.Conn {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.conn
}

func (c *Client) GetStop() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stop
}

func (c *Client) SetStop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.stop = true
}

func (c *Client) GetActiveTime() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.activeTime
}

func (c *Client) SetActiveTime() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.activeTime = time.Now().Unix()
}

func (c *Client) GetReceiveData(headerLen int32) (message []byte, exists bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.receiveQueue) < int(headerLen) {
		return
	}

	// 获取头部信息
	header := c.receiveQueue[:headerLen]

	// 将头部数据转换为内容的长度
	contentLength := BytesToInt32(header, true)

	// 判断长度是否满足
	if len(c.receiveQueue) < int(headerLen+contentLength) {
		return
	}

	// 提取消息内容
	message = c.receiveQueue[headerLen : headerLen+contentLength]

	// 将对应的数据截断，以得到新的内容
	c.receiveQueue = c.receiveQueue[headerLen+contentLength:]

	// 存在合理的数据
	exists = true

	return
}

func (c *Client) AppendReceiveQueue(message []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.receiveQueue = append(c.receiveQueue, message...)
}

func (c *Client) AppendSendQueue(message []byte) {
	c.sendCh <- message
}

func (c *Client) start() {
	Go(func(Stop chan struct{}) {
		defer func() {
			c.GetConn().Close()
			c.SetStop()
		}()
		for !c.GetStop() {
			readBytes := make([]byte, 1024)
			n, err := c.GetConn().Read(readBytes)
			if err != nil {
				if err != io.EOF {
					ErrorLog("读取消息错误：%s，本次读取的字节数为：%d", err, n)
				}
				break
			}
			c.AppendReceiveQueue(readBytes[:n])

			Go(func(Stop chan struct{}) {
				message, exists := c.GetReceiveData(tcpReceiveDataHeaderLen)
				if exists && tcpHandlerReceiveFunc != nil {
					tcpHandlerReceiveFunc(message, c.GetConn().RemoteAddr().String())
				}
			})
		}
	})
	Go(func(Stop chan struct{}) {
		for !c.GetStop() {
			select {
			case message := <-c.sendCh:
				Go(func(Stop chan struct{}) {
					header := Int32ToBytes(int32(len(message)), true)
					header = append(header, message...)
					_, err := c.GetConn().Write(header)
					if err != nil {
						ErrorLog("向客户端:%v发送数据出错,错误信息:%v", c.GetConn().RemoteAddr().String(), err)
						c.GetConn().Close()
						c.SetStop()
					}
				})
			case <-Stop:
				return
			}
		}
	})
}

func RegisterTcpClient(c *Client) {
	tcpClientMap.Store(c.GetConn().RemoteAddr().String(), c)
}

func clearExpireTcpClient() {
	Go(func(Stop chan struct{}) {
		t := time.NewTicker(5 * time.Second)
		for allForStopSignal == 0 {
			<-t.C
			removeClient := make([]string, 0)
			tcpClientMap.Range(func(key, value interface{}) bool {
				client := value.(*Client)
				if client.GetActiveTime()+clientExpireTime <= time.Now().Unix() {
					removeClient = append(removeClient, key.(string))
				}

				return true
			})

			// 移除过期的客户端
			callBackList := make([]string, 0)
			for _, key := range removeClient {
				value, exists := tcpClientMap.Load(key)
				if !exists {
					continue
				}

				// 再次判断是否过期，防止将要移除时有发生通信的事件
				client := value.(*Client)
				if client.GetActiveTime()+clientExpireTime > time.Now().Unix() {
					continue
				}

				// 移除过期客户端
				client.SetStop()
				client.GetConn().Close()
				tcpClientMap.Delete(key)
				callBackList = append(callBackList, key)
			}

			if len(callBackList) > 0 {
				InfoLog("移除过期客户端连接:%v", callBackList)
				if tcpClientExpireHandleFunc != nil {
					tcpClientExpireHandleFunc(callBackList)
				}
			}
		}
	})
}
