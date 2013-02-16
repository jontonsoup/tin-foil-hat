package main

import (
	"bufio"
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

	// THis was the workaround, Ted.
	//	_, err = kademlia.SendFindNodeAddr(kadem, kadem.NodeID, firstPeerStr)

	//	var foundNodes []kademlia.Contact
	//	if err == nil {
	foundNodes, err := kademlia.IterativeFindNode(kadem, kadem.NodeID)
	//	}

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
		runCommand(kadem, line)
	}
}

func runCommand(k *kademlia.Kademlia, s string) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return
	}

	switch fields[0] {

	case "whoami":
		fmt.Printf("%v\n", k.NodeID.AsString())
	case "local_find_value":
		if len(fields) != 2 {
			fmt.Println("usage: local_find_value key")
			return
		}

		key, err := kademlia.FromString(fields[1])

		if err != nil {
			fmt.Println("Invalid key: ", fields[1])
			return
		}

		value, ok := kademlia.LocalLookup(k, key)

		if ok {
			fmt.Printf("%v\n", string(value))

		} else {
			fmt.Printf("ERR\n")
		}
	case "get_contact":
		if len(fields) != 2 {
			fmt.Println("usage: get_contact ID")
			return
		}

		nodeID, err := kademlia.FromString(fields[1])

		if err != nil {
			fmt.Println("Invalid nodeID:", fields[1])
			return
		}

		c, ok := kademlia.LookupContact(k, nodeID)

		if ok {
			fmt.Printf("%v %v\n", c.Host, c.Port)
			return
		}
		fmt.Println("ERR")
	case "ping":
		var address string

		if len(fields) != 2 {
			fmt.Println("usage: ping [ip:port | NodeID]")
			return
		}
		localhostfmt := strings.Contains(fields[1], ":")
		if localhostfmt {
			address = fields[1]
		} else {
			id, err := kademlia.FromString(fields[1])
			if err != nil {
				fmt.Println("usage: ping [ip:port | NodeID]")
				return
			}
			if c, ok := kademlia.LookupContact(k, id); !ok {
				fmt.Println("Node not found")
				return
			} else {
				address = c.Host.String() + ":" +
					strconv.FormatUint(uint64(c.Port), 10)
			}

		}
		pong, err := kademlia.SendPing(k, address)
		if err != nil {
			fmt.Println("Ping error:", err)
			return
		}
		fmt.Printf("pong msgID: %v\n", pong.MsgID.AsString())
	case "iterativeFindNode":
		if len(fields) != 2 {
			fmt.Println("usage: iterativeFindNode key")
			return
		}

		id, err := kademlia.FromString(fields[1])
		if err != nil {
			fmt.Println("Invalid NodeID: ", fields[1])
			return
		}

		contacts, err := kademlia.IterativeFindNode(k, id)

		if err != nil {
			fmt.Println("Iterative find node error:", err)
			return
		}

		ids := make([]string, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID.AsString()
		}

		fmt.Printf("%v\n", ids)
	case "iterativeStore":
		if len(fields) != 3 {
			fmt.Println("usage: iterativeStore key value")
			return
		}

		key, err := kademlia.FromString(fields[1])
		if err != nil {
			fmt.Println("Invalid Key: ", fields[1])
			return
		}

		value := []byte(fields[2])

		lastID, err := kademlia.IterativeStore(k, key, value)

		if err == nil {
			fmt.Println(lastID.AsString())
		}
	case "iterativeFindValue":
		// TODO
		// fmt.Println("NOT IMPLEMENTED")
		if len(fields) != 2 {
			fmt.Println("usage: iterativeFindValue key")
		}

		searchID, err := kademlia.FromString(fields[1])
		if err != nil {
			fmt.Println("Invalid NodeID: ", fields[1])
			return
		}

		findValResult, err := kademlia.IterativeFindValue(k, searchID)

		if err != nil {
			fmt.Println("Iterative find value error:", err)
			return
		}

		// if value exists, print it and return
		if findValResult.Value != nil {
			fmt.Printf("%v %v\n", searchID.AsString(), string(findValResult.Value))
			return
		}

		contacts := findValResult.Nodes
		// otherwise print all the found nodes
		ids := make([]string, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID.AsString()
		}

		fmt.Printf("%v\n", ids)

	case "store":
		if len(fields) != 4 {
			fmt.Println("usage: store nodeID key value")
			return
		}
		nodeID, err := kademlia.FromString(fields[1])
		if err != nil {
			fmt.Println("Invalid nodeID: ", fields[1])
			return
		}
		key, err := kademlia.FromString(fields[2])
		if err != nil {
			fmt.Println("Invalid key: ", fields[2])
			return
		}
		value := []byte(fields[3])

		err = kademlia.SendStore(k, key, value, nodeID)
		if err == nil {
			fmt.Println()
		} else {
			fmt.Println("ERR")
		}

	case "find_node":
		if len(fields) != 3 {
			fmt.Println("usage: find_node nodeID searchID")
			return
		}
		nodeID, err := kademlia.FromString(fields[1])
		if err != nil {
			fmt.Println("Invalid nodeID:", fields[1])
			return
		}
		searchID, err := kademlia.FromString(fields[2])
		if err != nil {
			fmt.Println("Invalid nodeID:", fields[2])
			return
		}

		nodes, err := kademlia.SendFindNode(k, searchID, nodeID)
		if err != nil {
			fmt.Println("ERR:", err)
			return
		}
		ids := make([]string, len(nodes))

		for i, node := range nodes {
			ids[i] = node.NodeID.AsString()
		}

		fmt.Printf("%v\n", ids)

	case "find_value":
		if len(fields) != 3 {
			fmt.Println("usage: find_value nodeID key")
			return
		}
		nodeID, err := kademlia.FromString(fields[1])
		if err != nil {
			fmt.Println("Invalid nodeID:", fields[1])
			return
		}
		key, err := kademlia.FromString(fields[2])
		if err != nil {
			fmt.Println("Invalid nodeID:", fields[2])
			return
		}

		findValResult, err := kademlia.SendFindValue(k, key, nodeID)

		if err != nil {
			fmt.Println("error:", err)
			return
		}
		if findValResult.Value != nil {
			fmt.Printf("%v %v\n", nodeID.AsString(), string(findValResult.Value))
			return
		}

		contacts := findValResult.Nodes
		// otherwise print all the found nodes
		ids := make([]string, len(contacts))

		for i, c := range contacts {
			ids[i] = c.NodeID.AsString()
		}

		fmt.Printf("%v\n", ids)

	default:
		fmt.Println("Unrecognized command", fields[0])
	}
	return
}
