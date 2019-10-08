package crypto

import (
	"encoding/base64"
	"errors"
)

//Base64编码字符串
func EncodeString(plain string) (string, error) {
	return EncodeData([]byte(plain))
}

//Base64编码数据
func EncodeData(plain []byte) (string, error) {
	if nil == plain || len(plain) <= 0 {
		return "", errors.New("invalid plain text")
	}

	return base64.StdEncoding.EncodeToString(plain), nil
}

//Base64解码字符串
func DecodeString(crypt string) (string, error) {
	b, err := DecodeData(crypt)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

//Base64解码数据
func DecodeData(crypt string) ([]byte, error) {
	if "" == crypt {
		return nil, errors.New("invalid encrypted text")
	}

	return base64.StdEncoding.DecodeString(crypt)
}
