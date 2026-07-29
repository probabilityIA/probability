package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	probabilityWeightUnit    = "kg"
	probabilityDimensionUnit = "cm"
)

type providerUpsertMsg struct {
	BusinessID     uint     `json:"business_id"`
	IntegrationID  uint     `json:"integration_id"`
	SKU            string   `json:"sku"`
	Name           string   `json:"name"`
	TrackInventory bool     `json:"track_inventory"`
	Price          float64  `json:"price"`
	ExternalID     string   `json:"external_id"`
	Weight         *float64 `json:"weight,omitempty"`
	WeightUnit     string   `json:"weight_unit,omitempty"`
	Length         *float64 `json:"length,omitempty"`
	Width          *float64 `json:"width,omitempty"`
	Height         *float64 `json:"height,omitempty"`
	DimensionUnit  string   `json:"dimension_unit,omitempty"`
	ImageURL       string   `json:"image_url,omitempty"`
}

type reconcileData struct {
	cred        domain.Credential
	integration *domain.Integration
	probability []domain.ProductForSync
	vtex        []domain.VTEXSKU
	rules       []productmatch.Rule
	outcome     productmatch.Outcome
}

func vtexRefs(s domain.VTEXSKU) productmatch.ExternalRefs {
	return productmatch.ExternalRefs{
		ProductID: s.ID,
		VariantID: s.ID,
		SKU:       s.RefID,
		Barcode:   s.EAN,
	}
}

func probabilityItems(products []domain.ProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func vtexItems(skus []domain.VTEXSKU) []productmatch.Item {
	items := make([]productmatch.Item, len(skus))
	for i, s := range skus {
		items[i] = s.MatchItem()
	}
	return items
}

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func (uc *vtexUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileData, error) {
	integration, err := uc.integrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	cred, err := uc.resolveCredential(ctx, integration, integrationID)
	if err != nil {
		return nil, err
	}

	probabilityProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}

	vtexSKUs, err := uc.client.ListSKUs(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("listing vtex skus: %w", err)
	}

	rules := productmatch.Sanitize(integration.ProductMatchRules)
	return &reconcileData{
		cred:        cred,
		integration: integration,
		probability: probabilityProducts,
		vtex:        vtexSKUs,
		rules:       rules,
		outcome:     productmatch.Reconcile(rules, probabilityItems(probabilityProducts), vtexItems(vtexSKUs)),
	}, nil
}

func (uc *vtexUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	mapped, err := uc.productRepo.ListMappedItems(ctx, data.integration.ID)
	if err != nil {
		return nil, fmt.Errorf("listing mapped items: %w", err)
	}
	associated := make(map[string]bool, len(mapped))
	for _, m := range mapped {
		associated[normalizeSKU(m.SKU)] = true
	}

	result := &domain.ReconcileResult{
		ProbabilityNoSKU: data.outcome.ProbabilityUnmatchable,
		VTEXNoSKU:        data.outcome.ChannelUnmatchable,
		MatchRules:       data.rules,
	}

	for _, pair := range data.outcome.Pairs {
		prob := data.probability[pair.ProbabilityIndex]
		channel := data.vtex[pair.ChannelIndex]
		brief := domain.ProductBrief{
			SKU:          prob.SKU,
			Name:         channel.Name,
			MatchedBy:    pair.Rule.Key(),
			MatchedValue: channel.MatchItem().Values()[pair.Rule.Channel],
		}
		result.Matched++
		result.MatchedItems = append(result.MatchedItems, brief)
		if !associated[normalizeSKU(prob.SKU)] {
			result.MatchedNotAssociated = append(result.MatchedNotAssociated, brief)
		}
	}

	for _, idx := range data.outcome.OnlyInProbability {
		p := data.probability[idx]
		result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
	}
	for _, idx := range data.outcome.OnlyInChannel {
		sku := data.vtex[idx]
		result.OnlyInVTEX = append(result.OnlyInVTEX, domain.ProductBrief{SKU: sku.RefID, Name: sku.Name})
	}

	return result, nil
}

func (uc *vtexUseCase) SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	if err := uc.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID); err != nil {
		return err
	}
	return uc.AssociateProducts(ctx, integrationID, businessID, correlationID, nil)
}

func selectedSKUs(skus []string) map[string]bool {
	if len(skus) == 0 {
		return nil
	}
	set := make(map[string]bool, len(skus))
	for _, s := range skus {
		if key := normalizeSKU(s); key != "" {
			set[key] = true
		}
	}
	return set
}

func (uc *vtexUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	existing := make(map[string]bool, len(data.probability))
	for _, p := range data.probability {
		existing[normalizeSKU(p.SKU)] = true
	}

	only := selectedSKUs(skus)
	var toCreate []domain.VTEXSKU
	for _, sku := range data.vtex {
		key := normalizeSKU(sku.RefID)
		if key == "" || existing[key] || (only != nil && !only[key]) {
			continue
		}
		toCreate = append(toCreate, sku)
	}

	return uc.publishUpserts(ctx, businessID, data.integration.ID, correlationID, domain.ModeCreate, toCreate)
}

func (uc *vtexUseCase) UpdateProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	existing := make(map[string]bool, len(data.probability))
	for _, p := range data.probability {
		existing[normalizeSKU(p.SKU)] = true
	}

	only := selectedSKUs(skus)
	var toUpdate []domain.VTEXSKU
	for _, sku := range data.vtex {
		key := normalizeSKU(sku.RefID)
		if key == "" || !existing[key] || (only != nil && !only[key]) {
			continue
		}
		toUpdate = append(toUpdate, sku)
	}

	return uc.publishUpserts(ctx, businessID, data.integration.ID, correlationID, domain.ModeUpdate, toUpdate)
}

func (uc *vtexUseCase) upsertMsgFromVTEX(businessID, integrationID uint, sku domain.VTEXSKU) providerUpsertMsg {
	msg := providerUpsertMsg{
		BusinessID:     businessID,
		IntegrationID:  integrationID,
		SKU:            sku.RefID,
		Name:           sku.Name,
		TrackInventory: true,
		Price:          sku.Price,
		ExternalID:     sku.ID,
		Weight:         sku.Weight,
		Length:         sku.Length,
		Width:          sku.Width,
		Height:         sku.Height,
		ImageURL:       sku.ImageURL,
	}

	if sku.Weight != nil {
		msg.WeightUnit = probabilityWeightUnit
	}
	if sku.Length != nil || sku.Width != nil || sku.Height != nil {
		msg.DimensionUnit = probabilityDimensionUnit
	}

	return msg
}

func (uc *vtexUseCase) publishUpserts(ctx context.Context, businessID, integrationID uint, correlationID, mode string, items []domain.VTEXSKU) error {
	total := len(items)

	uc.emitSyncEvent(ctx, businessID, integrationID, "vtex.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      domain.DirectionToProbability,
		"mode":           mode,
		"total":          total,
	})

	if uc.rabbit == nil {
		uc.emitSyncEvent(ctx, businessID, integrationID, "vtex.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID,
			"direction":      domain.DirectionToProbability,
			"mode":           mode,
			"total":          total,
			"created":        0,
			"failed":         total,
		})
		return fmt.Errorf("RabbitMQ no disponible: no se pueden sincronizar productos hacia Probability")
	}

	if derr := uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true); derr != nil {
		return fmt.Errorf("declarando la cola de upsert de productos: %w", derr)
	}

	fails := &failedSKUs{}
	applied := 0

	for i, item := range items {
		msg := uc.upsertMsgFromVTEX(businessID, integrationID, item)

		payload, merr := json.Marshal(msg)
		if merr != nil {
			fails.add(item.RefID)
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, payload); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", item.RefID).Msg("Error al publicar producto hacia Probability")
			fails.add(item.RefID)
		} else {
			applied++
		}

		uc.maybeProductProgress(ctx, businessID, integrationID, correlationID, domain.DirectionToProbability, i+1, total, applied, 0, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, integrationID, "vtex.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      domain.DirectionToProbability,
		"mode":           mode,
		"total":          total,
		"created":        applied,
		"updated":        applied,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *vtexUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	filter := make(map[string]bool, len(skus))
	for _, s := range skus {
		filter[normalizeSKU(s)] = true
	}

	total := 0
	associated := 0
	fails := &failedSKUs{}

	for _, pair := range data.outcome.Pairs {
		p := data.probability[pair.ProbabilityIndex]
		key := normalizeSKU(p.SKU)
		if len(filter) > 0 && (key == "" || !filter[key]) {
			continue
		}

		total++
		refs := vtexRefs(data.vtex[pair.ChannelIndex])
		if err := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, data.integration.ID, refs); err != nil {
			uc.logger.Error(ctx).Err(err).Str("sku", p.SKU).Msg("Error al asociar producto con VTEX")
			fails.add(p.SKU)
			continue
		}
		associated++
	}

	uc.emitSyncEvent(ctx, businessID, data.integration.ID, "vtex.product.associate.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"total":          total,
		"associated":     associated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	uc.logger.Info(ctx).
		Int("total", total).
		Int("associated", associated).
		Int("failed", fails.count()).
		Msg("VTEX product association completed")

	return nil
}
