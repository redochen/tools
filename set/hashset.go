package set

import (
	"bytes"
	"fmt"
)

type Set interface {
	Add(e interface{}) bool
	Remove(e interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	String() string
}
type HashSet struct {
	m map[interface{}]bool
}

func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}
func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}
func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}
func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}
func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}
func (set *HashSet) Len() int {
	return len(set.m)
}
func (set *HashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("HashSet{")
	first := true
	for k := range set.m {
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
