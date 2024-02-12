package v1

import (
	"errors"
	"time"

	"github.com/maypok86/otter"
)

const mgClientCacheTTL = time.Hour * 1

var NegativeCapacity = errors.New("capacity cannot be less than 1")

type MGClientPool struct {
	cache *otter.CacheWithVariableTTL[string, *MgClient]
}

// NewMGClientPool initializes the client cache
func NewMGClientPool(capacity int) (*MGClientPool, error) {
	if capacity <= 0 {
		return nil, NegativeCapacity
	}

	cache, _ := otter.MustBuilder[string, *MgClient](capacity).WithVariableTTL().Build()
	return &MGClientPool{cache: &cache}, nil
}

func (m *MGClientPool) Get(token string, url string) *MgClient {
	if client, ok := m.cache.Get(token); ok {
		return client
	}

	client := New(url, token)
	m.cache.Set(token, client, mgClientCacheTTL)

	return client
}

func (m *MGClientPool) Remove(token string) {
	m.cache.Delete(token)
}

func (m *MGClientPool) Close() {
	m.cache.Close()
}
