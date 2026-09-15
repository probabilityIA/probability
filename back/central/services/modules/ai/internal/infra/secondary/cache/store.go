package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

func usageKey(userID uint) string {
	return fmt.Sprintf("%s:%d", usagePrefix, userID)
}

func introKey(userID uint) string {
	return fmt.Sprintf("%s:%d", introPrefix, userID)
}

func isNotFound(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "key not found")
}

func (s *Store) ConsumeMessage(ctx context.Context, userID uint, window time.Duration) (*entities.Usage, error) {
	key := usageKey(userID)
	count, err := s.redis.Incr(ctx, key)
	if err != nil {
		return nil, err
	}

	ttl, err := s.redis.TTL(ctx, key)
	if err != nil {
		return nil, err
	}
	if ttl <= 0 {
		if err := s.redis.Expire(ctx, key, window); err != nil {
			return nil, err
		}
		ttl = window
	}

	resetAt := time.Now().Add(ttl)
	return &entities.Usage{Count: int(count), ResetAt: &resetAt}, nil
}

func (s *Store) GetUsage(ctx context.Context, userID uint) (*entities.Usage, error) {
	key := usageKey(userID)
	raw, err := s.redis.Get(ctx, key)
	if isNotFound(err) {
		return &entities.Usage{}, nil
	}
	if err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(raw)
	if err != nil {
		return &entities.Usage{}, nil
	}

	usage := &entities.Usage{Count: count}
	if ttl, err := s.redis.TTL(ctx, key); err == nil && ttl > 0 {
		resetAt := time.Now().Add(ttl)
		usage.ResetAt = &resetAt
	}
	return usage, nil
}

func (s *Store) IsIntroSeen(ctx context.Context, userID uint) (bool, error) {
	raw, err := s.redis.Get(ctx, introKey(userID))
	if isNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return raw != "", nil
}

func (s *Store) MarkIntroSeen(ctx context.Context, userID uint) error {
	return s.redis.Set(ctx, introKey(userID), "1", introTTL)
}
