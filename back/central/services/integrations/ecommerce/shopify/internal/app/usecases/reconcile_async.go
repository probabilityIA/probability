package usecases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	channelLabel        = "Shopify"
	reconcileEventType  = "shopify.product.reconcile.completed"
)

type reconcileDetailItem struct {
	SKU   string `json:"sku"`
	Label string `json:"label"`
	Tone  string `json:"tone"`
	Group string `json:"group"`
}

type syncRunEnvelope struct {
	Type          string                 `json:"type"`
	BusinessID    uint                   `json:"business_id"`
	IntegrationID uint                   `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
}

func (uc *SyncOrdersUseCase) ReconcileProductsAsync(ctx context.Context, integrationID string, businessID, integIDUint uint, correlationID string) {
	uc.emitProductEvent(ctx, integIDUint, businessID, "shopify.product.reconcile.started", map[string]interface{}{"correlation_id": correlationID})

	result, err := uc.ReconcileProducts(ctx, integrationID, businessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Error al comparar el catalogo con Shopify")
		payload := map[string]interface{}{
			"correlation_id": correlationID,
			"error":          err.Error(),
		}
		uc.emitProductEvent(ctx, integIDUint, businessID, reconcileEventType, payload)
		uc.publishReconcileRun(ctx, businessID, integIDUint, payload)
		return
	}

	counts := map[string]interface{}{
		"correlation_id":      correlationID,
		"matched":             result.Matched,
		"not_associated":      len(result.MatchedNotAssociated),
		"only_in_probability": len(result.OnlyInProbability),
		"only_in_channel":     len(result.OnlyInShopify),
	}

	uc.emitProductEvent(ctx, integIDUint, businessID, reconcileEventType, counts)

	stored := make(map[string]interface{}, len(counts)+1)
	for key, value := range counts {
		stored[key] = value
	}
	stored["detail"] = reconcileDetail(result)
	uc.publishReconcileRun(ctx, businessID, integIDUint, stored)
}

func reconcileDetail(result *domain.ReconcileResult) []reconcileDetailItem {
	detail := make([]reconcileDetailItem, 0)

	add := func(items []domain.ProductBrief, label, tone, group string) {
		for _, item := range items {
			text := label
			if item.Name != "" {
				text = label + " · " + item.Name
			}
			detail = append(detail, reconcileDetailItem{SKU: item.SKU, Label: text, Tone: tone, Group: group})
		}
	}

	add(result.MatchedNotAssociated, "sin asociar", "warn", "not_associated")
	add(result.MatchedItems, "en ambos", "ok", "both")
	add(result.OnlyInProbability, "solo en Probability", "warn", "only_probability")
	add(result.OnlyInShopify, "solo en "+channelLabel, "warn", "only_channel")

	return detail
}

func (uc *SyncOrdersUseCase) publishReconcileRun(ctx context.Context, businessID, integrationID uint, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	if err := uc.rabbit.DeclareQueue(rabbitmq.QueueIntegrationSyncRuns, true); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al declarar la cola de resultados de sincronizacion")
		return
	}
	payload, err := json.Marshal(syncRunEnvelope{
		Type:          reconcileEventType,
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Timestamp:     time.Now(),
		Data:          data,
	})
	if err != nil {
		return
	}
	if err := uc.rabbit.Publish(ctx, rabbitmq.QueueIntegrationSyncRuns, payload); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al publicar el resultado de la comparacion de catalogo")
	}
}
