package kademlia

import (
	"crypto/aes"
	"crypto/sha256"
	"os"
)


//
// [2 bytes = length of padding ][32 bytes SHA hash of unencrypted file whole file][... X number of SHA hashes for file chunks]
//
//

// takes a file path and ecrypts that file, returning
// a hash representing a way to
func Encrypt(filePath string) (outStr string, err error) {
  // open the file
	file, err := os.Open("file.go") // For read access.
	if err != nil {
		outStr = "FATAL"
		return
	}

  //compute the SHA on the untouched file for sanity check

  //determine how much we need to pad unencrypted file to make it mod 256 bit (append to output key (for user))

  //encrypt file with SHA hash as key with AES and append to output key (for user)

  //split file into 256 bit chunks

  //compute SHA hash of chunk in order append to output key (for user)

  //randomly choose a chunk, send chunk out to store
  //do until all chunks are stored

  //return completed key to user
}
