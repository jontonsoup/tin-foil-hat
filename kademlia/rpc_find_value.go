package kademlia

import (
	"container/list"
	"errors"
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
		k.removeContact(nodeID)
		return
	}
	req := new(FindValueRequest)
	req.MsgID = NewRandomID()
	req.Sender = k.Self
	req.Key = key
	err = client.Call("Kademlia.FindValue", req, &ret)
	if err != nil {
		k.removeContact(nodeID)
		return
	}
	defer client.Close()
	if !ret.MsgID.Equals(ret.MsgID) {
		err = errors.New("FindValue MsgID didn't match SendFindValue MsgID")
	}
	if ret.Value != nil && CorrectHash(key[:], ret.Value) {
		err = errors.New("Bad hash")
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

func IterativeFindValue(k *Kademlia, searchID ID) (FindValueResult, error) {
	val, ok := LocalLookup(k, searchID)
	if ok {
		res := new(FindValueResult)
		res.Value = val
		return *res, nil
	}

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
		// shortlist will be nonempty bc nextSearchNodes is nonempty
		closestNode := shortList.Front().Value.(Contact)

		newNodesChan := k.goFindValue(nextSearchNodes, searchID)

		// send find_node rpcs to nextSearchNodes, add their nodes to shortList
		for _, _ = range nextSearchNodes {
			response := <-newNodesChan

			if response.Err != nil {
				removeFromSorted(shortList, response.searchNode.NodeID)
			}
			// if it's a value, just return it because we're done

			if response.FoundValueResult.Value != nil {
				return response.FoundValueResult, nil
			}

			// mark them alreadySeen
			alreadySeen[response.searchNode.NodeID] = true

			// aggregate the new contacts into the shortList, keeping
			// only the K closest
			newNodes := make([]Contact, 0)

			for _, node := range response.FoundValueResult.Nodes {
				newNodes = append(newNodes, foundNodeToContact(&node))
			}

			insertUnseenSorted(shortList, newNodes, isCloser, alreadySeen, MAX_BUCKET_SIZE)
		}

		front := shortList.Front()

		// shortList is empty, start over with stale nodes gone
		if front == nil {
			initNodes := k.closestContacts(searchID, k.NodeID, ALPHA)
			insertUnseenSorted(shortList, initNodes, isCloser, alreadySeen, MAX_BUCKET_SIZE)
			continue
		}

		newClosest := front.Value.(Contact)

		if newClosest.NodeID.Equals(closestNode.NodeID) {
			// Didn't find anything closer than the old closest so stop
			allUnseen := getUnseen(shortList, alreadySeen, shortList.Len())

			unseenNodes := k.goFindValue(allUnseen, searchID)

			// send find_node rpcs to nextSearchNodes, add their nodes to shortList
			for _, _ = range allUnseen {
				response := <-unseenNodes

				if response.Err != nil {
					removeFromSorted(shortList, response.searchNode.NodeID)
				}
				// if it's a value, just return it because we're done

				if response.FoundValueResult.Value != nil {
					return response.FoundValueResult, nil
				}
			}

			break
		}

	}
	closestNodes := make([]FoundNode, 0)

	for e := shortList.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		f := contactToFoundNode(c)
		closestNodes = append(closestNodes, f)
	}

	findValResult := new(FindValueResult)
	// nobody cares about this, so make it 0
	findValResult.MsgID = *new(ID)
	findValResult.Nodes = closestNodes
	return *findValResult, nil
}

type SignedFindValueResults struct {
	FoundValueResult FindValueResult
	Err              error
	searchNode       Contact
}

func (k *Kademlia) goFindValue(searchNodes []Contact, searchID ID) <-chan SignedFindValueResults {
	outChan := make(chan SignedFindValueResults)

	for _, curNode := range searchNodes {
		node := curNode
		go func() {
			findValResult, err := SendFindValue(k, searchID, node.NodeID)
			output := SignedFindValueResults{*findValResult, err, node}
			outChan <- output
		}()
	}

	return outChan
}
