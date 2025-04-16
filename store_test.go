package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	sOpts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(sOpts)
	data := bytes.NewReader([]byte("a jpg image2"))
	if err := s.writeStream("mypic", data); err != nil {
		t.Error(err)
	}
}
