package kademlia

import (
	"errors"
	"log"
	"net/rpc"
)

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
	log.Println("Handling store request from", req.Sender.Address())
	k.updateContact(req.Sender)
	res.MsgID = req.MsgID
	err := Verify(req.Key, req.Value)
	if err != nil {
		return err
	}
	k.Table[req.Key] = req.Value
	res.Err = nil
	return nil
}

func SendStore(k *Kademlia, key ID, value []byte, nodeID ID) error {
	c, ok := LookupContact(k, nodeID)
	if !ok {
		return errors.New("node not found")
	}
	address := c.Address()

	client, err := rpc.DialHTTP("tcp", address)
	log.Println("Sending Store rpc to", address)
	if err != nil {
		k.removeContact(c.NodeID)
		return nil
	}

	msgID := NewRandomID()
	req := StoreRequest{k.Self, msgID, key, value}

	var res StoreResult
	err = client.Call("Kademlia.Store", req, &res)
	if err != nil {
		return err
	}
	defer client.Close()

	return res.Err
}

func IterativeStore(k *Kademlia, key ID, value []byte) (lastID ID, err error) {

	k.Table[key] = value

	nodes, err := IterativeFindNode(k, key)

	if err != nil {
		return
	}

	if len(nodes) == 0 {
		return k.NodeID, nil
	}

	for _, node := range nodes {
		SendStore(k, key, value, node.NodeID)
	}

	log.Println("Found", len(nodes), nodes)
	lastID = nodes[len(nodes)-1].NodeID

	return
}
