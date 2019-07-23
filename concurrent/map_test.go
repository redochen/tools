package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestConcurrentMap(t *testing.T) {
	cm := NewConcurrentMap()

	cm.Add(111, "AAA")
	cm.Add("aaa", 888)
	cm.Add(true, true)
	cm.Add(time.Now(), nil)

	t.Log(fmt.Sprintf("Len:%d", cm.Len()))

	if !cm.ContainsKey("aaa") {
		t.Error("TestConcurrentMap Add failed")
	}

	if cm.GetInt("aaa", true) != 888 {
		t.Error("TestConcurrentMap GetInt failed")
	}

	if !cm.GetBool(true, false) {
		t.Error("TestConcurrentMap GetBool failed")
	}

	cm.Remove(111)

	if cm.ContainsKey(111) {
		t.Error("TestConcurrentMap Remove failed")
	}

	cm.Clear()

	if cm.ContainsKey(123) {
		t.Error("TestConcurrentMap Clear failed")
	}

	t.Log("TestConcurrentMap success")
}
