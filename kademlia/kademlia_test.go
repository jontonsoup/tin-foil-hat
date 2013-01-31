package kademlia

import (
	"fmt"
	"net"
	"testing"
)

func ExampleBucket() {
	c := new(Contact)
	c.NodeID = NewRandomID()
	c.Port = 7970
	c.Host = net.ParseIP("127.0.0.1")
	b := new(Bucket)
	b.contacts.PushBack(c)
	val := b.contacts.Back().Value
	cval := val.(*Contact)
	fmt.Println(cval.Host)
	fmt.Println(cval.Port)
	// Output:
	// 127.0.0.1
	// 7970
}

func TestBucket(*testing.T) {
	// b := new(Bucket)
	// fmt.Println(b)
}

func ExampleIndex() {
	k := NewKademlia("127.0.0.1:8890")
	k.NodeID = HalfHalfID()
	other := OnesID()
	index := k.index(other)
	fmt.Println(index)
	// Output:
	// 80
}

func TestAddContact(t *testing.T) {
	// create contact
	c := new(Contact)
	c.NodeID = NewRandomID()
	c.Port = 7970
	c.Host = net.ParseIP("127.0.0.1")
	// create kademlia
	k := NewKademlia("127.0.0.1:8890")
	k.NodeID = HalfHalfID()
	// add a bucket at the right place
	b := new(Bucket)
	index := k.index(c.NodeID)
	k.Buckets[index] = *b
	// add contact
	k.addContact(c)
	if k.Buckets[index].contacts.Len() == 0 {
		t.Fail()
	}
	if _, ok := LookupContact(k, c.NodeID); !ok {
		t.Fail()
	}
	return
}

func TestParseAddress(t *testing.T) {

}
