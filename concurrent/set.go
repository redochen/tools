package concurrent

import (
	"bytes"
	"fmt"
	"sync"
)

//线程安全的集合类
type ConcurrentSet struct {
	m     *sync.RWMutex
	items map[interface{}]bool
}

func NewConcurrentSet() *ConcurrentSet {
	return &ConcurrentSet{
		m:     new(sync.RWMutex),
		items: make(map[interface{}]bool),
	}
}

//添加项到集合
func (cs *ConcurrentSet) Add(item interface{}) bool {
	cs.m.Lock()
	defer cs.m.Unlock()

	if !cs.items[item] {
		cs.items[item] = true
		return true
	}

	return false
}

//从集合中移除项
func (cs *ConcurrentSet) Remove(item interface{}) {
	cs.m.Lock()
	defer cs.m.Unlock()

	delete(cs.items, item)
}

//清空集合
func (cs *ConcurrentSet) Clear() {
	cs.m.Lock()
	defer cs.m.Unlock()

	cs.items = make(map[interface{}]bool)
}

//测试集合是否包含项
func (cs *ConcurrentSet) Contains(e interface{}) bool {
	cs.m.RLock()
	defer cs.m.RUnlock()

	return cs.items[e]
}

//获取集合的项数量
func (cs *ConcurrentSet) Len() int {
	cs.m.RLock()
	defer cs.m.RUnlock()

	return len(cs.items)
}

func (cs *ConcurrentSet) String() string {
	cs.m.RLock()
	defer cs.m.RUnlock()

	var buf bytes.Buffer
	buf.WriteString("ConcurrentSet{")

	first := true
	for k := range cs.items {
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
