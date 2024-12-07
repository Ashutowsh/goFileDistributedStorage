package p2p

import (
	"net"
	"sync"
)

// TCPPeer represents a peer over a TCP connection.
type TCPPeer struct {
	net.Conn
	outbound bool
	wg       *sync.WaitGroup
}

// NewTCPPeer creates a new TCP peer instance.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

// CloseStream signals the end of a stream for this peer.
func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

// Send sends data to the peer.
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}
