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
	k.updateContact(req.Sender)
	res.MsgID = req.MsgID
	k.Table[req.MsgID] = req.Value
	res.Err = nil
	return nil
}

func (k *Kademlia) SendStore(key ID, value []byte, address string) error {
	client, err := rpc.DialHTTP("tcp", address)
	log.Print("Sending Store ")
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

func IterativeStore(k *Kademlia, key ID, value ID) error {
	return errors.New("not implemented")
}
