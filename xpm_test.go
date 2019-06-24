package xpm

import (
	"fmt"
	"testing"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func TestNum2charcode(t *testing.T) {
	if num2charcode(0, letters) != "a" {
		t.Errorf("wanted a, got %s", num2charcode(0, letters))
	}
	if num2charcode(1, letters) != "b" {
		t.Errorf("wanted b, got %s", num2charcode(1, letters))
	}
	if num2charcode(25, letters) != "z" {
		t.Errorf("wanted z, got %s", num2charcode(25, letters))
	}
	if num2charcode(26, letters) != "aa" {
		t.Errorf("wanted aa, got %s", num2charcode(26, letters))
	}
	if num2charcode(27, letters) != "ab" {
		t.Errorf("wanted ab, got %s", num2charcode(27, letters))
	}
	if num2charcode(51, letters) != "az" {
		t.Errorf("wanted az, got %s", num2charcode(51, letters))
	}
	if num2charcode(52, letters) != "ba" {
		t.Errorf("wanted ba, got %s", num2charcode(52, letters))
	}
	if num2charcode(53, letters) != "bb" {
		t.Errorf("wanted bb, got %s", num2charcode(53, letters))
	}
	if inc("zzzzzzzz", letters) != "aaaaaaaaa" {
		t.Errorf("wanted aaaaaaaaa, got %s", inc("zzzzzzzz", letters))
	}
}

func TestNum2charcode2(t *testing.T) {
	for x := 18000; x < 18700; x++ {
		//fmt.Println(num2charcode(0, letters,(x))
		num2charcode(x, letters)
	}
}

func ExampleHexify() {
	fmt.Println(hexify([]byte{0, 7, 0x80, 0xff}))
	// Output:
	// [0x00 0x07 0x80 0xff]
}
