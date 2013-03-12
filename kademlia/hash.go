package kademlia

import (
	"crypto/sha256"
)

func HashStore(k *Kademlia, value []byte) (hash []byte, err error) {
	hash = Hash(value)
	key, err := fromBytes(hash)
	if err != nil {
		return
	}
	_, err = IterativeStore(k, key, value)
	return
}

func Hash(bs []byte) []byte {
	h := sha256.New()
	return h.Sum(bs)
}
