package starGo

import "time"

type timerFunc func(time time.Time)

func RegisterOneMinuteFunc(funcName string, f timerFunc) {
	timerMutex.Lock()
	defer timerMutex.Unlock()
	if _, exists := oneMinuteFunc[funcName]; exists {
		ErrorLog("方法名:%v的一分钟执行方法已注册")
		return
	}
	oneMinuteFunc[funcName] = f
}

func RegisterFiveMinuteFunc(funcName string, f timerFunc) {
	timerMutex.Lock()
	defer timerMutex.Unlock()
	if _, exists := fiveMinuteFunc[funcName]; exists {
		ErrorLog("方法名:%v的五分钟执行方法已注册")
		return
	}
	fiveMinuteFunc[funcName] = f
}

func RegisterThirtyMinuteFunc(funcName string, f timerFunc) {
	timerMutex.Lock()
	defer timerMutex.Unlock()
	if _, exists := thirtyMinuteFunc[funcName]; exists {
		ErrorLog("方法名:%v的三十分钟执行方法已注册")
		return
	}
	thirtyMinuteFunc[funcName] = f
}

func timerStart() {
	Go(func(Stop chan struct{}) {
		t := time.NewTicker(1 * time.Second)
		for allForStopSignal == 0 {
			select {
			case <-Stop:
				return
			case <-t.C:
				nowTime := time.Now()

				// 整分钟数开始执行
				if nowTime.Second() == 0 {
					// 执行一分钟方法
					callOneMinuteFunc(nowTime)

					if nowTime.Minute()%5 == 0 {
						// 执行五分钟方法
						callFiveMinuteFunc(nowTime)
					}

					if nowTime.Minute()%30 == 0 {
						// 执行三十分方法
						callThirtyMinuteFunc(nowTime)
					}

					if nowTime.Minute() == 0 {
						// 整理日志文件
						reorganizeLog(nowTime)
					}
				}
			}
		}
	})
}

func callOneMinuteFunc(nowTime time.Time) {
	timerMutex.RLock()
	defer timerMutex.RUnlock()
	for name, f := range oneMinuteFunc {
		startTime := time.Now().UnixNano()
		InfoLog("开始执行%v方法", name)
		Try(func() { f(nowTime) }, nil)
		useTime := float64(time.Now().UnixNano() - startTime)
		InfoLog("%v方法执行完成,总耗时%v毫秒", name, useTime/float64(1000000))
	}
}

func callFiveMinuteFunc(nowTime time.Time) {
	timerMutex.RLock()
	defer timerMutex.RUnlock()
	for name, f := range fiveMinuteFunc {
		startTime := time.Now().UnixNano()
		InfoLog("开始执行%v方法", name)
		Try(func() { f(nowTime) }, nil)
		useTime := float64(time.Now().UnixNano() - startTime)
		InfoLog("%v方法执行完成,总耗时%v毫秒", name, useTime/float64(1000000))
	}
}

func callThirtyMinuteFunc(nowTime time.Time) {
	timerMutex.RLock()
	defer timerMutex.RUnlock()
	for name, f := range thirtyMinuteFunc {
		startTime := time.Now().UnixNano()
		InfoLog("开始执行%v方法", name)
		Try(func() { f(nowTime) }, nil)
		useTime := float64(time.Now().UnixNano() - startTime)
		InfoLog("%v方法执行完成,总耗时%v毫秒", name, useTime/float64(1000000))
	}
}
