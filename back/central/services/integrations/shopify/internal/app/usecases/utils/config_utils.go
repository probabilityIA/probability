package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

func NormalizeConfig(config interface{}, integrationName string) (map[string]interface{}, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok || configMap == nil {
		if integrationName != "" {
			return map[string]interface{}{
				"store_name": integrationName,
			}, nil
		}
		return nil, fmt.Errorf("invalid integration config format and Name is empty")
	}
	return configMap, nil
}

func ExtractStoreName(config map[string]interface{}, integrationName string) (string, error) {
	if url, ok := config["store_url"].(string); ok && url != "" {
		url = strings.TrimPrefix(url, "https://")
		url = strings.TrimPrefix(url, "http://")
		url = strings.TrimSuffix(url, "/")
		return url, nil
	}

	if name, ok := config["store_name"].(string); ok && name != "" {
		return name, nil
	}

	if integrationName != "" {
		return integrationName, nil
	}

	return "", fmt.Errorf("store_url, store_name not found in config and integration name is empty")
}

func GetAccessToken(ctx context.Context, integrationService domain.IIntegrationService, integrationID string) (string, error) {
	accessToken, err := integrationService.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return "", fmt.Errorf("failed to decrypt access_token: %w", err)
	}
	return accessToken, nil
}
