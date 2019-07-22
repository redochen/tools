package crypto

import (
	"encoding/base64"
	"errors"
)

var (
	CcBase64 = NewBase64Helper()
)

//Base64帮助类
type Base64Helper struct {
}

//获取一个新的Base64Helper实例
func NewBase64Helper() *Base64Helper {
	return &Base64Helper{}
}

//Base64编码
func (h *Base64Helper) Encode(plain []byte) (string, error) {
	if nil == plain || len(plain) <= 0 {
		return "", errors.New("invalid plain text")
	}
	return base64.StdEncoding.EncodeToString(plain), nil
}

//Base64解码
func (h *Base64Helper) Decode(crypt string) ([]byte, error) {
	if "" == crypt {
		return nil, errors.New("invalid encrypted text")
	}
	return base64.StdEncoding.DecodeString(crypt)
}
