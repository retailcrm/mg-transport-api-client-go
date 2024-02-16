package v1

import (
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/suite"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TokensBucketTest struct {
	suite.Suite
}

func TestTokensBucket(t *testing.T) {
	suite.Run(t, new(TokensBucketTest))
}

func (t *TokensBucketTest) Test_NewTokensBucket() {
	t.Assert().NotNil(NewTokensBucket(10, time.Hour, time.Hour))
}

func (t *TokensBucketTest) new(
	maxRPS uint32, unusedTokenTime, checkTokenTime time.Duration, sleeper sleeper) *TokensBucket {
	bucket := &TokensBucket{
		maxRPS:          maxRPS,
		unusedTokenTime: unusedTokenTime,
		checkTokenTime:  checkTokenTime,
		sleep:           sleeper,
	}
	runtime.SetFinalizer(bucket, destructBasket)
	return bucket
}

func (t *TokensBucketTest) Test_Obtain_NoThrottle() {
	tb := t.new(100, time.Hour, time.Minute, &realSleeper{})
	start := time.Now()
	for i := 0; i < 100; i++ {
		tb.Obtain("a")
	}
	t.Assert().True(time.Since(start) < time.Second) // check that rate limiter did not perform throttle.
}

func (t *TokensBucketTest) Test_Obtain_Sleep() {
	clock := &fakeSleeper{}
	tb := t.new(100, time.Hour, time.Minute, clock)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 301; i++ {
			tb.Obtain("a")
		}
		wg.Done()
	}()

	wg.Wait()
	t.Assert().Equal(3, int(clock.total.Load()))
}

func (t *TokensBucketTest) Test_Obtain_AddRPS() {
	clock := clockwork.NewFakeClock()
	tb := t.new(100, time.Hour, time.Minute, clock)
	go tb.deleteUnusedToken()
	tb.Obtain("a")
	clock.Advance(time.Minute * 2)

	item, found := tb.tokens.Load("a")
	t.Require().True(found)
	t.Assert().Equal(1, int(item.(*token).rps.Load()))
	tb.Obtain("a")
	t.Assert().Equal(2, int(item.(*token).rps.Load()))
}

type fakeSleeper struct {
	total atomic.Uint32
}

func (s *fakeSleeper) Sleep(time.Duration) {
	s.total.Add(1)
}
