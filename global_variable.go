package starGo

import (
	"github.com/nats-io/nats.go"
	"math"
	"math/big"
	"os"
	"sync"
)

var (
	maxBigInt64Edge  = big.NewInt(0).Add(big.NewInt(math.MaxInt64), big.NewInt(1))
	baseString       = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	allForStopSignal int32

	logDirPath   string
	logFileMap   map[logLv]*os.File
	logMutex     sync.Mutex
	logCh        chan *logObj
	logOnce      sync.Once
	logNowLv     logLv
	logLvNameMap = map[logLv]string{
		Debug: "debug",
		Info:  "info",
		Warn:  "warn",
		Error: "error",
		Fatal: "fatal",
	}

	waitAllGroup   sync.WaitGroup
	goCount        int32
	goId           int32
	stopChanForGo  = make(chan struct{})
	stopChanForLog = make(chan struct{})

	timerMutex       sync.RWMutex
	oneMinuteFunc    map[string]timerFunc
	fiveMinuteFunc   map[string]timerFunc
	thirtyMinuteFunc map[string]timerFunc

	systemExitFunc   []func()
	systemReloadFunc []func()

	tcpClientMap            sync.Map
	udpClientMap            sync.Map
	tcpReceiveDataHeaderLen int32
	udpReceiveDataHeaderLen int32
	tcpHandlerReceiveFunc   ClientCallBack
	udpHandlerReceiveFunc   ClientCallBack

	natChMap sync.Map
	natConn  *nats.Conn
)
