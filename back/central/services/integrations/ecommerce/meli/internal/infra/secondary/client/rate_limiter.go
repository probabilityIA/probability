package client

import (
	"context"
	"sync"
	"time"
)

type rateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64
	refillRate float64
	last       time.Time
}

func newRateLimiter(requestsPerMinute float64) *rateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 100
	}
	return &rateLimiter{
		tokens:     requestsPerMinute,
		capacity:   requestsPerMinute,
		refillRate: requestsPerMinute / 60.0,
		last:       time.Now(),
	}
}

func (r *rateLimiter) Wait(ctx context.Context) error {
	for {
		r.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(r.last).Seconds()
		r.last = now
		r.tokens += elapsed * r.refillRate
		if r.tokens > r.capacity {
			r.tokens = r.capacity
		}
		if r.tokens >= 1 {
			r.tokens--
			r.mu.Unlock()
			return nil
		}
		deficit := 1 - r.tokens
		wait := time.Duration(deficit / r.refillRate * float64(time.Second))
		r.mu.Unlock()

		if wait <= 0 {
			wait = 10 * time.Millisecond
		}
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
