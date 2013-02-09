package kademlia

import (
	"errors"
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

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
	// return value for FindValueRequest.Key; if not found, return k closest nodes
	// be sure to fill up res.MsgID with req.MsgID!

	// update contact!
	k.updateContact(&req.Sender)

	if val, ok := k.Table[req.Key]; ok {
		// return value
		res.MsgID = req.MsgID
		res.Value = val
		res.Err = nil
		res.Nodes = nil
	} else {
		// return closest nodes
		res.MsgID = req.MsgID
		res.Value = nil
		res.Err = nil
		// run FindNode rpc, by calling SendFindNode
		var find_node_res *FindNodeResult
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
