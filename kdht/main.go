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

	log.Printf("pong msgID: %s\n", pong.MsgID.AsString())

	_, err = kademlia.SendFindNode(kadem, kadem.NodeID, firstPeerStr)

	var foundNodes []kademlia.Contact
	if err == nil {
		foundNodes, err = kademlia.IterativeFindNode(kadem, kadem.NodeID)
	}

	if err != nil {
		log.Fatal("Bootstrap find_node error: ", err)
	}

	log.Println("Received", len(foundNodes), "nodes")
	for i, node := range foundNodes {
		log.Println("Node ", i, ": ", node.NodeID.AsString())
	}

	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		err = runCommand(kadem, line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func runCommand(k *kademlia.Kademlia, s string) (err error) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return nil
	}
	err = nil
	switch fields[0] {
	case "get_node_id":
		fmt.Printf("OK: %s\n", k.NodeID.AsString())
	case "ping":
		var address string

		if len(fields) != 2 {
			log.Println("usage: ping [ip:port | NodeID]")
			return
		}
		localhostfmt := strings.Contains(fields[1], ":")
		if localhostfmt {
			address = fields[1]
		} else {
			id, err := kademlia.FromString(fields[1])
			if err != nil {
				log.Println("usage: ping [ip:port | NodeID]")
				return nil
			}
			if c, ok := kademlia.LookupContact(k, id); !ok {
				log.Println("Node not found")
				return nil
			} else {
				log.Println("Found contact ", c)
				address = c.Host.String() + ":" +
					strconv.FormatUint(uint64(c.Port), 10)
			}

		}
		pong, err := kademlia.SendPing(k, address)
		if err != nil {
			return err
		}
		log.Printf("pong msgID: %s\n", pong.MsgID.AsString())
	case "find_node":
		if len(fields) != 2 {
			log.Println("usage: find_node key")
			return
		}

		id, err := kademlia.FromString(fields[1])
		if err != nil {
			log.Println("Invalid NodeID: ", fields[1])
			return nil
		}

		contacts, err := kademlia.IterativeFindNode(k, id)

		if err != nil {
			return err
		}

		log.Println(len(contacts), "contacts found")
		for _, node := range contacts {
			fmt.Println("Ok: ", node.NodeID.AsString())
		}
	default:
		fmt.Println("Unrecognized command", fields[0])
	}
	return
}
