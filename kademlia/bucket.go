package kademlia

import (
	"container/list"
	"log"
)

/*
 Implements the Kademlia k-bucket type
*/

const MAX_BUCKET_SIZE = 20

// a simple list to implement a k-bucket
type Bucket struct {
	contacts list.List
}

// if a contact with the given id is present, return that contact
func (b *Bucket) lookupContact(id ID) (e *list.Element, ok bool) {
	log.Print("Looking up ", id.AsString())
	for e = b.contacts.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Contact)
		log.Print("Checking: ", c.NodeID.AsString())
		if c.NodeID.Equals(id) {
			ok = true
			return
		} else {
			log.Println("not equal: ", c.NodeID.AsString(), id.AsString())
		}
	}
	ok = false
	return
}
