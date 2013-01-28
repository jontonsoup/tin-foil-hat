package kademlia

import (
	"fmt"
	"testing"
)

func ExampleBucket() {
	c := new(Contact)
	c.NodeID = "some freaking id"
	c.Port = "7970"
	c.IPAddress = "127.0.0.1"
	b := new(Bucket)
	b.Contacts.PushBack(c)
	val := b.Contacts.Back().Value
	cval := val.(*Contact)
	fmt.Println(cval.IPAddress)
	// Output:
	// 127.0.0.1
}

func TestBucket(*testing.T) {
	// b := new(Bucket)
	// fmt.Println(b)
}
