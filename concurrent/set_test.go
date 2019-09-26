package concurrent

import (
	"testing"
	"time"
)

func TestConcurrentSet(t *testing.T) {
	s := NewConcurrentSet()

	s.Add(123)
	s.Add("abc")
	s.Add(nil)
	s.Add(time.Now())

	t.Logf("Len:%d,String:%s", s.Len(), s.String())

	if !s.Contains("abc") {
		t.Error("TestConcurrentSet Add failed")
	}

	s.Remove("abc")

	if s.Contains("abc") {
		t.Error("TestConcurrentSet Remove failed")
	}

	s.Clear()

	if s.Contains(123) {
		t.Error("TestConcurrentSet Clear failed")
	}

	t.Log("TestConcurrentSet success")
}
