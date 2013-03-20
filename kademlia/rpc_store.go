package kademlia

import (
	"errors"
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
	k.updateContact(req.Sender)
	res.MsgID = req.MsgID

	if !CorrectHash(req.Key[:], req.Value) {
		return errors.New("Bad hash")
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

	lastID = nodes[len(nodes)-1].NodeID

	return
}
