package starGo

type logLv int

const (
	Debug logLv = iota //调试信息
	Info               //资讯讯息
	Warn               //警告状况发生
	Error              //一般错误，可能导致功能不正常
	Fatal              //严重错误，会导致进程退出
)
