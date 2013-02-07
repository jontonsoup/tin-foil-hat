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
	if err != nil {
		return
	}

	ping := new(Ping)
	ping.MsgID = NewRandomID()
	ping.Sender = k.Self
	err = client.Call("Kademlia.Ping", ping, &pong)
	if err != nil {
		return
	}
	defer client.Close()

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

func SendFindNode(k *Kademlia, nodeID ID, address string) ([]FoundNode, error) {
	// TODO
	// send a findNode rpc and return the k-closest nodes
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	msgID := NewRandomID()
	req := FindNodeRequest{k.Self, msgID, k.NodeID}

	var res *FindNodeResult
	err = client.Call("Kademlia.FindNode", req, &res)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	if !res.MsgID.Equals(req.MsgID) {
		err = errors.New("Response MsgID didn't match request MsgID")
	} else if res.Err != nil {
		err = res.Err
	}

	for _, node := range res.Nodes {
		c := foundNodeToContact(&node)
		k.addContact(&c)
	}
	return res.Nodes, err
}

func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// TODO: Implement.
	log.Println("Handling FindNode request:", req)
	// Find the k closest nodes to FindNodeRequest.NodeID and pack
	// them in FindNodeResult.Nodes
	k.addContact(&req.Sender)

	log.Println("Finding close nodes")
	nodes, err := k.closestNodes(req.NodeID, req.Sender.NodeID)

	if err != nil {
		res.Err = err
		return err
	}

	log.Println("Finishing up")
	//set the msg id
	res.MsgID = req.MsgID
	res.Nodes = nodes

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

func foundNodeToContact(f *FoundNode) Contact {
	c := new(Contact)
	ip := net.ParseIP(f.IPAddr)
	c.Host = ip
	c.NodeID = f.NodeID
	c.Port = f.Port
	return *c
}
