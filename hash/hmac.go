package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

//HmacSum get hmac summary
func HmacSum(key, plain string) string {
	if "" == plain {
		return ""
	}

	d := []byte(plain)

	c := hmac.New(md5.New, []byte(key))
	c.Write(d)
	hash := c.Sum([]byte(""))

	return hex.EncodeToString(hash)
}
