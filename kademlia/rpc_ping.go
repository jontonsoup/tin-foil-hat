package kademlia

import (
	"errors"
	"log"
	"net/rpc"
)

// PING
type Ping struct {
	Sender Contact
	MsgID  ID
}

type Pong struct {
	MsgID ID
}

func SendPing(k *Kademlia, address string) (pong *Pong, err error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return
	}

	ping := new(Ping)
	ping.MsgID = NewRandomID()
	ping.Sender = k.Self
	err = client.Call("Kademlia.Ping", ping, &pong)
	if err != nil {
		return
	}
	defer client.Close()

	if !pong.MsgID.Equals(ping.MsgID) {
		err = errors.New("Pong MsgID didn't match Ping MsgID")
	}

	return
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
	// This one's a freebie.
	if !ping.Sender.NodeID.Equals(k.NodeID) {
		log.Print("Adding contact: ", ping.Sender.NodeID.AsString())
		k.updateContact(&ping.Sender)
	}
	pong.MsgID = ping.MsgID
	return nil
}
