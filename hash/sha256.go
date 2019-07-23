package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	//"fmt"
)

var (
	CcSHA256 = NewSHA256Helper()
)

//SHA256帮助类
type SHA256Helper struct {
}

//获取一个新的SHA256Helper实例
func NewSHA256Helper() *SHA256Helper {
	return &SHA256Helper{}
}

//获取摘要
func (h *SHA256Helper) Sum(plain string) string {
	if "" == plain {
		return ""
	}

	d := []byte(plain)

	c := sha256.New()
	c.Write(d)
	hash := c.Sum([]byte(""))

	return hex.EncodeToString(hash)
}
