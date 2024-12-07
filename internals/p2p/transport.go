package p2p

// Transport defines the interface for any transport mechanism (e.g., TCP, UDP, WebSockets).
type Transport interface {
	// Addr returns the listening address of the transport.
	Addr() string

	// Dial connects to a remote address.
	Dial(addr string) error

	// ListenAndAccept starts listening for incoming connections.
	ListenAndAccept() error

	// Consume returns a read-only channel for incoming RPC messages.
	Consume() <-chan RPC

	// Close shuts down the transport and its underlying connections.
	Close() error
}
