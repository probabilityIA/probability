package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// InventoryCache implementa repository.IInventoryCache usando Redis
type InventoryCache struct {
	redis  redis.IRedis
	config env.IConfig
	log    log.ILogger
}

// NewInventoryCache crea una nueva instancia del servicio de caché de inventario
func NewInventoryCache(redisClient redis.IRedis, config env.IConfig, logger log.ILogger) repository.IInventoryCache {
	cache := &InventoryCache{
		redis:  redisClient,
		config: config,
		log:    logger.WithModule("inventory.cache"),
	}

	// Registrar prefijo en Redis para tracking
	if redisClient != nil {
		redisClient.RegisterCachePrefix(cache.getPrefix())
	}

	return cache
}

// getPrefix obtiene el prefijo de cache configurado
func (c *InventoryCache) getPrefix() string {
	prefix := c.config.Get("REDIS_INVENTORY_CACHE_PREFIX")
	if prefix == "" {
		prefix = "probability:inventory"
	}
	return prefix
}

// getTTL obtiene el TTL configurado para el caché de inventario
func (c *InventoryCache) getTTL() time.Duration {
	ttlStr := c.config.Get("REDIS_INVENTORY_CACHE_TTL")
	if ttlStr == "" {
		return 600 * time.Second // 10 minutos por defecto
	}

	ttlSeconds, err := strconv.Atoi(ttlStr)
	if err != nil {
		c.log.Warn(context.Background()).
			Str("ttl", ttlStr).
			Msg("Invalid REDIS_INVENTORY_CACHE_TTL, using default 600s")
		return 600 * time.Second
	}

	return time.Duration(ttlSeconds) * time.Second
}

// buildProductKey construye la key para inventario de un producto en un negocio
// Pattern: probability:inventory:product:{product_id}:{business_id}
func (c *InventoryCache) buildProductKey(productID string, businessID uint) string {
	return fmt.Sprintf("%s:product:%s:%d", c.getPrefix(), productID, businessID)
}

// buildLevelKey construye la key para un nivel de inventario individual
// Pattern: probability:inventory:level:{product_id}:{warehouse_id}
func (c *InventoryCache) buildLevelKey(productID string, warehouseID uint) string {
	return fmt.Sprintf("%s:level:%s:%d", c.getPrefix(), productID, warehouseID)
}

// GetProductLevels obtiene los niveles de inventario de un producto desde caché
func (c *InventoryCache) GetProductLevels(ctx context.Context, productID string, businessID uint) ([]entities.InventoryLevel, error) {
	if c.redis == nil {
		return nil, nil
	}

	key := c.buildProductKey(productID, businessID)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil // cache MISS
	}

	var levels []entities.InventoryLevel
	if err := json.Unmarshal([]byte(data), &levels); err != nil {
		c.log.Error(ctx).Err(err).Str("key", key).Msg("Error deserializing product levels from cache")
		return nil, nil
	}

	return levels, nil
}

// SetProductLevels guarda los niveles de inventario de un producto en caché
func (c *InventoryCache) SetProductLevels(ctx context.Context, productID string, businessID uint, levels []entities.InventoryLevel) error {
	if c.redis == nil {
		return nil
	}

	key := c.buildProductKey(productID, businessID)
	data, err := json.Marshal(levels)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("product_id", productID).Msg("Failed to marshal product levels for cache")
		return err
	}

	return c.redis.Set(ctx, key, string(data), c.getTTL())
}

// InvalidateProduct elimina los niveles de inventario de un producto del caché
func (c *InventoryCache) InvalidateProduct(ctx context.Context, productID string, businessID uint) error {
	if c.redis == nil {
		return nil
	}

	key := c.buildProductKey(productID, businessID)
	return c.redis.Delete(ctx, key)
}

// GetLevel obtiene un nivel de inventario individual desde caché
func (c *InventoryCache) GetLevel(ctx context.Context, productID string, warehouseID uint) (*entities.InventoryLevel, error) {
	if c.redis == nil {
		return nil, nil
	}

	key := c.buildLevelKey(productID, warehouseID)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil // cache MISS
	}

	var level entities.InventoryLevel
	if err := json.Unmarshal([]byte(data), &level); err != nil {
		c.log.Error(ctx).Err(err).Str("key", key).Msg("Error deserializing level from cache")
		return nil, nil
	}

	return &level, nil
}

// SetLevel guarda un nivel de inventario individual en caché
func (c *InventoryCache) SetLevel(ctx context.Context, productID string, warehouseID uint, level *entities.InventoryLevel) error {
	if c.redis == nil {
		return nil
	}
	if level == nil {
		return nil
	}

	key := c.buildLevelKey(productID, warehouseID)
	data, err := json.Marshal(level)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("product_id", productID).Msg("Failed to marshal level for cache")
		return err
	}

	return c.redis.Set(ctx, key, string(data), c.getTTL())
}

// InvalidateLevel elimina un nivel de inventario individual del caché
func (c *InventoryCache) InvalidateLevel(ctx context.Context, productID string, warehouseID uint) error {
	if c.redis == nil {
		return nil
	}

	key := c.buildLevelKey(productID, warehouseID)
	return c.redis.Delete(ctx, key)
}
