package token_bucket_filter

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	size          int
	lastRequestAt time.Time
	tokens        int
	sync.Mutex
}

func New(size int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		size:          size,
		lastRequestAt: time.Now(),
		tokens:        0,
	}
}

func (limiter *TokenBucketLimiter) GetToken() {
	limiter.Lock()
	defer limiter.Unlock()

	limiter.tokens += int(time.Now().Sub(limiter.lastRequestAt).Seconds())
	limiter.tokens = min(limiter.tokens, limiter.size)
	if limiter.tokens == 0 {
		time.Sleep(time.Second)
	} else {
		limiter.tokens--
	}

	limiter.lastRequestAt = time.Now()
	fmt.Println("Granting token at time: " + time.Now().GoString())
}
