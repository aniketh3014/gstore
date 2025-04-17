package main

import (
	"log"
	"time"

	"github.com/aniketh3014/distributed-file-storage/p2p"
)

func main() {

	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr: ":4000",
		Shakehands: p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
		//TODO: OnPeer
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "4000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}
	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
