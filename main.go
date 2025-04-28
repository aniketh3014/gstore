package main

import (
	"bytes"
	"log"
	"time"

	"github.com/aniketh3014/distributed-file-storage/p2p"
)

func makeServer(listenAddr string, root string, nodes ...string) *FileServer {

	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		Shakehands: p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
		//TODO: OnPeer
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootStapNodes:     nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s

}

func main() {
	s1 := makeServer(":4000", "ag2", "")
	s2 := makeServer(":5000", "ag3", ":4000")
	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	s.Stop()
	// }()

	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(2 * time.Second)
	go func() {
		log.Fatal(s2.Start())
	}()

	time.Sleep(time.Second * 2)
	data := bytes.NewReader([]byte("this is a big data file!"))

	err := s2.StoreData("key", data)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
