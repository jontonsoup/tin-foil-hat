package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

import (
	"kademlia-go/kademlia"
)

func main() {
	// By default, Go seeds its RNG with 1. This would cause every program to
	// generate the same sequence of IDs.
	rand.Seed(time.Now().UnixNano())

	// Get the bind and connect connection strings from command-line arguments.
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Must be invoked with exactly two arguments!\n")
	}
	listenStr := args[0]
	firstPeerStr := args[1]

	fmt.Printf("kademlia starting up!\n")
	kadem := kademlia.NewKademlia(listenStr)

	rpc.Register(kadem)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", listenStr)
	if err != nil {
		log.Fatal("Listen: ", err)
	}

	// Serve forever.
	go http.Serve(l, nil)

	// Confirm our server is up with a PING request and then exit.
	// Your code should loop forever, reading instructions from stdin and
	// printing their results to stdout. See README.txt for more details.
	pong, err := kademlia.SendPing(kadem, firstPeerStr)
	if err != nil {
		log.Fatal("Initial ping error: ", err)
	}

	fmt.Printf("pong msgID: %s\n", pong.MsgID.AsString())

	_, err = kademlia.SendFindNodeAddr(kadem, kadem.NodeID, firstPeerStr)

	var foundNodes []kademlia.Contact
	if err == nil {
		foundNodes, err = kademlia.IterativeFindNode(kadem, kadem.NodeID)
	}

	if err != nil {
		log.Fatal("Bootstrap find_node error: ", err)
	}

	fmt.Println("Received", len(foundNodes), "nodes")
	for i, node := range foundNodes {
		fmt.Println("Node ", i, ": ", node.NodeID.AsString())
	}

	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		str, err := runCommand(kadem, line)
		if err == nil {
			fmt.Println(str)
		} else {
			fmt.Println(err)
		}
	}
}

func runCommand(k *kademlia.Kademlia, s string) (outStr string, err error) {
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

		key, errCheck := kademlia.FromString(fields[1])

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid key: %v", fields[1]))
			return
		}

		value, ok := kademlia.LocalLookup(k, key)

		if ok {
			outStr = fmt.Sprintf("%v", string(value))
			return
		} else {
			err = errors.New("kademlia.LocalLookup error")
			return
		}
	case "get_contact":
		if len(fields) != 2 {
			err = errors.New("usage: get_contact ID")
			return
		}

		nodeID, errCheck := kademlia.FromString(fields[1])

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}

		c, ok := kademlia.LookupContact(k, nodeID)

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
			id, errCheck := kademlia.FromString(fields[1])
			if errCheck != nil {
				err = errors.New("usage: ping [ip:port | NodeID]")
				return
			}
			if c, ok := kademlia.LookupContact(k, id); !ok {
				err = errors.New("Node not found")
				return
			} else {
				address = c.Host.String() + ":" +
					strconv.FormatUint(uint64(c.Port), 10)
			}

		}
		pong, errCheck := kademlia.SendPing(k, address)
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

		id, errCheck := kademlia.FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid NodeID: %v", fields[1]))
			return
		}

		contacts, errCheck := kademlia.IterativeFindNode(k, id)

		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Iterative find node error: %v", err))
			return
		}

		ids := make([]kademlia.ID, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)
	case "iterativeStore":
		if len(fields) != 3 {
			err = errors.New("usage: iterativeStore key value")
			return
		}

		key, errCheck := kademlia.FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid Key: %v", fields[1]))
			return
		}

		value := []byte(fields[2])

		lastID, errCheck := kademlia.IterativeStore(k, key, value)

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

		searchID, errCheck := kademlia.FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid NodeID: %v", fields[1]))
			return
		}

		findValResult, errCheck := kademlia.IterativeFindValue(k, searchID)

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
		ids := make([]kademlia.ID, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)

	case "store":
		if len(fields) != 4 {
			err = errors.New("usage: store nodeID key value")
			return
		}
		nodeID, errCheck := kademlia.FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}
		key, errCheck := kademlia.FromString(fields[2])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid key: %v", fields[2]))
			return
		}
		value := []byte(fields[3])

		errCheck = kademlia.SendStore(k, key, value, nodeID)
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
		nodeID, errCheck := kademlia.FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}
		searchID, errCheck := kademlia.FromString(fields[2])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[2]))
			return
		}

		nodes, errCheck := kademlia.SendFindNode(k, searchID, nodeID)
		if errCheck != nil {
			return
		}
		ids := make([]kademlia.ID, len(nodes))

		for i, node := range nodes {
			ids[i] = node.NodeID
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)

	case "find_value":
		if len(fields) != 3 {
			err = errors.New("usage: find_value nodeID key")
			return
		}
		nodeID, errCheck := kademlia.FromString(fields[1])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[1]))
			return
		}
		key, errCheck := kademlia.FromString(fields[2])
		if errCheck != nil {
			err = errors.New(fmt.Sprintf("Invalid nodeID: %v", fields[2]))
			return
		}

		findValResult, errCheck := kademlia.SendFindValue(k, key, nodeID)

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
		ids := make([]kademlia.ID, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID
		}

		// TODO: make sure this is right, it looks really dumb
		outStr = fmt.Sprintf("%v", ids)

	default:
		err = errors.New(fmt.Sprintf("Unrecognized command: %v", fields[0]))
	}
	return
}
