package crypto

import (
	"encoding/base64"
	"errors"
)

//Base64EncodeString Base64编码字符串
func Base64EncodeString(plain string) (string, error) {
	return Base64EncodeData([]byte(plain))
}

//Base64EncodeData Base64编码数据
func Base64EncodeData(plain []byte) (string, error) {
	if nil == plain || len(plain) <= 0 {
		return "", errors.New("invalid plain text")
	}

	return base64.StdEncoding.EncodeToString(plain), nil
}

//Base64DecodeString Base64解码字符串
func Base64DecodeString(crypt string) (string, error) {
	b, err := Base64DecodeData(crypt)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

//Base64DecodeData Base64解码数据
func Base64DecodeData(crypt string) ([]byte, error) {
	if "" == crypt {
		return nil, errors.New("invalid encrypted text")
	}

	return base64.StdEncoding.DecodeString(crypt)
}
