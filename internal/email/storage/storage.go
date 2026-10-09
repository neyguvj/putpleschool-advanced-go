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

func (s *Storage) Get(hash string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	email, ok := s.data[hash]
	return email, ok
}
