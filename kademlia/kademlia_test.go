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
	b.Contacts.PushBack(c)
	val := b.Contacts.Back().Value
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
