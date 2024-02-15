package v1

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type token struct {
	rps     uint32
	lastUse time.Time
}

type TokensBucket struct {
	maxRPS          uint32
	mux             sync.Mutex
	tokens          map[string]*token
	unusedTokenTime time.Duration
	checkTokenTime  time.Duration
	cancel          atomic.Bool
}

func NewTokensBucket(maxRPS uint32, unusedTokenTime, checkTokenTime time.Duration) *TokensBucket {
	bucket := &TokensBucket{
		maxRPS:          maxRPS,
		tokens:          map[string]*token{},
		unusedTokenTime: unusedTokenTime,
		checkTokenTime:  checkTokenTime,
	}

	go bucket.deleteUnusedToken()
	runtime.SetFinalizer(bucket, destructBasket)
	return bucket
}

func (m *TokensBucket) Obtain(id string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, ok := m.tokens[id]; !ok {
		m.tokens[id] = &token{
			lastUse: time.Now(),
			rps:     1,
		}
		return
	}

	sleepTime := time.Second - time.Since(m.tokens[id].lastUse)
	if sleepTime < 0 {
		m.tokens[id].lastUse = time.Now()
		m.tokens[id].rps = 0
	} else if m.tokens[id].rps >= m.maxRPS {
		time.Sleep(sleepTime)
		m.tokens[id].lastUse = time.Now()
		m.tokens[id].rps = 0
	}
	m.tokens[id].rps++
}

func destructBasket(m *TokensBucket) {
	m.cancel.Store(true)
}

func (m *TokensBucket) deleteUnusedToken() {
	for {
		if m.cancel.Load() {
			return
		}
		m.mux.Lock()

		for id, token := range m.tokens {
			if time.Since(token.lastUse) >= m.unusedTokenTime {
				delete(m.tokens, id)
			}
		}
		m.mux.Unlock()

		time.Sleep(m.checkTokenTime)
	}
}
