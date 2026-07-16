package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const productSyncProgressBatch = 10

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

type providerUpsertMsg struct {
	BusinessID     uint    `json:"business_id"`
	IntegrationID  uint    `json:"integration_id"`
	SKU            string  `json:"sku"`
	Name           string  `json:"name"`
	TrackInventory bool    `json:"track_inventory"`
	Price          float64 `json:"price"`
	ExternalID     string  `json:"external_id"`
}

func resolveSellerID(integration *domain.Integration) (int64, error) {
	if integration.StoreID != "" {
		if id, err := strconv.ParseInt(strings.TrimSpace(integration.StoreID), 10, 64); err == nil && id > 0 {
			return id, nil
		}
	}
	if v, ok := integration.Config["seller_id"]; ok {
		switch val := v.(type) {
		case float64:
			return int64(val), nil
		case int64:
			return val, nil
		case int:
			return int64(val), nil
		case string:
			if id, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				return id, nil
			}
		}
	}
	return 0, domain.ErrSellerIDNotFound
}

func (uc *meliUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (accessToken string, sellerID int64, probProducts []domain.ProductForSync, meliProducts []domain.MeliProduct, err error) {
	integration, ierr := uc.service.GetIntegrationByID(ctx, integrationID)
	if ierr != nil {
		return "", 0, nil, nil, fmt.Errorf("getting integration: %w", ierr)
	}
	if integration == nil {
		return "", 0, nil, nil, domain.ErrIntegrationNotFound
	}

	sellerID, err = resolveSellerID(integration)
	if err != nil {
		return "", 0, nil, nil, err
	}

	accessToken, err = uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return "", 0, nil, nil, err
	}

	probProducts, err = uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return "", 0, nil, nil, fmt.Errorf("listing probability products: %w", err)
	}

	meliProducts, err = uc.client.GetProducts(ctx, accessToken, sellerID)
	if err != nil {
		return "", 0, nil, nil, fmt.Errorf("listing meli products: %w", err)
	}

	return accessToken, sellerID, probProducts, meliProducts, nil
}

func (uc *meliUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	_, _, probProducts, meliProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	mapped, merr := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if merr != nil {
		return nil, fmt.Errorf("listing mapped items: %w", merr)
	}
	associatedSKUs := make(map[string]bool)
	for _, m := range mapped {
		if k := normalizeSKU(m.SKU); k != "" {
			associatedSKUs[k] = true
		}
	}

	probBySKU := make(map[string]domain.ProductForSync)
	result := &domain.ReconcileResult{
		MatchedNotAssociated: []domain.ProductBrief{},
		OnlyInProbability:    []domain.ProductBrief{},
		OnlyInMeli:           []domain.ProductBrief{},
	}
	for _, p := range probProducts {
		key := normalizeSKU(p.SKU)
		if key == "" {
			result.ProbabilityNoSKU++
			continue
		}
		probBySKU[key] = p
	}

	meliSKUs := make(map[string]bool)
	for _, m := range meliProducts {
		key := normalizeSKU(m.SKU)
		if key == "" {
			result.MeliNoSKU++
			continue
		}
		meliSKUs[key] = true
		if _, ok := probBySKU[key]; ok {
			if associatedSKUs[key] {
				result.Matched++
			} else {
				result.MatchedNotAssociated = append(result.MatchedNotAssociated, domain.ProductBrief{SKU: m.SKU, Name: m.Name})
			}
		} else {
			result.OnlyInMeli = append(result.OnlyInMeli, domain.ProductBrief{SKU: m.SKU, Name: m.Name})
		}
	}

	for key, p := range probBySKU {
		if !meliSKUs[key] {
			result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
		}
	}

	return result, nil
}

func (uc *meliUseCase) ApplyProductsToMeli(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	accessToken, _, probProducts, meliProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	meliSKUs := make(map[string]bool)
	for _, m := range meliProducts {
		if key := normalizeSKU(m.SKU); key != "" {
			meliSKUs[key] = true
		}
	}

	missing := make([]domain.ProductForSync, 0)
	for _, p := range probProducts {
		key := normalizeSKU(p.SKU)
		if key == "" || meliSKUs[key] {
			continue
		}
		missing = append(missing, p)
	}

	total := len(missing)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_meli",
		"total":          total,
	})

	siteID, currencyID, listingTypeID := uc.resolveProductPublishConfig(ctx, integrationID)

	created, failed := 0, 0
	for i, p := range missing {
		newID, cerr := uc.client.CreateProduct(ctx, accessToken, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			SiteID:        siteID,
			CurrencyID:    currencyID,
			ListingTypeID: listingTypeID,
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en MercadoLibre")
			failed++
		} else {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), newID); merr != nil {
				uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en MeLi pero fallo el mapeo")
			}
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_meli",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
	return nil
}

func (uc *meliUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	_, _, probProducts, meliProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	probSKUs := make(map[string]bool)
	for _, p := range probProducts {
		if key := normalizeSKU(p.SKU); key != "" {
			probSKUs[key] = true
		}
	}

	missing := make([]domain.MeliProduct, 0)
	for _, m := range meliProducts {
		key := normalizeSKU(m.SKU)
		if key == "" || probSKUs[key] {
			continue
		}
		missing = append(missing, m)
	}

	total := len(missing)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
	})

	if uc.rabbit != nil {
		_ = uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true)
	}

	created, failed := 0, 0
	for i, m := range missing {
		msg := providerUpsertMsg{
			BusinessID:     businessID,
			IntegrationID:  uint(integIDUint),
			SKU:            m.SKU,
			Name:           m.Name,
			TrackInventory: true,
			Price:          m.Price,
			ExternalID:     m.ID,
		}
		data, merr := json.Marshal(msg)
		if merr != nil || uc.rabbit == nil {
			failed++
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, data); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", m.SKU).Msg("Error al publicar producto para crear en Probability")
			failed++
		} else {
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
	return nil
}

func (uc *meliUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	_, _, probProducts, meliProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	meliBySKU := make(map[string]string)
	for _, m := range meliProducts {
		if k := normalizeSKU(m.SKU); k != "" && m.ID != "" {
			meliBySKU[k] = m.ID
		}
	}
	probBySKU := make(map[string]domain.ProductForSync)
	for _, p := range probProducts {
		if k := normalizeSKU(p.SKU); k != "" {
			probBySKU[k] = p
		}
	}

	mapped, err := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return err
	}
	associated := make(map[string]bool)
	for _, m := range mapped {
		if k := normalizeSKU(m.SKU); k != "" {
			associated[k] = true
		}
	}

	targets := make([]string, 0)
	if len(skus) > 0 {
		for _, s := range skus {
			if k := normalizeSKU(s); k != "" {
				targets = append(targets, k)
			}
		}
	} else {
		for k := range probBySKU {
			if meliBySKU[k] != "" && !associated[k] {
				targets = append(targets, k)
			}
		}
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
	})

	created, failed := 0, 0
	for i, k := range targets {
		p, okP := probBySKU[k]
		ref, okM := meliBySKU[k]
		if !okP || !okM || ref == "" || associated[k] {
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
			continue
		}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), ref); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto a MercadoLibre")
			failed++
		} else {
			associated[k] = true
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "meli.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
	return nil
}

func (uc *meliUseCase) resolveProductPublishConfig(ctx context.Context, integrationID string) (string, string, string) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil || integration == nil {
		return "", "", ""
	}
	site, _ := integration.Config["meli_site_id"].(string)
	currency, _ := integration.Config["meli_currency_id"].(string)
	listing, _ := integration.Config["meli_listing_type_id"].(string)
	return site, currency, listing
}

func (uc *meliUseCase) maybeProgress(ctx context.Context, businessID, integrationID uint, correlationID string, processed, total, created, failed int) {
	if processed%productSyncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "meli.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
}

func (uc *meliUseCase) emitSyncEvent(ctx context.Context, businessID, integrationID uint, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "meli",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data:          data,
	})
}
