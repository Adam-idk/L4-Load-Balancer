package main

import (
	"sync"
	"sync/atomic"
)

type Backend struct {
	URL   string
	Alive bool
	mux   sync.RWMutex
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func (s *ServerPool) AddBackend(b *Backend) {
	s.backends = append(s.backends, b)
}

func (s *ServerPool) GetNextValidPeer() *Backend {
	l := uint64(len(s.backends))
	if l == 0 {
		return nil
	}

	next := atomic.AddUint64(&s.current, 1)

	for i := uint64(0); i < l; i++ {
		idx := (next + i) % l
		if s.backends[idx].IsAlive() {
			return s.backends[idx]
		}
	}

	return nil
}

func NewBackend(rawURL string) *Backend {
	return &Backend{
		URL:   rawURL,
		Alive: true,
	}
}

func NewServerPool() *ServerPool {
	return &ServerPool{
		backends: make([]*Backend, 0),
		current:  0,
	}
}
