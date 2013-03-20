package main

import (
	"bytes"
	"encoding/gob"
)

const MAX_FAKE_BYTE_RATIO = float64(0.1)

type tfhKey struct {
	Hash         []byte
	EncryptKey   []byte
	NumPadBytes  int
	PartKeys     [][]byte
	NumRealBytes int
}

func (tfhK tfhKey) serialize() (out []byte, err error) {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err = e.Encode(tfhK)
	return b.Bytes(), err
}

func unSerialize(b []byte) (tfhK *tfhKey, err error) {
	r := bytes.NewReader(b)
	d := gob.NewDecoder(r)
	var tk tfhKey
	err = d.Decode(&tk)
	return &tk, err
}
