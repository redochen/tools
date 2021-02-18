package set

import (
	"bytes"
	"fmt"
)

//HashSet 哈希集合类
type HashSet struct {
	m map[interface{}]bool
}

//NewHashSet get new instance of HashSet
func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}

//Add add one element to hashset
func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

//Remove remove one element form hashset
func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}

//Clear clear hashset
func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

//Contains check if hashset contains the element
func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

//Len get the length of hashset
func (set *HashSet) Len() int {
	return len(set.m)
}

//String format hashset to string
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
