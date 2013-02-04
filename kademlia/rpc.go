package kademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
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
	if !pong.MsgID.Equals(k.NodeID) {
		host, port, err := parseAddress(address)
		if err != nil {
			return pong, err
		}
		contact := &Contact{pong.MsgID, host, port}
		log.Print("Adding contact ", pong.MsgID.AsString())
		k.addContact(contact)
	} else {
		log.Print("Stop pinging yourself!")
	}

	return
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
	// This one's a freebie.
	if !ping.MsgID.Equals(k.NodeID) {
		log.Print("Adding contact: ", ping.Sender.NodeID.AsString())
		k.addContact(&ping.Sender)
	}
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
