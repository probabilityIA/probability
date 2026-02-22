package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// IConfigCache define la interfaz para el servicio de caché de configuraciones
type IConfigCache interface {
	Get(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error)
	Set(ctx context.Context, config *entities.InvoicingConfig) error
	Invalidate(ctx context.Context, integrationID uint) error
}

// ConfigCache implementa el servicio de caché de configuraciones usando Redis
type ConfigCache struct {
	redis  redis.IRedis
	config env.IConfig
	log    log.ILogger
}

// NewConfigCache crea una nueva instancia del servicio de caché
func NewConfigCache(redisClient redis.IRedis, config env.IConfig, logger log.ILogger) IConfigCache {
	return &ConfigCache{
		redis:  redisClient,
		config: config,
		log:    logger.WithModule("factus.config_cache"),
	}
}

// buildCacheKey construye la key de Redis para cachear una configuración
func (c *ConfigCache) buildCacheKey(integrationID uint) string {
	prefix := c.config.Get("REDIS_INVOICING_CONFIG_PREFIX")
	if prefix == "" {
		prefix = "probability:invoicing:config"
	}
	return fmt.Sprintf("%s:%d", prefix, integrationID)
}

func (c *ConfigCache) getTTL() time.Duration {
	return 3600 * time.Second
}

// Get obtiene una configuración desde Redis
func (c *ConfigCache) Get(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
	if c.redis == nil {
		return nil, nil
	}

	key := c.buildCacheKey(integrationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil
	}

	var config entities.InvoicingConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		c.log.Error(ctx).Err(err).Str("key", key).Msg("Error al deserializar config desde caché")
		return nil, nil
	}

	return &config, nil
}

// Set guarda una configuración en Redis con TTL
func (c *ConfigCache) Set(ctx context.Context, config *entities.InvoicingConfig) error {
	if c.redis == nil || config == nil {
		return nil
	}

	key := c.buildCacheKey(config.IntegrationID)

	data, err := json.Marshal(config)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_id", config.IntegrationID).Msg("Failed to marshal config for cache")
		return err
	}

	ttl := c.getTTL()
	if err := c.redis.Set(ctx, key, string(data), ttl); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Error al guardar config en caché")
		return err
	}

	return nil
}

// Invalidate elimina una configuración del caché
func (c *ConfigCache) Invalidate(ctx context.Context, integrationID uint) error {
	if c.redis == nil {
		return nil
	}

	key := c.buildCacheKey(integrationID)

	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Error al invalidar config en caché")
		return err
	}

	return nil
}
