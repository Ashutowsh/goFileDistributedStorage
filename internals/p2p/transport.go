package p2p

type Transport interface {
	Addr() string
	Dial(addr string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}

// Peer represents a node in the network.
type Peer interface {
	Send(b []byte) error
	CloseStream()
}

// RPC represents a Remote Procedure Call message.
type RPC struct {
	From   string
	Data   []byte
	Stream bool
}
