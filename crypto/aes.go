package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	CcAes         = NewAesHelper()
	defaultAesKey = "Weav3kDf5VxmuuwB"
	defaultAesIvs = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

//AES帮助类
type AesHelper struct {
}

//获取一个新的AesHelper实例
func NewAesHelper() *AesHelper {
	return &AesHelper{}
}

//AES加密
func (h *AesHelper) Encrypt(plain string) ([]byte, error) {
	return h.EncryptEx(plain, nil, nil)
}

//AES加密
func (h *AesHelper) EncryptEx(plain string, key, iv []byte) ([]byte, error) {
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

//AES解密
func (h *AesHelper) Decrypt(crypt []byte) (string, error) {
	return h.DecryptEx(crypt, nil, nil)
}

//AES解密
func (h *AesHelper) DecryptEx(crypt []byte, key, iv []byte) (string, error) {
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
	aesKey := key
	if nil == aesKey {
		aesKey = []byte(defaultAesKey)
	} else {
		keyLen := len(aesKey)
		if keyLen != 16 && keyLen != 24 && keyLen != 32 {
			return nil, errors.New("length of key must be 16, 24 or 32")
		}
	}
	return aesKey, nil
}

//获取AES加密向量
func getAesIv(iv []byte) []byte {
	// 16 bytes for AES-128, 24 bytes for AES-192, 32 bytes for AES-256
	aesIvs := iv
	if nil == aesIvs {
		ivTemps := []byte(defaultAesIvs)
		aesIvs = ivTemps[:aes.BlockSize]
	} else {
		aesIvs = iv[:aes.BlockSize]
	}
	return aesIvs
}
