package main

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := []byte("xDG8KiW148cwUQIPS23pHxPSVppD8VIH")
	fileContents := parseFile("crypto_test.go")
	encrypted, _ := encrypt(fileContents, key)
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
		t.Log("But it has size", len(encrypted))
		t.Fail()
	}
	return
}

func TestSplitBytes(t *testing.T) {
	key := []byte("zDglWUd5gbjArrcjxbg8t4Wfspszpxyp")
	fileContents := parseFile("crypto_test.go")
	encrypted, _ := encrypt(fileContents, key)
	split := splitBytes(encrypted)

	totalSplitSize := 0
	for _, bytes := range split {
		totalSplitSize += len(bytes)
	}
	if totalSplitSize != len(encrypted) {
		t.Log("Encrypted and split aren't the same size")
		t.Fail()
	}
	// each byte array in the split array should correspond to a
	// slice of the original array
	for i, bs := range split {
		if len(bs) != CHUNK_SIZE {
			t.Log("encrypted chunks should be", CHUNK_SIZE, "long")
			t.Fail()
		}
		for j, b := range bs {
			if uint(encrypted[i*CHUNK_SIZE+j]) != uint(b) {
				t.Log("Split bytes don't match original")
				t.Fail()
			}
		}
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
