package integration_cache

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// CachedIntegration representa los datos de integración en cache
// Debe coincidir con domain.CachedIntegration de IntegrationCore
type CachedIntegration struct {
	ID                  uint                   `json:"id"`
	Name                string                 `json:"name"`
	Code                string                 `json:"code"`
	Config              map[string]interface{} `json:"config"`
	IsActive            bool                   `json:"is_active"`
	IntegrationTypeCode string                 `json:"integration_type_code"`
	Category            string                 `json:"category"`
	IntegrationTypeID   uint                   `json:"integration_type_id"`
	BusinessID          *uint                  `json:"business_id"`
	StoreID             string                 `json:"store_id"`
	IsDefault           bool                   `json:"is_default"`
	Description         string                 `json:"description"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// CachedCredentials representa credenciales desencriptadas en cache
type CachedCredentials struct {
	IntegrationID uint                   `json:"integration_id"`
	Credentials   map[string]interface{} `json:"credentials"`
	CachedAt      time.Time              `json:"cached_at"`
}

// IIntegrationCacheClient define operaciones de lectura de cache
type IIntegrationCacheClient interface {
	GetIntegration(ctx context.Context, integrationID uint) (*CachedIntegration, error)
	GetCredential(ctx context.Context, integrationID uint, field string) (string, error)
	GetAllCredentials(ctx context.Context, integrationID uint) (map[string]interface{}, error)
}

// Client cliente para acceder al cache de IntegrationCore
type Client struct {
	redis redis.IRedis
	log   log.ILogger
}

// New crea una nueva instancia del cliente de cache
func New(redisClient redis.IRedis, logger log.ILogger) IIntegrationCacheClient {
	return &Client{
		redis: redisClient,
		log:   logger.WithModule("factus.integration_cache"),
	}
}

// GetIntegration lee metadata de integración desde cache
func (c *Client) GetIntegration(ctx context.Context, integrationID uint) (*CachedIntegration, error) {
	key := "integration:meta:" + strconv.Itoa(int(integrationID))

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		c.log.Warn(ctx).Err(err).Uint("integration_id", integrationID).Msg("Cache miss - integration")
		return nil, err
	}

	var integration CachedIntegration
	if err := json.Unmarshal([]byte(data), &integration); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal cached integration")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_id", integrationID).Msg("✅ Cache hit - integration")
	return &integration, nil
}

// GetCredential obtiene un campo específico de credentials
func (c *Client) GetCredential(ctx context.Context, integrationID uint, field string) (string, error) {
	creds, err := c.GetAllCredentials(ctx, integrationID)
	if err != nil {
		return "", err
	}

	value, exists := creds[field]
	if !exists {
		return "", errors.New("credential field not found: " + field)
	}

	strValue, ok := value.(string)
	if !ok {
		return "", errors.New("credential field is not a string: " + field)
	}

	c.log.Debug(ctx).
		Uint("integration_id", integrationID).
		Str("field", field).
		Msg("🔐 Credential retrieved from cache")

	return strValue, nil
}

// GetAllCredentials lee todas las credentials de cache
func (c *Client) GetAllCredentials(ctx context.Context, integrationID uint) (map[string]interface{}, error) {
	key := "integration:creds:" + strconv.Itoa(int(integrationID))

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		c.log.Warn(ctx).Err(err).Uint("integration_id", integrationID).Msg("Cache miss - credentials")
		return nil, err
	}

	var cached CachedCredentials
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal cached credentials")
		return nil, err
	}

	c.log.Debug(ctx).Uint("integration_id", integrationID).Msg("✅ Cache hit - credentials")
	return cached.Credentials, nil
}
