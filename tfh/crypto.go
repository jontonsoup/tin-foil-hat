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

func aesEncryptFile(msg []byte, inputkey string) {
	// some key, 32 Byte long
	key := []byte(inputkey)
	str := string(key)
	fmt.Println("key: ", str)
	// fmt.Println("message: ", msg)

	// // create the new cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	encrypted_file := make([]byte, len(msg))
	out := make([]byte, 16)

	for i := 0; i < len(msg) - 16 ; i = i + 16 {
		c.Encrypt(out, msg[i:i+15])
		encrypted_file = append(encrypted_file, out...)
	}

	fmt.Println("len of encrypted: ", len(encrypted_file))
	fmt.Println(">> ", encrypted_file)

	// now we decrypt our encrypted text
	plain := make([]byte, len(out))
	decrypt_block := make([]byte, 16)
	for i := 0; i < len(encrypted_file) - 16; i = i + 16 {
		c.Decrypt(decrypt_block, encrypted_file[i:i+15])
		plain = append(plain, decrypt_block...)
	}


	// fmt.Println("msg: ", plain)
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
