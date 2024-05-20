package token_bucket_filter

import (
	"errors"
	"fmt"
	"time"
)

type TokenBucketLimiterEn struct {
	size          int
	lastRequestAt time.Time
	tokens        int
	rate          int
	queue         chan chan struct{}
	queueSize     int
}

func NewEn(size, queueSize, rate int) *TokenBucketLimiterEn {
	limiter := TokenBucketLimiterEn{
		size:          size,
		lastRequestAt: time.Now(),
		tokens:        0,
		rate:          rate,
		queueSize:     queueSize,
		queue:         make(chan chan struct{}, queueSize),
	}
	go limiter.process()
	return &limiter
}

func (limiter *TokenBucketLimiterEn) GetToken(id int, cb chan struct{}) error {
	fmt.Println("Token request received for go-routine: ", id)
	select {
	case limiter.queue <- cb:
		return nil
	default:
		return errors.New("queue is full")
	}
}

func (limiter *TokenBucketLimiterEn) process() {
	for {
		select {
		case cb := <-limiter.queue:
			limiter.tokens += int(time.Now().Sub(limiter.lastRequestAt).Seconds() * float64(limiter.rate))
			limiter.tokens = min(limiter.tokens, limiter.size)
			if limiter.tokens == 0 {
				time.Sleep(time.Second / time.Duration(limiter.rate))
			} else {
				limiter.tokens--
			}

			limiter.lastRequestAt = time.Now()
			select {
			case cb <- struct{}{}:
			default:
			}
		}
	}
}
