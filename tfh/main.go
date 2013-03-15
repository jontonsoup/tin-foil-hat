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
	"time"
)

import (
	"kademlia-secure/kademlia"
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
	tfh := NewTFH(kadem)
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
		str, err := tfh.runCommand(line)
		if err == nil {
			fmt.Println("yolo")
			fmt.Println(str)
		} else {
			fmt.Println(err)
		}
	}
}
