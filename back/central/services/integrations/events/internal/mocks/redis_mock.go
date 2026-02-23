package mocks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisMock es el mock de redisclient.IRedis.
// Todas las funciones son inyectables; los métodos no configurados retornan valores cero.
type RedisMock struct {
	ConnectFn              func(ctx context.Context) error
	CloseFn                func() error
	ClientFn               func(ctx context.Context) *redis.Client
	PingFn                 func(ctx context.Context) error
	GetFn                  func(ctx context.Context, key string) (string, error)
	SetFn                  func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	DeleteFn               func(ctx context.Context, keys ...string) error
	ExistsFn               func(ctx context.Context, keys ...string) (int64, error)
	ExpireFn               func(ctx context.Context, key string, expiration time.Duration) error
	TTLFn                  func(ctx context.Context, key string) (time.Duration, error)
	KeysFn                 func(ctx context.Context, pattern string) ([]string, error)
	IncrFn                 func(ctx context.Context, key string) (int64, error)
	DecrFn                 func(ctx context.Context, key string) (int64, error)
	HGetFn                 func(ctx context.Context, key, field string) (string, error)
	HSetFn                 func(ctx context.Context, key string, values ...interface{}) error
	HGetAllFn              func(ctx context.Context, key string) (map[string]string, error)
	HDelFn                 func(ctx context.Context, key string, fields ...string) error
	RegisterCachePrefixFn  func(prefix string)
	RegisterChannelFn      func(channel string)
}

func (m *RedisMock) Connect(ctx context.Context) error {
	if m.ConnectFn != nil {
		return m.ConnectFn(ctx)
	}
	return nil
}

func (m *RedisMock) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}

func (m *RedisMock) Client(ctx context.Context) *redis.Client {
	if m.ClientFn != nil {
		return m.ClientFn(ctx)
	}
	return nil
}

func (m *RedisMock) Ping(ctx context.Context) error {
	if m.PingFn != nil {
		return m.PingFn(ctx)
	}
	return nil
}

func (m *RedisMock) Get(ctx context.Context, key string) (string, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, key)
	}
	return "", nil
}

func (m *RedisMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if m.SetFn != nil {
		return m.SetFn(ctx, key, value, expiration)
	}
	return nil
}

func (m *RedisMock) Delete(ctx context.Context, keys ...string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, keys...)
	}
	return nil
}

func (m *RedisMock) Exists(ctx context.Context, keys ...string) (int64, error) {
	if m.ExistsFn != nil {
		return m.ExistsFn(ctx, keys...)
	}
	return 0, nil
}

func (m *RedisMock) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if m.ExpireFn != nil {
		return m.ExpireFn(ctx, key, expiration)
	}
	return nil
}

func (m *RedisMock) TTL(ctx context.Context, key string) (time.Duration, error) {
	if m.TTLFn != nil {
		return m.TTLFn(ctx, key)
	}
	return 0, nil
}

func (m *RedisMock) Keys(ctx context.Context, pattern string) ([]string, error) {
	if m.KeysFn != nil {
		return m.KeysFn(ctx, pattern)
	}
	return nil, nil
}

func (m *RedisMock) Incr(ctx context.Context, key string) (int64, error) {
	if m.IncrFn != nil {
		return m.IncrFn(ctx, key)
	}
	return 0, nil
}

func (m *RedisMock) Decr(ctx context.Context, key string) (int64, error) {
	if m.DecrFn != nil {
		return m.DecrFn(ctx, key)
	}
	return 0, nil
}

func (m *RedisMock) HGet(ctx context.Context, key, field string) (string, error) {
	if m.HGetFn != nil {
		return m.HGetFn(ctx, key, field)
	}
	return "", nil
}

func (m *RedisMock) HSet(ctx context.Context, key string, values ...interface{}) error {
	if m.HSetFn != nil {
		return m.HSetFn(ctx, key, values...)
	}
	return nil
}

func (m *RedisMock) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if m.HGetAllFn != nil {
		return m.HGetAllFn(ctx, key)
	}
	return nil, nil
}

func (m *RedisMock) HDel(ctx context.Context, key string, fields ...string) error {
	if m.HDelFn != nil {
		return m.HDelFn(ctx, key, fields...)
	}
	return nil
}

func (m *RedisMock) RegisterCachePrefix(prefix string) {
	if m.RegisterCachePrefixFn != nil {
		m.RegisterCachePrefixFn(prefix)
	}
}

func (m *RedisMock) RegisterChannel(channel string) {
	if m.RegisterChannelFn != nil {
		m.RegisterChannelFn(channel)
	}
}
