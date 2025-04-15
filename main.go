package main

import (
	"fmt"
	"log"

	"github.com/aniketh3014/distributed-file-storage/p2p"
)

func onPeer(p p2p.Peer) error {
	p.Close()
	return nil
}
func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":4000",
		Shakehands: p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
		OnPeer:     onPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()
	select {}
}
