package client

import (
	"context"
	"strconv"
	"sync"
	"time"
)

const (
	assumedRateLimitPerMinute = 60
	minRequestInterval        = 50 * time.Millisecond
	maxRequestInterval        = 5 * time.Second
)

type storePacer struct {
	mu          sync.Mutex
	interval    time.Duration
	lastRequest time.Time
}

func newStorePacer() *storePacer {
	return &storePacer{interval: intervalForLimit(assumedRateLimitPerMinute)}
}

func intervalForLimit(perMinute int) time.Duration {
	if perMinute <= 0 {
		return intervalForLimit(assumedRateLimitPerMinute)
	}
	interval := time.Minute / time.Duration(perMinute)
	if interval < minRequestInterval {
		return minRequestInterval
	}
	if interval > maxRequestInterval {
		return maxRequestInterval
	}
	return interval
}

func (p *storePacer) wait(ctx context.Context) error {
	p.mu.Lock()
	interval := p.interval
	last := p.lastRequest
	now := time.Now()

	next := last.Add(interval)
	if last.IsZero() || !next.After(now) {
		p.lastRequest = now
		p.mu.Unlock()
		return nil
	}

	p.lastRequest = next
	p.mu.Unlock()

	timer := time.NewTimer(next.Sub(now))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (p *storePacer) observeLimit(header string) {
	if header == "" {
		return
	}
	limit, err := strconv.Atoi(header)
	if err != nil || limit <= 0 {
		return
	}

	interval := intervalForLimit(limit)

	p.mu.Lock()
	p.interval = interval
	p.mu.Unlock()
}

func (p *storePacer) backOff() {
	p.mu.Lock()
	doubled := p.interval * 2
	if doubled > maxRequestInterval {
		doubled = maxRequestInterval
	}
	p.interval = doubled
	p.mu.Unlock()
}

type pacerRegistry struct {
	mu     sync.Mutex
	pacers map[string]*storePacer
}

func newPacerRegistry() *pacerRegistry {
	return &pacerRegistry{pacers: make(map[string]*storePacer)}
}

func (r *pacerRegistry) forStore(key string) *storePacer {
	r.mu.Lock()
	defer r.mu.Unlock()

	pacer, ok := r.pacers[key]
	if !ok {
		pacer = newStorePacer()
		r.pacers[key] = pacer
	}
	return pacer
}
