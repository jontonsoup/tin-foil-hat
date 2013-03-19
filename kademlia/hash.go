package kademlia

import (
	"crypto/sha256"
	"log"
)

func HashStore(k *Kademlia, value []byte) (hash []byte, err error) {
	hash = Hash(value)
	key, err := FromBytes(hash)
	if err != nil {
		return
	}
	log.Println(key.AsString())
	_, err = IterativeStore(k, key, value)
	return
}

func Hash(bs []byte) []byte {
	h := sha256.New()
	h.Write(bs)

	return h.Sum(nil)
}

func CorrectHash(testHash []byte, value []byte) bool {
	hash := Hash(value)
	for i, b := range testHash {
		if uint(b) != uint(hash[i]) {
			return false
		}
	}
	return true
}
