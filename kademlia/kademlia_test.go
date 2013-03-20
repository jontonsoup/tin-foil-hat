package kademlia

import (
	"fmt"
	"log"
	"testing"
)

func ExampleIndex() {
	k, err := NewUnBootedKademlia("127.0.0.1:8890")
	if err != nil {
		log.Fatal(err)
		return
	}
	k.NodeID = HalfHalfID()
	other := OnesID()
	index := k.index(other)
	fmt.Println(index)
	// Output:
	// 128
}

// func TestAddContact(t *testing.T) {
// 	// create contact
// 	c := new(Contact)
// 	c.NodeID = NewRandomID()
// 	c.Port = 7210
// 	c.Host = net.ParseIP("127.0.0.1")
// 	// create kademlia
// 	k, err := NewUnBootedKademlia("127.0.0.1:8890")
// 	if err != nil {
// 		t.Log(err)
// 		t.Fail()
// 		return
// 	}
// 	//	fmt.Println(k.NodeID)
// 	k.NodeID = HalfHalfID()
// 	// add a bucket at the right place
// 	b := new(Bucket)
// 	index := k.index(c.NodeID)
// 	k.Buckets[index] = *b
// 	// add contact
// 	fmt.Println(index)
// 	k.updateContact(*c)
// 	if k.Buckets[index].contacts.Len() == 0 {
// 		t.Fail()
// 	}
// 	if _, ok := LookupContact(k, c.NodeID); !ok {
// 		t.Fail()
// 	}
// 	return
// }

// func ExampleContactToFoundNode() {
// 	c := new(Contact)
// 	c.NodeID = OnesID()
// 	c.Port = 7970
// 	c.Host = net.ParseIP("127.0.0.1")
// 	f := contactToFoundNode(*c)
// 	fmt.Println(f.IPAddr)
// 	fmt.Println(f.NodeID)
// 	fmt.Println(f.Port)
// 	// Output:
// 	// 127.0.0.1
// 	// [1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1]
// 	// 7970
// }

func TestOnesIndices(t *testing.T) {
	ones := OnesID().OnesIndices()
	fmt.Println(ones)
	for i := 0; i < NUM_BUCKETS; i++ {
		if i%8 == 0 {
			if !ones[i] {
				t.Error("SHould be 1:", i)
			}
		} else if ones[i] {
			t.Error("SHould be 0:", i)
		}
	}

	zeros := new(ID).OnesIndices()
	for i := 0; i < NUM_BUCKETS; i++ {
		if i != NUM_BUCKETS-1 {
			if zeros[i] {
				t.Error("SHould be 0:", i)
			}
		} else if !zeros[i] {
			t.Error("Should be 1:", i)
		}
	}
}

func TestDoInSearchOrder(t *testing.T) {
	k, err := NewUnBootedKademlia("127.0.0.1:8890")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	k.NodeID = HalfHalfID()

	other := OnesID()

	count := 0
	k.doInSearchOrder(other, func(index int) bool {
		if index < 0 || index >= NUM_BUCKETS {
			t.Fail()
		}
		count++
		return true
	})
	if count != NUM_BUCKETS {
		t.Fail()
	}

}

// func TestDoStore(t *testing.T) {
// 	bootupKademlia("127.0.0.1:8765", "127.0.0.1:8765")
// 	// val := []byte("foobar")
// 	// stoReq := StoreRequest{k.Self, NewRandomID(), NewRandomID(), val}
// 	// stoRes := new(StoreResult)
// 	// SendStore(k, )
// 	//	go bootupKademlia("127.0.0.1:8891", "127.0.0.1:8765")
// }
