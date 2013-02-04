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

// Correctly updates the bucket given that the contact given has just
// been observed
func (b *Bucket) updateContact(c *Contact) {
	e, ok := b.lookupContact(c.NodeID)
	if !ok {
		if b.contacts.Len() <= MAX_BUCKET_SIZE {
			log.Print("Adding contact to bucket: ", c.NodeID.AsString())
			b.contacts.PushBack(c)
		} else {
			// TODO: properly update a full bucket. Ping
			// the first contact in the list
		}
	} else {
		log.Print("Previously seen contact recently seen: ", c.NodeID.AsString())
		// Move contact to most recently seen in the bucket.
		b.contacts.MoveToBack(e)
	}
}
