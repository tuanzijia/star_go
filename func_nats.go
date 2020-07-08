package starGo

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type NatCallBack func(result *NatResult)

type NatResult struct {
	Message []byte
	Reply   string
}

func StartNatConn(addr, natName string) {
	option := make([]nats.Option, 0)
	option = append(option, nats.Name(natName),
		nats.PingInterval(30*time.Second),
		nats.MaxPingsOutstanding(5),
		nats.ReconnectWait(3*time.Second),
		nats.MaxReconnects(200),
		nats.DisconnectErrHandler(func(conn *nats.Conn, e error) {
			ErrorLog("nat连接断开,等待重新连接时间:%v分钟", 10)
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			ErrorLog("nat重新连接，URL:%v", conn.ConnectedUrl())
		}),
		nats.ClosedHandler(func(conn *nats.Conn) {
			ErrorLog("nat退出,err:%v", conn.LastError())
		}))

	nc, err := nats.Connect(addr, option...)

	if err != nil {
		ErrorLog("连接nats服务器失败，错误信息:%v", err)
		panic(fmt.Errorf("连接nats服务器失败"))
	}

	natConn = nc
}

// 管道模式订阅
func SubscribeChannel(channel string, channelCount int32, cb NatCallBack) {
	_, exists := natChMap.Load(channel)
	if exists {
		ErrorLog("%v通道已被订阅", channel)
	}

	ch := make(chan *nats.Msg, 64)
	sub, err := natConn.ChanSubscribe(channel, ch)
	if err != nil {
		ErrorLog("订阅%v错误,错误信息:%v", channel, err)
		return
	}

	natChMap.Store(channel, channelCount)

	for i := int32(0); i < channelCount; i++ {
		Go2(func() {
			defer func() {
				_ = sub.Unsubscribe()
				_ = sub.Drain()
			}()

			for allForStopSignal == 0 {
				msg := <-ch
				if cb != nil {
					result := &NatResult{
						Message: msg.Data,
						Reply:   msg.Reply,
					}

					cb(result)
				}
			}
		})
	}
}

// 异步模式订阅
func SubscribeAsync(channel string, cb NatCallBack) {
	_, err := natConn.Subscribe(channel, func(msg *nats.Msg) {
		Go2(func() {
			if cb != nil {
				result := &NatResult{
					Message: msg.Data,
					Reply:   msg.Reply,
				}

				cb(result)
			}
		})
	})

	if err != nil {
		ErrorLog("使用异步模式订阅%v错误,错误信息:%v", channel, err)
	}
}

// 队列模式订阅
func SubscribeQueue(channel, queue string, cb NatCallBack) {
	_, err := natConn.QueueSubscribe(channel, queue, func(msg *nats.Msg) {
		Go2(func() {
			if cb != nil {
				result := &NatResult{
					Message: msg.Data,
					Reply:   msg.Reply,
				}

				cb(result)
			}
		})
	})

	if err != nil {
		ErrorLog("使用队列:%v订阅:%v出错,错误信息:%v", queue, channel, err)
	}
}

// 发布消息
func Publish(channel string, data []byte) {
	err := natConn.Publish(channel, data)
	if err != nil {
		ErrorLog("发布消息出错,错误信息:%v", err)
	}

	err = natConn.Flush()
	if err != nil {
		ErrorLog("发布消息后刷新出错,错误信息:%v", err)
	}
}

// 远程rpc调用
func RpcCall(channel string, data []byte, timeout int32) ([]byte, error) {
	msg, err := natConn.Request(channel, data, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
