package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
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
		batch, err := uc.siigoClient.ListProducts(ctx, credentials, page, siigoProductsPageSize)
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

func (uc *invoicingUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) ([]dtos.ProductForSync, []dtos.ProductItem, error) {
	credentials, err := uc.resolveWebhookCredentials(ctx, integrationID)
	if err != nil {
		return nil, nil, err
	}
	siigoProducts, err := uc.listAllSiigoProducts(ctx, credentials)
	if err != nil {
		return nil, nil, fmt.Errorf("listing siigo products: %w", err)
	}
	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, nil, fmt.Errorf("listing probability products: %w", err)
	}
	return probProducts, siigoProducts, nil
}

func (uc *invoicingUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*dtos.ReconcileResult, error) {
	probProducts, siigoProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	result := &dtos.ReconcileResult{
		OnlyInProbability: []dtos.ProductBrief{},
		OnlyInSiigo:       []dtos.ProductBrief{},
	}

	probBySKU := make(map[string]dtos.ProductForSync)
	for _, p := range probProducts {
		key := normalizeSKU(p.SKU)
		if key == "" {
			result.ProbabilityNoSKU++
			continue
		}
		probBySKU[key] = p
	}

	siigoSKUs := make(map[string]bool)
	for _, s := range siigoProducts {
		key := normalizeSKU(s.Code)
		if key == "" {
			result.SiigoNoSKU++
			continue
		}
		siigoSKUs[key] = true
		if _, ok := probBySKU[key]; ok {
			result.Matched++
		} else {
			result.OnlyInSiigo = append(result.OnlyInSiigo, dtos.ProductBrief{SKU: s.Code, Name: s.Name})
		}
	}

	for key, p := range probBySKU {
		if !siigoSKUs[key] {
			result.OnlyInProbability = append(result.OnlyInProbability, dtos.ProductBrief{SKU: p.SKU, Name: p.Name})
		}
	}

	return result, nil
}

func (uc *invoicingUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	probProducts, siigoProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	probSKUs := make(map[string]bool)
	for _, p := range probProducts {
		if key := normalizeSKU(p.SKU); key != "" {
			probSKUs[key] = true
		}
	}

	missing := make([]dtos.ProductItem, 0)
	for _, s := range siigoProducts {
		key := normalizeSKU(s.Code)
		if key == "" || probSKUs[key] {
			continue
		}
		missing = append(missing, s)
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
