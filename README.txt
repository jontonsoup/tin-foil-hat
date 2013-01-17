************
* BUILDING *
************

Go's build tools depend on the value of the GOPATH environment variable. $GOPATH
should be the project root: the absolute path of the directory containing
{bin,pkg,src}.

Once you've set that, you should be able to build the skeleton and create an
executable at bin/main with:

    go install main

Running main as

    main localhost:7890 localhost:7890

will cause it to start up a server bound to localhost:7890 (the first argument)
and then connect as a client to itself (the second argument). All it does by
default is perform a PING RPC and exit.



**************************
* COMMAND-LINE INTERFACE *
**************************

As demonstrated above, your program must accept two positional arguments of the
form "host:port". The first tells it the bind address of its own server; the
second gives the first peer your client should connect to to join the network.

After setting up its server and establishing a connection to its first peer,
your executable should loop forever, reading commands from stdin, executing
them, and printing their results to stdout. Valid commands are defined below
(you may do anything with an invalid command):

// RPC commands.
// On an error - that is, when Call returns non-nil - you should
//     Printf("ERR: %v\n", err).
// On success, you should
//     Printf("OK: %v\n", theAppropriateReturnType)
ping host:port
ping nodeID
    Execute a PING RPC.

store key data
    Execute a STORE RPC. You may assume the data is human-readable and contains
    no whitespace.

find_node key
    Execute a FIND_NODE RPC. 

find_value key
    Execute a FIND_VALUE RPC.


// Local commands. Your code should serve these without interacting with the DHT
// at all.
get_node_id
    You should Printf("OK: %v\n", yourNodeID)

get_local_value key
    If your node has data for the given key, you should
        Printf("OK: %v\n", theAppropriateData)
    If your node does not have data for the given key, you should
        Printf("ERR\n")

