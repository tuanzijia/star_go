package starGo

import "net"

func StartTcpServer(addr string, handler CallBack, headerLen int32) error {
	DebugLog("开始监听地址:%v", addr)
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
			client := newClient(c)
			RegisterClient(client)
			client.start()

			DebugLog("收到客户端:%v的连接请求", c.RemoteAddr().String())
		}
	})

	// 注册回调方法
	handlerReceiveFunc = handler

	// 记录头部数据长度
	receiveDataHeaderLen = headerLen

	// 启动客户端处理协程
	clearExpireClient()

	return nil
}
