package kademlia

import (
	"log"
	"net"
)

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
	NodeID  ID
	Buckets [IDBytes*8 + 1]Bucket
	Self    Contact
}

func NewKademlia(address string) *Kademlia {
	ip, port, err := parseAddress(address)
	if err != nil {
		log.Fatal("bad address")
	}
	return newKademliaSplitAddress(ip, port)
}

func newKademliaSplitAddress(ip net.IP, port uint16) *Kademlia {
	k := new(Kademlia)
	k.NodeID = NewRandomID()
	k.Self = Contact{k.NodeID, ip, port}
	k.addContact(&k.Self)
	return k
}

func (k *Kademlia) index(id ID) int {
	// convert ID into index for Kademlia.Buckets array
	return k.NodeID.Xor(id).PrefixLen()
}

func LookupContact(k *Kademlia, id ID) (c *Contact, ok bool) {
	index := k.index(id)
	if index >= len(k.Buckets) {
		ok = false
		return
	}
	bucket := &k.Buckets[index]
	e, ok := bucket.lookupContact(id)
	if ok {
		c = e.Value.(*Contact)
	}
	return
}

func (k *Kademlia) addContact(c *Contact) {
	index := k.index(c.NodeID)
	bucket := &k.Buckets[index]
	bucket.updateContact(c)
}

func (k *Kademlia) indexSearchOrder(id ID) <-chan int {
	// go a goroutine that sends the correct indices out on the
	// channel and returns the channel

	return nil
}

func (k *Kademlia) produceIndexSearchOrder(id ID, outChan <-chan int) {
	// produce the indices for the closest k-buckets to the id

}
