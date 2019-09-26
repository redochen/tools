package concurrent

import "testing"

func TestConcurrentQueue(t *testing.T) {
	q := NewConcurrentQueue()

	q.Push(1)
	q.Push("chen")

	if q.Length() != 2 {
		t.Error("TestConcurrentQueue Push failed")
	}

	n := q.Pop().(int)
	if n != 1 {
		t.Error("TestConcurrentQueue Pop number failed")
	}

	v, _ := q.TryPop()
	if nil == v || v.(string) != "chen" {
		t.Error("TestConcurrentQueue TryPop string failed")
	}

	if q.Length() != 0 {
		t.Error("TestConcurrentQueue Pop failed")
	}
}
