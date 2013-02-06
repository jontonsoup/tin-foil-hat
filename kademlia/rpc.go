package kademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

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

func SendPing(k *Kademlia, address string) (pong *Pong, err error) {
	client, err := rpc.DialHTTP("tcp", address)
	defer client.Close()
	if err != nil {
		return
	}
	ping := new(Ping)
	ping.MsgID = k.NodeID
	ping.Sender = k.Self
	err = client.Call("Kademlia.Ping", ping, &pong)
	if err != nil {
		return
	}
	if !pong.MsgID.Equals(ping.MsgID) {
		err = errors.New("Pong MsgID didn't match Ping MsgID")
	}

	return
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
	// This one's a freebie.
	if !ping.Sender.NodeID.Equals(k.NodeID) {
		log.Print("Adding contact: ", ping.Sender.NodeID.AsString())
		k.addContact(&ping.Sender)
	}
	pong.MsgID = ping.MsgID
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

func parseAddress(address string) (ip net.IP, port uint16, err error) {
	hoststr, portstr, err := net.SplitHostPort(address)
	if err != nil {
		return
	}
	ip = net.ParseIP(hoststr)
	port64, err := strconv.ParseUint(portstr, 10, 16)
	if err != nil {
		log.Fatal("parseAddress failed!")
		return
	}
	port = uint16(port64)
	return ip, port, nil
}

func contactToFoundNode(c *Contact) FoundNode {
	// create a FoundNode and return it
	f := new(FoundNode)
	f.IPAddr = c.Host.String()
	f.NodeID = c.NodeID
	f.Port = c.Port
	return *f
}
