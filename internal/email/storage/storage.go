package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Storage struct {
	mu   sync.Mutex
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{path: path}
}

func (s *Storage) load() (map[string]string, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	data := make(map[string]string)
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Storage) save(data map[string]string) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Storage) Put(hash string, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}
	data[hash] = email
	return s.save(data)
}

func (s *Storage) Pop(hash string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return "", false, err
	}

	email, ok := data[hash]
	if !ok {
		return "", false, nil
	}

	delete(data, hash)
	if err := s.save(data); err != nil {
		return "", false, err
	}
	return email, true, nil
}
