package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

var (
	CcHmac = NewHmacHelper()
)

//Hmac帮助类
type HmacHelper struct {
}

//获取一个新的HmacHelper实例
func NewHmacHelper() *HmacHelper {
	return &HmacHelper{}
}

//获取摘要
func (h *HmacHelper) Sum(key, plain string) string {
	if "" == plain {
		return ""
	}

	d := []byte(plain)

	c := hmac.New(md5.New, []byte(key))
	c.Write(d)
	hash := c.Sum([]byte(""))

	return hex.EncodeToString(hash)
}
