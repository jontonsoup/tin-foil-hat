package kademlia

// Contains helpers for any functions/types

import (
	"container/list"
	"log"
	"net"
	"strconv"
)

func parseAddress(address string) (ip net.IP, port uint16, err error) {
	hoststr, portstr, err := net.SplitHostPort(address)
	if err != nil {
		return
	}
	ip = net.ParseIP(hoststr)
	port64, err := strconv.ParseUint(portstr, 10, 16)
	if err != nil {
		log.Fatal("parseAddress failed!")
		return
	}
	port = uint16(port64)
	return ip, port, nil
}

func contactToFoundNode(c *Contact) FoundNode {
	// create a FoundNode and return it
	f := new(FoundNode)
	f.IPAddr = c.Host.String()
	f.NodeID = c.NodeID
	f.Port = c.Port
	return *f
}

func foundNodeToContact(f *FoundNode) Contact {
	c := new(Contact)
	ip := net.ParseIP(f.IPAddr)
	c.Host = ip
	c.NodeID = f.NodeID
	c.Port = f.Port
	return *c
}

func insertSorted(inputlist *list.List, item *Contact, greaterThan func(*Contact, *Contact) bool) {
	for e := inputlist.Front(); e != nil; e = e.Next() {
		if greaterThan(e.Value.(*Contact), item) {
			inputlist.InsertBefore(item, e)
			return
		}
	}
	inputlist.PushBack(item)
}

// maxLength should be >= length of original inputList
func insertAllSorted(inputList *list.List, items [](Contact), greaterThan func(*Contact, *Contact) bool, maxLength int) {
	for _, c := range items {
		insertSorted(inputList, &c, greaterThan)
		if inputList.Len() == maxLength {
			inputList.Remove(inputList.Back())
		}
	}
}

// assumes the id is only in the list once
func removeFromSorted(l *list.List, id ID) {
	for e := l.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Contact)
		if c.NodeID.Equals(id) {
			l.Remove(e)
			return
		}
	}
}

func getUnseen(l *list.List, alreadySeen map[ID]bool, max int) []Contact {
	unseen := make([]Contact, 0)
	for e := l.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Contact)
		if !alreadySeen[c.NodeID] {
			unseen = append(unseen, *c)
			if len(unseen) == max {
				break
			}
		}
	}
	return unseen
}
