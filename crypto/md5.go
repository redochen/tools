package crypto

import (
	"crypto/md5"
	"fmt"
)

var (
	CcMd5 = NewMD5Helper()
)

//MD5帮助类
type MD5Helper struct {
}

//获取一个新的MD5Helper实例
func NewMD5Helper() *MD5Helper {
	return &MD5Helper{}
}

//获取摘要
func (h *MD5Helper) Sum(plain string) string {
	if "" == plain {
		return ""
	}
	data := []byte(plain)
	return fmt.Sprintf("%x", md5.Sum(data))
}
