package crypto

import (
	"encoding/base64"
	"errors"
)

var (
	//CcBase64 Base64帮助类实例
	CcBase64 = newBase64Helper()
)

//Base64Helper Base64帮助类
type Base64Helper struct {
}

//获取一个新的Base64Helper实例
func newBase64Helper() *Base64Helper {
	return &Base64Helper{}
}

//EncodeString Base64编码字符串
func (h *Base64Helper) EncodeString(plain string) (string, error) {
	return h.EncodeData([]byte(plain))
}

//EncodeData Base64编码数据
func (h *Base64Helper) EncodeData(plain []byte) (string, error) {
	if nil == plain || len(plain) <= 0 {
		return "", errors.New("invalid plain text")
	}

	return base64.StdEncoding.EncodeToString(plain), nil
}

//DecodeString Base64解码字符串
func (h *Base64Helper) DecodeString(crypt string) (string, error) {
	b, err := h.DecodeData(crypt)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

//DecodeData Base64解码数据
func (h *Base64Helper) DecodeData(crypt string) ([]byte, error) {
	if "" == crypt {
		return nil, errors.New("invalid encrypted text")
	}

	return base64.StdEncoding.DecodeString(crypt)
}
