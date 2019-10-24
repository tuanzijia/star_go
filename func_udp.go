package starGo

import (
	"net"
)

func StartUdpServer(addr string, handler ClientCallBack, headerLen int32) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		ErrorLog("监听Udp地址:%v失败,错误信息:%v", addr, err)
		return err
	}
	conn, err1 := net.ListenUDP("udp", udpAddr)
	if err1 != nil {
		ErrorLog("监听Udp地址:%v失败,错误信息:%v", addr, err)
		return err1
	}

	// 监听
	listen(conn)

	// 注册回调方法
	udpHandlerReceiveFunc = handler

	// 记录头部数据长度
	udpReceiveDataHeaderLen = headerLen

	return nil
}

func listen(conn *net.UDPConn) {
	// 设置缓冲区
	conn.SetReadBuffer(1024)
	conn.SetWriteBuffer(1024)

	Go(func(Stop chan struct{}) {
		c := make(chan struct{})
		Go(func(Stop1 chan struct{}) {
			select {
			case <-Stop1:
			case <-c:
			}
			conn.Close()
		})

		listenTrue(conn)

		close(c)
	})
}

func listenTrue(conn *net.UDPConn) {
	data := make([]byte, 1024)
	for allForStopSignal == 0 {
		n, udpAddr, err := conn.ReadFromUDP(data)

		if err != nil {
			if err.(net.Error).Timeout() {
				continue
			}
			break
		}

		if n <= 0 {
			continue
		}

		addr := udpAddr.String()
		clientInfo, exists := udpClientMap.Load(addr)
		if !exists {
			client := newUdpClient(conn, udpAddr)
			client.start()
			udpClientMap.Store(addr, client)
		}

		client := clientInfo.(*UdpClient)
		client.receiveCh <- data
	}
}

func GetUdpClient(addr string) *UdpClient {
	client, exists := udpClientMap.Load(addr)
	if !exists {
		return nil
	}

	return client.(*UdpClient)
}
