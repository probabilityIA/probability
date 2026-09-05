package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type IntegrationCache struct {
	redis redis.IRedis
	log   log.ILogger
}

func New(redisClient redis.IRedis, logger log.ILogger) domain.IIntegrationCache {
	return &IntegrationCache{
		redis: redisClient,
		log:   logger.WithModule("integration.cache"),
	}
}

func (c *IntegrationCache) SetIntegration(ctx context.Context, integration *domain.CachedIntegration) error {
	data, err := json.Marshal(integration)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to marshal integration")
		return err
	}

	key := integrationKey(integration.ID)
	if err := c.redis.Set(ctx, key, string(data), ttlMetadata); err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_id", integration.ID).Msg("Failed to cache integration")
		return err
	}

	if integration.Code != "" {
		codeIdx := codeKey(integration.Code)
		if integration.IsActive {
			if err := c.redis.Set(ctx, codeIdx, strconv.Itoa(int(integration.ID)), ttlMetadata); err != nil {
				c.log.Warn(ctx).Err(err).Str("code", integration.Code).Msg("Failed to cache code index")
			}
		} else {
			if err := c.redis.Delete(ctx, codeIdx); err != nil {
				c.log.Warn(ctx).Err(err).Str("code", integration.Code).Msg("Failed to remove stale code index")
			}
		}
	}

	if integration.BusinessID != nil {
		bizTypeIdx := businessTypeIndexKey(*integration.BusinessID, integration.IntegrationTypeID)
		if integration.IsActive {
			if err := c.redis.Set(ctx, bizTypeIdx, strconv.Itoa(int(integration.ID)), ttlMetadata); err != nil {
				c.log.Warn(ctx).Err(err).Msg("Failed to cache business+type index")
			}
		} else {
			if err := c.redis.Delete(ctx, bizTypeIdx); err != nil {
				c.log.Warn(ctx).Err(err).Msg("Failed to remove stale business+type index")
			}
		}
	}

	if integration.StoreID != "" {
		storeIdx := storeTypeIndexKey(integration.StoreID, integration.IntegrationTypeID)
		if integration.IsActive {
			if err := c.redis.Set(ctx, storeIdx, strconv.Itoa(int(integration.ID)), ttlMetadata); err != nil {
				c.log.Warn(ctx).Err(err).Str("store_id", integration.StoreID).Msg("Failed to cache store+type index")
			}
		} else {
			if err := c.redis.Delete(ctx, storeIdx); err != nil {
				c.log.Warn(ctx).Err(err).Str("store_id", integration.StoreID).Msg("Failed to remove stale store+type index")
			}
		}
	}

	for _, field := range indexedConfigFields {
		value := configValueAsString(integration.Config[field])
		if value == "" {
			continue
		}
		if integration.IsActive {
			_ = c.SetConfigValueIndex(ctx, integration.IntegrationTypeID, field, value, integration.ID)
		} else {
			_ = c.InvalidateConfigValueIndex(ctx, integration.IntegrationTypeID, field, value)
		}
	}

	c.log.Debug(ctx).Uint("integration_id", integration.ID).Msg("Integration cached")
	return nil
}

func (c *IntegrationCache) GetIntegration(ctx context.Context, integrationID uint) (*domain.CachedIntegration, error) {
	key := integrationKey(integrationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var integration domain.CachedIntegration
	if err := json.Unmarshal([]byte(data), &integration); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal cached integration")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_id", integrationID).Msg("✅ Cache hit - metadata")
	return &integration, nil
}

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

func (c *IntegrationCache) GetCredentials(ctx context.Context, integrationID uint) (*domain.CachedCredentials, error) {
	key := credentialsKey(integrationID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var creds domain.CachedCredentials
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal cached credentials")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_id", integrationID).Msg("✅ Cache hit - credentials")
	return &creds, nil
}

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

func (c *IntegrationCache) SetPlatformCredentials(ctx context.Context, integrationTypeID uint, creds map[string]interface{}) error {
	data, err := json.Marshal(creds)
	if err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Failed to marshal platform credentials")
		return err
	}

	key := platformCredentialsKey(integrationTypeID)
	if err := c.redis.Set(ctx, key, string(data), ttlPlatformCredentials); err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Failed to cache platform credentials")
		return err
	}

	c.log.Debug(ctx).Uint("integration_type_id", integrationTypeID).Msg("Platform credentials cached (no TTL)")
	return nil
}

func (c *IntegrationCache) GetPlatformCredentials(ctx context.Context, integrationTypeID uint) (map[string]interface{}, error) {
	key := platformCredentialsKey(integrationTypeID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var creds map[string]interface{}
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		c.log.Error(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Failed to unmarshal cached platform credentials")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_type_id", integrationTypeID).Msg("✅ Cache hit - platform credentials")
	return creds, nil
}

func (c *IntegrationCache) InvalidateIntegration(ctx context.Context, integrationID uint) error {
	if err := c.redis.Delete(ctx, integrationKey(integrationID)); err != nil {
		c.log.Warn(ctx).Err(err).Msg("Failed to delete metadata cache")
	}

	if err := c.redis.Delete(ctx, credentialsKey(integrationID)); err != nil {
		c.log.Warn(ctx).Err(err).Msg("Failed to delete credentials cache")
	}

	c.log.Info(ctx).Uint("integration_id", integrationID).Msg("🗑️ Cache invalidated")
	return nil
}

func (c *IntegrationCache) InvalidatePlatformCredentials(ctx context.Context, integrationTypeID uint) error {
	key := platformCredentialsKey(integrationTypeID)
	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).Err(err).Uint("integration_type_id", integrationTypeID).Msg("Failed to delete platform credentials cache")
		return err
	}
	c.log.Info(ctx).Uint("integration_type_id", integrationTypeID).Msg("Platform credentials cache invalidated")
	return nil
}

func (c *IntegrationCache) InvalidateMetadata(ctx context.Context, integrationID uint) error {
	if err := c.redis.Delete(ctx, integrationKey(integrationID)); err != nil {
		c.log.Warn(ctx).Err(err).Msg("Failed to delete metadata cache")
	}

	c.log.Info(ctx).Uint("integration_id", integrationID).Msg("🗑️ Metadata cache invalidated")
	return nil
}

func (c *IntegrationCache) InvalidateBusinessTypeIndex(ctx context.Context, businessID, integrationTypeID uint) error {
	key := businessTypeIndexKey(businessID, integrationTypeID)
	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Failed to delete business+type index cache")
		return err
	}
	return nil
}

func (c *IntegrationCache) InvalidateCodeIndex(ctx context.Context, code string) error {
	if code == "" {
		return nil
	}
	key := codeKey(code)
	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).Err(err).Str("code", code).Msg("Failed to delete code index cache")
		return err
	}
	return nil
}

func (c *IntegrationCache) SetBusinessTypeIndex(ctx context.Context, businessID, integrationTypeID, integrationID uint) error {
	key := businessTypeIndexKey(businessID, integrationTypeID)
	if err := c.redis.Set(ctx, key, strconv.Itoa(int(integrationID)), ttlMetadata); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Failed to set business+type index cache")
		return err
	}
	return nil
}

func (c *IntegrationCache) GetByCode(ctx context.Context, code string) (*domain.CachedIntegration, error) {
	idxKey := codeKey(code)
	idStr, err := c.redis.Get(ctx, idxKey)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("code", code).Msg("Failed to parse cached ID")
		return nil, err
	}

	return c.GetIntegration(ctx, uint(id))
}

func (c *IntegrationCache) GetByStoreAndType(ctx context.Context, storeID string, integrationTypeID uint) (*domain.CachedIntegration, error) {
	idxKey := storeTypeIndexKey(storeID, integrationTypeID)
	idStr, err := c.redis.Get(ctx, idxKey)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("store_id", storeID).Msg("Failed to parse cached ID")
		return nil, err
	}

	return c.GetIntegration(ctx, uint(id))
}

func (c *IntegrationCache) InvalidateStoreTypeIndex(ctx context.Context, storeID string, integrationTypeID uint) error {
	if storeID == "" {
		return nil
	}
	key := storeTypeIndexKey(storeID, integrationTypeID)
	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Failed to delete store+type index cache")
		return err
	}
	return nil
}

func (c *IntegrationCache) GetByBusinessAndType(ctx context.Context, businessID, integrationTypeID uint) (*domain.CachedIntegration, error) {
	idxKey := businessTypeIndexKey(businessID, integrationTypeID)
	idStr, err := c.redis.Get(ctx, idxKey)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to parse cached ID")
		return nil, err
	}

	return c.GetIntegration(ctx, uint(id))
}

var indexedConfigFields = []string{"phone_number_id"}

func (c *IntegrationCache) GetByConfigValue(ctx context.Context, integrationTypeID uint, field, value string) (*domain.CachedIntegration, error) {
	if field == "" || value == "" {
		return nil, fmt.Errorf("config index requires field and value")
	}

	idStr, err := c.redis.Get(ctx, configValueIndexKey(integrationTypeID, field, value))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.log.Error(ctx).Err(err).Str("field", field).Msg("Failed to parse cached ID from config index")
		return nil, err
	}

	return c.GetIntegration(ctx, uint(id))
}

func (c *IntegrationCache) SetConfigValueIndex(ctx context.Context, integrationTypeID uint, field, value string, integrationID uint) error {
	if field == "" || value == "" {
		return nil
	}
	key := configValueIndexKey(integrationTypeID, field, value)
	if err := c.redis.Set(ctx, key, strconv.Itoa(int(integrationID)), ttlMetadata); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Failed to cache config value index")
		return err
	}
	return nil
}

func (c *IntegrationCache) InvalidateConfigValueIndex(ctx context.Context, integrationTypeID uint, field, value string) error {
	if field == "" || value == "" {
		return nil
	}
	key := configValueIndexKey(integrationTypeID, field, value)
	if err := c.redis.Delete(ctx, key); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Failed to delete config value index")
		return err
	}
	return nil
}

func configValueAsString(raw interface{}) string {
	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case json.Number:
		return v.String()
	}
	return ""
}

func (c *IntegrationCache) InvalidateConfigValueIndexes(ctx context.Context, integrationTypeID uint, config map[string]interface{}) error {
	for _, field := range indexedConfigFields {
		value := configValueAsString(config[field])
		if value == "" {
			continue
		}
		if err := c.InvalidateConfigValueIndex(ctx, integrationTypeID, field, value); err != nil {
			return err
		}
	}
	return nil
}
