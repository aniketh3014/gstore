package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const DefaultRootFolderName = "agstore"

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	// Root is the folder that contain all the folders/files.
	Root              string
	PathTransformFunc PathTransformFunc
}

type PathKey struct {
	PathName string
	FileName string
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 10
	scliceLen := len(hashStr) / blocksize
	paths := make([]string, scliceLen)

	for i := range scliceLen {
		start, end := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[start:end]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {

	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = DefaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	fileKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, fileKey.PathName)
	fullPath := pathNameWithRoot + "/" + fileKey.FileName
	return os.Open(fullPath)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	firstFolder := strings.Split(pathKey.PathName, "/")[0]
	firstFolderNameWithRoot := fmt.Sprintf("%s/%s", s.Root, firstFolder)
	defer func() {
		log.Printf("deleted %s from disk", pathKey.FileName)
	}()
	fmt.Println(firstFolderNameWithRoot)
	return os.RemoveAll(firstFolderNameWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fileName := pathKey.FileName
	fullPathName := pathNameWithRoot + "/" + fileName

	f, err := os.Create(fullPathName)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullPathName)

	return nil
}
