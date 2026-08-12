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

// MercadoLibre permite 1.500 llamadas por minuto POR VENDEDOR. Se deja margen
// para el resto del trafico de esa cuenta (ordenes, webhooks, envios).
const defaultRequestsPerMinute = 1000

// limiterRegistry entrega un limitador por vendedor. Uno global seria incorrecto:
// el limite de MercadoLibre es por cuenta, asi que un negocio sincronizando no
// debe frenar a los demas.
type limiterRegistry struct {
	mu       sync.Mutex
	limiters map[string]*rateLimiter
}

func newLimiterRegistry() *limiterRegistry {
	return &limiterRegistry{limiters: make(map[string]*rateLimiter)}
}

func (r *limiterRegistry) forKey(key string) *rateLimiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if l, ok := r.limiters[key]; ok {
		return l
	}
	l := newRateLimiter(defaultRequestsPerMinute)
	r.limiters[key] = l
	return l
}

func newRateLimiter(requestsPerMinute float64) *rateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = defaultRequestsPerMinute
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
