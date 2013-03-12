package main

import (
	"fmt"
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
