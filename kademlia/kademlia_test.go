package kademlia

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
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

func TestSendPing(t *testing.T) {
	nodeStr := "127.0.0.1:8890"
	listenStr := "127.0.0.1:8890"
	// create a kademlia instance, register it
	kadem := NewKademlia(listenStr)
	rpc.Register(kadem)
	// create an rpc server
	rpc.HandleHTTP()
	// start listening at listenStr
	l, err := net.Listen("tcp", listenStr)
	if err != nil {
		log.Fatal("Listen: ", err)
	}
	// start serving
	go http.Serve(l, nil)
	// send a ping to server at listenStr, using nodeStr
	// if there's an error, we'll fail this test
	_, err = SendPing(kadem, nodeStr)
	if err != nil {
		t.Fail()
	}
}

func ExampleContactToFoundNode() {
	c := new(Contact)
	c.NodeID = OnesID()
	c.Port = 7970
	c.Host = net.ParseIP("127.0.0.1")
	f := contactToFoundNode(c)
	fmt.Println(f.IPAddr)
	fmt.Println(f.NodeID)
	fmt.Println(f.Port)
	// Output:
	// 127.0.0.1
	// [1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1]
	// 7970
}

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
