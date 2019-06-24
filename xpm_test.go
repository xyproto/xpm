package xpm

import (
	"fmt"
	"testing"
)

func TestNum2charcode(t *testing.T) {
	if num2charcode(0, 'a', 'z') != "a" {
		t.Errorf("wanted a, got %s", num2charcode(0, 'a', 'z'))
	}
	if num2charcode(1, 'a', 'z') != "b" {
		t.Errorf("wanted b, got %s", num2charcode(1, 'a', 'z'))
	}
	if num2charcode(25, 'a', 'z') != "z" {
		t.Errorf("wanted z, got %s", num2charcode(25, 'a', 'z'))
	}
	if num2charcode(26, 'a', 'z') != "aa" {
		t.Errorf("wanted aa, got %s", num2charcode(26, 'a', 'z'))
	}
	if num2charcode(27, 'a', 'z') != "ab" {
		t.Errorf("wanted ab, got %s", num2charcode(27, 'a', 'z'))
	}
	if num2charcode(51, 'a', 'z') != "az" {
		t.Errorf("wanted az, got %s", num2charcode(51, 'a', 'z'))
	}
	if num2charcode(52, 'a', 'z') != "ba" {
		t.Errorf("wanted ba, got %s", num2charcode(52, 'a', 'z'))
	}
	if num2charcode(53, 'a', 'z') != "bb" {
		t.Errorf("wanted bb, got %s", num2charcode(53, 'a', 'z'))
	}
	if inc("zzzzzzzz", 'a', 'z') != "aaaaaaaaa" {
		t.Errorf("wanted aaaaaaaaa, got %s", inc("zzzzzzzz", 'a', 'z'))
	}
}

func TestNum2charcode2(t *testing.T) {
	for x := 18000; x < 18700; x++ {
		//fmt.Println(num2charcode(0, 'a', 'z',(x))
		num2charcode(x, 'a', 'z')
	}
}

func ExampleHexify() {
	fmt.Println(hexify([]byte{0, 7, 0x80, 0xff}))
	// Output:
	// [0x00 0x07 0x80 0xff]
}
