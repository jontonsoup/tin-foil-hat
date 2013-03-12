package main

import (
	"crypto/aes"
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
func Encrypt(filePath string, key string) (outStr string, err error) {

	fileContents := parseFile(filePath)
	outStr, _ = hashFile(fileContents)
	aesEncryptFile(fileContents, key)

	//determine how much we need to pad unencrypted file to make it mod 256 bit (append to output key (for user))

	//encrypt file with SHA hash as key with AES and append to output key (for user)

	//split file into 256 bit chunks

	//compute SHA hash of chunk in order append to output key (for user)

	//randomly choose a chunk, send chunk out to store
	//do until all chunks are stored

	//return completed key to user
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

func aesEncryptFile(msg []byte, inputkey string) {
	// some key, 32 Byte long
	key := []byte(inputkey)
	str := string(key)
	fmt.Println("key: ", str)

	// // create the new cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	out := make([]byte, len(msg))

	c.Encrypt(msg, out)

	fmt.Println("len of encrypted: ", len(out))
	fmt.Println(">> ", out)

	// now we decrypt our encrypted text
	plain := make([]byte, len(out))
	c.Decrypt(out, plain)

	fmt.Println("msg: ", string(plain))
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
