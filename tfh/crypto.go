package main

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"tin-foil-hat/kademlia"
)

const CHUNK_SIZE = 32

//this is the main function that encrypts the file and stores it across kademlia
func (tfh *TFH) encryptAndStore(filePath, keyFilePath, key string) (returnPath string, err error) {
	fileContents, err := parseFile(filePath)
	if err != nil {
		return
	}
	tk := new(tfhKey)
	tk.EncryptKey = []byte(key)
	tk.Hash, err = hashFile(fileContents)
	if err != nil {
		return
	}
	var encryptedBytes []byte
	encryptedBytes, tk.NumPadBytes = encrypt(fileContents, tk.EncryptKey)

	parts := splitBytes(encryptedBytes)
	tk.NumRealBytes = len(parts)

	ratio := randomRatio(MAX_FAKE_BYTE_RATIO)
	numFakeBytes := int(math.Ceil(float64(tk.NumRealBytes) * ratio))
	parts = addJunk(parts, numFakeBytes)
	tk.PartKeys, err = tfh.storeAll(parts)
	if err != nil {
		return
	}

	// append the keys to the outstr
	decryptKey, err := tk.serialize()
	var decryptKeyStr string
	if err == nil {
		decryptKeyStr = hex.EncodeToString(decryptKey)
	}
	tfh.storeDecryptKeyString(decryptKeyStr, keyFilePath)
	returnPath = keyFilePath
	return
}

// stores the key at specified path
func (tfh *TFH) storeDecryptKeyString(decryptKeyStr string, filePath string) {
	// create and open the file, if none exists; overwrite if it does
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// write string to file
	_, err = file.WriteString(decryptKeyStr)
	if err != nil {
		panic(err)
	}
	return
}

// stores the file at the path
func (tfh *TFH) writeFile(filepath string, data []byte) {
	// create and open the file, if none exists; overwrite if it does
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// write string to file
	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
	return
}

func (tfh *TFH) retrieveDecryptKeyString(filePath string) (decryptKeyStr string, err error) {
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = errors.New(fmt.Sprintf("no such file or directory: %s", filePath))
		return
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	decryptKeyStr = string(content)
	return
}

//this is the main function (run by the cli) that both gets the file and decrypts it
func (tfh *TFH) decryptAndGet(pathToKey string, pathToFile string) (outStr string, err error) {
	//deconstruct key into parts
	key, err := tfh.retrieveDecryptKeyString(pathToKey)
	if err != nil {
		return
	}
	keybytes, _ := hex.DecodeString(key)
	tfhkey, _ := unSerialize(keybytes)

	allBytes, err := tfh.findAll(tfhkey.PartKeys)
	if err != nil {
		return
	}
	bytes := allBytes[:tfhkey.NumRealBytes]
	flattened_bytes := flatten(bytes)

	decryptBytes := decrypt(flattened_bytes, tfhkey.EncryptKey)
	decryptBytes = trimDecryptedFile(decryptBytes, tfhkey.NumPadBytes)
	hash, err := hashFile(decryptBytes)
	if err != nil {
		return
	}
	if !kademlia.CorrectHash(hash, decryptBytes) {
		err = errors.New("Reassembled file has bad hash. Aborting!!!")
		return
	}
	tfh.writeFile(pathToFile, decryptBytes)
	outStr = pathToFile

	return
}

func encrypt(msg []byte, inputkey []byte) (encrypted_file []byte, numPadBytes int) {
	numPadBytes = numBytesToPad(msg)
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

func parseFile(filePath string) (fileContents []byte, err error) {
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = errors.New(fmt.Sprintf("No such file or directory: %s", filePath))
		return
	}

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
	order := rand.Perm(len(vals))

	for _, i := range order {
		keys[i], err = kademlia.HashStore(tfh.kadem, vals[i])
		if err != nil {
			return
		}
	}

	return
}

// returns the keys in the not the same order they were
func (tfh *TFH) findAll(keys [][]byte) (values [][]byte, err error) {
	values = make([][]byte, len(keys))
	order := rand.Perm(len(keys))

	for _, i := range order {
		id, _ := kademlia.FromBytes(keys[i])
		findValResult, errCheck := kademlia.IterativeFindValue(tfh.kadem, id)
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Iterative find value error: %v", err))
			return
		}
		values[i] = findValResult.Value
	}

	return
}

//flattens a [][]byte array into []byte
func flatten(byte_array [][]byte) (flattened_bytes []byte) {
	flattened_bytes = make([]byte, 0)
	for i := 0; i < len(byte_array); i++ {
		flattened_bytes = append(flattened_bytes, byte_array[i]...)
	}
	return
}

func trimDecryptedFile(file []byte, numToTrim int) (trimFile []byte) {
	trimFile = file[:len(file)-numToTrim]
	return
}

func addJunk(bs ([][]byte), numBytes int) (newBs [][]byte) {
	newBs = bs[:]
	for i := 0; i < numBytes; i++ {
		junk := makeRandKey(numBytes)
		newBs = append(newBs, junk)
	}
	return
}

func randomRatio(maxRat float64) float64 {
	return rand.Float64() * maxRat
}
