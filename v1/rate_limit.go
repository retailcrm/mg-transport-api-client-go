package v1

import (
	"hash/fnv"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var NoopLimiter Limiter = &noopLimiter{}

type token struct {
	rps     uint32
	lastUse int64 // Unix timestamp in nanoseconds
}

// Limiter interface for rate-limiting.
type Limiter interface {
	Obtain(id string)
}

// TokensBucket implements a sharded rate limiter with fixed window and tokens.
type TokensBucket struct {
	maxRPS          uint32
	unusedTokenTime int64 // in nanoseconds
	checkTokenTime  time.Duration
	shards          []*tokenShard
	shardCount      uint32
	cancel          atomic.Bool
	sleep           sleeper
}

type tokenShard struct {
	tokens map[string]*token
	mu     sync.Mutex
}

// NewTokensBucket creates a sharded token bucket limiter.
func NewTokensBucket(maxRPS uint32, unusedTokenTime, checkTokenTime time.Duration) Limiter {
	shardCount := uint32(runtime.NumCPU() * 2) // Use double the CPU count for sharding
	shards := make([]*tokenShard, shardCount)
	for i := range shards {
		shards[i] = &tokenShard{tokens: make(map[string]*token)}
	}

	bucket := &TokensBucket{
		maxRPS:          maxRPS,
		unusedTokenTime: unusedTokenTime.Nanoseconds(),
		checkTokenTime:  checkTokenTime,
		shards:          shards,
		shardCount:      shardCount,
		sleep:           realSleeper{},
	}

	go bucket.cleanupRoutine()
	runtime.SetFinalizer(bucket, destructBucket)
	return bucket
}

// Obtain request hit. Will throttle RPS.
func (m *TokensBucket) Obtain(id string) {
	shard := m.getShard(id)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	item, exists := shard.tokens[id]
	now := time.Now().UnixNano()

	if !exists {
		shard.tokens[id] = &token{
			rps:     1,
			lastUse: now,
		}
		return
	}

	sleepTime := int64(time.Second) - (now - item.lastUse)
	if sleepTime <= 0 {
		item.lastUse = now
		atomic.StoreUint32(&item.rps, 1)
	} else if atomic.LoadUint32(&item.rps) >= m.maxRPS {
		m.sleep.Sleep(time.Duration(sleepTime))
		item.lastUse = time.Now().UnixNano()
		atomic.StoreUint32(&item.rps, 1)
	} else {
		atomic.AddUint32(&item.rps, 1)
	}
}

func (m *TokensBucket) getShard(id string) *tokenShard {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(id))
	return m.shards[hash.Sum32()%m.shardCount]
}

func (m *TokensBucket) cleanupRoutine() {
	ticker := time.NewTicker(m.checkTokenTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if m.cancel.Load() {
				return
			}
			now := time.Now().UnixNano()
			for _, shard := range m.shards {
				shard.mu.Lock()
				for id, token := range shard.tokens {
					if now-token.lastUse >= m.unusedTokenTime {
						delete(shard.tokens, id)
					}
				}
				shard.mu.Unlock()
			}
		}
	}
}

func destructBucket(m *TokensBucket) {
	m.cancel.Store(true)
}

type noopLimiter struct{}

func (l *noopLimiter) Obtain(string) {}

type sleeper interface {
	Sleep(time.Duration)
}

type realSleeper struct{}

func (s realSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}
