package crypto

import (
	"bytes"
)

//PKCS5Padding PKCS5补码
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

//PKCS5UnPadding PKCS5反补码
func PKCS5UnPadding(padText []byte) []byte {
	length := len(padText)
	// 去掉最后一个字节 unPadding 次
	unPadding := int(padText[length-1])
	return padText[:(length - unPadding)]
}
