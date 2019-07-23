package crypto

import (
	"crypto/sha1"
	"encoding/hex"
)

var (
	CcSHA1 = NewSHA1Helper()
)

//SHA1帮助类
type SHA1Helper struct {
}

//获取一个新的SHA1Helper实例
func NewSHA1Helper() *SHA1Helper {
	return &SHA1Helper{}
}

//获取摘要
func (h *SHA1Helper) Sum(plain string) string {
	if "" == plain {
		return ""
	}

	d := []byte(plain)

	c := sha1.New()
	c.Write(d)
	hash := c.Sum([]byte(""))

	return hex.EncodeToString(hash)
}
