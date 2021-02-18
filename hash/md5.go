package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5Sum get md5 summary
func MD5Sum(plain string) string {
	if "" == plain {
		return ""
	}

	d := []byte(plain)
	//return fmt.Sprintf("%x", md5.Sum(d))

	c := md5.New()
	c.Write(d)
	hash := c.Sum([]byte(""))

	return hex.EncodeToString(hash)
}
