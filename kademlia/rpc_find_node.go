package kademlia

import (
	"errors"
	"log"
	"net/rpc"
)

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
