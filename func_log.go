package starGo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type logObj struct {
	lv      logLv
	logInfo string
}

func logStart() {
	goForLog(func(Stop chan struct{}) {
		for logForStopSignal == 0 {
			select {
			case <-Stop:
				return
			case logItem := <-logCh:
				writeLog(logItem)
			}
		}
	})
}

func writeLog(log *logObj) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// 日志级别不存在
	_, logLvExists := logLvNameMap[log.lv]
	if !logLvExists {
		return
	}

	// 日志文件未开启
	file, fileExists := logFileMap[log.lv]
	if !fileExists || file == nil {
		return
	}

	// 记录日志
	file.WriteString(log.logInfo)
	fmt.Printf("%v", log.logInfo)
}

func reorganizeLog(nowTime time.Time) {
	logMutex.Lock()
	defer logMutex.Unlock()

	for _, file := range logFileMap {
		// 优先获取文件名
		fileName := file.Name()

		// 将文件流关闭
		file.Close()
		file = nil

		// 重命名文件
		nameList := strings.Split(fileName, ".")
		newFileName := fmt.Sprintf("%v_%v.log", nameList[0], ToDateTimeString(nowTime.Add(-1*time.Hour)))
		err := os.Rename(fileName, newFileName)
		if err != nil {
			ErrorLog("重命名文件:%v失败，错误信息:%v", fileName, err)
		}
	}

	// 重新开启日志文件流
	for logLv, logName := range logLvNameMap {
		if logLv < logNowLv {
			continue
		}
		// 得到最终的文件绝对路径
		fileName := fmt.Sprintf("%v.log", logName)
		fileAbsolutePath := filepath.Join(logDirPath, fileName)

		// 打开文件(如果文件存在就以写模式打开，并追加写入；如果文件不存在就创建，然后以写模式打开。)
		f, err := os.OpenFile(fileAbsolutePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm|os.ModeTemporary)
		if err != nil {
			fmt.Printf("打开%v日志文件错误，错误信息:%v", fileName, err)
		}

		// 将文件流保存
		logFileMap[logLv] = f
	}
}

func StartLog(dirPatch string, lv logLv) {
	logOnce.Do(func() {
		logNowLv = lv
		logDirPath = dirPatch
		// 文件夹路径不存在就创建
		if !IsDirExists(dirPatch) {
			if err := os.MkdirAll(dirPatch, os.ModePerm|os.ModeTemporary); err != nil {
				fmt.Printf("创建日志文件夹错误，错误信息:%v", err)
			}
		}

		for logLv, logName := range logLvNameMap {
			// 低于当前等级的日志不打开文件
			if logLv < lv {
				continue
			}

			// 得到最终的文件绝对路径
			fileName := fmt.Sprintf("%v.log", logName)
			fileAbsolutePath := filepath.Join(dirPatch, fileName)

			// 打开文件(如果文件存在就以写模式打开，并追加写入；如果文件不存在就创建，然后以写模式打开。)
			f, err := os.OpenFile(fileAbsolutePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm|os.ModeTemporary)
			if err != nil {
				fmt.Printf("打开%v日志文件错误，错误信息:%v", fileName, err)
			}

			// 将文件流保存
			logFileMap[logLv] = f
		}
	})
}

func Stack() {
	buf := make([]byte, 1<<12)
	log(Error, string(buf[:runtime.Stack(buf, false)]))
}

func SimpleTack() string {
	_, file, line, _ := runtime.Caller(2)
	i := strings.LastIndex(file, "/") + 1
	i = strings.LastIndex((string)(([]byte(file))[:i-1]), "/") + 1

	return fmt.Sprintf("%s:%d", (string)(([]byte(file))[i:]), line)
}

func log(lv logLv, v ...interface{}) {
	if lv < logNowLv {
		return
	}

	// 判断日志文件是否存在
	lvName, exists := logLvNameMap[lv]
	if !exists {
		return
	}

	// 记录日志
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}

	i := strings.LastIndex(file, "/") + 1
	logContent := fmt.Sprintf("[%s][%s][%s:%d]:", lvName, time.Now().Format("2006-01-02 15:04:05"), (string)(([]byte(file))[i:]), line)
	if len(v) > 1 {
		logContent += fmt.Sprintf(v[0].(string), v[1:]...)
	} else {
		logContent += fmt.Sprint(v[0])
	}
	logContent += GetNewLineString()

	logCh <- &logObj{
		lv:      lv,
		logInfo: logContent,
	}
}

func DebugLog(v ...interface{}) {
	debugLog(v...)
}

func InfoLog(v ...interface{}) {
	infoLog(v...)
}

func WarnLog(v ...interface{}) {
	warnLog(v...)
}

func ErrorLog(v ...interface{}) {
	errorLog(v...)
}

func FatalLog(v ...interface{}) {
	fatalLog(v...)
}

func debugLog(v ...interface{}) {
	log(Debug, v...)
}

func infoLog(v ...interface{}) {
	log(Info, v...)
}

func warnLog(v ...interface{}) {
	log(Warn, v...)
}

func errorLog(v ...interface{}) {
	log(Error, v...)
}

func fatalLog(v ...interface{}) {
	log(Fatal, v...)
}
