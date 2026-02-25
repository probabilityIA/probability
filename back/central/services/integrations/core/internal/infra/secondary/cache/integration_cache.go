package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type IntegrationCache struct {
	redis redis.IRedis
	log   log.ILogger
}

// New crea una nueva instancia del cache de integraciones
func New(redisClient redis.IRedis, logger log.ILogger) domain.IIntegrationCache {
	return &IntegrationCache{
		redis: redisClient,
		log:   logger.WithModule("integration.cache"),
	}
}

// SetIntegration cachea metadata + secondary indexes
func (c *IntegrationCache) SetIntegration(ctx context.Context, integration *domain.CachedIntegration) error {
	// 1. Serializar metadata
	data, err := json.Marshal(integration)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to marshal integration")
		return err
	}

	// 2. Set metadata key
	key := integrationKey(integration.ID)
	if err := c.redis.Set(ctx, key, string(data), ttlMetadata); err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_id", integration.ID).Msg("Failed to cache integration")
		return err
	}

	// 3. Set secondary indexes
	if integration.Code != "" {
		codeIdx := codeKey(integration.Code)
		if err := c.redis.Set(ctx, codeIdx, strconv.Itoa(int(integration.ID)), ttlMetadata); err != nil {
			c.log.Warn(ctx).Err(err).Str("code", integration.Code).Msg("Failed to cache code index")
		}
	}

	if integration.BusinessID != nil {
		bizTypeIdx := businessTypeIndexKey(*integration.BusinessID, integration.IntegrationTypeID)
		if err := c.redis.Set(ctx, bizTypeIdx, strconv.Itoa(int(integration.ID)), ttlMetadata); err != nil {
			c.log.Warn(ctx).Err(err).Msg("Failed to cache business+type index")
		}
	}

	c.log.Debug(ctx).Uint("integration_id", integration.ID).Msg("✅ Integration cached")
	return nil
}

// GetIntegration lee metadata desde cache
func (c *IntegrationCache) GetIntegration(ctx context.Context, integrationID uint) (*domain.CachedIntegration, error) {
	key := integrationKey(integrationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err // Cache miss
	}

	var integration domain.CachedIntegration
	if err := json.Unmarshal([]byte(data), &integration); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal cached integration")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_id", integrationID).Msg("✅ Cache hit - metadata")
	return &integration, nil
}

// SetCredentials cachea credentials desencriptadas (TTL corto - 1h)
func (c *IntegrationCache) SetCredentials(ctx context.Context, creds *domain.CachedCredentials) error {
	creds.CachedAt = time.Now()

	data, err := json.Marshal(creds)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to marshal credentials")
		return err
	}

	key := credentialsKey(creds.IntegrationID)
	if err := c.redis.Set(ctx, key, string(data), ttlCredentials); err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_id", creds.IntegrationID).Msg("Failed to cache credentials")
		return err
	}

	c.log.Debug(ctx).Uint("integration_id", creds.IntegrationID).Msg("✅ Credentials cached (TTL: 1h)")
	return nil
}

// GetCredentials lee credentials desencriptadas desde cache
func (c *IntegrationCache) GetCredentials(ctx context.Context, integrationID uint) (*domain.CachedCredentials, error) {
	key := credentialsKey(integrationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err // Cache miss
	}

	var creds domain.CachedCredentials
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal cached credentials")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_id", integrationID).Msg("✅ Cache hit - credentials")
	return &creds, nil
}

// GetCredentialField obtiene un campo específico de credentials
func (c *IntegrationCache) GetCredentialField(ctx context.Context, integrationID uint, field string) (string, error) {
	creds, err := c.GetCredentials(ctx, integrationID)
	if err != nil {
		return "", err
	}

	value, exists := creds.Credentials[field]
	if !exists {
		return "", fmt.Errorf("credential field not found: %s", field)
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("credential field is not a string: %s", field)
	}

	c.log.Debug(ctx).
		Uint("integration_id", integrationID).
		Str("field", field).
		Msg("🔐 Credential field retrieved from cache")

	return strValue, nil
}

// SetPlatformCredentials cachea las credenciales de plataforma de un tipo de integración (TTL: 24h)
func (c *IntegrationCache) SetPlatformCredentials(ctx context.Context, integrationTypeID uint, creds map[string]interface{}) error {
	data, err := json.Marshal(creds)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Failed to marshal platform credentials")
		return err
	}

	key := platformCredentialsKey(integrationTypeID)
	if err := c.redis.Set(ctx, key, string(data), ttlMetadata); err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Failed to cache platform credentials")
		return err
	}

	c.log.Debug(ctx).Uint("integration_type_id", integrationTypeID).Msg("✅ Platform credentials cached (TTL: 24h)")
	return nil
}

// GetPlatformCredentials lee las credenciales de plataforma de un tipo de integración desde cache
func (c *IntegrationCache) GetPlatformCredentials(ctx context.Context, integrationTypeID uint) (map[string]interface{}, error) {
	key := platformCredentialsKey(integrationTypeID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err // Cache miss
	}

	var creds map[string]interface{}
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Failed to unmarshal cached platform credentials")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_type_id", integrationTypeID).Msg("✅ Cache hit - platform credentials")
	return creds, nil
}

// InvalidateIntegration elimina metadata + credentials de cache
func (c *IntegrationCache) InvalidateIntegration(ctx context.Context, integrationID uint) error {
	// Delete metadata
	if err := c.redis.Delete(ctx, integrationKey(integrationID)); err != nil {
		c.log.Warn(ctx).Err(err).Msg("Failed to delete metadata cache")
	}

	// Delete credentials (seguridad)
	if err := c.redis.Delete(ctx, credentialsKey(integrationID)); err != nil {
		c.log.Warn(ctx).Err(err).Msg("Failed to delete credentials cache")
	}

	c.log.Info(ctx).Uint("integration_id", integrationID).Msg("🗑️ Cache invalidated")
	return nil
}

// GetByCode busca por código usando index
func (c *IntegrationCache) GetByCode(ctx context.Context, code string) (*domain.CachedIntegration, error) {
	// 1. Get ID from index
	idxKey := codeKey(code)
	idStr, err := c.redis.Get(ctx, idxKey)
	if err != nil {
		return nil, err // Cache miss
	}

	// 2. Parse ID
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("code", code).Msg("Failed to parse cached ID")
		return nil, err
	}

	// 3. Get full metadata
	return c.GetIntegration(ctx, uint(id))
}

// GetByBusinessAndType busca por business+type usando index
func (c *IntegrationCache) GetByBusinessAndType(ctx context.Context, businessID, integrationTypeID uint) (*domain.CachedIntegration, error) {
	idxKey := businessTypeIndexKey(businessID, integrationTypeID)
	idStr, err := c.redis.Get(ctx, idxKey)
	if err != nil {
		return nil, err // Cache miss
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to parse cached ID")
		return nil, err
	}

	return c.GetIntegration(ctx, uint(id))
}
