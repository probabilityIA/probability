package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	productSyncProgressBatch = 10
	siigoProductsPageSize    = 100
	siigoProductsMaxPages    = 5000
)

type providerUpsertMsg struct {
	BusinessID     uint    `json:"business_id"`
	IntegrationID  uint    `json:"integration_id"`
	SKU            string  `json:"sku"`
	Name           string  `json:"name"`
	TrackInventory bool    `json:"track_inventory"`
	Price          float64 `json:"price"`
	ExternalID     string  `json:"external_id"`
}

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func (uc *invoicingUseCase) listAllSiigoProducts(ctx context.Context, credentials dtos.Credentials) ([]dtos.ProductItem, error) {
	all := make([]dtos.ProductItem, 0)
	for page := 1; page <= siigoProductsMaxPages; page++ {
		var batch []dtos.ProductItem
		var err error
		for attempt := 1; attempt <= 3; attempt++ {
			batch, err = uc.siigoClient.ListProducts(ctx, credentials, page, siigoProductsPageSize)
			if err == nil {
				break
			}
			uc.log.Warn(ctx).Err(err).Int("page", page).Int("attempt", attempt).Msg("Reintentando listar productos de Siigo")
		}
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if len(batch) < siigoProductsPageSize {
			break
		}
	}
	return all, nil
}

type reconcileContext struct {
	probProducts  []dtos.ProductForSync
	siigoProducts []dtos.ProductItem
	rules         []productmatch.Rule
	outcome       productmatch.Outcome
}

func probabilityItems(products []dtos.ProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func siigoItems(products []dtos.ProductItem) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = productmatch.Item{
			SKU:        p.Code,
			Barcode:    p.Barcode,
			ExternalID: p.ID,
			Name:       p.Name,
		}
	}
	return items
}

func (uc *invoicingUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileContext, error) {
	integration, err := uc.integrationCore.GetIntegrationByID(ctx, integrationID)
	if err != nil || integration == nil {
		return nil, fmt.Errorf("integracion no encontrada")
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, fmt.Errorf("la integracion no pertenece al negocio")
	}

	credentials, err := uc.resolveWebhookCredentials(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	siigoProducts, err := uc.listAllSiigoProducts(ctx, credentials)
	if err != nil {
		return nil, fmt.Errorf("listing siigo products: %w", err)
	}
	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}

	rules := productmatch.Sanitize(integration.ProductMatchRules)
	return &reconcileContext{
		probProducts:  probProducts,
		siigoProducts: siigoProducts,
		rules:         rules,
		outcome:       productmatch.Reconcile(rules, probabilityItems(probProducts), siigoItems(siigoProducts)),
	}, nil
}

func (uc *invoicingUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*dtos.ReconcileResult, error) {
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	associated, err := uc.productRepo.ListAssociatedSKUs(ctx, businessID, uint(integIDUint))
	if err != nil {
		return nil, fmt.Errorf("listing associated skus: %w", err)
	}

	result := &dtos.ReconcileResult{
		MatchedNotAssociated: []dtos.ProductBrief{},
		OnlyInProbability:    []dtos.ProductBrief{},
		OnlyInSiigo:          []dtos.ProductBrief{},
		ProbabilityNoSKU:     rc.outcome.ProbabilityUnmatchable,
		SiigoNoSKU:           rc.outcome.ChannelUnmatchable,
		MatchRules:           rc.rules,
	}

	for _, pair := range rc.outcome.Pairs {
		prob := rc.probProducts[pair.ProbabilityIndex]
		channel := rc.siigoProducts[pair.ChannelIndex]
		if associated[normalizeSKU(prob.SKU)] {
			result.Matched++
			continue
		}
		result.MatchedNotAssociated = append(result.MatchedNotAssociated, dtos.ProductBrief{
			SKU:          channel.Code,
			Name:         channel.Name,
			MatchedBy:    pair.Rule.Key(),
			MatchedValue: siigoItems([]dtos.ProductItem{channel})[0].Values()[pair.Rule.Channel],
		})
	}

	for _, idx := range rc.outcome.OnlyInChannel {
		s := rc.siigoProducts[idx]
		result.OnlyInSiigo = append(result.OnlyInSiigo, dtos.ProductBrief{SKU: s.Code, Name: s.Name})
	}
	for _, idx := range rc.outcome.OnlyInProbability {
		p := rc.probProducts[idx]
		result.OnlyInProbability = append(result.OnlyInProbability, dtos.ProductBrief{SKU: p.SKU, Name: p.Name})
	}

	return result, nil
}

func (uc *invoicingUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	rc, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "siigo.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID,
			"direction":      "to_probability",
			"total":          0,
			"created":        0,
			"updated":        0,
			"failed":         0,
			"error":          err.Error(),
		})
		return err
	}

	missing := make([]dtos.ProductItem, 0)
	if len(skus) > 0 {
		want := make(map[string]bool, len(skus))
		for _, s := range skus {
			if key := normalizeSKU(s); key != "" {
				want[key] = true
			}
		}
		for _, s := range rc.siigoProducts {
			if key := normalizeSKU(s.Code); key != "" && want[key] {
				missing = append(missing, s)
			}
		}
	} else {
		for _, idx := range rc.outcome.OnlyInChannel {
			s := rc.siigoProducts[idx]
			if normalizeSKU(s.Code) == "" {
				continue
			}
			missing = append(missing, s)
		}
	}

	total := len(missing)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "siigo.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
	})

	if uc.rabbit != nil {
		_ = uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true)
	}

	created, failed := 0, 0
	for i, s := range missing {
		msg := providerUpsertMsg{
			BusinessID:     businessID,
			IntegrationID:  uint(integIDUint),
			SKU:            s.Code,
			Name:           s.Name,
			TrackInventory: s.StockControl,
			Price:          s.Price,
			ExternalID:     s.ID,
		}
		data, merr := json.Marshal(msg)
		if merr != nil || uc.rabbit == nil {
			failed++
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, data); perr != nil {
			uc.log.Error(ctx).Err(perr).Str("sku", s.Code).Msg("Error al publicar producto para crear en Probability")
			failed++
		} else {
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, 0, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "siigo.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
	return nil
}

func (uc *invoicingUseCase) maybeProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, created, updated, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "siigo.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
}

func (uc *invoicingUseCase) emitSyncEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "siigo",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data:          data,
	})
}
