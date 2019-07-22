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

//获取类型
func (h *ObjectHelper) GetType(v interface{}) reflect.Type {
	return h.GetTypeEx(v, false)
}

//获取类型
func (h *ObjectHelper) GetTypeEx(v interface{}, getUnderlyingType bool) reflect.Type {
	if nil == v {
		return nil
	}

	t := reflect.TypeOf(v)

	if getUnderlyingType {
		return h.GetUnderlyingType(t)
	} else {
		return t
	}
}

//获取底层类型
func (h *ObjectHelper) GetUnderlyingType(t reflect.Type) reflect.Type {
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

//获取类型
func (h *ObjectHelper) GetKind(v interface{}) reflect.Kind {
	if nil == v {
		return reflect.Invalid
	}

	return h.GetType(v).Kind()
}

//获取值
func (h *ObjectHelper) GetValue(v interface{}) reflect.Value {
	if nil == v {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

//根据类型创建新的对象（返回值为新对象的地址）
func (h *ObjectHelper) New(t reflect.Type) interface{} {
	if nil == t {
		return nil
	}

	return reflect.New(t).Interface()
}

//根据类型的底层类型创建新的对象（返回值为新对象的地址）
func (h *ObjectHelper) NewEx(t reflect.Type, byUnderlyingType bool) interface{} {
	if byUnderlyingType {
		return h.New(CcObject.GetUnderlyingType(t))
	} else {
		return h.New(t)
	}
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

//Object类
type Object struct {
	object interface{}
}

//创建一个新的Object实例
func NewObject(v interface{}) *Object {
	if nil == v {
		return nil
	}

	return &Object{
		object: v,
	}
}

//获取类型
func (o *Object) GetType() reflect.Type {
	return CcObject.GetType(o.object)
}

func (o *Object) GetUnderlyingType() reflect.Type {
	return CcObject.GetTypeEx(o.object, true)
}

//获取类型
func (o *Object) GetKind() reflect.Kind {
	return CcObject.GetKind(o.object)
}

//获取值
func (o *Object) GetValue() reflect.Value {
	return CcObject.GetValue(o.object)
}

//获取接口
func (o *Object) Interface() interface{} {
	return o.object
}

//根据类型创建新的对象（返回值为新对象的地址）
func (h *ObjectHelper) NewObject(t reflect.Type) *Object {
	return NewObject(h.New(t))
}

//根据类型的底层类型创建新的对象（返回值为新对象的地址）
func (h *ObjectHelper) NewObjectEx(t reflect.Type, byUnderlyingType bool) *Object {
	return NewObject(h.NewEx(t, byUnderlyingType))
}
