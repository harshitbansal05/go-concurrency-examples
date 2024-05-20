package token_bucket_filter

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucketLimiterCond struct {
	size          int
	lastRequestAt time.Time
	tokens        int
	rate          int
	cond          sync.Cond
}

func NewCond(size, rate int) *TokenBucketLimiterCond {
	limiter := TokenBucketLimiterCond{
		size:          size,
		lastRequestAt: time.Now(),
		tokens:        0,
		rate:          rate,
		cond:          sync.Cond{L: &sync.Mutex{}},
	}
	return &limiter
}

func (limiter *TokenBucketLimiterCond) GetToken(id int) {
	fmt.Println("Token request received for go-routine: ", id)
	limiter.cond.L.Lock()
	defer limiter.cond.Signal()
	defer limiter.cond.L.Unlock()

	limiter.tokens += int(time.Now().Sub(limiter.lastRequestAt).Seconds() * float64(limiter.rate))
	limiter.tokens = min(limiter.tokens, limiter.size)
	if limiter.tokens == 0 {
		time.Sleep(time.Second / time.Duration(limiter.rate))
	} else {
		limiter.tokens--
	}

	limiter.lastRequestAt = time.Now()
	fmt.Printf("Granting token for id: %d at time: %s\n", id, time.Now().GoString())
}
