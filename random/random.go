package random

import (
	"crypto/rand"
)

const (
	lettersAlpha  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	lettersNumber = "0123456789"
	lettersAll    = lettersNumber + lettersAlpha
)

//GetRandomNumber 获取随机数字字符串
func GetRandomNumber(length int) string {
	return GetRandomText(length, false, true)
}

//GetRandomString 获取随机字母字符串
func GetRandomString(length int) string {
	return GetRandomText(length, true, false)
}

//GetRandomText 获取随机字符串。includeAlpha－是否包含字母；includeNumber－则否包含数字
func GetRandomText(length int, includeAlpha, includeNumber bool) string {
	var letters string

	if includeAlpha && includeNumber {
		letters = lettersAll
	} else if includeNumber {
		letters = lettersNumber
	} else if includeAlpha {
		letters = lettersAlpha
	} else {
		return ""
	}
	var bytes = make([]byte, length)
	rand.Read(bytes)

	for k, v := range bytes {
		bytes[k] = letters[v%byte(len(letters))]
	}

	return string(bytes)
}
