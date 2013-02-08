package kademlia

import (
	"net"
	"strconv"
)

// Host identification.
type Contact struct {
	NodeID ID
	Host   net.IP
	Port   uint16
}

func (c *Contact) Address() string {
	return c.Host.String() + ":" + strconv.FormatUint(uint64(c.Port), 10)
}
