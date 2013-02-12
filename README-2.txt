                                                                     
                                                                     
                                                                     
                                             
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
them, and printing their results to stdout. All data should be printed with
the %v specifier and should be followed by exactly one newline. You may assume
values are alphanumeric and are no more than 4095 B. All operations should
complete within 10 seconds.

Implement the following commands:

whoami
    Print your node ID.

local_find_value key
    If your node has data for the given key, print it.
    If your node does not have data for the given key, you should print "ERR".

get_contact ID
    If your buckets contain a node with the given ID,
        printf("%v %v\n", theNode.addr, theNode.port)
    If your buckers do not contain any such node, print "ERR".

iterativeStore key value
    Perform the iterativeStore operation and then print the ID of the node that
    received the final STORE operation.

iterativeFindNode ID
    Print a list of ≤ k closest nodes and print their IDs. You should collect
    the IDs in a slice and print that.

iterativeFindValue key
    printf("%v %v\n", ID, value), where ID refers to the node that finally
    returned the value. If you do not find a value, print "ERR".

// The following four commands cause your code to invoke the appropriate RPC on
// another node, specified by the nodeID argument.
ping nodeID
ping host:port
    Perform a ping. 

store nodeID key value 
    Perform a store and print a blank line.

find_node nodeID key
    Perform a find_node and print its results as for iterativeFindNode.

find_value nodeID key
    Perform a find_value. If it returns nodes, print them as for find_node. If
    it returns a value, print the value as in iterativeFindValue.

