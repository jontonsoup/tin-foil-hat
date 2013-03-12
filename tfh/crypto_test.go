package main

import (
	"fmt"
	"testing"
)

func ExampleNumBytesToPad1() {
	x := []byte{0, 0, 0, 0, 0}
	fmt.Println(numBytesToPad(x))
	// Output:
	// 27
}

func ExampleNumBytesToPad2() {
	x := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fmt.Println(numBytesToPad(x))
	// Output:
	// 16
}

func ExampleNumBytesToPad3() {
	x := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fmt.Println(numBytesToPad(x))
	// Output:
	// 15
}

func TestPadFile(t *testing.T) {
	x := []byte{4, 2}
	y := padFile(x, 2)
	// padded file should have the right length
	if len(y) != 4 {
		t.Fail()
	}
	// padded file should have the original file's content
	for i := 0; i < len(x); i++ {
		if y[i] != x[i] {
			t.Fail()
		}
	}
}
