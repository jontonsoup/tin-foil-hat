package main

import (
	"crypto/aes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
)

const CHUNK_SIZE = 256

func (tfh *TFH) encryptAndStore(filePath, key string) (outStr string, err error) {
	fileContents := parseFile(filePath)
	outStr, err = hashFile(fileContents)
	if err != nil {
		fmt.Println("Hash Failed somehow")
	}
	//	encryptedBytes := encrypt(fileContents, key)

	//	parts := splitBytes(encryptedBytes)

	return
}

//
// [2 bytes = length of padding ][32 bytes SHA hash of unencrypted file whole file][... X number of SHA hashes for file chunks]
//
//

// takes a file path and ecrypts that file, returning
// a hash representing a way to
/*
func Encrypt(filePath string, key string) (outStr string, err error) {

	fileContents := parseFile(filePath)
	outStr, _ = hashFile(fileContents)
	encrypted_file := aesEncryptFile(fileContents, key)
	plain_text := aesDecryptFile(encrypted_file, key)
	fmt.Println("before: ", string(fileContents))
	fmt.Println("after: ", string(plain_text))
	//determine how much we need to pad unencrypted file to make it mod 256 bit (append to output key (for user))

	//encrypt file with SHA hash as key with AES and append to output key (for user)

	//split file into 256 bit chunks

	//compute SHA hash of chunk in order append to output key (for user)

	//randomly choose a chunk, send chunk out to store
	//do until all chunks are stored

	//return completed key to user
	return
}*/

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

func encrypt(msg []byte, inputkey string) (encrypted_file []byte) {
	numPadBytes := numBytesToPad(msg)
	msg = padFile(msg, numPadBytes)
	// some key, 32 Byte long
	key := []byte(inputkey)

	// create the new cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	encrypted_file = make([]byte, 0)
	encrypt_block := make([]byte, c.BlockSize())

	for i := 0; i != len(msg); i = i + c.BlockSize() {
		c.Encrypt(encrypt_block, msg[i:i+c.BlockSize()])
		encrypted_file = append(encrypted_file, encrypt_block...)
	}
	return encrypted_file

}

func decrypt(encrypted_file []byte, inputkey string) []byte {
	key := []byte(inputkey)
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	// now we decrypt our encrypted text
	decrypted_bytes := make([]byte, 0)
	decrypt_block := make([]byte, c.BlockSize())
	for i := 0; i != len(encrypted_file); i = i + c.BlockSize() {
		c.Decrypt(decrypt_block, encrypted_file[i:i+c.BlockSize()])
		decrypted_bytes = append(decrypted_bytes, decrypt_block...)
	}

	return decrypted_bytes
}

func padFile(fileContents []byte, numBytesPadding int) (paddedFile []byte) {
	paddedFile = make([]byte, len(fileContents)+numBytesPadding)
	for i := 0; i < len(fileContents); i++ {
		paddedFile[i] = fileContents[i]
	}
	return
}

// deconstructs the user's key string to find out things like padding, sha key, etc
func destructureKeyString(key string) (padding []byte, unencHash []byte, chunks []byte) {
	b := []byte(key)
	padding = b[:2]
	unencHash = b[2:34]
	chunks = b[34:]
	return
}

// Returns number of bytes to pad to make given file mod 32 byte
func numBytesToPad(fileContents []byte) (numBytes int) {
	numBytes = 32 - (len(fileContents) % 32)
	return
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
