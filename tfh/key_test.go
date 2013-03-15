package main

import (
	"bytes"
	"testing"
)

func TestSerialize(t *testing.T) {
	tk := new(tfhKey)
	tk.Hash = []byte("testatasdlf;kjasd;lfj")
	tk.EncryptKey = []byte("xDG8KiW148cwUQIPS23pHxPSVppD8VIH")
	tk.NumPadBytes = 5
	tk.PartKeys = make([][]byte, 2)
	tk.PartKeys[0] = []byte("whathasdfa")
	tk.PartKeys[1] = []byte("as;dlfjkas;ldkfj")

	out, err := tk.serialize()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(out)
	decoded, err := unSerialize(out)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	if bytes.Compare(tk.Hash, decoded.Hash) != 0 {
		t.Log("Hashes are different", tk.Hash, decoded.Hash)
		t.Fail()
	}
	if bytes.Compare(tk.EncryptKey, decoded.EncryptKey) != 0 {
		t.Log("Encrypt keys are different", tk.EncryptKey, decoded.EncryptKey)
		t.Fail()
	}
	if tk.NumPadBytes != decoded.NumPadBytes {
		t.Log("Num pad bytes different", tk.NumPadBytes, decoded.NumPadBytes)
		t.Fail()
	}
	for i, _ := range tk.PartKeys {
		if bytes.Compare(tk.PartKeys[i], decoded.PartKeys[i]) != 0 {
			t.Log("Hashes are different", tk.PartKeys[i], decoded.PartKeys[i])
			t.Fail()
		}
	}
}
