package xpm

import (
	"fmt"
	"testing"
)

func TestNum2charcode(t *testing.T) {
	if num2charcode(0) != "a" {
		t.Errorf("wanted a, got %s", num2charcode(0))
	}
	if num2charcode(1) != "b" {
		t.Errorf("wanted b, got %s", num2charcode(1))
	}
	if num2charcode(25) != "z" {
		t.Errorf("wanted z, got %s", num2charcode(25))
	}
	if num2charcode(26) != "aa" {
		t.Errorf("wanted aa, got %s", num2charcode(26))
	}
	if num2charcode(27) != "ab" {
		t.Errorf("wanted ab, got %s", num2charcode(27))
	}
	if num2charcode(51) != "az" {
		t.Errorf("wanted az, got %s", num2charcode(51))
	}
	if num2charcode(52) != "ba" {
		t.Errorf("wanted ba, got %s", num2charcode(52))
	}
	if num2charcode(53) != "bb" {
		t.Errorf("wanted bb, got %s", num2charcode(53))
	}
	if inc("zzzzzzzz") != "aaaaaaaaa" {
		t.Errorf("wanted aaaaaaaaa, got %s", inc("zzzzzzzz"))
	}
}

func TestNum2charcode2(t *testing.T) {
	for x := 18000; x < 18700; x++ {
		//fmt.Println(num2charcode(x))
		num2charcode(x)
	}
}

func ExampleHexify() {
	fmt.Println(hexify([]byte{0, 7, 0x80, 0xff}))
	// Output:
	// [0x00 0x07 0x80 0xff]
}
