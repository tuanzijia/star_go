package starGo

import "os"

func init() {
	logCh = make(chan *logObj, 1024)
	logFileMap = make(map[logLv]*os.File)
	logStart()

	oneMinuteFunc = make(map[string]timerFunc)
	fiveMinuteFunc = make(map[string]timerFunc)
	thirtyMinuteFunc = make(map[string]timerFunc)
	timerStart()
}
