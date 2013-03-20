package main

import (
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
)

import (
	"tin-foil-hat/kademlia"
)

const KEY_SIZE = 32

type TFH struct {
	cipher cipher.Block
	kadem  *kademlia.Kademlia
}

func NewTFH(k *kademlia.Kademlia) *TFH {
	tfh := new(TFH)
	tfh.kadem = k
	key := makeRandKey(KEY_SIZE)
	cipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	tfh.cipher = cipher
	return tfh
}

func makeRandKey(keySize int) []byte {
	key := make([]byte, keySize)
	for i := 0; i < keySize; i++ {
		key[i] = uint8(rand.Intn(256))
	}
	return key
}
