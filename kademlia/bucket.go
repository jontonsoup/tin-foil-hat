package kademlia

import (
	"container/list"
	"log"
	"time"
)

/*
 Implements the Kademlia k-bucket type
*/

const MAX_BUCKET_SIZE = 20

const TREFRESH = 3600 * time.Second

//const TREFRESH =  * time.Second

// a simple list to implement a k-bucket
type Bucket struct {
	contacts    list.List
	k           *Kademlia
	index       int
	refreshChan chan bool
}

func (b *Bucket) start() {
	refreshChan := make(chan bool)
	b.refreshChan = refreshChan

	go b.waitRefresh()
}

func (b *Bucket) waitRefresh() {
	if b.index == len(b.k.Buckets)-1 {
		for {
			<-b.refreshChan
		}
	}
	for {
		timeOut := time.After(TREFRESH)
		select {
		case <-b.refreshChan:
			// restart timer
		case <-timeOut:
			// refresh bucket
			if b.contacts.Len() != 0 {
				b.refresh()
			}
		}
	}
}

func (b *Bucket) refresh() {
	// selects a random number in the bucket's range and do an
	// iterativeFindNode using that number as key

	c := b.contacts.Front().Value.(Contact)
	idInBucket := c.NodeID
	// TODO: make this random
	// idInBucket := NewRandomWithPrefix(b.k.NodeID, b.index)
	if index := b.k.index(idInBucket); b.index != index {
		// TODO: remove this once it's better tested
		log.Fatal("Supposed to do bucket", b.index, "got bucket", index)
	}
	IterativeFindNode(b.k, idInBucket)
}

func (b *Bucket) removeContact(id ID) {
	e, ok := b.lookupContact(id)
	if ok {
		b.contacts.Remove(e)
	}
}

func (b *Bucket) updateContact(c Contact) {
	e, ok := b.lookupContact(c.NodeID)

	b.refreshChan <- true
	if !ok {
		if b.contacts.Len() <= MAX_BUCKET_SIZE {
			b.contacts.PushBack(c)
		} else {
			// ping the least recently seen node
			firstEl := b.contacts.Front()
			first := firstEl.Value.(Contact)
			_, err := SendPing(b.k, first.Address())
			if err != nil {
				// first is now most recently seen
				b.contacts.MoveToBack(firstEl)
			} else {
				b.contacts.Remove(firstEl)
				b.contacts.PushBack(c)
			}
		}
	} else {
		// Move contact to most recently seen in the bucket.
		b.contacts.MoveToBack(e)
	}
}

// if a contact with the given id is present, return that contact
func (b *Bucket) lookupContact(id ID) (e *list.Element, ok bool) {

	for e = b.contacts.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		if c.NodeID.Equals(id) {
			ok = true
			return
		}
	}

	ok = false
	return
}
