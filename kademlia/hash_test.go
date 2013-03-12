package kademlia

import (
	"crypto/sha256"
	"testing"
)

func TestSha(t *testing.T) {
	h := sha256.New()
	if h.Size() != IDBytes {
		t.Log("SHA is wrong!", h.Size())
		t.Fail()
	}
}

// make sure the hash can be cast to an ID
func TestHash(t *testing.T) {
	b1 := []byte("tasdf")
	hash := Hash(b1)
	if len(hash) != IDBytes {
		t.Log("Wrong sized hash: ", len(hash))
		t.Fail()
	}
	_, err := fromBytes(hash)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	return
}

func TestVerify(t *testing.T) {
	b1 := []byte("test")
	hash := Hash(b1)
	key, err := fromBytes(hash)
	if err != nil {
		t.Log("fromBytes err")
		t.Fail()
	}
	err = Verify(key, b1)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
