package token_bucket_filter

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucketLimiterAlter struct {
	size   int
	tokens int
	sync.Mutex
}

func NewAlter(size int) *TokenBucketLimiterAlter {
	limiter := TokenBucketLimiterAlter{
		size:   size,
		tokens: 0,
	}
	go func() {
		for {
			limiter.Lock()
			limiter.tokens = min(limiter.size, limiter.tokens+1)
			limiter.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()
	return &limiter
}

func (limiter *TokenBucketLimiterAlter) GetToken() {
	limiter.Lock()
	defer limiter.Unlock()

	for limiter.tokens == 0 {
		limiter.Unlock()

		limiter.Lock()
	}

	limiter.tokens--
	fmt.Println("Granting token at time: " + time.Now().GoString())
}
