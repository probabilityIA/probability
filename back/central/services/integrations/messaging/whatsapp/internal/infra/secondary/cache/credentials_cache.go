package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

const (
	// WhatsApp integration_type_id en la DB
	whatsAppTypeID = 2
)

// Claves Redis generadas por integrations/core warmup
// Referencia: core/internal/infra/secondary/cache/keys.go
func businessTypeIdxKey(businessID uint) string {
	return fmt.Sprintf("integration:idx:biz:%d:type:%d", businessID, whatsAppTypeID)
}

func integrationCredsKey(integrationID uint) string {
	return fmt.Sprintf("integration:creds:%d", integrationID)
}

func integrationMetaKey(integrationID uint) string {
	return fmt.Sprintf("integration:meta:%d", integrationID)
}

func platformCredsKey() string {
	return fmt.Sprintf("integration:platform_creds:%d", whatsAppTypeID)
}

// cachedIntegrationMeta replica parcialmente domain.CachedIntegration de core
type cachedIntegrationMeta struct {
	ID      uint   `json:"id"`
	BaseURL string `json:"base_url"`
}

// cachedCredentials replica parcialmente domain.CachedCredentials de core
type cachedCredentials struct {
	IntegrationID uint                   `json:"integration_id"`
	Credentials   map[string]interface{} `json:"credentials"`
}

// credentialsCache lee credenciales de WhatsApp desde Redis (claves de integrations/core)
type credentialsCache struct {
	redis redisclient.IRedis
	log   log.ILogger
}

// newCredentialsCache crea una nueva instancia del cache de credenciales
func newCredentialsCache(redis redisclient.IRedis, logger log.ILogger) ports.ICredentialsCache {
	return &credentialsCache{
		redis: redis,
		log:   logger.WithModule("whatsapp-credentials-cache"),
	}
}

// GetWhatsAppConfig obtiene credenciales de WhatsApp para un business desde Redis.
// Lee las claves que integrations/core calienta al startup.
func (c *credentialsCache) GetWhatsAppConfig(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error) {
	// 1. Obtener integration_id desde índice business+type
	integrationIDStr, err := c.redis.Get(ctx, businessTypeIdxKey(businessID))
	if err != nil {
		return nil, fmt.Errorf("no se encontró integración WhatsApp para business %d en cache: %w", businessID, err)
	}

	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("integration_id inválido en cache: %s", integrationIDStr)
	}

	// 2. Obtener credenciales desencriptadas
	credsJSON, err := c.redis.Get(ctx, integrationCredsKey(uint(integrationID)))
	if err != nil {
		return nil, fmt.Errorf("credenciales no encontradas en cache para integration %d: %w", integrationID, err)
	}

	var creds cachedCredentials
	if err := json.Unmarshal([]byte(credsJSON), &creds); err != nil {
		return nil, fmt.Errorf("error deserializando credenciales: %w", err)
	}

	// 3. Obtener metadata (BaseURL)
	metaJSON, err := c.redis.Get(ctx, integrationMetaKey(uint(integrationID)))
	if err != nil {
		c.log.Warn(ctx).
			Err(err).
			Uint("integration_id", uint(integrationID)).
			Msg("Metadata no encontrada en cache, continuando sin BaseURL")
	}

	var meta cachedIntegrationMeta
	if metaJSON != "" {
		json.Unmarshal([]byte(metaJSON), &meta)
	}

	// 4. Construir WhatsAppConfig
	config, err := buildWhatsAppConfig(creds.Credentials, uint(integrationID), meta.BaseURL)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetWhatsAppDefaultConfig obtiene credenciales globales de plataforma desde Redis
func (c *credentialsCache) GetWhatsAppDefaultConfig(ctx context.Context) (*ports.WhatsAppConfig, error) {
	platJSON, err := c.redis.Get(ctx, platformCredsKey())
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma WhatsApp no encontradas en cache: %w", err)
	}

	var platCreds map[string]interface{}
	if err := json.Unmarshal([]byte(platJSON), &platCreds); err != nil {
		return nil, fmt.Errorf("error deserializando credenciales de plataforma: %w", err)
	}

	config, err := buildWhatsAppConfig(platCreds, 0, "")
	if err != nil {
		return nil, err
	}

	return config, nil
}

// buildWhatsAppConfig construye WhatsAppConfig desde un map de credenciales
func buildWhatsAppConfig(creds map[string]interface{}, integrationID uint, baseURL string) (*ports.WhatsAppConfig, error) {
	config := &ports.WhatsAppConfig{
		IntegrationID: integrationID,
		WhatsAppURL:   baseURL,
	}

	// Extraer phone_number_id
	if phoneID, ok := creds["phone_number_id"]; ok {
		switch v := phoneID.(type) {
		case float64:
			config.PhoneNumberID = uint(v)
		case string:
			parsed, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("phone_number_id inválido: %s", v)
			}
			config.PhoneNumberID = uint(parsed)
		}
	}

	// Extraer access_token
	if token, ok := creds["access_token"].(string); ok {
		config.AccessToken = token
	}

	// Extraer whatsapp_url si viene en credenciales (platform_creds)
	if url, ok := creds["whatsapp_url"].(string); ok && config.WhatsAppURL == "" {
		config.WhatsAppURL = url
	}

	if config.PhoneNumberID == 0 {
		return nil, fmt.Errorf("phone_number_id no encontrado en credenciales")
	}
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access_token no encontrado en credenciales")
	}

	return config, nil
}
