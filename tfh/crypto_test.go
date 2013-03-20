package main

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := []byte("xDG8KiW148cwUQIPS23pHxPSVppD8VIH")
	fileContents, err := parseFile("crypto_test.go")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
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
	fileContents, err := parseFile("crypto_test.go")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
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

func TestNumBytesToPad(t *testing.T) {
	// less than 16 bytes...
	x := []byte{0, 0, 0, 0, 0}
	numBytes := numBytesToPad(x)
	if numBytes != 27 {
		t.Log("Expected", 27, ", but got ", numBytes)
		t.Fail()
	}
	// 16 bytes...
	x = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	numBytes = numBytesToPad(x)
	if numBytes != 16 {
		t.Log("Expected", 16, ", but got ", numBytes)
		t.Fail()
	}
	// more than 16 bytes...
	x = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	numBytes = numBytesToPad(x)
	if numBytes != 15 {
		t.Log("Expected", 15, ", but got ", numBytes)
		t.Fail()
	}
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

func TestaddJunk(t *testing.T) {
	bs := make([][]byte, 2)
	bs[0] = []byte{1, 2}
	bs[1] = []byte{3}
	newBs := addJunk(bs, 2)
	if len(newBs) != 4 {
		t.Fail()
		return
	}
	for i, b := range bs {
		for j, b2 := range b {
			if uint(b2) != uint(newBs[i][j]) {
				t.Fail()
				return
			}
		}
	}
}
