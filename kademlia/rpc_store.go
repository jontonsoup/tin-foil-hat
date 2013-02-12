package kademlia

import (
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
<<<<<<< HEAD
	k.updateContact(&req.Sender)
=======
	log.Println("Handling store request from", req.Sender.Address())
	k.updateContact(req.Sender)
>>>>>>> 7cdba3c2b6b9f427243582771905375407810327
	res.MsgID = req.MsgID
	k.Table[req.Key] = req.Value
	res.Err = nil
	return nil
}

func (k *Kademlia) sendStore(key ID, value []byte, address string) error {
	client, err := rpc.DialHTTP("tcp", address)
	log.Println("Sending Store rpc to", address)
	if err != nil {
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

func IterativeStore(k *Kademlia, key ID, value []byte) error {

	k.Table[key] = value

	nodes, err := IterativeFindNode(k, key)

	if err != nil {
		return err
	}

	for _, node := range nodes {
		k.sendStore(key, value, node.Address())
	}

	return nil
}
