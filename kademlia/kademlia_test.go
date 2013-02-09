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
	k.updateContact(c)
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

func TestDoInSearchOrder(t *testing.T) {
	k := NewKademlia("127.0.0.1:8890")
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

func bootupKademlia(listenStr string, firstPeerStr string) Kademlia {
	// By default, Go seeds its RNG with 1. This would cause every program to
	// generate the same sequence of IDs.
	rand.Seed(time.Now().UnixNano())

	kadem := kademlia.NewKademlia(listenStr)

	rpc.Register(kadem)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", listenStr)
	if err != nil {
		log.Fatal("Listen: ", err)
	}

	// Serve forever.
	go http.Serve(l, nil)

	// Confirm our server is up with a PING request and then exit.
	// Your code should loop forever, reading instructions from stdin and
	// printing their results to stdout. See README.txt for more details.
	pong, err := kademlia.SendPing(kadem, firstPeerStr)
	if err != nil {
		log.Fatal("Initial ping error: ", err)
	}

	log.Printf("pong msgID: %s\n", pong.MsgID.AsString())

	foundNodes, err := kademlia.SendFindNode(kadem, kadem.NodeID, firstPeerStr)

	if err != nil {
		log.Fatal("Bootstrap find_node error: ", err)
	}

	log.Println("Received", len(foundNodes), "nodes")
	for i, node := range foundNodes {
		log.Println("Node ", i, ": ", node.NodeID.AsString())
	}
	return kadem
}

func TestDoStore (t *testing.T){

	k := bootupKademlia("127.0.0.1:8890", "127.0.0.1:8090")
	k2 := bootupKademlia("127.0.0.1:8891", "127.0.0.1:8890")

}


