package starGo

import (
	"net"
	"sync"
)

type UdpClient struct {
	addr      *net.UDPAddr
	conn      *net.UDPConn
	stop      bool
	receiveCh chan []byte
	sendCh    chan []byte
	mutex     sync.RWMutex
}

func newUdpClient(conn *net.UDPConn, addr *net.UDPAddr) *UdpClient {
	return &UdpClient{
		addr:      addr,
		conn:      conn,
		stop:      false,
		receiveCh: make(chan []byte, 1024),
		sendCh:    make(chan []byte, 1024),
	}
}

func (c *UdpClient) GetAddr() *net.UDPAddr {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.addr
}

func (c *UdpClient) GetStop() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stop
}

func (c *UdpClient) SetStop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.stop = true
}

func (c *UdpClient) GetReceiveData(headerLen int32, data []byte) (message []byte, exists bool) {
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

func (c *UdpClient) SetSendQueue(data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.sendCh <- data
}

func (c *UdpClient) read() {
	Go(func(Stop chan struct{}) {
		defer func() {
			c.SetStop()
		}()

		for !c.GetStop() {
			select {
			case <-Stop:
				return
			case receiveData := <-c.receiveCh:
				Go(func(Stop1 chan struct{}) {
					select {
					case <-Stop1:
						return
					default:
						message, exists := c.GetReceiveData(udpReceiveDataHeaderLen, receiveData)
						if exists {
							udpHandlerReceiveFunc(message, c.GetAddr().String())
						}
					}
				})
			}
		}
	})
}

func (c *UdpClient) write() {
	Go(func(Stop chan struct{}) {
		defer func() {
			c.SetStop()
		}()

		if !c.GetStop() {
			select {
			case <-Stop:
				return
			case sendData := <-c.sendCh:
				Go(func(Stop1 chan struct{}) {
					select {
					case <-Stop1:
						return
					default:
						c.conn.WriteToUDP(sendData, c.GetAddr())
					}
				})
			}
		}
	})
}

func (c *UdpClient) start() {
	c.read()
	c.write()
}
