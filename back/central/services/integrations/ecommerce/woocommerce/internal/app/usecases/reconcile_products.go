package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func wooExternalRef(w domain.WooProduct) string {
	if w.ParentID != "" {
		return w.ParentID + ":" + w.ID
	}
	return w.ID
}

func fullImageURL(imageURL string) string {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" || strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return imageURL
	}
	base := strings.TrimRight(os.Getenv("URL_BASE_DOMAIN_S3"), "/")
	if base == "" {
		return ""
	}
	return base + "/" + strings.TrimLeft(imageURL, "/")
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

func (uc *wooCommerceUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (storeURL, ck, cs string, probProducts []domain.ProductForSync, wooProducts []domain.WooProduct, err error) {
	storeURL, ck, cs, err = uc.resolveStoreCreds(ctx, integrationID)
	if err != nil {
		return "", "", "", nil, nil, err
	}
	probProducts, err = uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return "", "", "", nil, nil, fmt.Errorf("listing probability products: %w", err)
	}
	wooProducts, err = uc.client.GetProducts(ctx, storeURL, ck, cs)
	if err != nil {
		return "", "", "", nil, nil, fmt.Errorf("listing woocommerce products: %w", err)
	}
	return storeURL, ck, cs, probProducts, wooProducts, nil
}

func (uc *wooCommerceUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	_, _, _, probProducts, wooProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	probBySKU := make(map[string]domain.ProductForSync)
	result := &domain.ReconcileResult{
		OnlyInProbability: []domain.ProductBrief{},
		OnlyInWoo:         []domain.ProductBrief{},
	}
	for _, p := range probProducts {
		key := normalizeSKU(p.SKU)
		if key == "" {
			result.ProbabilityNoSKU++
			continue
		}
		probBySKU[key] = p
	}

	wooSKUs := make(map[string]bool)
	for _, w := range wooProducts {
		key := normalizeSKU(w.SKU)
		if key == "" {
			result.WooNoSKU++
			continue
		}
		wooSKUs[key] = true
		if _, ok := probBySKU[key]; ok {
			result.Matched++
		} else {
			result.OnlyInWoo = append(result.OnlyInWoo, domain.ProductBrief{SKU: w.SKU, Name: w.Name})
		}
	}

	for key, p := range probBySKU {
		if !wooSKUs[key] {
			result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
		}
	}

	return result, nil
}

func (uc *wooCommerceUseCase) ApplyProductsToWoo(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	storeURL, ck, cs, probProducts, wooProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	wooBySKU := make(map[string]string)
	for _, w := range wooProducts {
		if key := normalizeSKU(w.SKU); key != "" && w.ID != "" {
			wooBySKU[key] = wooExternalRef(w)
		}
	}

	targets := make([]domain.ProductForSync, 0)
	for _, p := range probProducts {
		if normalizeSKU(p.SKU) == "" {
			continue
		}
		targets = append(targets, p)
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_woo",
		"total":          total,
	})

	created, updated, failed := 0, 0, 0
	for i, p := range targets {
		if wooID, ok := wooBySKU[normalizeSKU(p.SKU)]; ok && wooID != "" {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), wooID); merr != nil {
				uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al mapear producto existente de WooCommerce")
				failed++
			} else {
				updated++
			}
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
			continue
		}

		newID, cerr := uc.client.CreateProduct(ctx, storeURL, ck, cs, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			ManageStock:   p.TrackInventory,
			ImageURL:      fullImageURL(p.ImageURL),
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en WooCommerce")
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
			continue
		}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), newID); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en Woo pero fallo el mapeo")
			failed++
			uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
			continue
		}
		created++
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, updated, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_woo",
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
	return nil
}

func (uc *wooCommerceUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	_, _, _, probProducts, wooProducts, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	probSKUs := make(map[string]bool)
	for _, p := range probProducts {
		if key := normalizeSKU(p.SKU); key != "" {
			probSKUs[key] = true
		}
	}

	missing := make([]domain.WooProduct, 0)
	for _, w := range wooProducts {
		key := normalizeSKU(w.SKU)
		if key == "" || probSKUs[key] {
			continue
		}
		missing = append(missing, w)
	}

	total := len(missing)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
	})

	if uc.rabbit != nil {
		_ = uc.rabbit.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true)
	}

	created, failed := 0, 0
	for i, w := range missing {
		msg := providerUpsertMsg{
			BusinessID:     businessID,
			IntegrationID:  uint(integIDUint),
			SKU:            w.SKU,
			Name:           w.Name,
			TrackInventory: true,
			Price:          w.Price,
			ExternalID:     wooExternalRef(w),
		}
		data, merr := json.Marshal(msg)
		if merr != nil || uc.rabbit == nil {
			failed++
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, data); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", w.SKU).Msg("Error al publicar producto para crear en Probability")
			failed++
		} else {
			created++
		}
		uc.maybeProgress(ctx, businessID, uint(integIDUint), correlationID, i+1, total, created, 0, failed)
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "woocommerce.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "to_probability",
		"total":          total,
		"created":        created,
		"updated":        0,
		"failed":         failed,
	})
	return nil
}
