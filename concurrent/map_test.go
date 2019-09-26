package concurrent

import (
	"testing"
	"time"
)

func TestConcurrentMap(t *testing.T) {
	m := NewConcurrentMap()

	m.Add(111, "AAA")
	m.Add("aaa", 888)
	m.Add(true, true)
	m.Add(time.Now(), nil)

	t.Logf("Len:%d", m.Len())

	if !m.ContainsKey("aaa") {
		t.Error("TestConcurrentMap Add failed")
	}

	if m.GetInt("aaa", true) != 888 {
		t.Error("TestConcurrentMap GetInt failed")
	}

	if !m.GetBool(true, false) {
		t.Error("TestConcurrentMap GetBool failed")
	}

	m.Remove(111)

	if m.ContainsKey(111) {
		t.Error("TestConcurrentMap Remove failed")
	}

	m.Clear()

	if m.ContainsKey(123) {
		t.Error("TestConcurrentMap Clear failed")
	}

	t.Log("TestConcurrentMap success")
}
