package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

import (
	"kademlia-secure/kademlia"
)

func (tfh *TFH) runCommand(s string) (outStr string, err error) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		err = errors.New("You need some fields for runCommand")
		return
	}

	switch fields[0] {

	case "encrypt":
		if len(fields) != 3 {
			err = errors.New("usage: Path to file | 32 character key")
			return
		}
		if len(fields[2]) != 32 {
			err = errors.New("ERROR: you must have a 32 character key")
			return
		}
		outStr, err = Encrypt(fields[1], fields[2])
		return

	case "decrypt":
		if len(fields) != 2 {
			err = errors.New("usage: key")
			return
		}
		return

	case "whoami":
		outStr = fmt.Sprintf("%v", tfh.kadem.NodeID.AsString())
		return

	case "ping":
		var address string

		if len(fields) != 2 {
			err = errors.New("usage: ping [ip:port | NodeID]")
			return
		}
		localhostfmt := strings.Contains(fields[1], ":")
		if localhostfmt {
			address = fields[1]
		} else {
			id, errCheck := kademlia.FromString(fields[1])
			if errCheck != nil {
				err = errors.New("usage: ping [ip:port | NodeID]")
				return
			}
			if c, ok := kademlia.LookupContact(tfh.kadem, id); !ok {
				err = errors.New("Node not found")
				return
			} else {
				address = c.Host.String() + ":" +
					strconv.FormatUint(uint64(c.Port), 10)
			}

		}
		pong, errCheck := kademlia.SendPing(tfh.kadem, address)
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Ping error: %v", err))
			return
		} else {
			outStr = fmt.Sprintf("pong msgID: %v", pong.MsgID.AsString())
		}

	default:
		err = errors.New(fmt.Sprintf("Unrecognized command: %v", fields[0]))
	}
	return
}