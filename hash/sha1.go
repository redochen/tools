package crypto

import (
	"crypto/sha1"
	"encoding/hex"
)

//SHA1Sum get sha1 summary
func SHA1Sum(plain string) string {
	if "" == plain {
		return ""
	}

	d := []byte(plain)

	c := sha1.New()
	c.Write(d)
	hash := c.Sum([]byte(""))

	return hex.EncodeToString(hash)
}
