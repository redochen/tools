package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestConcurrentSet(t *testing.T) {
	cs := NewConcurrentSet()

	cs.Add(123)
	cs.Add("abc")
	cs.Add(nil)
	cs.Add(time.Now())

	t.Log(fmt.Sprintf("Len:%d,String:%s", cs.Len(), cs.String()))

	if !cs.Contains("abc") {
		t.Error("TestConcurrentSet Add failed")
	}

	cs.Remove("abc")

	if cs.Contains("abc") {
		t.Error("TestConcurrentSet Remove failed")
	}

	cs.Clear()

	if cs.Contains(123) {
		t.Error("TestConcurrentSet Clear failed")
	}

	t.Log("TestConcurrentSet success")
}
