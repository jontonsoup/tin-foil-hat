package kademlia

import (
	"crypto/sha256"
	"errors"
	"log"
)

func HashStore(k *Kademlia, value []byte) (hash []byte, err error) {
	hash = Hash(value)
	key, err := fromBytes(hash)
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

func Verify(key ID, value []byte) error {
	bs := key[:]
	hash := Hash(bs)
	for i, b := range hash {
		if uint(b) != uint(hash[i]) {
			return errors.New("verify failure")
		}
	}
	return nil
}
