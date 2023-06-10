package string

import (
	"testing"
)

func TestString(t *testing.T) {
	num := RomanToInt("MCMXCIV")
	if num != 1994 {
		t.Error("RomanToInt failed")
	}

	roman := IntToRoman(1994)
	if roman != "MCMXCIV" {
		t.Error("IntToRoman failed")
	}
}
