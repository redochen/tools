package crypto

import (
	"testing"
)

func TestBase64(t *testing.T) {
	var str1 = "test base64"

	data, err := CcBase64.EncodeString(str1)
	if err != nil {
		t.Error("TestBase64 failed EncodeString:", err)
		return
	}

	str2, err := CcBase64.DecodeString(data)
	if err != nil {
		t.Error("TestBase64 failed DecodeString:", err)
		return
	}

	if str2 == str1 {
		t.Log("TestBase64 success")
	} else {
		t.Error("TestBase64 failed")
	}
}
