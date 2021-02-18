package set

//Set 接口
type Set interface {
	Add(e interface{}) bool
	Remove(e interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	String() string
}
