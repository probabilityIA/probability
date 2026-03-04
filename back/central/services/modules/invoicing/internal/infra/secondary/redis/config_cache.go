package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// ═══════════════════════════════════════════════════════════════
// CONFIG CACHE SERVICE - Servicio de caché para configuraciones
// ═══════════════════════════════════════════════════════════════

// ConfigCache implementa el servicio de caché de configuraciones usando Redis
type ConfigCache struct {
	redis  redis.IRedis
	config env.IConfig
	log    log.ILogger
}

// NewConfigCache crea una nueva instancia del servicio de caché
func NewConfigCache(redisClient redis.IRedis, config env.IConfig, logger log.ILogger) ports.IConfigCache {
	return &ConfigCache{
		redis:  redisClient,
		config: config,
		log:    logger.WithModule("invoicing.config_cache"),
	}
}

// buildCacheKey construye la key de Redis para cachear una configuración por integration_id
// Pattern: probability:invoicing:config:{integration_id}
func (c *ConfigCache) buildCacheKey(integrationID uint) string {
	prefix := c.config.Get("REDIS_INVOICING_CONFIG_PREFIX")
	if prefix == "" {
		prefix = "probability:invoicing:config"
	}
	return fmt.Sprintf("%s:%d", prefix, integrationID)
}

// buildBusinessCacheKey construye la key de Redis para cachear la config activa de un negocio
// Pattern: probability:invoicing:config:business:{business_id}
func (c *ConfigCache) buildBusinessCacheKey(businessID uint) string {
	prefix := c.config.Get("REDIS_INVOICING_CONFIG_PREFIX")
	if prefix == "" {
		prefix = "probability:invoicing:config"
	}
	return fmt.Sprintf("%s:business:%d", prefix, businessID)
}

// getTTL obtiene el TTL configurado para el caché de configuraciones
// Default: 1 hora (3600 segundos)
func (c *ConfigCache) getTTL() time.Duration {
	ttlStr := c.config.Get("REDIS_INVOICING_CONFIG_TTL")
	if ttlStr == "" {
		return 3600 * time.Second // 1 hora por defecto
	}

	ttlSeconds, err := strconv.Atoi(ttlStr)
	if err != nil {
		c.log.Warn(context.Background()).
			Str("ttl", ttlStr).
			Msg("Invalid REDIS_INVOICING_CONFIG_TTL, using default 3600s")
		return 3600 * time.Second
	}

	return time.Duration(ttlSeconds) * time.Second
}

// Get obtiene una configuración desde Redis
// Retorna nil si no existe en caché (cache MISS) - NO es un error
func (c *ConfigCache) Get(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
	// Si Redis no está disponible, retornar nil (cache MISS resiliente)
	if c.redis == nil {
		return nil, nil
	}

	key := c.buildCacheKey(integrationID)

	// Leer desde Redis
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		// Cache MISS - no es un error
		return nil, nil
	}

	// Deserializar JSON
	var config entities.InvoicingConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("key", key).
			Msg("Error al deserializar config desde caché")
		return nil, nil // Retornar nil para forzar fallback a BD
	}

	return &config, nil
}

// Set guarda una configuración en Redis con TTL bajo la clave integrationID
func (c *ConfigCache) Set(ctx context.Context, integrationID uint, config *entities.InvoicingConfig) error {
	// Si Redis no está disponible, no hacer nada (resiliente)
	if c.redis == nil {
		return nil
	}

	if config == nil {
		return nil
	}

	key := c.buildCacheKey(integrationID)

	// Serializar a JSON
	data, err := json.Marshal(config)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("integration_id", integrationID).
			Msg("Failed to marshal config for cache")
		return err
	}

	// Guardar en Redis con TTL
	ttl := c.getTTL()
	if err := c.redis.Set(ctx, key, string(data), ttl); err != nil {
		c.log.Warn(ctx).
			Err(err).
			Str("key", key).
			Msg("Error al guardar config en caché")
		return err
	}

	return nil
}

// Invalidate elimina una configuración del caché
func (c *ConfigCache) Invalidate(ctx context.Context, integrationID uint) error {
	// Si Redis no está disponible, no hacer nada (resiliente)
	if c.redis == nil {
		return nil
	}

	key := c.buildCacheKey(integrationID)

	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).
			Err(err).
			Str("key", key).
			Msg("Error al invalidar config en caché")
		return err
	}

	return nil
}

// ═══════════════════════════════════════════════════════════════
// CACHE POR BUSINESS ID - Config activa de un negocio
// ═══════════════════════════════════════════════════════════════

// GetByBusinessID obtiene la configuración activa de un negocio desde Redis
func (c *ConfigCache) GetByBusinessID(ctx context.Context, businessID uint) (*entities.InvoicingConfig, error) {
	if c.redis == nil {
		return nil, nil
	}

	key := c.buildBusinessCacheKey(businessID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, nil // Cache MISS
	}

	var config entities.InvoicingConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("key", key).
			Msg("Error al deserializar config de negocio desde caché")
		return nil, nil
	}

	return &config, nil
}

// SetByBusinessID guarda la configuración activa de un negocio en Redis con TTL
func (c *ConfigCache) SetByBusinessID(ctx context.Context, businessID uint, config *entities.InvoicingConfig) error {
	if c.redis == nil || config == nil {
		return nil
	}

	key := c.buildBusinessCacheKey(businessID)

	data, err := json.Marshal(config)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("business_id", businessID).
			Msg("Failed to marshal config for business cache")
		return err
	}

	ttl := c.getTTL()
	if err := c.redis.Set(ctx, key, string(data), ttl); err != nil {
		c.log.Warn(ctx).
			Err(err).
			Str("key", key).
			Msg("Error al guardar config de negocio en caché")
		return err
	}

	return nil
}

// InvalidateByBusinessID elimina la configuración activa de un negocio del caché
func (c *ConfigCache) InvalidateByBusinessID(ctx context.Context, businessID uint) error {
	if c.redis == nil {
		return nil
	}

	key := c.buildBusinessCacheKey(businessID)

	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).
			Err(err).
			Str("key", key).
			Msg("Error al invalidar config de negocio en caché")
		return err
	}

	return nil
}
