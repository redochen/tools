package json

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/redochen/tools/object"
)

//GetString 转换成字符串
func GetString(v interface{}) string {
	s, _ := Serialize(v)
	return s
}

//Serialize Json序列化
func Serialize(v interface{}) (string, error) {
	if nil == v {
		return "", errors.New("invalid parameter")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

//FromString 从字符串解析
func FromString(s string, t reflect.Type) interface{} {
	v, _ := Deserialize(s, t)
	return v
}

//Deserialize Json反序列化
func Deserialize(s string, t reflect.Type) (interface{}, error) {
	if "" == s {
		return nil, errors.New("invalid parameter")
	}

	//object.NewEx返回对象的指针
	ptr := object.NewEx(t, true)

	err := json.Unmarshal([]byte(s), ptr)
	if err != nil {
		return nil, err
	}

	return ptr, nil
}
