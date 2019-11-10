package starGo

import "net"

func StartTcpServer(addr string, handler ClientCallBack, clientExpireHandler ClientExpireCallBack, headerLen int32) error {
	InfoLog("开始监听Tcp地址:%v", addr)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		ErrorLog("tcp监听地址%v出错,错误信息:%v", addr, err)
		return err
	}
	Go(func(Stop chan struct{}) {
		c := make(chan struct{})
		Go(func(Stop1 chan struct{}) {
			select {
			case <-Stop1:
			case <-c:
			}
			listen.Close()
		})

		for allForStopSignal == 0 {
			c, err := listen.Accept()
			if err != nil {
				ErrorLog("接收客户端连接失败,错误信息:%v", err)
				break
			}

			// 新注册客户端
			client := newTcpClient(c)
			registerTcpClient(client)
			client.start()

			InfoLog("收到客户端:%v的连接请求", c.RemoteAddr().String())
		}
	})

	// 注册回调方法
	tcpHandlerReceiveFunc = handler
	tcpClientExpireHandleFunc = clientExpireHandler

	// 记录头部数据长度
	tcpReceiveDataHeaderLen = headerLen

	// 启动客户端处理协程
	clearExpireTcpClient()

	return nil
}

func GetTcpClient(addr string) *Client {
	client, exists := tcpClientMap.Load(addr)
	if !exists {
		return nil
	}

	return client.(*Client)
}
