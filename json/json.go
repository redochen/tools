package json

import (
	"encoding/json"
	"errors"
	. "github.com/redochen/tools/object"
	"reflect"
)

var (
	CcJson = NewJsonHelper()
)

//Json帮助类
type JsonHelper struct {
}

//获取一个新的JsonHelper实例
func NewJsonHelper() *JsonHelper {
	return &JsonHelper{}
}

//转换成字符串
func (h *JsonHelper) GetString(v interface{}) string {
	s, _ := h.Serialize(v)
	return s
}

//Json序列化
func (h *JsonHelper) Serialize(v interface{}) (string, error) {
	if nil == v {
		return "", errors.New("parameter is nil")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

//从字符串解析
func (h *JsonHelper) FromString(s string, t reflect.Type) interface{} {
	v, _ := h.Deserialize(s, t)
	return v
}

//Json反序列化
func (h *JsonHelper) Deserialize(s string, t reflect.Type) (interface{}, error) {
	if "" == s {
		return nil, errors.New("invalid parameter")
	}

	//CcObject.NewEx返回对象的指针
	ptr := CcObject.NewEx(t, true)

	err := json.Unmarshal([]byte(s), ptr)
	if err != nil {
		return nil, err
	}

	return ptr, nil
}
