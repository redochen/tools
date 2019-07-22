package object

import (
	"reflect"
)

var (
	CcObject = NewObjectHelper()
)

//interface{}帮助类
type ObjectHelper struct {
}

//获取一个新的ObjectHelper实例
func NewObjectHelper() *ObjectHelper {
	return &ObjectHelper{}
}

//深拷贝对象
func (h *ObjectHelper) DeepCopy(src, dst interface{}) {
	sVal := reflect.ValueOf(src).Elem()
	dVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < sVal.NumField(); i++ {
		value := sVal.Field(i)
		name := sVal.Type().Field(i).Name

		dValue := dVal.FieldByName(name)
		if dValue.IsValid() == false {
			continue
		}
		dValue.Set(value) //这里默认共同成员的类型一样，否则这个地方可能导致 panic，需要简单修改一下。
	}
}

//获取集合的长度
func (h *ObjectHelper) GetLengthOfCollection(item interface{}) (length int) {
	value := reflect.ValueOf(item)
	if value.IsNil() {
		length = 0
	}

	kind := value.Kind()
	if reflect.Array == kind ||
		reflect.Slice == kind ||
		reflect.Map == kind {
		length = value.Len()
	}

	return
}
