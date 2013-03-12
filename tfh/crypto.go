package main

import (
	//"crypto/aes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
)

func (tfh *TFH) encryptFile() {

}

//
// [2 bytes = length of padding ][32 bytes SHA hash of unencrypted file whole file][... X number of SHA hashes for file chunks]
//
//

// takes a file path and ecrypts that file, returning
// a hash representing a way to
func Encrypt(filePath string) (outStr string, err error) {

	fileContents := parseFile(filePath)
	outStr, _ = hashFile(fileContents)

	//determine how much we need to pad unencrypted file to make it mod 256 bit (append to output key (for user))

	//encrypt file with SHA hash as key with AES and append to output key (for user)

	//split file into 256 bit chunks

	//compute SHA hash of chunk in order append to output key (for user)

	//randomly choose a chunk, send chunk out to store
	//do until all chunks are stored

	//return completed key to user
	return
}

// Returns number of bytes to pad to make given file mod 32 byte
func numBytesToPad(fileContents []byte) (numBytes int) {
	x := len(fileContents) % 32
	if x <= 16 {
		x = 32 - x
	}
	numBytes = x
	return
}

func Decrypt(key string) (outStr string, err error) {
	//deconstruct key into parts

	//randomly call find on parts

	//check to see if each part matches its SHA key (for file integrity)

	//order parts

	//decrypt whole ordered file

	//remove padding from file

	//return file bytes to user
	return
}

func aesEncryptFile(msg []byte) {
	// some key, 16 Byte long
	key := []byte{0x2c, 0x88, 0x25, 0x1a, 0xaa, 0xae, 0xc2, 0xa2, 0xaf, 0xe7, 0x84, 0x8a, 0x10, 0xcf, 0xe3, 0x2a}

	Println("len of message: ", len(msg))
	Println("len of key: ", len(key))
	// create the new cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		Println("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	// now we convert our string into 32 long byte array
	msgbuf := strings.Bytes(msg)
	out := make([]byte, len(msg))

	c.Encrypt(msgbuf[0:16], out[0:16])   // encrypt the first half
	c.Encrypt(msgbuf[16:32], out[16:32]) // encrypt the second half

	Println("len of encrypted: ", len(out))
	Println(">> ", out)

	// now we decrypt our encrypted text
	plain := make([]byte, len(out))
	c.Decrypt(out[0:16], plain[0:16])   // decrypt the first half
	c.Decrypt(out[16:32], plain[16:32]) // decrypt the second half

	Println("msg: ", string(plain))
}

func hashFile(fileContents []byte) (outStr string, err error) {

	//compute the SHA on the untouched file for sanity check
	h := sha256.New()
	h.Write(fileContents)
	shaSum := h.Sum(nil)
	outStr = fmt.Sprintf("% x", shaSum)
	return
}

func parseFile(filePath string) (fileContents []byte) {
	// open the file
	file, err := os.Open(filePath) // For read access.
	if err != nil {
		err = errors.New("FATAL: Error opening file")
		return
	}
	defer file.Close()

	// read file into a string
	fstat, err := file.Stat()
	if err != nil {
		err = errors.New("FATAL: fstat fail")
		return
	}
	fileContents = make([]byte, fstat.Size())
	_, err = file.Read(fileContents)
	if err != nil {
		err = errors.New("FATAL: READ ERROR")
		return
	}

	return
}
