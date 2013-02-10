package kademlia

import (
	"container/list"
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
	log.Println("Sending FindNode rpc to ", address)
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

	nodes := k.closestNodes(req.NodeID, req.Sender.NodeID, MAX_BUCKET_SIZE)

	//set the msg id
	res.MsgID = req.MsgID
	res.Nodes = nodes

	return nil
}

func IterativeFindNode(k *Kademlia, searchID ID) ([]Contact, error) {
	shortList := list.New()
	alreadySeen := make(map[ID]bool)
	initNodes := k.closestContacts(searchID, k.NodeID, ALPHA)
	log.Println(len(initNodes), "closest contacts found", initNodes)
	isCloser := func(c1 *Contact, c2 *Contact) int {
		d1 := c1.NodeID.Xor(searchID)
		d2 := c2.NodeID.Xor(searchID)
		return d1.Compare(d2)
	}

	insertUnseenSorted(shortList, initNodes, isCloser, alreadySeen, MAX_BUCKET_SIZE)
	log.Println(shortList.Len(), "in the shortList")
	for e := shortList.Front(); e != nil; e = e.Next() {
		log.Println(e.Value.(*Contact).NodeID.AsString(), "is in the shortlist")
	}

	for {

		nextSearchNodes := getUnseen(shortList, alreadySeen, ALPHA)

		if len(nextSearchNodes) == 0 {
			break
		}

		newNodesChan := k.goFindNodes(nextSearchNodes, searchID)

		// send find_node rpcs to nextSearchNodes, add their nodes to shortList
		for _, _ = range nextSearchNodes {
			response := <-newNodesChan
			if response.Err != nil {
				removeFromSorted(shortList, response.searchNode.NodeID)
			}
			// mark them alreadySeen
			alreadySeen[response.searchNode.NodeID] = true

			// aggregate the new contacts into the shortList, keeping
			// only the K closest
			newNodes := make([]Contact, 0)

			for _, node := range response.FoundNodes {
				newNodes = append(newNodes, foundNodeToContact(&node))
			}

			insertUnseenSorted(shortList, newNodes, isCloser, alreadySeen, MAX_BUCKET_SIZE)
		}

	}

	closestNodes := make([]Contact, 0)

	for e := shortList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Contact)
		closestNodes = append(closestNodes, *c)
	}

	return closestNodes, nil
}

type SignedFoundNodes struct {
	FoundNodes []FoundNode
	Err        error
	searchNode *Contact
}

func (k *Kademlia) goFindNodes(searchNodes []Contact, searchID ID) <-chan SignedFoundNodes {
	outChan := make(chan SignedFoundNodes)

	for _, node := range searchNodes {
		log.Println("sending rpc to ", node.NodeID)
		go func() {
			foundNodes, err := SendFindNode(k, searchID, node.Address())
			output := SignedFoundNodes{foundNodes, err, &node}
			outChan <- output
		}()
	}

	return outChan
}
