package storage

import "sync"

type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Put(hash string, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[hash] = email
}

func (s *Storage) Pop(hash string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	email, ok := s.data[hash]
	if !ok {
		return "", false
	}
	delete(s.data, hash)
	return email, true
}
