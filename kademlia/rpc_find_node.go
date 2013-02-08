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
		k.updateContact(&c)
	}
	return res.Nodes, err
}

func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// TODO: Implement.
	log.Println("Handling FindNode request:", req)
	// Find the k closest nodes to FindNodeRequest.NodeID and pack
	// them in FindNodeResult.Nodes
	k.updateContact(&req.Sender)

	log.Println("Finding close nodes")
	nodes, err := k.closestNodes(req.NodeID, req.Sender.NodeID, MAX_BUCKET_SIZE)

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

func IterativeFindNode(k *Kademlia, searchID nodeID) ([]Contact, error) {
	shortList := new(list.List)
	alreadySeen := new(map[*Contact]bool)
	initNodes := closestContacts(searchID, k.NodeID, ALPHA)

	closerNode := func(c1 *Contact, c2 *Contact) bool {
		d1 := c1.Xor(searchID)
		d2 := c2.Xor(searchID)
		return d1.Compare(d2) == 1
	}

	for _, node := range initNodes {
		InsertSorted(shortList, node, closerNode)
	}

	for {
		nextSearchNodes := getUnseen(shortList, alreadySeen)

		if len(nextSearchNodes) == 0 {
			break
		}

		// send find_node rpcs to nextSearchNodes

		// mark them alreadySeen
		for _, node := range nextSearchNodes {
			alreadySeen[node] = true
		}

		// aggregate the new contacts into the shortList, keeping
		// only the K closest
	}

	return closestNodes, nil
}
