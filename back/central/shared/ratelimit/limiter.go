package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
	"golang.org/x/time/rate"
)

type Decision struct {
	Allowed     bool
	Blacklisted bool
	RetryAfter  time.Duration
}

type Limiter interface {
	Check(ctx context.Context, key string) Decision
}

type Config struct {
	RatePerSec  float64
	Burst       int
	Threshold   int
	Escalation  []time.Duration
	RedisPrefix string
}

type entry struct {
	lim         *rate.Limiter
	violations  int
	blockUntil  time.Time
	blCheckedAt time.Time
	lastSeen    time.Time
}

type limiter struct {
	cfg   Config
	redis redis.IRedis
	log   log.ILogger
	mu    sync.Mutex
	items map[string]*entry
}

func New(cfg Config, r redis.IRedis, logger log.ILogger) Limiter {
	if cfg.RatePerSec <= 0 {
		cfg.RatePerSec = 1
	}
	if cfg.Burst <= 0 {
		cfg.Burst = 5
	}
	if cfg.Threshold <= 0 {
		cfg.Threshold = 5
	}
	if len(cfg.Escalation) == 0 {
		cfg.Escalation = []time.Duration{time.Minute, 15 * time.Minute, time.Hour}
	}
	if cfg.RedisPrefix == "" {
		cfg.RedisPrefix = "ratelimit"
	}
	l := &limiter{
		cfg:   cfg,
		redis: r,
		log:   logger.WithModule("ratelimit"),
		items: make(map[string]*entry),
	}
	go l.janitor()
	return l
}

func (l *limiter) getEntry(key string, now time.Time) *entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.items[key]
	if !ok {
		e = &entry{lim: rate.NewLimiter(rate.Limit(l.cfg.RatePerSec), l.cfg.Burst)}
		l.items[key] = e
	}
	e.lastSeen = now
	return e
}

func (l *limiter) Check(ctx context.Context, key string) Decision {
	now := time.Now()
	e := l.getEntry(key, now)

	l.mu.Lock()
	if now.Before(e.blockUntil) {
		ra := e.blockUntil.Sub(now)
		l.mu.Unlock()
		return Decision{Blacklisted: true, RetryAfter: ra}
	}
	needRedisCheck := e.blCheckedAt.IsZero() || now.Sub(e.blCheckedAt) > 5*time.Second
	l.mu.Unlock()

	if needRedisCheck && l.redis != nil {
		if ttl, err := l.redis.TTL(ctx, l.blKey(key)); err == nil && ttl > 0 {
			l.mu.Lock()
			e.blockUntil = now.Add(ttl)
			e.blCheckedAt = now
			l.mu.Unlock()
			return Decision{Blacklisted: true, RetryAfter: ttl}
		}
		l.mu.Lock()
		e.blCheckedAt = now
		l.mu.Unlock()
	}

	if e.lim.Allow() {
		l.mu.Lock()
		e.violations = 0
		l.mu.Unlock()
		return Decision{Allowed: true}
	}

	l.mu.Lock()
	e.violations++
	v := e.violations
	l.mu.Unlock()

	if v == l.cfg.Threshold {
		dur := l.blacklist(ctx, key)
		l.mu.Lock()
		e.blockUntil = now.Add(dur)
		e.violations = 0
		l.mu.Unlock()
		return Decision{Blacklisted: true, RetryAfter: dur}
	}
	if v > l.cfg.Threshold {
		return Decision{Allowed: false, Blacklisted: true, RetryAfter: l.cfg.Escalation[0]}
	}

	res := e.lim.Reserve()
	ra := res.Delay()
	res.Cancel()
	return Decision{Allowed: false, RetryAfter: ra}
}

func (l *limiter) blacklist(ctx context.Context, key string) time.Duration {
	offense := 1
	if l.redis != nil {
		if n, err := l.redis.Incr(ctx, l.offKey(key)); err == nil {
			offense = int(n)
			_ = l.redis.Expire(ctx, l.offKey(key), 24*time.Hour)
		}
	}
	idx := offense - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(l.cfg.Escalation) {
		idx = len(l.cfg.Escalation) - 1
	}
	dur := l.cfg.Escalation[idx]
	if l.redis != nil {
		_ = l.redis.Set(ctx, l.blKey(key), "1", dur)
	}
	l.log.Warn(ctx).Str("key", key).Dur("duration", dur).Int("offense", offense).Msg("rate limit: clave en lista negra")
	return dur
}

func (l *limiter) blKey(key string) string  { return l.cfg.RedisPrefix + ":bl:" + key }
func (l *limiter) offKey(key string) string { return l.cfg.RedisPrefix + ":off:" + key }

func (l *limiter) janitor() {
	t := time.NewTicker(2 * time.Minute)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		cut := now.Add(-10 * time.Minute)
		l.mu.Lock()
		for k, e := range l.items {
			if e.lastSeen.Before(cut) && now.After(e.blockUntil) {
				delete(l.items, k)
			}
		}
		l.mu.Unlock()
	}
}
