package app

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	channelLabel       = "Siigo"
	reconcileEventType = "siigo.product.reconcile.completed"
)

type reconcileDetailItem struct {
	SKU          string `json:"sku"`
	Label        string `json:"label"`
	Tone         string `json:"tone"`
	Group        string `json:"group"`
	MatchedBy    string `json:"matched_by,omitempty"`
	MatchedValue string `json:"matched_value,omitempty"`

	CounterpartSKU  string `json:"counterpart_sku,omitempty"`
	CounterpartName string `json:"counterpart_name,omitempty"`
	ChannelQty      *int   `json:"channel_qty,omitempty"`
	OwnQty          *int   `json:"own_qty,omitempty"`
	FixSide         string `json:"fix_side,omitempty"`
	Pattern         string `json:"pattern,omitempty"`
}

type syncRunEnvelope struct {
	Type          string                 `json:"type"`
	BusinessID    uint                   `json:"business_id"`
	IntegrationID uint                   `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
}

func (uc *invoicingUseCase) ReconcileProductsAsync(ctx context.Context, integrationID string, businessID, integIDUint uint, correlationID string) {
	uc.emitSyncEvent(ctx, businessID, integIDUint, "siigo.product.reconcile.started", map[string]interface{}{
		"correlation_id": correlationID,
	})

	result, err := uc.ReconcileProducts(ctx, integrationID, businessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Error al comparar el catalogo con Siigo")
		payload := map[string]interface{}{
			"correlation_id": correlationID,
			"error":          err.Error(),
		}
		uc.emitSyncEvent(ctx, businessID, integIDUint, reconcileEventType, payload)
		uc.publishReconcileRun(ctx, businessID, integIDUint, payload)
		return
	}

	uc.recordReconcileRun(ctx, businessID, integIDUint, correlationID, result)
}

func (uc *invoicingUseCase) ReconcileProductsAndRecord(ctx context.Context, integrationID string, businessID uint) (*dtos.ReconcileResult, error) {
	result, err := uc.ReconcileProducts(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	uc.recordReconcileRun(ctx, businessID, uint(integIDUint), "", result)
	return result, nil
}

func (uc *invoicingUseCase) recordReconcileRun(ctx context.Context, businessID, integIDUint uint, correlationID string, result *dtos.ReconcileResult) {
	counts := map[string]interface{}{
		"correlation_id":      correlationID,
		"matched":             result.Matched,
		"not_associated":      len(result.MatchedNotAssociated),
		"only_in_probability": len(result.OnlyInProbability),
		"only_in_channel":     len(result.OnlyInSiigo),
		"probability_no_sku":  result.ProbabilityNoSKU,
		"channel_no_sku":      result.SiigoNoSKU,
		"sku_typo":            len(result.TypoSuspects) - productmatch.CountPattern(result.TypoSuspects, productmatch.PatternSpacing),
		"sku_spacing":         productmatch.CountPattern(result.TypoSuspects, productmatch.PatternSpacing),
		"match_rules":         result.MatchRules,
	}

	uc.emitSyncEvent(ctx, businessID, integIDUint, reconcileEventType, counts)

	stored := make(map[string]interface{}, len(counts)+1)
	for key, value := range counts {
		stored[key] = value
	}
	stored["detail"] = reconcileDetail(result)
	uc.publishReconcileRun(ctx, businessID, integIDUint, stored)
}

func reconcileDetail(result *dtos.ReconcileResult) []reconcileDetailItem {
	detail := make([]reconcileDetailItem, 0)

	add := func(items []dtos.ProductBrief, label, tone, group string) {
		for _, item := range items {
			text := label
			if item.Name != "" {
				text = label + " · " + item.Name
			}
			detail = append(detail, reconcileDetailItem{
				SKU:          item.SKU,
				Label:        text,
				Tone:         tone,
				Group:        group,
				MatchedBy:    item.MatchedBy,
				MatchedValue: item.MatchedValue,
			})
		}
	}

	add(result.MatchedNotAssociated, "sin asociar", "warn", "not_associated")
	add(result.MatchedItems, "en ambos", "ok", "both")
	add(result.OnlyInProbability, "solo en Probability", "warn", "only_probability")
	add(result.OnlyInSiigo, "solo en "+channelLabel, "warn", "only_channel")

	for _, d := range productmatch.TypoDetails(result.TypoSuspects, channelLabel) {
		detail = append(detail, reconcileDetailItem{
			SKU:             d.SKU,
			Label:           d.Label,
			Tone:            d.Tone,
			Group:           d.Group,
			MatchedBy:       d.SKU,
			MatchedValue:    d.CounterpartSKU,
			CounterpartSKU:  d.CounterpartSKU,
			CounterpartName: d.CounterpartName,
			ChannelQty:      d.ChannelQty,
			OwnQty:          d.OwnQty,
			FixSide:         d.FixSide,
			Pattern:         d.Pattern,
		})
	}

	return detail
}

func (uc *invoicingUseCase) publishReconcileRun(ctx context.Context, businessID, integrationID uint, data map[string]interface{}) {
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
