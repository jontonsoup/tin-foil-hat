package kademlia

import (
	"testing"
)

// make sure the hash can be cast to an ID
func Testhash(t *testing.T) {
	b1 := []byte("test")
	hash := Hash(b1)
	if len(hash) != IDBytes {
		t.Fail()
	}
	_, err := fromBytes(hash)
	if err != nil {
		t.Fail()
	}
	return
}
