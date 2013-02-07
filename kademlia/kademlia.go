package kademlia

import (
	"container/list"
	"log"
	"net"
)

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

const NUM_BUCKETS = IDBytes*8 + 1

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
	k.addContact(&k.Self)
	k.Table = make(map[ID][]byte)
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
	bucketIndexes, doneChan := k.indexSearchOrder(searchID)
	contacts = make([]Contact, 0)
indicesLoop:
	for i := range bucketIndexes {
		// add as many contacts from bucket i as possible,

		currentBucket := k.Buckets[i].contacts
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

				if len(contacts) == amount {
					break indicesLoop
				}
			}
		}
	}
	log.Println("Done now")
	doneChan <- true
	return
}

func InsertSorted(inputlist *list.List, item *Contact, greaterThan func(*Contact, *Contact) bool) {
	for e := inputlist.Front(); e != nil; e = e.Next() {
		if greaterThan(e.Value.(*Contact), item) {
			inputlist.InsertBefore(item, e)
			return
		}
	}
	inputlist.PushBack(item)
}

func (k *Kademlia) indexSearchOrder(id ID) (<-chan int, chan<- bool) {
	// go a goroutine that sends the correct indices out on the
	// channel and returns the channel
	indicesChan := make(chan int)
	doneChan := make(chan bool)
	go k.produceIndexSearchOrder(id, indicesChan, doneChan)
	return indicesChan, doneChan
}

func (k *Kademlia) produceIndexSearchOrder(id ID, outChan chan<- int, doneChan <-chan bool) {
	// produce the indices for the closest k-buckets to the id
	ones := k.NodeID.Xor(id).OnesIndices()

	i := 0
	searchOnes := true
	for {
		select {
		case <-doneChan:
			return
		default:
			// not done, produce more!
		}

		// check if we're searching back or forward
		if searchOnes {
			// first loop right looking for ones
			for ; i < NUM_BUCKETS; i++ {
				if ones[i] {
					outChan <- i
					i++
					break
				}
			}
			if i == NUM_BUCKETS {
				searchOnes = false
				i--
			}
		} else {
			// then loop back looking for zeros
			for ; i >= 0; i-- {
				if !ones[i] {
					outChan <- i
					i--
					break
				}
			}
			if i < 0 {
				break
			}
		}
	}
	close(outChan)
	<-doneChan
}
