package json

import (
	"reflect"
	"testing"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestJson(t *testing.T) {
	p1 := &Person{
		Name: "Json.Test",
		Age:  18,
	}

	s, err := Serialize(p1)
	if err != nil {
		t.Error("TestJson failed Serialize:", err)
		return
	}

	//t.Log(s)

	v, err := Deserialize(s, reflect.TypeOf(p1))
	if err != nil {
		t.Error("TestJson failed Deserialize:", err)
		return
	}

	p2 := v.(*Person)

	if p2.Name == p1.Name && p2.Age == p1.Age {
		t.Log("TestJson success")
	} else {
		t.Error("TestJson failed")
	}
}
