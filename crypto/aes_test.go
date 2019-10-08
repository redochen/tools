package crypto

import (
	"testing"
)

func TestAes(t *testing.T) {
	var str1 = "test aes"

	data, err := EncryptString(str1)
	if err != nil {
		t.Error("TestAes failed EncryptString:", err)
		return
	}

	str2, err := DecryptString(data)
	if err != nil {
		t.Error("TestAes failed DecryptString:", err)
		return
	}

	if str2 == str1 {
		t.Log("TestAes success")
	} else {
		t.Error("TestAes failed")
	}
}
