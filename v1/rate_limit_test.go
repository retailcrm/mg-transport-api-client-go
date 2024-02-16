package v1

import (
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/suite"
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

func (t *TokensBucketTest) Test_Obtain_NoThrottle() {
	tb := NewTokensBucket(100, time.Hour, time.Minute)
	start := time.Now()
	for i := 0; i < 100; i++ {
		tb.Obtain("a")
	}
	t.Assert().True(time.Since(start) < time.Second) // check that rate limiter did not perform throttle.
}

func (t *TokensBucketTest) Test_Obtain_Sleep() {
	clock := &fakeSleeper{}
	tb := NewTokensBucket(100, time.Hour, time.Minute)
	tb.cancel.Store(true) // prevent unused token removal.
	tb.sleep = clock

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
	tb := NewTokensBucket(100, time.Hour, time.Minute)
	tb.sleep = clock
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
