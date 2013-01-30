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
	kadem := kademlia.NewKademlia()

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
	client, err := rpc.DialHTTP("tcp", firstPeerStr)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
	ping := new(kademlia.Ping)
	ping.MsgID = kadem.NodeID
	var pong kademlia.Pong
	err = client.Call("Kademlia.Ping", ping, &pong)
	if err != nil {
		log.Fatal("Call: ", err)
	}

	log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
	log.Printf("pong msgID: %s\n", pong.MsgID.AsString())

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
		client, err := rpc.DialHTTP("tcp", fields[1])
		if err != nil {
			log.Fatal("DialHTTP: ", err)
		}
		ping := new(kademlia.Ping)
		ping.MsgID = k.NodeID
		var pong kademlia.Pong
		err = client.Call("Kademlia.Ping", ping, &pong)
		log.Printf("pong msgID: %s\n", pong.MsgID.AsString())
	default:
		fmt.Println("Unrecognized command", fields[0])
	}
	return
}
