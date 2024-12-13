package transport

import (
	"errors"
	"log"
	"net"
	"sync"
)

// TCPTransport implements the Transport interface for TCP communication.
type TCPTransport struct {
	opts     TransportOpts
	listener net.Listener
	rpcCh    chan RPC
}

// NewTCPTransport creates a new instance of TCPTransport.
func NewTCPTransport(opts TransportOpts) *TCPTransport {
	return &TCPTransport{
		opts:  opts,
		rpcCh: make(chan RPC, 1024),
	}
}

// Dial connects to a remote address via TCP.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConnection(conn, true)
	return nil
}

// ListenAndAccept starts the TCP listener and accepts incoming connections.
func (t *TCPTransport) ListenAndAccept() error {
	listener, err := net.Listen("tcp", t.opts.ListenAddr)
	if err != nil {
		return err
	}
	t.listener = listener
	log.Printf("Listening on %s", t.opts.ListenAddr)
	go t.acceptConnections()
	return nil
}

// Close stops the transport and closes the listener.
func (t *TCPTransport) Close() error {
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}

// Consume returns the channel for incoming RPC messages.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

// Addr returns the listening address of the transport.
func (t *TCPTransport) Addr() string {
	return t.opts.ListenAddr
}

func (t *TCPTransport) acceptConnections() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go t.handleConnection(conn, false)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn, outbound bool) {
	defer conn.Close()

	peer := NewTCPPeer(conn, outbound)
	if t.opts.HandshakeFunc != nil {
		if err := t.opts.HandshakeFunc(peer); err != nil {
			log.Printf("Handshake failed: %s", err)
			return
		}
	}

	if t.opts.OnPeer != nil {
		if err := t.opts.OnPeer(peer); err != nil {
			log.Printf("Peer rejected: %s", err)
			return
		}
	}

	for {
		rpc := RPC{}
		if err := t.opts.Decoder.Decode(conn, &rpc); err != nil {
			log.Printf("Error decoding message: %s", err)
			return
		}
		rpc.From = conn.RemoteAddr().String()
		t.rpcCh <- rpc
	}
}

// TCPPeer implements the Peer interface for TCP connections.
type TCPPeer struct {
	conn     net.Conn
	outbound bool
	wg       sync.WaitGroup
}

// NewTCPPeer creates a new TCPPeer instance.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Send writes data to the peer's connection.
func (p *TCPPeer) Send(data []byte) error {
	_, err := p.conn.Write(data)
	return err
}

// CloseStream marks the stream as completed.
func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}
