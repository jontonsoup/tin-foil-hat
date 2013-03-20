package kademlia

import (
	"container/list"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

const NUM_BUCKETS = IDBytes*8 + 1
const ALPHA = 3

// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
	NodeID  ID
	Buckets [NUM_BUCKETS]Bucket // TODO: refreshes
	Self    Contact
	Table   map[ID][]byte // TODO: republishes
}

func NewKademlia(address string, firstPeerAddr string) (k *Kademlia, err error) {
	k, err = NewUnBootedKademlia(address, firstPeerAddr)
	if err != nil {
		return
	}
	err = k.bootUp(address, firstPeerAddr)

	return
}

func NewUnBootedKademlia(listenAddr, peerAddr string) (k *Kademlia, err error) {
	ip, port, err := parseAddress(listenAddr)
	if err != nil {
		return
	}
	rand.Seed(time.Now().UnixNano())

	k = new(Kademlia)
	k.NodeID = NewRandomID()
	k.Self = Contact{k.NodeID, ip, port}
	for i, _ := range k.Buckets {
		b := &k.Buckets[i]
		b.k = k
		b.index = i
		b.start()
	}
	k.updateContact(k.Self)
	k.Table = make(map[ID][]byte)
	return
}

func LocalLookup(k *Kademlia, key ID) ([]byte, bool) {
	val, ok := k.Table[key]
	return val, ok
}

func (k *Kademlia) bootUp(listenAddr string, peerAddr string) (err error) {
	rpc.Register(k)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("Listen: ", err)
		return
	}

	// Serve forever.
	go http.Serve(l, nil)

	_, err = SendPing(k, peerAddr)
	if err != nil {
		log.Fatal("Initial ping error: ", err)
		return
	}

	_, err = IterativeFindNode(k, k.NodeID)
	if err != nil {
		log.Fatal("Bootstrap find_node error: ", err)
	}
	return
}

func (k *Kademlia) removeContact(id ID) {
	index := k.index(id)
	b := &k.Buckets[index]
	b.removeContact(id)
}

// Correctly updates the bucket given that the contact given has just
// been observed
func (k *Kademlia) updateContact(c Contact) {
	index := k.index(c.NodeID)
	b := &k.Buckets[index]
	b.updateContact(c)
}

func (k *Kademlia) index(id ID) int {
	// convert ID into index for Kademlia.Buckets array
	return k.NodeID.Xor(id).PrefixLen()
}

func (k *Kademlia) closestNodes(searchID ID, excludedID ID, amount int) []FoundNode {
	cs := k.closestContacts(searchID, excludedID, amount)

	nodes := make([]FoundNode, len(cs))

	for i, c := range cs {
		nodes[i] = contactToFoundNode(c)
	}
	return nodes
}

func (k *Kademlia) closestContacts(searchID ID, excludedID ID, amount int) (contacts []Contact) {
	contacts = make([]Contact, 0)

	k.doInSearchOrder(searchID, func(index int) bool {
		// add as many contacts from bucket i as possible,

		currentBucket := k.Buckets[index].contacts
		sortedList := new(list.List)

		//sort that list |suspect|
		for e := currentBucket.Front(); e != nil; e = e.Next() {
			insertSorted(sortedList, e.Value.(Contact), func(first Contact, second Contact) int {
				firstDistance := first.NodeID.Xor(searchID)
				secondDistance := second.NodeID.Xor(searchID)
				return firstDistance.Compare(secondDistance)
			})
		}

		// (^._.^)~ kirby says add as much as you can to output slice
		for e := sortedList.Front(); e != nil; e = e.Next() {
			c := e.Value.(Contact)
			if !c.NodeID.Equals(excludedID) {
				contacts = append(contacts, c)
				// if the slice is full, break
				if len(contacts) == amount {
					return false
				}
			}
		}
		// slice isn't full, do on the next index
		return true
	})
	return
}

func LookupContact(k *Kademlia, id ID) (c Contact, ok bool) {
	index := k.index(id)
	if index >= len(k.Buckets) {
		ok = false
		return
	}
	bucket := &k.Buckets[index]
	e, ok := bucket.lookupContact(id)
	if ok {
		c = e.Value.(Contact)
	}
	return
}

func (k *Kademlia) doInSearchOrder(id ID, usrFunc func(int) bool) {
	// produce the indices for the closest k-buckets to the id
	ones := k.NodeID.Xor(id).OnesIndices()

	for i := 0; i < NUM_BUCKETS; i++ {
		if ones[i] {
			if !usrFunc(i) {
				return
			}
		}
	}

	for i := NUM_BUCKETS - 1; i >= 0; i-- {
		if !ones[i] {
			if !usrFunc(i) {
				return
			}
		}
	}

}
