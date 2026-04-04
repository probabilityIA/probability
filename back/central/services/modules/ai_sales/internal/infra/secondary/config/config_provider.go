package config

import (
	"context"
	"encoding/json"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

const (
	// platformCredsKey es la clave Redis donde se guardan las platform_credentials del IntegrationType WhatsApp (id=2)
	platformCredsKey = "integration:platform_creds:2"
)

func (c *configProvider) GetAIConfig(ctx context.Context) (*domain.AIConfig, error) {
	data, err := c.redis.Get(ctx, platformCredsKey)
	if err != nil || data == "" {
		c.log.Warn(ctx).Msg("No se encontraron platform_creds para WhatsApp en Redis")
		return defaultConfig(), nil
	}

	var creds map[string]interface{}
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error parseando platform_creds de Redis")
		return defaultConfig(), nil
	}

	config := &domain.AIConfig{
		Enabled:           getBool(creds, "ai_sales_enabled", false),
		ModelID:           getString(creds, "ai_sales_model_id", "amazon.nova-micro-v1:0"),
		SessionTTLMinutes: getInt(creds, "ai_sales_session_ttl_minutes", 20),
		MaxToolIterations: getInt(creds, "ai_sales_max_tool_iterations", 5),
		DemoBusinessID:    getUint(creds, "ai_sales_demo_business_id", 1),
	}

	return config, nil
}

func defaultConfig() *domain.AIConfig {
	return &domain.AIConfig{
		Enabled:           false,
		ModelID:           "amazon.nova-micro-v1:0",
		SessionTTLMinutes: 20,
		MaxToolIterations: 5,
		DemoBusinessID:    1,
	}
}

func getBool(m map[string]interface{}, key string, def bool) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}

func getString(m map[string]interface{}, key string, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

func getInt(m map[string]interface{}, key string, def int) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return def
}

func getUint(m map[string]interface{}, key string, def uint) uint {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			if n > 0 {
				return uint(n)
			}
		}
	}
	return def
}

