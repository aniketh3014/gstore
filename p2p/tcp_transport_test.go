package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOps := TCPTransportOpts{
		ListenAddr: ":4000",
		Decoder:    DefaultDecoder{},
		Shakehands: NOPHandshakeFunc,
	}
	tr := NewTCPTransport(tcpOps)

	assert.Equal(t, tr.ListenAddr, ":4000")
	assert.Nil(t, tr.ListenAndAccept())
}
