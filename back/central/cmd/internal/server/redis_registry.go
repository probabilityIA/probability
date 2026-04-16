package server

import (
	"sync"
)

// RedisRegistry mantiene un registro de prefijos de caché y canales pub/sub
type RedisRegistry struct {
	cachePrefixes []string
	channels      []string
	mu            sync.Mutex
}

// NewRedisRegistry crea un nuevo registro de Redis
func NewRedisRegistry() *RedisRegistry {
	return &RedisRegistry{
		cachePrefixes: make([]string, 0),
		channels:      make([]string, 0),
	}
}

// RegisterCachePrefix registra un prefijo de caché usado
func (r *RedisRegistry) RegisterCachePrefix(prefix string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Evitar duplicados
	for _, p := range r.cachePrefixes {
		if p == prefix {
			return
		}
	}

	r.cachePrefixes = append(r.cachePrefixes, prefix)
}

// RegisterChannel registra un canal pub/sub activo
func (r *RedisRegistry) RegisterChannel(channel string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Evitar duplicados
	for _, c := range r.channels {
		if c == channel {
			return
		}
	}

	r.channels = append(r.channels, channel)
}

// GetCachePrefixes retorna todos los prefijos registrados
func (r *RedisRegistry) GetCachePrefixes() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]string, len(r.cachePrefixes))
	copy(result, r.cachePrefixes)
	return result
}

// GetChannels retorna todos los canales registrados
func (r *RedisRegistry) GetChannels() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]string, len(r.channels))
	copy(result, r.channels)
	return result
}
