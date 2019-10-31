package starGo

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
)

func WaitForSystemExit() {
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sign
	DebugLog("收到退出信号")
	systemExit()

	waitAllGroup.Wait()
	if !atomic.CompareAndSwapInt32(&logForStopSignal, 0, 1) {
		return
	}
	close(stopChanForLog)
	waitLogGroup.Wait()
	fmt.Println("服务器已关闭")
}

func RegisterSystemExitFunc(f func()) {
	systemExitFunc = append(systemExitFunc, f)
}

func RegisterSystemReloadFunc(f func()) {
	systemReloadFunc = append(systemReloadFunc)
}

func systemExit() {
	DebugLog("调用退出时方法")
	for _, f := range systemExitFunc {
		f()
	}
	DebugLog("更新停止信号")
	// 更新停止信号
	if !atomic.CompareAndSwapInt32(&allForStopSignal, 0, 1) {
		return
	}
	close(stopChanForGo)
	DebugLog("关闭所有连接")
	// 关闭所有tcp连接
	tcpClientMap.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		client.SetStop()
		client.GetConn().Close()
		tcpClientMap.Delete(key)

		return true
	})

	// 关闭所有udp连接
	udpClientMap.Range(func(key, value interface{}) bool {
		udpClientMap.Delete(key)
		return true
	})

	// 关闭所有webSocket连接
	wsClientMap.Range(func(key, value interface{}) bool {
		client := value.(*WebSocketClient)
		client.SetStop()
		client.GetConn().Close()
		wsClientMap.Delete(key)
		return true
	})
	DebugLog("系统退出方法调用完成")
}

func systemReload() {
	for _, f := range systemReloadFunc {
		f()
	}
}

func Daemon(skip ...string) {
	if os.Getppid() != 1 {
		filePath, _ := filepath.Abs(os.Args[0])
		newCmd := []string{os.Args[0]}
		add := 0
		for _, v := range os.Args[1:] {
			if add == 1 {
				add = 0
				continue
			} else {
				add = 0
			}
			for _, s := range skip {
				if strings.Contains(v, s) {
					if strings.Contains(v, "--") {
						add = 2
					} else {
						add = 1
					}
					break
				}
			}
			if add == 0 {
				newCmd = append(newCmd, v)
			}
		}
		DebugLog("后台运行参数:%v", newCmd)
		cmd := exec.Command(filePath)
		cmd.Args = newCmd
		cmd.Start()
	}
}
