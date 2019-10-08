package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	defaultAesKey = "oO3OEFhB7ALGGzAm"
	defaultAesIvs = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

//AES加密字符串
func EncryptString(plain string) ([]byte, error) {
	return EncryptEx(plain, nil, nil)
}

//AES加密
func Encrypt(plain, key string, iv []byte) ([]byte, error) {
	return EncryptEx(plain, []byte(key), iv)
}

//AES加密
func EncryptEx(plain string, key, iv []byte) ([]byte, error) {
	if "" == plain {
		return nil, errors.New("invalid plain text")
	}

	aesKey, err := getAesKey(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, getAesIv(iv))

	content := []byte(plain)
	content = PKCS5Padding(content, block.BlockSize())

	crypt := make([]byte, len(content))
	cbc.CryptBlocks(crypt, content)

	return crypt, nil
}

//AES解密字符串
func DecryptString(crypt []byte) (string, error) {
	return DecryptEx(crypt, nil, nil)
}

//AES解密
func Decrypt(crypt []byte, key string, iv []byte) (string, error) {
	return DecryptEx(crypt, []byte(key), iv)
}

//AES解密
func DecryptEx(crypt []byte, key, iv []byte) (string, error) {
	if nil == crypt || len(crypt) <= 0 {
		return "", errors.New("invalid crypt bytes")
	}

	aesKey, err := getAesKey(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	cbc := cipher.NewCBCDecrypter(block, getAesIv(iv))

	plain := make([]byte, len(crypt))
	cbc.CryptBlocks(plain, crypt)

	plain = PKCS5UnPadding(plain)
	content := string(plain)

	return content, nil
}

//获取AES加密KEY
func getAesKey(key []byte) ([]byte, error) {
	if nil == key {
		return []byte(defaultAesKey), nil
	}

	aesKey := key
	keyLen := len(aesKey)

	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		// 16 bytes for AES-128, 24 bytes for AES-192, 32 bytes for AES-256
		return nil, errors.New("length of key must be 16, 24 or 32")
	}

	return aesKey, nil
}

//获取AES加密向量
func getAesIv(iv []byte) []byte {
	if iv != nil && len(iv) >= aes.BlockSize {
		return iv[:aes.BlockSize]
	}

	return defaultAesIvs[:aes.BlockSize]
}
