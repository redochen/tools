package compress

import (
	"testing"
)

func TestBestZip(t *testing.T) {
	var str1 = "test BestZip"

	data, err := CcGzip.BestZip(str1)
	if err != nil {
		t.Error("TestingBestZip failed BestZip:", err)
		return
	}

	str2, err := CcGzip.Unzip(data)
	if err != nil {
		t.Error("TestingBestZip failed Unzip:", err)
		return
	}

	if str2 == str1 {
		t.Log("TestingBestZip success")
	} else {
		t.Error("TestingBestZip failed")
	}
}

func TestFastZip(t *testing.T) {
	var str1 = "test FastZip"

	data, err := CcGzip.FastZip(str1)
	if err != nil {
		t.Error("TestingFastZip failed FastZip:", err)
		return
	}

	str2, err := CcGzip.Unzip(data)
	if err != nil {
		t.Error("TestingFastZip failed Unzip:", err)
		return
	}

	if str2 == str1 {
		t.Log("TestingFastZip success")
	} else {
		t.Error("TestingFastZip failed")
	}
}
