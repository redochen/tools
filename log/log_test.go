package log

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("TestLog failed:%v", err)
		}
	}()

	Debug("test debug log")
	Debugf("test formatted debug log: %v", time.Now())

	Info("test info log")
	Infof("test formatted info log: %v", time.Now())

	Warn("test warning log")
	Warnf("test formatted warning log: %v", time.Now())

	Error("test error log")
	Errorf("test formatted error log: %v", time.Now())

	Fatal("test fatal log")
	Fatalf("test formatted fatal log: %v", time.Now())

	t.Log("TestLog success")
}
