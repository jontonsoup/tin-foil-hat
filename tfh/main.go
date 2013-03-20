// Copyright (c) 2013, gnarmis, jontonsoup, maxsnew
// All rights reserved.

// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

// Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
// Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
			fmt.Println(str)
		} else {
			fmt.Println(err)
		}
	}
}
