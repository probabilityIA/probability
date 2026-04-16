package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// CacheRegistryCallback es un callback para registrar prefijos de caché usados
type CacheRegistryCallback func(prefix string)

// ChannelRegistryCallback es un callback para registrar canales pub/sub activos
type ChannelRegistryCallback func(channel string)

// IRedis define la interfaz para la conexión a Redis
type IRedis interface {
	Connect(ctx context.Context) error
	Close() error
	Client(ctx context.Context) *redis.Client
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// Métodos de registro para tracking
	RegisterCachePrefix(prefix string)
	RegisterChannel(channel string)
}

// redisClient implementa la interfaz IRedis
type redisClient struct {
	client          *redis.Client
	log             log.ILogger
	config          env.IConfig
	cacheRegistry   CacheRegistryCallback
	channelRegistry ChannelRegistryCallback
}

// New crea una nueva instancia de Redis y conecta automáticamente
func New(logger log.ILogger, config env.IConfig) IRedis {
	r := &redisClient{
		log:    logger,
		config: config,
	}

	// Conectar automáticamente a Redis
	if err := r.Connect(context.Background()); err != nil {
		logger.Warn(context.Background()).
			Err(err).
			Msg("Error al conectar a Redis - la aplicación continuará sin cache")
		// No hacemos fatal porque Redis es opcional (cache)
	}

	return r
}

// SetCacheRegistry establece un callback para registrar prefijos de caché
func (r *redisClient) SetCacheRegistry(callback CacheRegistryCallback) {
	r.cacheRegistry = callback
}

// SetChannelRegistry establece un callback para registrar canales pub/sub
func (r *redisClient) SetChannelRegistry(callback ChannelRegistryCallback) {
	r.channelRegistry = callback
}

// RegisterCachePrefix registra un prefijo de caché usado (para mostrar en startup logs)
func (r *redisClient) RegisterCachePrefix(prefix string) {
	if r.cacheRegistry != nil {
		r.cacheRegistry(prefix)
	}
}

// RegisterChannel registra un canal pub/sub activo (para mostrar en startup logs)
func (r *redisClient) RegisterChannel(channel string) {
	if r.channelRegistry != nil {
		r.channelRegistry(channel)
	}
}

// Connect establece la conexión con Redis
func (r *redisClient) Connect(ctx context.Context) error {
	host := r.config.Get("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := r.config.Get("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := r.config.Get("REDIS_PASSWORD")
	db := 0 // Default DB

	addr := fmt.Sprintf("%s:%s", host, port)

	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Verificar conexión
	if err := r.Ping(ctx); err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("addr", addr).
			Msg("Error al conectar a Redis")
		return fmt.Errorf("error al conectar a Redis: %w", err)
	}

	// Redis info mostrada en LogStartupInfo() - no duplicar aquí
	return nil
}

// Close cierra la conexión con Redis
func (r *redisClient) Close() error {
	if r.client != nil {
		err := r.client.Close()
		if err != nil {
			r.log.Error(context.Background()).
				Err(err).
				Msg("Error al cerrar conexión Redis")
			return err
		}
		r.log.Info(context.Background()).Msg("Conexión Redis cerrada")
	}
	return nil
}

// Client retorna el cliente de Redis
func (r *redisClient) Client(ctx context.Context) *redis.Client {
	return r.client
}

// Ping verifica la conexión con Redis
func (r *redisClient) Ping(ctx context.Context) error {
	start := time.Now()
	err := r.client.Ping(ctx).Err()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Dur("duration", duration).
			Msg("Redis Ping failed")
		return err
	}

	// Ping exitoso - no logueamos para reducir ruido en debug logs
	return nil
}

// Get obtiene un valor de Redis
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	val, err := r.client.Get(ctx, key).Result()
	duration := time.Since(start)

	if err == redis.Nil {
		r.log.Debug(ctx).
			Str("key", key).
			Dur("duration", duration).
			Msg("Redis Get - key not found")
		return "", fmt.Errorf("key not found: %s", key)
	}

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("duration", duration).
			Msg("Redis Get failed")
		return "", err
	}

	// Get exitoso - no logueamos para reducir ruido en debug logs
	return val, nil
}

// Set guarda un valor en Redis
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()
	err := r.client.Set(ctx, key, value, expiration).Err()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("expiration", expiration).
			Dur("duration", duration).
			Msg("Redis Set failed")
		return err
	}

	// Set exitoso - no logueamos para reducir ruido en debug logs
	return nil
}

// Delete elimina una o más claves de Redis
func (r *redisClient) Delete(ctx context.Context, keys ...string) error {
	start := time.Now()
	err := r.client.Del(ctx, keys...).Err()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Strs("keys", keys).
			Dur("duration", duration).
			Msg("Redis Delete failed")
		return err
	}

	// Delete exitoso - no logueamos para reducir ruido en debug logs
	return nil
}

// Exists verifica si una o más claves existen en Redis
func (r *redisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	start := time.Now()
	count, err := r.client.Exists(ctx, keys...).Result()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Strs("keys", keys).
			Dur("duration", duration).
			Msg("Redis Exists failed")
		return 0, err
	}

	// Exists exitoso - no logueamos para reducir ruido en debug logs
	return count, nil
}

// Expire establece el tiempo de expiración de una clave
func (r *redisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	start := time.Now()
	err := r.client.Expire(ctx, key, expiration).Err()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("expiration", expiration).
			Dur("duration", duration).
			Msg("Redis Expire failed")
		return err
	}

	// Expire exitoso - no logueamos para reducir ruido en debug logs
	return nil
}

// TTL obtiene el tiempo restante de expiración de una clave
func (r *redisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	start := time.Now()
	ttl, err := r.client.TTL(ctx, key).Result()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("duration", duration).
			Msg("Redis TTL failed")
		return 0, err
	}

	// TTL exitoso - no logueamos para reducir ruido en debug logs
	return ttl, nil
}

// Keys obtiene todas las claves que coinciden con un patrón
func (r *redisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	start := time.Now()
	keys, err := r.client.Keys(ctx, pattern).Result()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("pattern", pattern).
			Dur("duration", duration).
			Msg("Redis Keys failed")
		return nil, err
	}

	// Keys exitoso - no logueamos para reducir ruido en debug logs
	return keys, nil
}

// Incr incrementa el valor de una clave numérica
func (r *redisClient) Incr(ctx context.Context, key string) (int64, error) {
	start := time.Now()
	val, err := r.client.Incr(ctx, key).Result()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("duration", duration).
			Msg("Redis Incr failed")
		return 0, err
	}

	// Incr exitoso - no logueamos para reducir ruido en debug logs
	return val, nil
}

// Decr decrementa el valor de una clave numérica
func (r *redisClient) Decr(ctx context.Context, key string) (int64, error) {
	start := time.Now()
	val, err := r.client.Decr(ctx, key).Result()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("duration", duration).
			Msg("Redis Decr failed")
		return 0, err
	}

	// Decr exitoso - no logueamos para reducir ruido en debug logs
	return val, nil
}

// HGet obtiene el valor de un campo en un hash
func (r *redisClient) HGet(ctx context.Context, key, field string) (string, error) {
	start := time.Now()
	val, err := r.client.HGet(ctx, key, field).Result()
	duration := time.Since(start)

	if err == redis.Nil {
		r.log.Debug(ctx).
			Str("key", key).
			Str("field", field).
			Dur("duration", duration).
			Msg("Redis HGet - field not found")
		return "", fmt.Errorf("field not found: %s in %s", field, key)
	}

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Str("field", field).
			Dur("duration", duration).
			Msg("Redis HGet failed")
		return "", err
	}

	// HGet exitoso - no logueamos para reducir ruido en debug logs
	return val, nil
}

// HSet establece el valor de uno o más campos en un hash
func (r *redisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	start := time.Now()
	err := r.client.HSet(ctx, key, values...).Err()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("duration", duration).
			Msg("Redis HSet failed")
		return err
	}

	// HSet exitoso - no logueamos para reducir ruido en debug logs
	return nil
}

// HGetAll obtiene todos los campos y valores de un hash
func (r *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	start := time.Now()
	val, err := r.client.HGetAll(ctx, key).Result()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Dur("duration", duration).
			Msg("Redis HGetAll failed")
		return nil, err
	}

	// HGetAll exitoso - no logueamos para reducir ruido en debug logs
	return val, nil
}

// HDel elimina uno o más campos de un hash
func (r *redisClient) HDel(ctx context.Context, key string, fields ...string) error {
	start := time.Now()
	err := r.client.HDel(ctx, key, fields...).Err()
	duration := time.Since(start)

	if err != nil {
		r.log.Error(ctx).
			Err(err).
			Str("key", key).
			Strs("fields", fields).
			Dur("duration", duration).
			Msg("Redis HDel failed")
		return err
	}

	// HDel exitoso - no logueamos para reducir ruido en debug logs
	return nil
}
