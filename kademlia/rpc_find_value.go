package kademlia

import (
	"errors"
	"log"
	"net/rpc"
)

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

func SendFindValue(k *Kademlia, key ID, nodeID ID) (ret *FindValueResult, err error) {
	contact, _ := LookupContact(k, nodeID)
	client, err := rpc.DialHTTP("tcp", contact.Address())
	if err != nil {
		return
	}
	req := new(FindValueRequest)
	req.MsgID = NewRandomID()
	req.Sender = k.Self
	req.Key = key
	err = client.Call("Kademlia.FindValue", req, &ret)
	if err != nil {
		return
	}
	defer client.Close()
	if !ret.MsgID.Equals(ret.MsgID) {
		err = errors.New("FindValue MsgID didn't match SendFindValue MsgID")
	}
	return
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
	// return value for FindValueRequest.Key; if not found, return k closest nodes
	// be sure to fill up res.MsgID with req.MsgID!

	// update contact!
	k.updateContact(req.Sender)

	if val, ok := k.Table[req.Key]; ok {
		// return value
		res.MsgID = req.MsgID
		res.Value = val
		res.Err = nil
		res.Nodes = nil
	} else {
		log.Println("No value for key:", req.Key.AsString())
		// return closest nodes
		res.MsgID = req.MsgID
		res.Value = nil
		res.Err = nil
		// run FindNode rpc, by calling SendFindNode
		// var find_node_res *FindNodeResult
		find_node_res := new(FindNodeResult)
		msgID := NewRandomID()
		find_node_req := FindNodeRequest{req.Sender, msgID, req.Key}
		err := k.FindNode(find_node_req, find_node_res)
		if err != nil {
			res.Err = errors.New("FindNode rpc call returned an error")
		}
		res.Nodes = find_node_res.Nodes
	}

	return nil
}
