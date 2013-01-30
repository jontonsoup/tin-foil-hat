package kademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import "net"

// Host identification.
type Contact struct {
	NodeID ID
	Host   net.IP
	Port   uint16
}

// PING
type Ping struct {
	Sender Contact
	MsgID  ID
}

type Pong struct {
	MsgID ID
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
	// This one's a freebie.
	pong.MsgID = k.NodeID
	return nil
}

// STORE
type StoreRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
	Value  []byte
}

type StoreResult struct {
	MsgID ID
	Err   error
}

func (k *Kademlia) Store(req StoreRequest, res *StoreResult) error {
	// TODO: Implement.
	return nil
}

// FIND_NODE
type FindNodeRequest struct {
	Sender Contact
	MsgID  ID
	NodeID ID
}

type FoundNode struct {
	IPAddr string
	Port   uint16
	NodeID ID
}

type FindNodeResult struct {
	MsgID ID
	Nodes []FoundNode
	Err   error
}

func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// TODO: Implement.
	return nil
}

// FIND_VALUE
type FindValueRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
	MsgID ID
	Value []byte
	Nodes []FoundNode
	Err   error
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
	// TODO: Implement.
	return nil
}
