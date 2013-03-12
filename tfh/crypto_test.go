package main

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := "xDG8KiW148cwUQIPS23pHxPSVppD8VIH"
	fileContents := parseFile("crypto_test.go")
	encrypted := encrypt(fileContents, key)
	decrypted := decrypt(encrypted, key)
	for i, b := range fileContents {
		if uint(b) != uint(decrypted[i]) {
			t.Log("Decryption doesn't match original file")
			t.Log("original: ", string(encrypted), "decrypted")
			t.Fail()
			return
		}
	}
	if len(encrypted)%CHUNK_SIZE != 0 {
		t.Log("Encrypted file should be splittable into", CHUNK_SIZE, "byte chunks")
		t.Fail()
	}
	return
}

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
