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

func CorrectHash(key ID, value []byte) bool {
	bs := key[:]
	hash := Hash(bs)
	for i, b := range hash {
		if uint(b) != uint(hash[i]) {
			return false
		}
	}
	return true
}
