package kademlia

import "container/list"

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
	NodeID  ID
	Buckets list.List
}

// a simple list to implement a k-bucket
type Bucket struct {
	Contacts list.List
}

func NewKademlia() *Kademlia {
	// TODO: Assign yourself a random ID and prepare other state here.
	k := new(Kademlia)
	k.NodeID = NewRandomID()
	return k
}
