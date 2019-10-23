package starGo

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

// Start ...启动信号管理器
func Start() {
	Go(func(Stop chan struct{}) {
		sign := make(chan os.Signal)
		signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

		for allForStopSignal == 0 {
			// 准备接收信息
			tempSig := <-sign

			// 输出信号
			DebugLog("收到信号:%v", sign)

			if tempSig == syscall.SIGHUP {
				DebugLog("收到重启的信号，准备重新加载配置")

				// 重新加载
				systemReload()

				DebugLog("收到重启的信号，重新加载配置完成")
			} else {
				fmt.Println("收到退出程序的信号，开始退出……")

				// 调用退出的方法
				systemExit()
				close(sign)

				fmt.Println("收到退出程序的信号，退出完成……")

				// 一旦收到信号，则表明管理员希望退出程序，则先保存信息，然后退出
				os.Exit(0)
			}
		}
	})
}

func WaitForSystemExit() {
	waitAllGroup.Wait()
}

func RegisterSystemExitFunc(f func()) {
	systemExitFunc = append(systemExitFunc, f)
}

func RegisterSystemReloadFunc(f func()) {
	systemReloadFunc = append(systemReloadFunc)
}

func systemExit() {
	for _, f := range systemExitFunc {
		f()
	}

	// 关闭所有tcp连接
	clientMap.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		client.SetStop()
		client.GetConn().Close()
		clientMap.Delete(key)

		return true
	})

	// 通知所有协程退出
	stopChanForLog <- struct{}{}
	for i := int32(0); i < goCount; i++ {
		stopChanForLog <- struct{}{}
	}

	// 更新停止信号
	if !atomic.CompareAndSwapInt32(&allForStopSignal, 0, 1) {
		return
	}

	// 关闭所有通道
	close(stopChanForGo)
	close(stopChanForLog)
	close(logCh)

	fmt.Println("服务器停止")
}

func systemReload() {
	for _, f := range systemReloadFunc {
		f()
	}
}
