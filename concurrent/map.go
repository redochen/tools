package concurrent

import (
	"sync"
)

//ConcurrentMap 线程安全的字典类
type ConcurrentMap struct {
	lock  *sync.RWMutex
	items map[interface{}]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		lock:  new(sync.RWMutex),
		items: make(map[interface{}]interface{}),
	}
}

//Add 添加项到字典
func (m *ConcurrentMap) Add(key, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.items[key] = value
}

//Get 从字典中获取键对应的值
func (m *ConcurrentMap) Get(key interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.getWithoutLock(key, false)
}

//Remove 从字典中删除键并返回值
func (m *ConcurrentMap) Remove(key interface{}) interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.getWithoutLock(key, true)
}

//Clear 清空字典
func (m *ConcurrentMap) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.items = make(map[interface{}]interface{})
}

//ContainsKey 测试字典是否包含键
func (m *ConcurrentMap) ContainsKey(key interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	_, ok := m.items[key]
	return ok
}

// Len 获取字典的项数量
func (m *ConcurrentMap) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(m.items)
}

//GetKeys 获取字典中的的所有键集合
func (m *ConcurrentMap) GetKeys() interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	keys := make([]interface{}, 0)

	for _, v := range m.items {
		keys = append(keys, v)
	}

	return keys
}

//GetValues 获取字典中的的所有值集合
func (m *ConcurrentMap) GetValues() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	values := make([]interface{}, 0)

	for _, v := range m.items {
		values = append(values, v)
	}

	return values
}

//getWithoutLock 读取键值对（不加锁，慎用）
func (m *ConcurrentMap) getWithoutLock(key interface{}, del bool) interface{} {
	if item, ok := m.items[key]; ok {
		if del {
			delete(m.items, key)
		}

		return item
	}

	return nil
}

//GetValue 获取键的值
func (m *ConcurrentMap) GetValue(key interface{}, deleteOnGet bool) interface{} {
	if deleteOnGet {
		return m.Remove(key)
	} else {
		return m.Get(key)
	}
}

//GetString 获取字符串值
func (m *ConcurrentMap) GetString(key interface{}, deleteOnGet bool) string {
	val := m.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(string)
	} else {
		return ""
	}
}

//GetInt 获取整数值
func (m *ConcurrentMap) GetInt(key interface{}, deleteOnGet bool) int {
	val := m.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(int)
	} else {
		return 0
	}
}

//GetBool 获取布尔值
func (m *ConcurrentMap) GetBool(key interface{}, deleteOnGet bool) bool {
	val := m.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(bool)
	} else {
		return false
	}
}

//GetFloat 获取小数值
func (m *ConcurrentMap) GetFloat(key interface{}, deleteOnGet bool) float32 {
	val := m.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(float32)
	} else {
		return 0.0
	}
}
