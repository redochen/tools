package concurrent

import (
	"bytes"
	"fmt"
	"sync"
)

//ConcurrentSet 线程安全的集合类
type ConcurrentSet struct {
	lock  *sync.RWMutex
	items map[interface{}]bool
}

func NewConcurrentSet() *ConcurrentSet {
	return &ConcurrentSet{
		lock:  new(sync.RWMutex),
		items: make(map[interface{}]bool),
	}
}

//Add 添加项到集合
func (s *ConcurrentSet) Add(item interface{}) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.items[item] {
		s.items[item] = true
		return true
	}

	return false
}

//Remove 从集合中移除项
func (s *ConcurrentSet) Remove(item interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.items, item)
}

//Clear 清空集合
func (s *ConcurrentSet) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = make(map[interface{}]bool)
}

//Contains 测试集合是否包含项
func (s *ConcurrentSet) Contains(e interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.items[e]
}

//获取集合的项数量
func (s *ConcurrentSet) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.items)
}

//String 转换成字符串
func (s *ConcurrentSet) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var buf bytes.Buffer
	buf.WriteString("ConcurrentSet{")

	first := true
	for k := range s.items {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}

		buf.WriteString(fmt.Sprintf("%v", k))
	}

	buf.WriteString("}")

	return buf.String()
}
