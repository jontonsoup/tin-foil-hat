package kademlia

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func RunCommand(k *Kademlia, s string) (outStr string, err error) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		err = errors.New("You need some fields for runCommand")
		return
	}

	switch fields[0] {

	case "whoami":
		outStr = fmt.Sprintf("%v", k.NodeID.AsString())
		return
	case "local_find_value":
		if len(fields) != 2 {
			err = errors.New("usage: local_find_value key")
			return
		}

		key, errCheck := FromString(fields[1])

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid key: %v", fields[1]))
			return
		}

		value, ok := LocalLookup(k, key)

		if ok {
			outStr = fmt.Sprintf("%v", string(value))
			return
		} else {
			err = errors.New("LocalLookup error")
			return
		}
	case "get_contact":
		if len(fields) != 2 {
			err = errors.New("usage: get_contact ID")
			return
		}

		nodeID, errCheck := FromString(fields[1])

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}

		c, ok := LookupContact(k, nodeID)

		if ok {
			outStr = fmt.Sprintf("%v %v", c.Host, c.Port)
			return
		} else {
			err = errors.New("kademlia.LookupContact error")
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
			id, errCheck := FromString(fields[1])
			if errCheck != nil {
				err = errors.New("usage: ping [ip:port | NodeID]")
				return
			}
			if c, ok := LookupContact(k, id); !ok {
				err = errors.New("Node not found")
				return
			} else {
				address = c.Host.String() + ":" +
					strconv.FormatUint(uint64(c.Port), 10)
			}

		}
		pong, errCheck := SendPing(k, address)
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Ping error: %v", err))
			return
		} else {
			outStr = fmt.Sprintf("pong msgID: %v", pong.MsgID.AsString())
		}
	case "iterativeFindNode":
		if len(fields) != 2 {
			err = errors.New("usage: iterativeFindNode key")
			return
		}

		id, errCheck := FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid NodeID: %v", fields[1]))
			return
		}

		contacts, errCheck := IterativeFindNode(k, id)

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Iterative find node error: %v", err))
			return
		}

		ids := make([]string, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID.AsString()
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)
	case "iterativeStore":
		if len(fields) != 3 {
			err = errors.New("usage: iterativeStore key value")
			return
		}

		key, errCheck := FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid Key: %v", fields[1]))
			return
		}

		value := []byte(fields[2])

		lastID, errCheck := IterativeStore(k, key, value)

		if errCheck == nil {
			outStr = lastID.AsString()
			return
		}
	case "iterativeFindValue":
		// TODO
		// fmt.Println("NOT IMPLEMENTED")
		if len(fields) != 2 {
			err = errors.New("usage: iterativeFindValue key")
			return
		}

		searchID, errCheck := FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid NodeID: %v", fields[1]))
			return
		}

		findValResult, errCheck := IterativeFindValue(k, searchID)

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Iterative find value error: %v", err))
			return
		}

		// if value exists, print it and return
		if findValResult.Value != nil {
			outStr = fmt.Sprintf("%v %v", searchID.AsString(), string(findValResult.Value))
			return
		}

		contacts := findValResult.Nodes
		// otherwise print all the found nodes
		ids := make([]string, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID.AsString()
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)

	case "store":
		if len(fields) != 4 {
			err = errors.New("usage: store nodeID key value")
			return
		}
		nodeID, errCheck := FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}
		key, errCheck := FromString(fields[2])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid key: %v", fields[2]))
			return
		}
		value := []byte(fields[3])

		errCheck = SendStore(k, key, value, nodeID)
		if errCheck == nil {
			// store shouldn't return anything, as per spec
			return
		} else {
			err = errors.New("ERR")
			return
		}

	case "find_node":
		if len(fields) != 3 {
			err = errors.New("usage: find_node nodeID searchID")
			return
		}
		nodeID, errCheck := FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}
		searchID, errCheck := FromString(fields[2])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[2]))
			return
		}

		nodes, errCheck := SendFindNode(k, searchID, nodeID)
		if errCheck != nil {
			return
		}
		ids := make([]string, len(nodes))

		for i, node := range nodes {
			ids[i] = node.NodeID.AsString()
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)

	case "find_value":
		if len(fields) != 3 {
			err = errors.New("usage: find_value nodeID key")
			return
		}
		nodeID, errCheck := FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}
		key, errCheck := FromString(fields[2])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[2]))
			return
		}

		findValResult, errCheck := SendFindValue(k, key, nodeID)

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("error: %v", err))
			return
		}
		if findValResult.Value != nil {
			outStr = fmt.Sprintf("%v %v", nodeID.AsString(), string(findValResult.Value))
			return
		}

		contacts := findValResult.Nodes
		// otherwise print all the found nodes
		ids := make([]string, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID.AsString()
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)

	default:
		err = errors.New(fmt.Sprintf("Unrecognized command: %v", fields[0]))
	}
	return
}
