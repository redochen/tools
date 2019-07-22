package json

import (
	"encoding/json"
	"errors"
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

//Json序列化
func (h *JsonHelper) Serialize(i interface{}) (string, error) {
	if nil == i {
		return "", errors.New("parameter is nil")
	}

	v, err := json.Marshal(i)
	if err != nil {
		return "", err
	} else {
		return string(v), nil
	}
}

//序列化成字符串
func (h *JsonHelper) GetString(i interface{}) string {
	value, _ := h.Serialize(i)
	return value
}
