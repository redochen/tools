package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

//SHA256Sum get sha256 summary
func SHA256Sum(plain string) string {
	if "" == plain {
		return ""
	}

	d := []byte(plain)

	c := sha256.New()
	c.Write(d)
	hash := c.Sum([]byte(""))

	return hex.EncodeToString(hash)
}
