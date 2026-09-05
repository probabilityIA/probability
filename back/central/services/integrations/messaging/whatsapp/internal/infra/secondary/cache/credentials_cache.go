package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

	return c.GetWhatsAppConfigByIntegration(ctx, integrationID)
}

func (c *credentialsCache) GetWhatsAppConfigByIntegration(ctx context.Context, integrationID uint) (*ports.WhatsAppConfig, error) {
	resolver := c.getResolver()
	if resolver == nil {
		return nil, fmt.Errorf("whatsapp resolver not configured")
	}

	if integrationID == 0 {
		return nil, fmt.Errorf("integración WhatsApp no encontrada")
	}

	integrationConfig, integrationCreds, err := resolver.GetIntegrationConfigAndCredentials(ctx, integrationID)
	if err != nil {
		c.log.Warn(ctx).Err(err).
			Uint("integration_id", integrationID).
			Msg("no se pudo leer la configuración propia de la integración, se usa el número de la plataforma")
		return c.platformConfigFor(ctx, integrationID)
	}

	merged := mergeCredentialMaps(integrationConfig, integrationCreds)

	if usePlatformToken(merged) {
		return c.platformConfigFor(ctx, integrationID)
	}

	phoneNumberID := stringValue(merged["phone_number_id"])
	if phoneNumberID == "" {
		return c.platformConfigFor(ctx, integrationID)
	}

	config, err := buildWhatsAppConfig(merged, integrationID, "")
	if err != nil {
		return nil, fmt.Errorf("integración WhatsApp %d con número propio mal configurada: %w", integrationID, err)
	}

	if config.WhatsAppURL == "" {
		if platform, platErr := c.GetWhatsAppDefaultConfig(ctx); platErr == nil {
			config.WhatsAppURL = platform.WhatsAppURL
		}
	}

	config.OwnNumber = true

	c.log.Debug(ctx).
		Uint("integration_id", integrationID).
		Uint("phone_number_id", config.PhoneNumberID).
		Str("waba_id", config.WABAID).
		Msg("configuración WhatsApp propia del negocio resuelta")

	return config, nil
}

func (c *credentialsCache) platformConfigFor(ctx context.Context, integrationID uint) (*ports.WhatsAppConfig, error) {
	config, err := c.GetWhatsAppDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("credenciales de plataforma no encontradas para WhatsApp (integración %d): %w", integrationID, err)
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

func (c *credentialsCache) ResolveByPhoneNumberID(ctx context.Context, phoneNumberID string) (*ports.IntegrationOwner, error) {
	resolver := c.getResolver()
	if resolver == nil {
		return nil, fmt.Errorf("whatsapp resolver not configured")
	}

	phoneNumberID = strings.TrimSpace(phoneNumberID)
	if phoneNumberID == "" {
		return nil, nil
	}

	integrationID, businessID, err := resolver.FindIntegrationByConfigValue(ctx, whatsAppTypeID, "phone_number_id", phoneNumberID)
	if err != nil {
		return nil, err
	}
	if integrationID == 0 {
		return nil, nil
	}

	return &ports.IntegrationOwner{IntegrationID: integrationID, BusinessID: businessID}, nil
}

func mergeCredentialMaps(config, credentials map[string]any) map[string]any {
	merged := make(map[string]any, len(config)+len(credentials))
	for key, value := range config {
		merged[key] = value
	}
	for key, value := range credentials {
		if value == nil {
			continue
		}
		if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
			continue
		}
		merged[key] = value
	}
	return merged
}

func usePlatformToken(merged map[string]any) bool {
	switch v := merged["use_platform_token"].(type) {
	case bool:
		return v
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(v))
		return err == nil && parsed
	}
	return false
}

func stringValue(raw any) string {
	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	}
	return ""
}

func buildWhatsAppConfig(creds map[string]any, integrationID uint, baseURL string) (*ports.WhatsAppConfig, error) {
	config := &ports.WhatsAppConfig{
		IntegrationID: integrationID,
		WhatsAppURL:   baseURL,
	}

	if phoneID := stringValue(creds["phone_number_id"]); phoneID != "" {
		parsed, err := strconv.ParseUint(phoneID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("phone_number_id inválido: %s", phoneID)
		}
		config.PhoneNumberID = uint(parsed)
	}

	if token, ok := creds["access_token"].(string); ok {
		config.AccessToken = strings.TrimSpace(token)
	}

	config.WABAID = stringValue(creds["waba_id"])

	if url, ok := creds["whatsapp_url"].(string); ok && config.WhatsAppURL == "" {
		config.WhatsAppURL = strings.TrimSpace(url)
	}

	if config.PhoneNumberID == 0 {
		return nil, fmt.Errorf("phone_number_id no encontrado en credenciales")
	}
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access_token no encontrado en credenciales")
	}

	return config, nil
}
