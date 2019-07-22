package cache

import (
	"sync"
	"time"
)

const (
	neverExpiry time.Duration = -1
)

var (
	mc = &memoryCache{
		lock:  new(sync.RWMutex),
		items: make(map[interface{}]*cacheItem),
	}
)

//缓存项
type cacheItem struct {
	data   interface{}
	expiry time.Time
}

func (ct *cacheItem) isExpired() bool {
	if ct.expiry.IsZero() {
		return false
	}

	return ct.expiry.Before(time.Now())
}

func (ct *cacheItem) setExpiry(duration time.Duration) {
	if neverExpiry == duration {
		var t time.Time
		ct.expiry = t
	} else {
		ct.expiry = time.Now().Add(duration)
	}
}

type memoryCache struct {
	lock  *sync.RWMutex
	items map[interface{}]*cacheItem
}

//设置缓存
func (mc *memoryCache) set(key interface{}, value interface{}, duration time.Duration) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	item := new(cacheItem)
	item.data = value
	item.setExpiry(duration)

	mc.items[key] = item
}

//读取缓存
func (mc *memoryCache) get(key interface{}) interface{} {
	mc.lock.RLock()
	defer mc.lock.RUnlock()

	return mc.getWithoutLock(key, false)
}

//删除缓存并返回值
func (mc *memoryCache) delete(key interface{}) interface{} {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	return mc.getWithoutLock(key, true)
}

//读取缓存（不加锁，请慎用）
func (mc *memoryCache) getWithoutLock(key interface{}, del bool) interface{} {
	if item, ok := mc.items[key]; ok {
		if del {
			delete(mc.items, key)
		}

		if !item.isExpired() {
			return item
		}
	}

	return nil
}

//设置缓存
func SetString(key interface{}, value string, duration ...time.Duration) {
	SetAnything(key, value, duration...)
}

//设置缓存
func SetInt(key interface{}, value int, duration ...time.Duration) {
	SetAnything(key, value, duration...)
}

//设置缓存
func SetBool(key interface{}, value bool, duration ...time.Duration) {
	SetAnything(key, value, duration...)
}

//设置缓存
func SetFloat(key interface{}, value float32, duration ...time.Duration) {
	SetAnything(key, value, duration...)
}

//设置缓存
func SetAnything(key interface{}, value interface{}, duration ...time.Duration) {
	dur := neverExpiry

	if duration != nil && len(duration) > 0 {
		for _, d := range duration {
			if d > 0 {
				dur = d
				break
			}
		}
	}

	mc.set(key, value, dur)
}

//读取缓存
func GetString(key interface{}, deleteOnGet ...bool) string {
	val := GetAnything(key, deleteOnGet...)
	if val != nil {
		return val.(string)
	} else {
		return ""
	}
}

//读取缓存
func GetInt(key interface{}, deleteOnGet ...bool) int {
	val := GetAnything(key, deleteOnGet...)
	if val != nil {
		return val.(int)
	} else {
		return 0
	}
}

//读取缓存
func GetBool(key interface{}, deleteOnGet ...bool) bool {
	val := GetAnything(key, deleteOnGet...)
	if val != nil {
		return val.(bool)
	} else {
		return false
	}
}

//读取缓存
func GetFloat(key interface{}, deleteOnGet ...bool) float32 {
	val := GetAnything(key, deleteOnGet...)
	if val != nil {
		return val.(float32)
	} else {
		return 0.0
	}
}

//读取缓存
func GetAnything(key interface{}, deleteOnGet ...bool) interface{} {
	del := false

	if deleteOnGet != nil && len(deleteOnGet) > 0 {
		for _, del = range deleteOnGet {
			if del {
				break
			}
		}
	}

	if del {
		return mc.delete(key)
	} else {
		return mc.get(key)
	}
}
