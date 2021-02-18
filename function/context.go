package function

import (
	"fmt"
	"sync"
)

//FuncContext context of function
type FuncContext struct {
	name string    //func name
	stop chan bool //stop chan
}

//FuncContextMap map for FuncContext
type FuncContextMap struct {
	lock  *sync.RWMutex
	funcs map[string]*FuncContext
}

//NewFuncContextMap returns new instance of FuncContextMap
func NewFuncContextMap() *FuncContextMap {
	m := &FuncContextMap{
		lock: new(sync.RWMutex),
	}

	m.funcs = make(map[string]*FuncContext)

	return m
}

//register func context
func (m *FuncContextMap) register(name string) (*FuncContext, error) {
	r := &FuncContext{name: name}

	r.stop = make(chan bool, 1)

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.funcs[r.name]; ok {
		return nil, fmt.Errorf("func %s already exists", name)
	}

	m.funcs[r.name] = r

	return r, nil
}

//unregister func context
func (m *FuncContextMap) unregister(name string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.funcs[name]; !ok {
		return fmt.Errorf("func %s not found", name)
	}

	delete(m.funcs, name)

	return nil
}
