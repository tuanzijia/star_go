package starGo

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestInfoLog(t *testing.T) {
	StartLog("log", Debug)

	InfoLog("qwe")
	go func() {
		time.Sleep(61 * time.Second)
		stopChanForLog <- struct{}{}

		for i := int32(0); i < goCount; i++ {
			stopChanForGo <- struct{}{}
		}
		//atomic.CompareAndSwapInt32(&allForStopSignal, 0, 1)
	}()
	go func() {
		for {
			InfoLog("输出当前时间:%v", time.Now())
			InfoLog("输出当前时间:%v", time.Now())
			WarnLog("输出当前时间:%v", time.Now())
			ErrorLog("输出当前时间:%v", time.Now())
			FatalLog("输出当前时间:%v", time.Now())
			time.Sleep(1 * time.Second)
		}
	}()
	WaitForSystemExit()
}

func TestCsv_UnMarshalFile(t *testing.T) {
	StartLog("log", Debug)
	type abc struct {
		Id   int32
		Name string
		Year string
	}
	inf := make([]abc, 0)
	cs := NewCsvReader()
	cs.UnMarshalFile("csv/test.csv", &inf)
	InfoLog(inf)
	WaitForSystemExit()
}

func TestNatPublish(t *testing.T) {
	StartNatConn("127.0.0.1:4222")
	StartLog("log", Debug)

	SubscribeQueue("hello", "h1", func(message []byte) {
		InfoLog("这是队列模式h1.1收到的消息,消息:%v", string(message))
		//fmt.Println(fmt.Sprintf("这是队列模式h1.1收到的消息,消息:%v", string(message)))
	})

	SubscribeQueue("hello", "h1", func(message []byte) {
		InfoLog("这是队列模式h1.2收到的消息,消息:%v", string(message))
		//fmt.Println(fmt.Sprintf("这是队列模式h1.2收到的消息,消息:%v", string(message)))
	})

	SubscribeQueue("hello", "h1", func(message []byte) {
		InfoLog("这是队列模式h1.3收到的消息,消息:%v", string(message))
		//fmt.Println(fmt.Sprintf("这是队列模式h1.3收到的消息,消息:%v", string(message)))
	})

	SubscribeQueue("hello", "h2", func(message []byte) {
		//fmt.Println(fmt.Sprintf("这是队列模式h2.1收到的消息,消息:%v", string(message)))
		InfoLog("这是队列模式h2.1收到的消息,消息:%v", string(message))
	})

	SubscribeQueue("hello", "h2", func(message []byte) {
		//fmt.Println(fmt.Sprintf("这是队列模式h2.2收到的消息,消息:%v", string(message)))
		InfoLog("这是队列模式h2.2收到的消息,消息:%v", string(message))
	})

	SubscribeQueue("hello", "h2", func(message []byte) {
		//fmt.Println(fmt.Sprintf("这是队列模式h2.3收到的消息,消息:%v", string(message)))
		InfoLog("这是队列模式h2.3收到的消息,消息:%v", string(message))
	})

	SubscribeAsync("hello", func(message []byte) {
		//fmt.Println(fmt.Sprintf("这是异步模式收到的消息,消息:%v", string(message)))
		InfoLog("这是异步模式收到的消息,消息:%v", string(message))
	})
	SubscribeChannel("hello", 1, func(message []byte) {
		//fmt.Println(fmt.Sprintf("这是管道模式收到的消息,消息:%v", string(message)))
		InfoLog("这是管道模式收到的消息,消息:%v", string(message))
	})

	Publish("hello", []byte("你好呀"))

	go func() {
		time.Sleep(10 * time.Second)
		stopChanForLog <- struct{}{}

		//stopChanForGo <- struct{}{}

		for i := int32(0); i < goCount; i++ {
			fmt.Println(i)
			stopChanForGo <- struct{}{}
		}
	}()
	WaitForSystemExit()
}

func TestChannel(t *testing.T) {
	var wg sync.WaitGroup
	c := make(chan struct{})
	fmt.Println(1)
	wg.Add(1)
	go func() {
		<-c
		wg.Done()
		fmt.Println(2)
	}()

	go func() {
		time.Sleep(3 * time.Second)
		close(c)
	}()
	fmt.Println(3)
	wg.Wait()
}

func TestWebSocketClient(t *testing.T) {
	StartLog("log", Debug)

	err := StartTcpServer("127.0.0.1:9999", nil, 4)
	if err != nil {
		ErrorLog(err)
		return
	}

	//Start()

	//go func() {
	//	time.Sleep(5 * time.Second)
	//	fmt.Println("开始关闭通道")
	//	systemExit()
	//}()

	WaitForSystemExit()
}
