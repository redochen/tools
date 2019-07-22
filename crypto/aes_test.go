package crypto

import (
	"testing"
)

func TestAes(t *testing.T) {
	var str1 = "test aes"

	data, err := CcAes.EncryptString(str1)
	if err != nil {
		t.Error("TestAes failed EncryptString:", err)
		return
	}

	str2, err := CcAes.DecryptString(data)
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
