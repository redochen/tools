package cache

import (
	"testing"
)

func TestString(t *testing.T) {
	key := "key"
	value := "value"

	SetString(key, value)
	val1 := GetString(key, true)
	if val1 == value {
		t.Log("TestString success")
	} else {
		t.Error("TestString failed")
	}

	val2 := GetString(key)
	if val2 == "" {
		t.Log("TestString success")
	} else {
		t.Error("TestString failed")
	}
}

func TestInt(t *testing.T) {
	key := "key"
	value := 1111

	SetInt(key, value)
	val := GetInt(key)

	if val == value {
		t.Log("TestInt success")
	} else {
		t.Error("TestInt failed")
	}
}

func TestBool(t *testing.T) {
	key := "key"
	value := true

	SetBool(key, value)
	val := GetBool(key)

	if val == value {
		t.Log("TestBool success")
	} else {
		t.Error("TestBool failed")
	}
}

func TestFloat(t *testing.T) {
	key := "key"
	value := float32(13.14)

	SetFloat(key, value)
	val := GetFloat(key)

	if val == value {
		t.Log("TestFloat success")
	} else {
		t.Error("TestFloat failed")
	}
}
