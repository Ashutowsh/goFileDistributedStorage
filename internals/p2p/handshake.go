// /internals/p2p/handshake.go
package p2p

import (
	"/internals/p2p/transport" // Import the transport package
)

// HandshakeFunc defines a custom handshake logic between peers.
type HandshakeFunc func(transport.Peer) error

// NOPHandshakeFunc is a no-op function for handshakes, meaning no specific handshake logic.
func NOPHandshakeFunc(transport.Peer) error {
	return nil
}
