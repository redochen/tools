package object

import (
	"reflect"
)

//GetType 获取类型
func GetType(v interface{}) reflect.Type {
	return GetTypeEx(v, false)
}

//GetTypeEx 获取类型
func GetTypeEx(v interface{}, getUnderlyingType bool) reflect.Type {
	if nil == v {
		return nil
	}

	t := reflect.TypeOf(v)

	if getUnderlyingType {
		return GetUnderlyingType(t)
	} else {
		return t
	}
}

//GetUnderlyingType 获取底层类型
func GetUnderlyingType(t reflect.Type) reflect.Type {
	if nil == t {
		return nil
	}

	switch t.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		t = t.Elem()
	}

	return t
}

//GetKind 获取类型
func GetKind(v interface{}) reflect.Kind {
	if nil == v {
		return reflect.Invalid
	}

	return GetType(v).Kind()
}

//GetValue 获取值
func GetValue(v interface{}) reflect.Value {
	if nil == v {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

//New 根据类型创建新的对象（返回值为新对象的地址）
func New(t reflect.Type) interface{} {
	if nil == t {
		return nil
	}

	return reflect.New(t).Interface()
}

//NewEx 根据类型的底层类型创建新的对象（返回值为新对象的地址）
func NewEx(t reflect.Type, byUnderlyingType bool) interface{} {
	if byUnderlyingType {
		return New(GetUnderlyingType(t))
	} else {
		return New(t)
	}
}

//DeepCopy 深拷贝对象
func DeepCopy(src, dst interface{}) {
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

//GetLengthOfCollection 获取集合的长度
func GetLengthOfCollection(item interface{}) (length int) {
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

//Object类
type Object struct {
	object interface{}
}

//newObject 创建一个新的Object实例
func newObject(v interface{}) *Object {
	if nil == v {
		return nil
	}

	return &Object{
		object: v,
	}
}

//GetType 获取类型
func (o *Object) GetType() reflect.Type {
	return GetType(o.object)
}

//GetUnderlyingType 获取底层类型
func (o *Object) GetUnderlyingType() reflect.Type {
	return GetTypeEx(o.object, true)
}

//GetKind 获取类型
func (o *Object) GetKind() reflect.Kind {
	return GetKind(o.object)
}

//GetValue 获取值
func (o *Object) GetValue() reflect.Value {
	return GetValue(o.object)
}

//Interface 获取接口
func (o *Object) Interface() interface{} {
	return o.object
}

//NewObject 根据类型创建新的对象（返回值为新对象的地址）
func NewObject(t reflect.Type) *Object {
	return newObject(New(t))
}

//NewObjectEx 根据类型的底层类型创建新的对象（返回值为新对象的地址）
func NewObjectEx(t reflect.Type, byUnderlyingType bool) *Object {
	return newObject(NewEx(t, byUnderlyingType))
}
