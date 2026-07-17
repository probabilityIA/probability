package usecases

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

var protocolPrefix = regexp.MustCompile(`^https?://`)

var vtexSuffixes = []string{
	".myvtex.com",
	".vtexcommercestable.com.br",
	".vtexcommercebeta.com.br",
	".vtexlocal.com.br",
	".vtexcommerce.com.br",
	".vtex.com",
}

func CleanAccountName(raw string) string {
	account := strings.ToLower(strings.TrimSpace(raw))
	account = protocolPrefix.ReplaceAllString(account, "")
	account = strings.Split(account, "/")[0]
	account = strings.Split(account, "?")[0]
	account = strings.Split(account, ":")[0]

	for _, suffix := range vtexSuffixes {
		if strings.HasSuffix(account, suffix) {
			account = strings.TrimSuffix(account, suffix)
			break
		}
	}

	parts := strings.Split(account, ".")
	if len(parts) > 0 {
		account = parts[0]
	}
	return account
}

func extractString(m map[string]interface{}, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("missing field: %s", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("field %s must be a non-empty string", key)
	}
	return s, nil
}

func extractBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}

func (uc *vtexUseCase) resolveCredential(ctx context.Context, integration *domain.Integration, integrationID string) (domain.Credential, error) {
	accountName, err := extractString(integration.Config, "account_name")
	if err != nil {
		if integration.StoreID == "" {
			return domain.Credential{}, domain.ErrMissingAccountName
		}
		accountName = integration.StoreID
	}

	accountName = CleanAccountName(accountName)
	if accountName == "" {
		return domain.Credential{}, domain.ErrMissingAccountName
	}

	appKey, err := uc.service.DecryptCredential(ctx, integrationID, "app_key")
	if err != nil {
		return domain.Credential{}, fmt.Errorf("decrypting app_key: %w", err)
	}
	if appKey == "" {
		return domain.Credential{}, domain.ErrMissingAppKey
	}

	appToken, err := uc.service.DecryptCredential(ctx, integrationID, "app_token")
	if err != nil {
		return domain.Credential{}, fmt.Errorf("decrypting app_token: %w", err)
	}
	if appToken == "" {
		return domain.Credential{}, domain.ErrMissingAppToken
	}

	return domain.Credential{
		AccountName: accountName,
		AppKey:      appKey,
		AppToken:    appToken,
	}, nil
}

func (uc *vtexUseCase) inventoryConfigFrom(config map[string]interface{}) domain.InventoryConfig {
	cfg := domain.InventoryConfig{
		Enabled:  extractBool(config, "inventory_sync_enabled"),
		IsSeller: extractBool(config, "is_seller"),
	}

	raw, ok := config["vtex_warehouse_mappings"].([]interface{})
	if !ok {
		return cfg
	}

	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		mapping := domain.WarehouseMapping{}
		if v, ok := m["vtex_warehouse_id"].(string); ok {
			mapping.VTEXWarehouseID = strings.TrimSpace(v)
		}
		switch v := m["internal_warehouse_id"].(type) {
		case float64:
			mapping.InternalWarehouseID = uint(v)
		case int:
			mapping.InternalWarehouseID = uint(v)
		}
		if mapping.VTEXWarehouseID == "" || mapping.InternalWarehouseID == 0 {
			continue
		}
		cfg.WarehouseMappings = append(cfg.WarehouseMappings, mapping)
	}

	return cfg
}
