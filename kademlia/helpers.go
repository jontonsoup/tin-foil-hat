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

func InsertSorted(inputlist *list.List, item *Contact, greaterThan func(*Contact, *Contact) bool) {
	for e := inputlist.Front(); e != nil; e = e.Next() {
		if greaterThan(e.Value.(*Contact), item) {
			inputlist.InsertBefore(item, e)
			return
		}
	}
	inputlist.PushBack(item)
}
