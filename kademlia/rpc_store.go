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
	k.updateContact(req.Sender)
	res.MsgID = req.MsgID
	k.Table[req.MsgID] = req.Value
	res.Err = nil
	return nil
}

func (k *Kademlia) SendStore(req StoreRequest, address string) (res *StoreResult){
	client, err := rpc.DialHTTP("tcp", address)
	log.Print("Sending Store ")
	if err != nil {
		return nil
	}
	res = new(StoreResult)

	err = client.Call("Kademlia.Store", req, &res)
	if err != nil {
		return
	}
	defer client.Close()
	return
}
