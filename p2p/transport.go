package p2p

// Peer is an Interface that represents the remote node.
type Peer interface {
	Close() error
}

// Transport is anything that handles the communication between the nodes in the network.
// This can of the form TCP, UDP, Websokets,...
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
