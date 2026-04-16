package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

const (
	compareCachePrefix = "invoicing:compare"
	compareCacheTTL    = 5 * time.Minute
)

// CompareCache almacena resultados de comparación de facturas en Redis con TTL corto.
// Proporciona un mecanismo de entrega alternativo a SSE (belt + suspenders).
type CompareCache struct {
	redis  redisclient.IRedis
	logger log.ILogger
}

// NewCompareCache crea una nueva instancia del cache de comparación
func NewCompareCache(redisClient redisclient.IRedis, logger log.ILogger) ports.ICompareCache {
	return &CompareCache{
		redis:  redisClient,
		logger: logger.WithModule("invoicing.compare_cache"),
	}
}

// buildKey construye la key de Redis: invoicing:compare:{correlationID}
func (c *CompareCache) buildKey(correlationID string) string {
	return fmt.Sprintf("%s:%s", compareCachePrefix, correlationID)
}

// StoreCompareResult almacena el resultado de una comparación en Redis con TTL de 5 minutos
func (c *CompareCache) StoreCompareResult(ctx context.Context, correlationID string, data *dtos.CompareResponseData) error {
	if c.redis == nil {
		return nil
	}

	if data == nil {
		return nil
	}

	key := c.buildKey(correlationID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to marshal compare result for cache")
		return err
	}

	if err := c.redis.Set(ctx, key, string(jsonData), compareCacheTTL); err != nil {
		c.logger.Warn(ctx).
			Err(err).
			Str("key", key).
			Msg("Failed to store compare result in Redis")
		return err
	}

	c.logger.Info(ctx).
		Str("correlation_id", correlationID).
		Str("key", key).
		Dur("ttl", compareCacheTTL).
		Msg("Compare result stored in Redis")

	return nil
}

// GetCompareResult recupera el resultado de una comparación de Redis.
// Retorna nil, nil si no existe (aún no listo o expirado).
func (c *CompareCache) GetCompareResult(ctx context.Context, correlationID string) (*dtos.CompareResponseData, error) {
	if c.redis == nil {
		return nil, nil
	}

	key := c.buildKey(correlationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		// Key not found - not ready yet or expired
		return nil, nil
	}

	var result dtos.CompareResponseData
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to unmarshal compare result from cache")
		return nil, nil
	}

	return &result, nil
}

// ÍTEMS / PRODUCTOS

const (
	itemCompareCachePrefix = "invoicing:items_compare"
)

// StoreItemCompareResult almacena el resultado de una comparación de ítems en Redis con TTL de 5 minutos
func (c *CompareCache) StoreItemCompareResult(ctx context.Context, correlationID string, data *dtos.ItemCompareResponseData) error {
	if c.redis == nil || data == nil {
		return nil
	}

	key := fmt.Sprintf("%s:%s", itemCompareCachePrefix, correlationID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to marshal item compare result for cache")
		return err
	}

	if err := c.redis.Set(ctx, key, string(jsonData), compareCacheTTL); err != nil {
		c.logger.Warn(ctx).
			Err(err).
			Str("key", key).
			Msg("Failed to store item compare result in Redis")
		return err
	}

	c.logger.Info(ctx).
		Str("correlation_id", correlationID).
		Str("key", key).
		Dur("ttl", compareCacheTTL).
		Msg("Item compare result stored in Redis")

	return nil
}

// GetItemCompareResult recupera el resultado de una comparación de ítems de Redis.
// Retorna nil, nil si no existe (aún no listo o expirado).
func (c *CompareCache) GetItemCompareResult(ctx context.Context, correlationID string) (*dtos.ItemCompareResponseData, error) {
	if c.redis == nil {
		return nil, nil
	}

	key := fmt.Sprintf("%s:%s", itemCompareCachePrefix, correlationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil
	}

	var result dtos.ItemCompareResponseData
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to unmarshal item compare result from cache")
		return nil, nil
	}

	return &result, nil
}

// CUENTAS BANCARIAS

const (
	bankAccountsCachePrefix = "invoicing:bank_accounts"
)

// StoreBankAccountsResult almacena el resultado de cuentas bancarias en Redis con TTL de 5 minutos
func (c *CompareCache) StoreBankAccountsResult(ctx context.Context, correlationID string, data *dtos.BankAccountsResponseData) error {
	if c.redis == nil || data == nil {
		return nil
	}

	key := fmt.Sprintf("%s:%s", bankAccountsCachePrefix, correlationID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to marshal bank accounts result for cache")
		return err
	}

	if err := c.redis.Set(ctx, key, string(jsonData), compareCacheTTL); err != nil {
		c.logger.Warn(ctx).
			Err(err).
			Str("key", key).
			Msg("Failed to store bank accounts result in Redis")
		return err
	}

	c.logger.Info(ctx).
		Str("correlation_id", correlationID).
		Str("key", key).
		Dur("ttl", compareCacheTTL).
		Msg("Bank accounts result stored in Redis")

	return nil
}

// GetBankAccountsResult recupera el resultado de cuentas bancarias de Redis.
// Retorna nil, nil si no existe (aún no listo o expirado).
func (c *CompareCache) GetBankAccountsResult(ctx context.Context, correlationID string) (*dtos.BankAccountsResponseData, error) {
	if c.redis == nil {
		return nil, nil
	}

	key := fmt.Sprintf("%s:%s", bankAccountsCachePrefix, correlationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil
	}

	var result dtos.BankAccountsResponseData
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to unmarshal bank accounts result from cache")
		return nil, nil
	}

	return &result, nil
}

// noopCompareCache es una implementación no-op para cuando Redis no está disponible
type noopCompareCache struct{}

// NewNoopCompareCache crea un cache que no hace nada (para cuando Redis no está disponible)
func NewNoopCompareCache() ports.ICompareCache {
	return &noopCompareCache{}
}

func (n *noopCompareCache) StoreCompareResult(_ context.Context, _ string, _ *dtos.CompareResponseData) error {
	return nil
}

func (n *noopCompareCache) GetCompareResult(_ context.Context, _ string) (*dtos.CompareResponseData, error) {
	return nil, nil
}

func (n *noopCompareCache) StoreItemCompareResult(_ context.Context, _ string, _ *dtos.ItemCompareResponseData) error {
	return nil
}

func (n *noopCompareCache) GetItemCompareResult(_ context.Context, _ string) (*dtos.ItemCompareResponseData, error) {
	return nil, nil
}

func (n *noopCompareCache) StoreBankAccountsResult(_ context.Context, _ string, _ *dtos.BankAccountsResponseData) error {
	return nil
}

func (n *noopCompareCache) GetBankAccountsResult(_ context.Context, _ string) (*dtos.BankAccountsResponseData, error) {
	return nil, nil
}

// Compile-time interface checks
var _ ports.ICompareCache = (*CompareCache)(nil)
var _ ports.ICompareCache = (*noopCompareCache)(nil)
