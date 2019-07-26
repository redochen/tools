package routine

import (
	. "github.com/redochen/tools/log"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var (
	CcFunctions = NewFuncContextMap()
)

/**
* 获取函数名称
 */
func GetFunctionName(i interface{}, seps ...rune) string {
	return getFuncName(reflect.ValueOf(i).Pointer(), seps...)
}

/**
* 获取调用者名称
 */
func GetCallerFuncName(seps ...rune) string {
	pc, _, _, _ := runtime.Caller(3)
	return getFuncName(pc, seps...)
}

/**
* 获取当前方法名称
 */
func GetCurrentFuncName(seps ...rune) string {
	pc, _, _, _ := runtime.Caller(2)
	return getFuncName(pc, seps...)
}

/**
* 获取函数名称
 */
func getFuncName(pc uintptr, seps ...rune) string {
	fn := runtime.FuncForPC(pc).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}

	return ""
}

/**
* 检查错误
 */
func CheckPanic() {
	name := GetCallerFuncName()
	if err := recover(); err != nil {
		Logger.FatalEx(name, err)
	}
}

/**
* 调用函数：fn要调用的方法；args方法参数
 */
func InvokeFunc(fn interface{}, args ...interface{}) interface{} {
	defer CheckPanic()

	if len(args) > 1 {
		f, ok := fn.(func(...interface{}) interface{})
		if ok {
			return f(args)
		}

		fn.(func(...interface{}))(args)
		return nil
	}

	if len(args) == 1 {
		f, ok := fn.(func(interface{}) interface{})
		if ok {
			return f(args[0])
		}

		fn.(func(interface{}))(args[0])
		return nil
	}

	f, ok := fn.(func() interface{})
	if ok {
		return f()
	}

	fn.(func())()

	return nil
}

/**
*定时调用函数：interval时间间隔；runRightNow是否立即执行一次；name方法名称；fn要调用的方法；args方法参数
 */
func InvokeTicker(interval time.Duration, runRightNow bool, name string, fn interface{}, args ...interface{}) {
	stop := make(chan bool, 1)
	deadline := time.Now().AddDate(100, 0, 0)
	InvokeTickerEx(interval, deadline, stop, runRightNow, name, fn, args...)
}

/**
*定时调用函数：interval时间间隔；deadline过期时间；stop停止执行控制通道；runRightNow是否立即执行一次；name方法名称；fn要调用的方法；args方法参数
 */
func InvokeTickerEx(interval time.Duration, deadline time.Time, stop <-chan bool,
	runRightNow bool, name string, fn interface{}, args ...interface{}) {
	defer CheckPanic()

	if time.Now().After(deadline) {
		Logger.ErrorEx("InvokeTickerEx", "%s invalid deadline: %v", name, deadline)
		return
	}

	if runRightNow {
		InvokeFunc(fn, args...)
	}

	t := time.NewTicker(interval)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			if time.Now().After(deadline) {
				Logger.DebugEx("InvokeTickerEx", "%s was dead", name)
				break Loop
			} else {
				InvokeFunc(fn, args...)
			}
		case <-stop:
			Logger.DebugEx("InvokeTickerEx", "%s was stopped", name)
			break Loop
		}
	}
}

/**
* 启用routine调用函数：result执行结果；name方法名称；fn要调用的方法；args方法参数
 */
func (rm *FuncContextMap) goFunc(result chan<- interface{}, name string, fn interface{}, args ...interface{}) *FuncContext {
	defer CheckPanic()

	c, e := rm.register(name)
	if e != nil {
		return nil
	}

	go func(r chan<- interface{}, n string, f interface{}, a ...interface{}) {
		defer rm.unregister(n)
		defer CheckPanic()

		result <- InvokeFunc(f, a...)
	}(result, name, fn, args...)

	return c
}

//启用routine调用Ticker：interval时间间隔；deadline过期时间；stop停止执行控制通道；runRightNow是否立即执行一次；name方法名称；fn要调用的方法；args方法参数
func (rm *FuncContextMap) startTicker(interval time.Duration, deadline time.Time, stop <-chan bool,
	runRightNow bool, name string, fn interface{}, args ...interface{}) *FuncContext {
	defer CheckPanic()

	c, e := rm.register(name)
	if e != nil {
		return nil
	}

	defer rm.unregister(name)

	go func(dur time.Duration, dead time.Time, stp <-chan bool, run bool, n string, f interface{}, a ...interface{}) {
		defer CheckPanic()

		InvokeTickerEx(dur, dead, stp, run, n, f, a...)
	}(interval, deadline, stop, runRightNow, name, fn, args...)

	return c
}

/**
* 调用函数：name方法名称；fn要调用的方法；args方法参数
 */
func (rm *FuncContextMap) InvokeFunc(name string, fn interface{}, args ...interface{}) interface{} {
	return rm.InvokeFuncWithTimeout(999999*time.Hour, name, fn, args...)
}

/**
* 调用函数：timeout超时时间；name方法名称；fn要调用的方法；args方法参数
 */
func (rm *FuncContextMap) InvokeFuncWithTimeout(timeout time.Duration, name string, fn interface{}, args ...interface{}) interface{} {
	defer CheckPanic()

	if timeout < 0 {
		Logger.ErrorEx("InvokeFuncWithDeadline", "%s invalid timeout: %d", name, timeout)
		return nil
	}

	ch := make(chan interface{}, 1)
	r := rm.goFunc(ch, name, fn, args...)

	var result interface{}
	result = nil

Loop:
	for {
		select {
		case result = <-ch:
			Logger.DebugEx("InvokeFuncWithDeadline", "%s ran to done", name)
			break Loop
		case <-r.stop:
			Logger.DebugEx("InvokeFuncWithDeadline", "%s has been stopped", name)
			break Loop
		case <-time.After(timeout):
			Logger.DebugEx("InvokeFuncWithDeadline", "%s has timed out", name)
			break Loop
		}
	}

	return result
}

/**
*定时执行任务：interval时间间隔；runRightNow是否执行一次；name方法名称；fn要调用的方法；args方法参数
 */
func (rm *FuncContextMap) StartTicker(interval time.Duration, runRightNow bool, name string, fn interface{}, args ...interface{}) {
	deadline := time.Now().AddDate(100, 0, 0)
	rm.StartTickerWithDeadline(interval, deadline, runRightNow, name, fn, args...)
}

/**
*定时执行任务：interval时间间隔；deadline过期时间；runRightNow是否执行一次；name方法名称；fn要调用的方法；args方法参数
 */
func (rm *FuncContextMap) StartTickerWithDeadline(interval time.Duration, deadline time.Time,
	runRightNow bool, name string, fn interface{}, args ...interface{}) {
	defer CheckPanic()

	if time.Now().After(deadline) {
		Logger.ErrorEx("StartTickerWithDeadline", "%s invalid deadline: %v", name, deadline)
		return
	}

	s := make(chan bool, 1)
	r := rm.startTicker(interval, deadline, s, runRightNow, name, fn, args...)

Loop:
	for {
		select {
		case <-r.stop:
			s <- true
			Logger.DebugEx("StartTickerWithDeadline", "%s has been stopped", name)
			break Loop
		}
	}
}

/**
* 停止routine：name方法名称
 */
func (rm *FuncContextMap) StopFunc(name string) {
	defer CheckPanic()

	r, ok := rm.funcs[name]
	if !ok {
		Logger.ErrorEx("StopFunc", "routine %s not found", name)
		return
	}

	r.stop <- true
}
