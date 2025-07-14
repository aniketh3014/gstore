package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
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

func (s *Store) Has(key string) bool {
	fileKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, fileKey.PathName)
	fullPath := pathNameWithRoot + "/" + fileKey.FileName

	_, err := os.Stat(fullPath)
	return !errors.Is(err, os.ErrNotExist)
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

func (s *Store) Read(key string) (int64, io.ReadCloser, error) {
	return s.readStream(key)
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	n, err := s.writeStream(key, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	fileKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, fileKey.PathName)
	fullPath := pathNameWithRoot + "/" + fileKey.FileName
	file, err := os.Open(fullPath)
	if err != nil {
		return 0, nil, err
	}

	f, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return f.Size(), file, nil
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

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	fileName := pathKey.FileName
	fullPathName := pathNameWithRoot + "/" + fileName

	f, err := os.Create(fullPathName)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}
