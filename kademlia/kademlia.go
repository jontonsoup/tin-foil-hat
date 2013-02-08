package kademlia

import (
	"container/list"
	"log"
	"net"
)

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

const NUM_BUCKETS = IDBytes*8 + 1
const ALPHA = 3

// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
	NodeID  ID
	Buckets [NUM_BUCKETS]Bucket
	Self    Contact
	Table   map[ID][]byte
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
	k.updateContact(&k.Self)
	k.Table = make(map[ID][]byte)
	return k
}

// Correctly updates the bucket given that the contact given has just
// been observed
func (k *Kademlia) updateContact(c *Contact) {
	index := k.index(c.NodeID)
	b := &k.Buckets[index]

	e, ok := b.lookupContact(c.NodeID)
	if !ok {
		if b.contacts.Len() <= MAX_BUCKET_SIZE {
			log.Print("Adding contact to bucket: ", c.NodeID.AsString())
			b.contacts.PushBack(c)
		} else {
			// ping the least recently seen node
			firstEl := b.contacts.Front()
			first := firstEl.Value.(*Contact)
			_, err := SendPing(k, first.Address())
			if err != nil {
				log.Println("Old node responded, ignoring new contact")
				// first is now most recently seen
				b.contacts.MoveToBack(firstEl)
			} else {
				log.Println("Old node did not respond, evicting and adding new contact")
				b.contacts.Remove(firstEl)
				b.contacts.PushBack(c)
			}
		}
	} else {
		log.Print("Previously seen contact recently seen: ", c.NodeID.AsString())
		// Move contact to most recently seen in the bucket.
		b.contacts.MoveToBack(e)
	}
}

func (k *Kademlia) index(id ID) int {
	// convert ID into index for Kademlia.Buckets array
	return k.NodeID.Xor(id).PrefixLen()
}

func (k *Kademlia) closestNodes(searchID ID, excludedID ID, amount int) ([]FoundNode, error) {
	cs, err := k.closestContacts(searchID, excludedID, amount)
	if err != nil {
		return nil, err
	}
	nodes := make([]FoundNode, len(cs))

	for i, c := range cs {
		nodes[i] = contactToFoundNode(&c)
	}
	return nodes, nil
}

func (k *Kademlia) closestContacts(searchID ID, excludedID ID, amount int) (contacts []Contact, err error) {
	contacts = make([]Contact, 0)

	k.doInSearchOrder(searchID, func(index int) bool {
		// add as many contacts from bucket i as possible,

		currentBucket := k.Buckets[index].contacts
		sortedList := new(list.List)

		//sort that list |suspect|
		for e := currentBucket.Front(); e != nil; e = e.Next() {
			InsertSorted(sortedList, e.Value.(*Contact), func(first *Contact, second *Contact) bool {
				firstDistance := first.NodeID.Xor(searchID)
				secondDistance := second.NodeID.Xor(searchID)
				return firstDistance.Compare(secondDistance) == 1
			})
		}

		// (^._.^)~ kirby says add as much as you can to output slice
		for e := sortedList.Front(); e != nil; e = e.Next() {
			c := e.Value.(*Contact)
			if !c.NodeID.Equals(excludedID) {
				contacts = append(contacts, *c)
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

func (k *Kademlia) doInSearchOrder(id ID, usrFunc func(int) bool) {
	// produce the indices for the closest k-buckets to the id
	ones := k.NodeID.Xor(id).OnesIndices()

	log.Println("Searching for ones")
	for i := 0; i < NUM_BUCKETS; i++ {
		if ones[i] {
			if !usrFunc(i) {
				return
			}
		}
	}

	log.Println("Searching for zeros")
	for i := NUM_BUCKETS - 1; i >= 0; i-- {
		if !ones[i] {
			if !usrFunc(i) {
				return
			}
		}
	}

	log.Println("Done searching")
}
