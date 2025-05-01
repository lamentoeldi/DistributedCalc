package memory

import (
	"context"
	"sync"
)

type Blacklist struct {
	m  map[string]struct{}
	mu sync.RWMutex
}

func NewBlacklist() *Blacklist {
	return &Blacklist{
		m: make(map[string]struct{}),
	}
}

func (b *Blacklist) Add(_ context.Context, tokenID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.m[tokenID] = struct{}{}
	return nil
}

func (b *Blacklist) Remove(_ context.Context, tokenID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.m, tokenID)
	return nil
}

func (b *Blacklist) IsBlackListed(_ context.Context, tokenID string) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, ok := b.m[tokenID]
	return ok, nil
}
