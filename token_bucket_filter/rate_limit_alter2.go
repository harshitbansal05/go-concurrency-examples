package token_bucket_filter

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucketLimiterAlter2 struct {
	size   int
	tokens int
	cond   sync.Cond
}

func NewAlter2(size int) *TokenBucketLimiterAlter2 {
	limiter := TokenBucketLimiterAlter2{
		size:   size,
		tokens: 0,
		cond:   sync.Cond{L: &sync.Mutex{}},
	}
	go func() {
		for {
			limiter.cond.L.Lock()
			limiter.tokens = min(limiter.size, limiter.tokens+1)
			limiter.cond.L.Unlock()
			limiter.cond.Broadcast()
			time.Sleep(1 * time.Second)
		}
	}()
	return &limiter
}

func (limiter *TokenBucketLimiterAlter2) GetToken() {
	limiter.cond.L.Lock()
	defer limiter.cond.L.Unlock()

	for limiter.tokens == 0 {
		limiter.cond.Wait()
	}

	limiter.tokens--
	fmt.Println("Granting token at time: " + time.Now().GoString())
}
