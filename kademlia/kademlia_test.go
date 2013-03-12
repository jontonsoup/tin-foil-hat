package kademlia

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"testing"
	"time"
)

func ExampleIndex() {
	k := NewKademlia("127.0.0.1:8890")
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
// 	k := bootupKademlia("127.0.0.1:8890", "127.0.0.1:8890")
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

// func TestParseAddress(t *testing.T) {

// }

// func TestSendPing(t *testing.T) {
// 	nodeStr := "127.0.0.1:8890"
// 	listenStr := "127.0.0.1:8890"
// 	// create a kademlia instance, register it
// 	kadem := NewKademlia(listenStr)
// 	rpc.Register(kadem)
// 	// create an rpc server
// 	rpc.HandleHTTP()
// 	// start listening at listenStr
// 	l, err := net.Listen("tcp", listenStr)
// 	if err != nil {
// 		log.Fatal("Listen: ", err)
// 	}
// 	// start serving
// 	go http.Serve(l, nil)
// 	// send a ping to server at listenStr, using nodeStr
// 	// if there's an error, we'll fail this test
// 	_, err = SendPing(kadem, nodeStr)
// 	if err != nil {
// 		t.Fail()
// 	}
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

func bootupKademlia(listenStr string, firstPeerStr string) *Kademlia {
	// By default, Go seeds its RNG with 1. This would cause every program to
	// generate the same sequence of IDs.
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("kademlia starting up!\n")
	kadem := NewKademlia(listenStr)

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
	pong, err := SendPing(kadem, firstPeerStr)
	if err != nil {
		log.Fatal("Initial ping error: ", err)
	}

	fmt.Printf("pong msgID: %s\n", pong.MsgID.AsString())

	_, err = SendFindNodeAddr(kadem, kadem.NodeID, firstPeerStr)

	var foundNodes []Contact
	if err == nil {
		foundNodes, err = IterativeFindNode(kadem, kadem.NodeID)
	}

	if err != nil {
		log.Fatal("Bootstrap find_node error: ", err)
	}

	fmt.Println("Received", len(foundNodes), "nodes")
	for i, node := range foundNodes {
		fmt.Println("Node ", i, ": ", node.NodeID.AsString())
	}
	return kadem
}

// func TestDoStore(t *testing.T) {
// 	bootupKademlia("127.0.0.1:8765", "127.0.0.1:8765")
// 	// val := []byte("foobar")
// 	// stoReq := StoreRequest{k.Self, NewRandomID(), NewRandomID(), val}
// 	// stoRes := new(StoreResult)
// 	// SendStore(k, )
// 	//	go bootupKademlia("127.0.0.1:8891", "127.0.0.1:8765")
// }
