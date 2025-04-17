package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestStore(t *testing.T) {
	sOpts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(sOpts)
	for i := range 50 {
		key := fmt.Sprintf("ag_%d", i)
		data := []byte("a jpg image2")

		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("expected %s got %s", data, b)
		}
		fmt.Println(string(b))
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if err := s.Clear(); err != nil {
			t.Error(err)
		}
	}
}
