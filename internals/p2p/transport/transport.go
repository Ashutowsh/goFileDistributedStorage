package transport

import (
	"fmt"
	"io"
)

// Transport is the generalized interface for all transport protocols.
type Transport interface {
	Dial(addr string) error
	ListenAndAccept() error
	Close() error
	Consume() <-chan RPC
	Addr() string
}

// Decoder defines the interface for decoding messages.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

// HandshakeFunc is a function type used for handshakes.
type HandshakeFunc func(Peer) error

// Peer represents a peer in the network, which is protocol-agnostic.
type Peer interface {
	Send([]byte) error
	CloseStream()
}

// RPC represents a message being passed between peers.
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}

// TransportOpts holds common options for different transport protocols.
type TransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

func NewTransport(opts TransportOpts, protocol string) (Transport, error) {
	switch protocol {
	case "tcp":
		return NewTCPTransport(opts), nil
	// case "udp":
	// 	return NewUDPTransport(opts), nil
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}
