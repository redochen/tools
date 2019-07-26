package concurrent

import (
	"sync"
)

/**
* 线程安全的字典类
 */
type ConcurrentMap struct {
	m     *sync.RWMutex
	items map[interface{}]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		m:     new(sync.RWMutex),
		items: make(map[interface{}]interface{}),
	}
}

/**
* 添加项到字典
 */
func (cm *ConcurrentMap) Add(key, value interface{}) {
	cm.m.Lock()
	defer cm.m.Unlock()

	cm.items[key] = value
}

/**
* 从字典中获取键对应的值
 */
func (cm *ConcurrentMap) Get(key interface{}) interface{} {
	cm.m.RLock()
	defer cm.m.RUnlock()

	return cm.getWithoutLock(key, false)
}

/**
* 从字典中删除键并返回值
 */
func (cm *ConcurrentMap) Remove(key interface{}) interface{} {
	cm.m.Lock()
	defer cm.m.Unlock()

	return cm.getWithoutLock(key, true)
}

/**
* 清空字典
 */
func (cm *ConcurrentMap) Clear() {
	cm.m.Lock()
	defer cm.m.Unlock()

	cm.items = make(map[interface{}]interface{})
}

/**
* 测试字典是否包含键
 */
func (cm *ConcurrentMap) ContainsKey(key interface{}) bool {
	cm.m.RLock()
	defer cm.m.RUnlock()

	_, ok := cm.items[key]
	return ok
}

/**
* 获取字典的项数量
 */
func (cm *ConcurrentMap) Len() int {
	cm.m.RLock()
	defer cm.m.RUnlock()

	return len(cm.items)
}

/**
* 获取字典中的的所有键集合
 */
func (cm *ConcurrentMap) GetKeys() interface{} {
	cm.m.RLock()
	defer cm.m.RUnlock()

	keys := make([]interface{}, 0)

	for _, v := range cm.items {
		keys = append(keys, v)
	}

	return keys
}

/**
* 获取字典中的的所有值集合
 */
func (cm *ConcurrentMap) GetValues() []interface{} {
	cm.m.RLock()
	defer cm.m.RUnlock()

	values := make([]interface{}, 0)

	for _, v := range cm.items {
		values = append(values, v)
	}

	return values
}

/**
* 读取键值对（不加锁，请慎用）
 */
func (cm *ConcurrentMap) getWithoutLock(key interface{}, del bool) interface{} {
	if item, ok := cm.items[key]; ok {
		if del {
			delete(cm.items, key)
		}

		return item
	}

	return nil
}

/**
* 获取键的值
 */
func (cm *ConcurrentMap) GetValue(key interface{}, deleteOnGet bool) interface{} {
	if deleteOnGet {
		return cm.Remove(key)
	} else {
		return cm.Get(key)
	}
}

/**
* 获取字符串值
 */
func (cm *ConcurrentMap) GetString(key interface{}, deleteOnGet bool) string {
	val := cm.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(string)
	} else {
		return ""
	}
}

/**
* 获取整数值
 */
func (cm *ConcurrentMap) GetInt(key interface{}, deleteOnGet bool) int {
	val := cm.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(int)
	} else {
		return 0
	}
}

/**
* 获取布尔值
 */
func (cm *ConcurrentMap) GetBool(key interface{}, deleteOnGet bool) bool {
	val := cm.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(bool)
	} else {
		return false
	}
}

/**
* 获取小数值
 */
func (cm *ConcurrentMap) GetFloat(key interface{}, deleteOnGet bool) float32 {
	val := cm.GetValue(key, deleteOnGet)
	if val != nil {
		return val.(float32)
	} else {
		return 0.0
	}
}
