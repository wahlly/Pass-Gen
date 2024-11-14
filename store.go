package main

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
)


const filename = "profile.bin"

var (
	ErrInvalidArgs = errors.New("invalid args")
	ErrNotFound = errors.New("not found")
)

type store struct {
	sync.RWMutex
	data map[string]*profile
}

func newStore() (*store, error) {
	s := &store{
		data: make(map[string]*profile),
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *store) load() error {
	flags := os.O_CREATE | os.O_RDONLY
	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		return nil
	}

	return gob.NewDecoder(file).Decode(&s.data)
}

func (s *store) save() error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(s.data)
}

func (s *store) find(platform string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	p, ok := s.data[platform]
	if !ok {
		return "", ErrNotFound
	}

	if err := p.decrypt(); err != nil {
		return "", err
	}

	return p.password, nil
}

func (s *store) add(platform, password string) error {
	if platform == "" {
		return ErrInvalidArgs
	}

	p := &profile{
		Platform: platform,
		password: password,
	}
	
	if err := p.encrypt(); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	s.data[platform] = p

	return s.save()
}