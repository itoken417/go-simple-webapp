package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const (
	idleTimeout     = 30 * time.Minute
	cleanupInterval = 1 * time.Minute
)

type entry struct {
	data      map[string]any
	expiresAt time.Time
}

type store struct {
	mu      sync.Mutex
	entries map[string]*entry
}

var globalStore = newStore()

func newStore() *store {
	s := &store{entries: make(map[string]*entry)}
	go s.cleanup()
	return s
}

func (s *store) cleanup() {
	t := time.NewTicker(cleanupInterval)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		s.mu.Lock()
		for id, e := range s.entries {
			if now.After(e.expiresAt) {
				delete(s.entries, id)
			}
		}
		s.mu.Unlock()
	}
}

func (s *store) get(id string) (*entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[id]
	if !ok || time.Now().After(e.expiresAt) {
		delete(s.entries, id)
		return nil, false
	}
	e.expiresAt = time.Now().Add(idleTimeout)
	return e, true
}

func (s *store) create() (string, *entry) {
	id := generateID()
	e := &entry{
		data:      make(map[string]any),
		expiresAt: time.Now().Add(idleTimeout),
	}
	s.mu.Lock()
	s.entries[id] = e
	s.mu.Unlock()
	return id, e
}

func (s *store) delete(id string) {
	s.mu.Lock()
	delete(s.entries, id)
	s.mu.Unlock()
}

func generateID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("session: failed to generate ID: " + err.Error())
	}
	return hex.EncodeToString(b)
}
