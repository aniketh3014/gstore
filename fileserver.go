package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aniketh3014/distributed-file-storage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootStapNodes     []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.RWMutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstarpNetwork()
	s.loop()

	return nil
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	buf := new(bytes.Buffer)

	tee := io.TeeReader(r, buf)

	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}

	msgbuf := new(bytes.Buffer)
	msg := Message{
		MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	if err := gob.NewEncoder(msgbuf).Encode(msg); err != nil {
		return err
	}
	for _, peer := range s.peers {
		if err := peer.Send(msgbuf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(3 * time.Second)
	for _, peer := range s.peers {
		n, err := io.Copy(peer, buf)
		if err != nil {
			return nil
		}
		fmt.Println(n)
	}

	return nil
	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)
	//
	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }
	//
	// p := &DataMessage{
	// 	Key:  key,
	// 	Data: buf.Bytes(),
	// }
	//
	// fmt.Println(p)
	// return s.broadcast(&Message{
	// 	From:    "TODO",
	// 	Payload: p,
	// })
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	fmt.Printf("connected with remote peer %s\n", p.RemoteAddr())

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) bootstarpNetwork() error {
	for _, addr := range s.BootStapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("stopping file server due to user quit signal")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println(err)
				return
			}
			if err := s.handleMessage(rpc.From.String(), &msg); err != nil {
				log.Println(err)
				return
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, &v)
		// default:
		// 	fmt.Printf("%+v", v)
	}
	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg *MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer %s was not found", from)
	}
	fmt.Println(msg.Size)
	if _, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); err != nil {
		return err
	}
	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
}
