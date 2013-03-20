package kademlia

import (
	"container/list"
	"errors"
	"net/rpc"
	"time"
)

const TRPC_WAIT = 30 * time.Second

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

func SendFindNode(k *Kademlia, searchNodeID ID, recipID ID) (nodes []FoundNode, err error) {
	c, ok := LookupContact(k, recipID)
	if !ok {
		err = errors.New("unkown recipient nodeID")
		return
	}
	nodes, err = SendFindNodeAddr(k, searchNodeID, c.Address())
	return
}

// this was capitalized for the workaround
func SendFindNodeAddr(k *Kademlia, nodeID ID, address string) (foundNodes []FoundNode, err error) {
	recvd := make(chan bool, 1)
	go func() {
		foundNodes, err = sendFindNodeAddr(k, nodeID, address)
		recvd <- true
	}()
	select {
	case <-recvd:
	case <-time.After(TRPC_WAIT):
		foundNodes, err = nil, errors.New("timeout")
	}

	if err != nil {
		k.removeContact(nodeID)
	}
	return
}

func sendFindNodeAddr(k *Kademlia, nodeID ID, address string) ([]FoundNode, error) {
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
		k.updateContact(c)
	}
	return res.Nodes, err
}

func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// Find the k closest nodes to FindNodeRequest.NodeID and pack
	// them in FindNodeResult.Nodes
	k.updateContact(req.Sender)

	nodes := k.closestNodes(req.NodeID, req.Sender.NodeID, MAX_BUCKET_SIZE)

	//set the msg id
	res.MsgID = req.MsgID
	res.Nodes = nodes

	return nil
}

func IterativeFindNode(k *Kademlia, searchID ID) ([]Contact, error) {
	shortList := list.New()
	alreadySeen := make(map[ID]bool)

	isCloser := func(c1 Contact, c2 Contact) int {
		d1 := c1.NodeID.Xor(searchID)
		d2 := c2.NodeID.Xor(searchID)
		return d1.Compare(d2)
	}

	initNodes := k.closestContacts(searchID, k.NodeID, ALPHA)
	insertUnseenSorted(shortList, initNodes, isCloser, alreadySeen, MAX_BUCKET_SIZE)

	for {
		nextSearchNodes := getUnseen(shortList, alreadySeen, ALPHA)
		if len(nextSearchNodes) == 0 {
			break
		}

		closestNode := shortList.Front().Value.(Contact)

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
		front := shortList.Front()

		// the shortList is empty, start over, all stale
		// contacts must have been removed by this point
		if front == nil {
			initNodes := k.closestContacts(searchID, k.NodeID, ALPHA)
			insertUnseenSorted(shortList, initNodes, isCloser, alreadySeen, MAX_BUCKET_SIZE)
			continue
		}

		newClosest := front.Value.(Contact)

		if newClosest.NodeID.Equals(closestNode.NodeID) {
			// Didn't find anything closer than the old
			// closest so contact everyone that's left and
			// make sure they're there
			allUnseen := getUnseen(shortList, alreadySeen, shortList.Len())

			unseenNodes := k.goFindNodes(allUnseen, searchID)

			// send find_node rpcs to nextSearchNodes, add their nodes to shortList
			for _, _ = range allUnseen {
				response := <-unseenNodes
				if response.Err != nil {
					removeFromSorted(shortList, response.searchNode.NodeID)
				}
			}
			break
		}
	}
	closestNodes := make([]Contact, 0)

	for e := shortList.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		closestNodes = append(closestNodes, c)
	}

	return closestNodes, nil
}

type SignedFoundNodes struct {
	FoundNodes []FoundNode
	Err        error
	searchNode Contact
}

func (k *Kademlia) goFindNodes(searchNodes []Contact, searchID ID) <-chan SignedFoundNodes {
	outChan := make(chan SignedFoundNodes)

	for _, curNode := range searchNodes {
		node := curNode
		go func() {
			foundNodes, err := SendFindNode(k, searchID, node.NodeID)
			output := SignedFoundNodes{foundNodes, err, node}
			outChan <- output
		}()
	}

	return outChan
}
