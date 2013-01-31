package kademlia

import (
	"container/list"
	"log"
)

/* 
 Implements the Kademlia k-bucket type
*/

// a simple list to implement a k-bucket
type Bucket struct {
	contacts list.List
}

// if a contact with the given id is present, return that contact
func (b *Bucket) lookupContact(id ID) (c *Contact, ok bool) {
	log.Print("Looking up ", id.AsString())
	for e := b.contacts.Front(); e != nil; e = e.Next() {
		c = e.Value.(*Contact)
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
	// TODO: check to see if the list is full or the contact is
	// already there
	b.contacts.PushBack(c)
}
