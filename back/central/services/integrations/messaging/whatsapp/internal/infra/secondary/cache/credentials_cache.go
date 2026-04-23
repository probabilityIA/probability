package cache

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const whatsAppTypeID uint = 2

type credentialsCache struct {
	log      log.ILogger
	mu       sync.RWMutex
	resolver ports.IPlatformCredentialsGetter
}

func newCredentialsCache(logger log.ILogger) *credentialsCache {
	return &credentialsCache{
		log: logger.WithModule("whatsapp-credentials-cache"),
	}
}

func (c *credentialsCache) SetResolver(resolver ports.IPlatformCredentialsGetter) {
	c.mu.Lock()
	c.resolver = resolver
	c.mu.Unlock()
}

func (c *credentialsCache) getResolver() ports.IPlatformCredentialsGetter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.resolver
}

func (c *credentialsCache) GetWhatsAppConfig(ctx context.Context, businessID uint) (*ports.WhatsAppConfig, error) {
	resolver := c.getResolver()
	if resolver == nil {
		return nil, fmt.Errorf("whatsapp resolver not configured")
	}

	integrationID, err := resolver.GetIntegrationIDByBusinessAndType(ctx, businessID, whatsAppTypeID)
	if err != nil {
		return nil, fmt.Errorf("no se encontró integración WhatsApp para business %d: %w", businessID, err)
	}

	config, err := c.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no encontradas para WhatsApp (business %d): %w", businessID, err)
	}

	config.IntegrationID = integrationID
	return config, nil
}

func (c *credentialsCache) GetWhatsAppDefaultConfig(ctx context.Context) (*ports.WhatsAppConfig, error) {
	resolver := c.getResolver()
	if resolver == nil {
		return nil, fmt.Errorf("whatsapp resolver not configured")
	}

	creds, err := resolver.GetCachedPlatformCredentials(ctx, whatsAppTypeID)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma WhatsApp no disponibles: %w", err)
	}

	return buildWhatsAppConfig(creds, 0, "")
}

func buildWhatsAppConfig(creds map[string]any, integrationID uint, baseURL string) (*ports.WhatsAppConfig, error) {
	config := &ports.WhatsAppConfig{
		IntegrationID: integrationID,
		WhatsAppURL:   baseURL,
	}

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

	if token, ok := creds["access_token"].(string); ok {
		config.AccessToken = token
	}

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
