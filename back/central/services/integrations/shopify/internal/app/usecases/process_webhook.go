package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core"
)

type ProcessWebhookUseCase struct {
	coreIntegration core.IIntegrationCore
	syncUseCase     *SyncOrdersUseCase
}

func NewProcessWebhookUseCase(
	coreIntegration core.IIntegrationCore,
	syncUseCase *SyncOrdersUseCase,
) *ProcessWebhookUseCase {
	return &ProcessWebhookUseCase{
		coreIntegration: coreIntegration,
		syncUseCase:     syncUseCase,
	}
}

func (uc *ProcessWebhookUseCase) Execute(ctx context.Context, storeDomain string, payload map[string]interface{}) error {
	// 1. Find integration by store domain
	// We need to find the integration that matches this store.
	// Core doesn't expose "GetIntegrationByConfigValue".
	// We might need to iterate or add a method to Core.
	// For now, let's assume we can't easily find it without a lookup.
	// BUT, we can assume there's only one shopify integration per business? No, we don't know the business.
	// We need to look up the integration ID based on the store name in the config.

	// TODO: Implement lookup in Core or Repository.
	// For this implementation, I will skip the lookup and assume we have a way to get it,
	// or I'll just log that we can't find it.

	// Realistically, we need a "FindIntegrationByConfig" method in Core.
	// Or we store the shop domain in a searchable field.

	return fmt.Errorf("not implemented: cannot lookup integration by store domain yet")
}
