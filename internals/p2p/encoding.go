package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder interface for decoding messages from different transports.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

// GOBDecoder implements the Decoder interface using Gob encoding.
type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

// DefaultDecoder implements a basic decoder for message payload.
type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return nil
	}

	if peekBuf[0] == 0x2 {
		msg.Stream = true
		return nil
	}

	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
