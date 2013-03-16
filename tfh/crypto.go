package main

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"kademlia-secure/kademlia"
	"os"
)

const CHUNK_SIZE = 32

func (tfh *TFH) encryptAndStore(filePath, key string) (decryptKeyStr string, err error) {
	fileContents := parseFile(filePath)
	tk := new(tfhKey)
	tk.EncryptKey = []byte(key)
	tk.Hash, err = hashFile(fileContents)
	if err != nil {
		return
	}
	var encryptedBytes []byte
	encryptedBytes, tk.NumPadBytes = encrypt(fileContents, tk.EncryptKey)

	parts := splitBytes(encryptedBytes)
	tk.PartKeys, err = tfh.storeAll(parts)
	if err != nil {
		return
	}

	// append the keys to the outstr
	decryptKey, err := tk.serialize()
	if err == nil {
		decryptKeyStr = hex.EncodeToString(decryptKey)
	}
	return
}

//
// [2 bytes = length of padding ][32 bytes SHA hash of unencrypted file whole file][... X number of SHA hashes for file chunks]
//
//

// takes a file path and ecrypts that file, returning
// a hash representing a way to

//determine how much we need to pad unencrypted file to make it mod 256 bit (append to output key (for user))

//encrypt file with SHA hash as key with AES and append to output key (for user)

//split file into 256 bit chunks

//compute SHA hash of chunk in order append to output key (for user)

//randomly choose a chunk, send chunk out to store
//do until all chunks are stored

//return completed key to user

func (tfh *TFH) decryptAndGet(key string) (outStr string, err error) {
	//deconstruct key into parts
	keybytes, _ := hex.DecodeString(key)
	tfhkey, _ := unSerialize(keybytes)

	fmt.Println("Padding: ", tfhkey.NumPadBytes)
	fmt.Println("unencHash: ", tfhkey.Hash)
	fmt.Println("chunks: ", tfhkey.PartKeys)
	//randomly call find on parts

	//check to see if each part matches its SHA key (for file integrity)

	//order parts

	//decrypt whole ordered file

	//remove padding from file

	//return file bytes to user
	return
}

func encrypt(msg []byte, inputkey []byte) (encrypted_file []byte, numPadBytes int) {
	numPadBytes = numBytesToPad(msg)
	fmt.Println("pad:", numPadBytes)
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
	return

}

func decrypt(encrypted_file []byte, key []byte) []byte {
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

func hashFile(fileContents []byte) (hash []byte, err error) {

	//compute the SHA on the untouched file for sanity check
	h := sha256.New()
	h.Write(fileContents)
	hash = h.Sum(nil)
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

func splitBytes(b []byte) (split [][]byte) {
	for i := 0; i < len(b); i += CHUNK_SIZE {
		split = append(split, b[i:i+CHUNK_SIZE])
	}
	return
}

// returns the keys in the same order they were
func (tfh *TFH) storeAll(vals [][]byte) (keys [][]byte, err error) {
	keys = make([][]byte, len(vals))
	order := randomOrder(len(vals))

	for _, i := range order {
		keys[i], err = kademlia.HashStore(tfh.kadem, vals[i])
		if err != nil {
			return
		}
	}

	return
}

// TODO: make actually random
func randomOrder(length int) (order []int) {
	order = make([]int, length)
	for i := 0; i < length; i++ {
		order[i] = i
	}
	return
}
