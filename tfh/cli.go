package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

import (
	"kademlia-secure/kademlia"
)

func (tfh *TFH) runCommand(s string) (outStr string, err error) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return
	}

	switch fields[0] {

	case "encrypt":
		if len(fields) != 3 {
			err = errors.New("usage: Path to file | path to keyfile")
			return
		}
		randomKey := makeRandKey(KEY_SIZE)
		outStr, err = tfh.encryptAndStore(fields[1], fields[2], string(randomKey))

		return

	case "decrypt":
		if len(fields) != 3 {
			err = errors.New("usage: path to key | path to output file")
			return
		}
		outStr, err = tfh.decryptAndGet(fields[1], fields[2])
		return

	case "whoami":
		outStr = fmt.Sprintf("%v", tfh.kadem.NodeID.AsString())
		return
	case "local_find_value":
		if len(fields) != 2 {
			err = errors.New("usage: local_find_value key")
			return
		}

		key, errCheck := kademlia.FromString(fields[1])

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid key: %v", fields[1]))
			return
		}

		value, ok := kademlia.LocalLookup(tfh.kadem, key)

		if ok {
			outStr = fmt.Sprintf("%v", string(value))
			return
		} else {
			err = errors.New("LocalLookup error")
			return
		}

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

	case "hashStore":
		if len(fields) != 2 {
			err = errors.New("usage: hashStore val")
			return
		}

		value := []byte(fields[1])

		var key []byte
		key, err = kademlia.HashStore(tfh.kadem, value)
		if err != nil {
			return
		}

		log.Println("Stored value", string(value), "with key", hex.EncodeToString(key))
	default:
		err = errors.New(fmt.Sprintf("Unrecognized command: %v", fields[0]))
	}
	return

}
