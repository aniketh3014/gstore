package main

import (
	// "bytes"
	// "bytes"
	"fmt"
	// "fmt"
	"io"
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
	//
	// for i := 0; i < 3; i++ {
	// 	data := bytes.NewReader([]byte("this is a big data file!"))
	//
	// 	err := s2.StoreData(fmt.Sprintf("myprivatekey_%d", i), data)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	time.Sleep(time.Millisecond * 5)
	// }

	r, err := s2.GetData("myprivatekey_1")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	select {}
}
