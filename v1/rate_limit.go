package v1

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// NoopLimiter implements Limiter but doesn't limit anything.
var NoopLimiter Limiter = &noopLimiter{}

type token struct {
	rps     atomic.Uint32
	lastUse atomic.Value
}

// Limiter implements some form of rate limiting.
type Limiter interface {
	// Obtain the right to send a request. Should lock the execution if current goroutine needs to wait.
	Obtain(string)
}

// TokensBucket implements basic Limiter with fixed window and fixed amount of tokens per window.
type TokensBucket struct {
	maxRPS          uint32
	tokens          sync.Map
	unusedTokenTime time.Duration
	checkTokenTime  time.Duration
	cancel          atomic.Bool
	sleep           sleeper
}

// NewTokensBucket constructs TokensBucket with provided parameters.
func NewTokensBucket(maxRPS uint32, unusedTokenTime, checkTokenTime time.Duration) Limiter {
	bucket := &TokensBucket{
		maxRPS:          maxRPS,
		unusedTokenTime: unusedTokenTime,
		checkTokenTime:  checkTokenTime,
		sleep:           realSleeper{},
	}

	go bucket.deleteUnusedToken()
	runtime.SetFinalizer(bucket, destructBasket)
	return bucket
}

func (m *TokensBucket) Obtain(id string) {
	val, ok := m.tokens.Load(id)
	if !ok {
		token := &token{}
		token.lastUse.Store(time.Now())
		token.rps.Store(1)
		m.tokens.Store(id, token)
		return
	}

	token := val.(*token)
	sleepTime := time.Second - time.Since(token.lastUse.Load().(time.Time))
	if sleepTime <= 0 {
		token.lastUse.Store(time.Now())
		token.rps.Store(0)
	} else if token.rps.Load() >= m.maxRPS {
		m.sleep.Sleep(sleepTime)
		token.lastUse.Store(time.Now())
		token.rps.Store(0)
	}
	token.rps.Add(1)
}

func destructBasket(m *TokensBucket) {
	m.cancel.Store(true)
}

func (m *TokensBucket) deleteUnusedToken() {
	for {
		if m.cancel.Load() {
			return
		}

		m.tokens.Range(func(key, value any) bool {
			id, token := key.(string), value.(*token)
			if time.Since(token.lastUse.Load().(time.Time)) >= m.unusedTokenTime {
				m.tokens.Delete(id)
			}
			return false
		})

		m.sleep.Sleep(m.checkTokenTime)
	}
}

type noopLimiter struct{}

func (l *noopLimiter) Obtain(string) {}

// sleeper sleeps. This thing is necessary for tests.
type sleeper interface {
	Sleep(time.Duration)
}

type realSleeper struct{}

func (s realSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}
