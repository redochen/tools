package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
)

const (
	LogSection string = "log"
)

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

func getLevelName(level int) string {
	switch level {
	case LevelError:
		return "error"
	case LevelAlert:
		return "alert"
	case LevelCritical:
		return "critical"
	case LevelWarn:
		return "warn"
	case LevelNotice:
		return "notice"
	case LevelInformational:
		return "info"
	case LevelDebug:
		return "debug"
	default:
		return "no"
	}

}

var (
	machineName  string
	Logger       *BeeLogger
	isUseMongodb bool //是否使用mongodb做为日志库
	mongodbUrl   string
	category     int    //日志类别
	businessNo   int    //业务编号
	tableName    string //mongodb日志表名称
)

func init() {
	registerFile()
	registerMongo()
	Logger = NewLogger(100000)
}

// Legacy loglevel constants to ensure backwards compatibility.
//
// Deprecated: will be removed in 1.5.0.
const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

type loggerType func() LoggerInterface

// LoggerInterface defines the behavior of a log provider.
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]loggerType)

// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, log loggerType) {
	if log == nil {
		fmt.Println("logs: Register provide is nil")
		return
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

// BeeLogger is default logger in beego application.
// it can contain several providers and log message into all providers.
type BeeLogger struct {
	lock                sync.Mutex
	level               int
	enableFuncCallDepth bool
	loggerFuncCallDepth int
	msg                 chan *logMsg
	outputs             map[string]LoggerInterface
}

type logMsg struct {
	level int
	msg   string
}

// NewLogger returns a new BeeLogger.
// channellen means the number of messages in chan.
// if the buffering chan is full, logger adapters write to file or other way.
func NewLogger(channellen int64) *BeeLogger {
	bl := new(BeeLogger)
	bl.level = LevelDebug
	bl.loggerFuncCallDepth = 2
	bl.msg = make(chan *logMsg, channellen)
	bl.outputs = make(map[string]LoggerInterface)
	var err error
	machineName, err = os.Hostname()
	if err != nil {
		panic(err)
	}
	bl.SetLogger("file", `{"level":7}`)
	bl.SetLogger("mongo", "4")
	go bl.startLogger()
	return bl
}

// SetLogger provides a given logger adapter into BeeLogger with config string.
// config need to be correct JSON as string: {"interval":360}.
func (bl *BeeLogger) SetLogger(adaptername string, config string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if log, ok := adapters[adaptername]; ok {
		lg := log()
		err := lg.Init(config)
		bl.outputs[adaptername] = lg
		if err != nil {
			fmt.Println("logs.BeeLogger.SetLogger: " + err.Error())
			return err
		}
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
	return nil
}

// remove a logger adapter in BeeLogger.
func (bl *BeeLogger) DelLogger(adaptername string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if lg, ok := bl.outputs[adaptername]; ok {
		lg.Destroy()
		delete(bl.outputs, adaptername)
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

func (bl *BeeLogger) writerMsg(loglevel int, msg string) error {
	if loglevel > bl.level {
		return nil
	}
	lm := new(logMsg)
	lm.level = loglevel
	if bl.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(bl.loggerFuncCallDepth)
		if ok {
			_, filename := path.Split(file)
			lm.msg = fmt.Sprintf("[%s:%d] %s", filename, line, msg)
		} else {
			lm.msg = msg
		}
	} else {
		lm.msg = msg
	}
	bl.msg <- lm
	return nil
}

// Set log message level.
//
// If message level (such as LevelDebug) is higher than logger level (such as LevelWarning),
// log providers will not even be sent the message.
func (bl *BeeLogger) SetLevel(l int) {
	bl.level = l
}

// set log funcCallDepth
func (bl *BeeLogger) SetLogFuncCallDepth(d int) {
	bl.loggerFuncCallDepth = d
}

// enable log funcCallDepth
func (bl *BeeLogger) EnableFuncCallDepth(b bool) {
	bl.enableFuncCallDepth = b
}

// start logger chan reading.
// when chan is not empty, write logs.
func (bl *BeeLogger) startLogger() {
	for {
		select {
		case bm := <-bl.msg:
			for _, l := range bl.outputs {
				err := l.WriteMsg(bm.msg, bm.level)
				if err != nil {
					fmt.Println("ERROR, unable to WriteMsg:", err)
				}
			}
		}
	}
}

// Log EMERGENCY level message.
func (bl *BeeLogger) Emerge(message string) {
	bl.writerMsg(LevelEmergency, message)
}

// Log EMERGENCY level message.
func (bl *BeeLogger) EmergeEx(title, format string, a ...interface{}) {
	bl.Emerge(getFmtMessage(title, format, a...))
}

// Log CRITICAL level message.
func (bl *BeeLogger) Critical(message string) {
	bl.writerMsg(LevelCritical, message)
}

// Log CRITICAL level message.
func (bl *BeeLogger) CriticalEx(title, format string, a ...interface{}) {
	bl.Critical(getFmtMessage(title, format, a...))
}

// Log ERROR level message.
func (bl *BeeLogger) Error(message string) {
	bl.writerMsg(LevelError, message)
}

// Log ERROR level message.
func (bl *BeeLogger) ErrorEx(title, format string, a ...interface{}) {
	bl.Error(getFmtMessage(title, format, a...))
}

// Log WARNING level message.
func (bl *BeeLogger) Warning(message string) {
	bl.writerMsg(LevelWarning, message)
}

// Log WARNING level message.
func (bl *BeeLogger) WarningEx(title, format string, a ...interface{}) {
	bl.Warning(getFmtMessage(title, format, a...))
}

// Log INFORMATIONAL level message.
func (bl *BeeLogger) Info(message string) {
	bl.writerMsg(LevelInformational, message)
}

// Log INFORMATIONAL level message.
func (bl *BeeLogger) InfoEx(title, format string, a ...interface{}) {
	bl.Info(getFmtMessage(title, format, a...))
}

// Log DEBUG level message.
func (bl *BeeLogger) Debug(message string) {
	bl.writerMsg(LevelDebug, message)
}

// Log DEBUG level message.
func (bl *BeeLogger) DebugEx(title, format string, a ...interface{}) {
	bl.Debug(getFmtMessage(title, format, a...))
}

func getFmtMessage(title, format string, a ...interface{}) string {
	return fmt.Sprintf("[%s]%s", title, fmt.Sprintf(format, a...))
}

// flush all chan data.
func (bl *BeeLogger) Flush() {
	for _, l := range bl.outputs {
		l.Flush()
	}
}

// close logger, flush all chan data and destroy all adapters in BeeLogger.
func (bl *BeeLogger) Close() {
	for {
		if len(bl.msg) > 0 {
			bm := <-bl.msg
			for _, l := range bl.outputs {
				err := l.WriteMsg(bm.msg, bm.level)
				if err != nil {
					fmt.Println("ERROR, unable to WriteMsg (while closing logger):", err)
				}
			}
		} else {
			break
		}
	}
	for _, l := range bl.outputs {
		l.Flush()
		l.Destroy()
	}
}
