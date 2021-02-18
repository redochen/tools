package log

import (
	"fmt"
	"runtime"
	"runtime/debug"

	log "github.com/redochen/log4go"
)

func init() {
	log.LoadConfiguration("./log.json")
}

//Debug 输出DEBUG日志
func Debug(msg ...string) {
	log.Log(log.DEBUG, getSource(), getMessage(msg...))
	log.Flush()
}

//Debugf 输出DEBUG日志
func Debugf(format string, args ...interface{}) {
	log.Log(log.DEBUG, getSource(), fmt.Sprintf(format, args...))
	log.Flush()
}

//Info 输出INFORMATION日志
func Info(msg ...string) {
	log.Log(log.INFO, getSource(), getMessage(msg...))
	log.Flush()
}

//Infof 输出INFORMATION日志
func Infof(format string, args ...interface{}) {
	log.Log(log.INFO, getSource(), fmt.Sprintf(format, args...))
	log.Flush()
}

//Warn 输出WARNING日志
func Warn(msg ...string) {
	log.Log(log.WARNING, getSource(), getMessage(msg...))
	log.Flush()
}

//Warnf 输出WARNING日志
func Warnf(format string, args ...interface{}) {
	log.Log(log.WARNING, getSource(), fmt.Sprintf(format, args...))
	log.Flush()
}

//Error 输出ERROR日志
func Error(msg ...string) {
	log.Log(log.ERROR, getSource(), getMessage(msg...))
	log.Flush()
}

//Errorf 输出ERROR日志
func Errorf(format string, args ...interface{}) {
	log.Log(log.ERROR, getSource(), fmt.Sprintf(format, args...))
	log.Flush()
}

//Fatal 输出FATAL日志
func Fatal(msg ...string) {
	log.Log(log.FATAL, getSource(), getMessage(msg...))
	log.Flush()
}

//Fatalf 输出FATAL日志
func Fatalf(format string, args ...interface{}) {
	log.Log(log.FATAL, getSource(), fmt.Sprintf(format, args...))
	log.Flush()
}

//FatalEx 输出FATAL日志
func FatalEx(err error) {
	log.Log(log.FATAL, getSource(), fmt.Sprintf("%v; stack: %s", err, debug.Stack()))
	log.Flush()
}

//getMessage 获取消息
func getMessage(msg ...string) string {
	if nil == msg || len(msg) <= 0 {
		return ""
	}

	var message string
	for _, m := range msg {
		message += m
	}

	return message
}

//getSource returns caller's info
func getSource() string {
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	return src
}
