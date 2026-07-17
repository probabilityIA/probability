package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	DirectionToJumpseller  = "to_jumpseller"
	DirectionToProbability = "to_probability"

	ModeCreate = "create"
	ModeUpdate = "update"
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
}

type jumpsellerSKU struct {
	SKU        string
	Name       string
	Price      float64
	ExternalID string
	ProductID  int64
	VariantID  int64
	Weight     float64
	Height     float64
	Width      float64
	Length     float64
}

type reconcileData struct {
	cred            domain.Credential
	storeWeightUnit string
	probability     []domain.ProductForSync
	jumpseller      []jumpsellerSKU
}

func normalizeSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func flattenProductSKUs(products []domain.JumpsellerProduct) []jumpsellerSKU {
	flat := make([]jumpsellerSKU, 0, len(products))
	for _, product := range products {
		if len(product.Variants) == 0 {
			flat = append(flat, jumpsellerSKU{
				SKU:        product.SKU,
				Name:       product.Name,
				Price:      product.Price,
				ExternalID: strconv.FormatInt(product.ID, 10),
				ProductID:  product.ID,
				Weight:     product.Weight,
				Height:     product.Height,
				Width:      product.Width,
				Length:     product.Length,
			})
			continue
		}
		for _, variant := range product.Variants {
			flat = append(flat, jumpsellerSKU{
				SKU:        variant.SKU,
				Name:       product.Name,
				Price:      variant.Price,
				ExternalID: strconv.FormatInt(product.ID, 10) + ":" + strconv.FormatInt(variant.ID, 10),
				ProductID:  product.ID,
				VariantID:  variant.ID,
				Weight:     product.Weight,
				Height:     product.Height,
				Width:      product.Width,
				Length:     product.Length,
			})
		}
	}
	return flat
}

func indexJumpsellerBySKU(products []jumpsellerSKU) map[string]jumpsellerSKU {
	index := make(map[string]jumpsellerSKU, len(products))
	for _, product := range products {
		if key := normalizeSKU(product.SKU); key != "" {
			index[key] = product
		}
	}
	return index
}

func indexProbabilityBySKU(products []domain.ProductForSync) map[string]domain.ProductForSync {
	index := make(map[string]domain.ProductForSync, len(products))
	for _, product := range products {
		if key := normalizeSKU(product.SKU); key != "" {
			index[key] = product
		}
	}
	return index
}

func (uc *jumpsellerUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileData, error) {
	_, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("obteniendo informacion de la tienda Jumpseller: %w", err)
	}

	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}

	jsProducts, err := uc.client.GetProducts(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("listing jumpseller products: %w", err)
	}

	return &reconcileData{
		cred:            cred,
		storeWeightUnit: storeInfo.WeightUnit,
		probability:     probProducts,
		jumpseller:      flattenProductSKUs(jsProducts),
	}, nil
}

func (uc *jumpsellerUseCase) associatedSKUs(ctx context.Context, integrationID uint) (map[string]bool, error) {
	mapped, err := uc.productRepo.ListMappedItems(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("listing mapped items: %w", err)
	}
	associated := make(map[string]bool, len(mapped))
	for _, item := range mapped {
		if key := normalizeSKU(item.SKU); key != "" {
			associated[key] = true
		}
	}
	return associated, nil
}

func (uc *jumpsellerUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	associated, err := uc.associatedSKUs(ctx, uint(integIDUint))
	if err != nil {
		return nil, err
	}

	result := &domain.ReconcileResult{
		MatchedNotAssociated: []domain.ProductBrief{},
		OnlyInProbability:    []domain.ProductBrief{},
		OnlyInJumpseller:     []domain.ProductBrief{},
	}

	probBySKU := make(map[string]domain.ProductForSync, len(data.probability))
	for _, p := range data.probability {
		key := normalizeSKU(p.SKU)
		if key == "" {
			result.ProbabilityNoSKU++
			continue
		}
		probBySKU[key] = p
	}

	jsSKUs := make(map[string]bool, len(data.jumpseller))
	for _, j := range data.jumpseller {
		key := normalizeSKU(j.SKU)
		if key == "" {
			result.JumpsellerNoSKU++
			continue
		}
		jsSKUs[key] = true

		if _, ok := probBySKU[key]; !ok {
			result.OnlyInJumpseller = append(result.OnlyInJumpseller, domain.ProductBrief{SKU: j.SKU, Name: j.Name})
			continue
		}
		if associated[key] {
			result.Matched++
		} else {
			result.MatchedNotAssociated = append(result.MatchedNotAssociated, domain.ProductBrief{SKU: j.SKU, Name: j.Name})
		}
	}

	for key, p := range probBySKU {
		if !jsSKUs[key] {
			result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
		}
	}

	return result, nil
}

func (uc *jumpsellerUseCase) probabilityWeightForStore(ctx context.Context, p domain.ProductForSync, storeUnit string) *float64 {
	if p.Weight == nil || *p.Weight <= 0 {
		return nil
	}

	weightKg := *p.Weight
	if unit := strings.ToLower(strings.TrimSpace(p.WeightUnit)); unit != "" && unit != probabilityWeightUnit {
		factor, known := weightFactor(unit)
		if !known {
			uc.logger.Warn(ctx).
				Str("sku", p.SKU).
				Str("weight_unit", p.WeightUnit).
				Msg("Unidad de peso desconocida en el producto de Probability, no se envia el peso a Jumpseller")
			return nil
		}
		weightKg = *p.Weight * factor
	}

	converted, ok := convertKgToStoreUnit(weightKg, storeUnit)
	if !ok {
		uc.logger.Warn(ctx).
			Str("sku", p.SKU).
			Str("store_weight_unit", storeUnit).
			Msg("La tienda Jumpseller no reporta una unidad de peso conocida, no se envia el peso")
		return nil
	}
	return positive(converted)
}

func (uc *jumpsellerUseCase) upsertMsgFromJumpseller(ctx context.Context, businessID, integrationID uint, j jumpsellerSKU, storeUnit string) providerUpsertMsg {
	msg := providerUpsertMsg{
		BusinessID:     businessID,
		IntegrationID:  integrationID,
		SKU:            j.SKU,
		Name:           j.Name,
		TrackInventory: true,
		Price:          j.Price,
		ExternalID:     j.ExternalID,
		Length:         positive(j.Length),
		Width:          positive(j.Width),
		Height:         positive(j.Height),
	}

	if msg.Length != nil || msg.Width != nil || msg.Height != nil {
		msg.DimensionUnit = probabilityDimensionUnit
	}

	if j.Weight > 0 {
		weightKg, ok := normalizeWeightToKg(j.Weight, storeUnit)
		if !ok {
			uc.logger.Warn(ctx).
				Str("sku", j.SKU).
				Str("store_weight_unit", storeUnit).
				Msg("La tienda Jumpseller no reporta una unidad de peso conocida, el peso no se importa")
		} else {
			msg.Weight = positive(weightKg)
			msg.WeightUnit = probabilityWeightUnit
		}
	}

	return msg
}

func (uc *jumpsellerUseCase) ApplyProductsToJumpseller(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	jsBySKU := indexJumpsellerBySKU(data.jumpseller)

	targets := make([]domain.ProductForSync, 0, len(data.probability))
	for _, p := range data.probability {
		if normalizeSKU(p.SKU) != "" {
			targets = append(targets, p)
		}
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeCreate,
		"total":          total,
	})

	fails := &failedSKUs{}
	created, updated := 0, 0
	for i, p := range targets {
		key := normalizeSKU(p.SKU)

		if existing, ok := jsBySKU[key]; ok {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), existing.ExternalID); merr != nil {
				uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al mapear producto existente de Jumpseller")
				fails.add(p.SKU)
			} else {
				updated++
			}
			uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, created, updated, fails.count())
			continue
		}

		newID, cerr := uc.client.CreateProduct(ctx, data.cred, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			ManageStock:   p.TrackInventory,
			Weight:        uc.probabilityWeightForStore(ctx, p, data.storeWeightUnit),
			Height:        p.Height,
			Width:         p.Width,
			Length:        p.Length,
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en Jumpseller")
			fails.add(p.SKU)
			uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, created, updated, fails.count())
			continue
		}

		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), newID); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en Jumpseller pero fallo el mapeo")
			fails.add(p.SKU)
		} else {
			created++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, created, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeCreate,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	probBySKU := indexProbabilityBySKU(data.probability)

	missing := make([]jumpsellerSKU, 0)
	for _, j := range data.jumpseller {
		key := normalizeSKU(j.SKU)
		if key == "" {
			continue
		}
		if _, exists := probBySKU[key]; exists {
			continue
		}
		missing = append(missing, j)
	}

	return uc.publishUpserts(ctx, businessID, uint(integIDUint), correlationID, ModeCreate, missing, data.storeWeightUnit)
}

func (uc *jumpsellerUseCase) UpdateProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	probBySKU := indexProbabilityBySKU(data.probability)

	existing := make([]jumpsellerSKU, 0)
	for _, j := range data.jumpseller {
		key := normalizeSKU(j.SKU)
		if key == "" {
			continue
		}
		if _, ok := probBySKU[key]; !ok {
			continue
		}
		existing = append(existing, j)
	}

	return uc.publishUpserts(ctx, businessID, uint(integIDUint), correlationID, ModeUpdate, existing, data.storeWeightUnit)
}

func (uc *jumpsellerUseCase) publishUpserts(ctx context.Context, businessID, integrationID uint, correlationID, mode string, items []jumpsellerSKU, storeUnit string) error {
	total := len(items)
	uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToProbability,
		"mode":           mode,
		"total":          total,
	})

	if uc.rabbit == nil {
		uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.completed", map[string]interface{}{
			"correlation_id": correlationID,
			"direction":      DirectionToProbability,
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
	for i, j := range items {
		msg := uc.upsertMsgFromJumpseller(ctx, businessID, integrationID, j, storeUnit)

		payload, merr := json.Marshal(msg)
		if merr != nil {
			fails.add(j.SKU)
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, payload); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", j.SKU).Msg("Error al publicar producto hacia Probability")
			fails.add(j.SKU)
		} else {
			applied++
		}
		uc.maybeProductProgress(ctx, businessID, integrationID, correlationID, DirectionToProbability, i+1, total, applied, 0, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToProbability,
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

func (uc *jumpsellerUseCase) UpdateProductsToJumpseller(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	jsBySKU := indexJumpsellerBySKU(data.jumpseller)

	targets := make([]domain.ProductForSync, 0)
	for _, p := range data.probability {
		key := normalizeSKU(p.SKU)
		if key == "" {
			continue
		}
		if _, ok := jsBySKU[key]; ok {
			targets = append(targets, p)
		}
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeUpdate,
		"total":          total,
	})

	fails := &failedSKUs{}
	updated, skipped := 0, 0
	touchedParents := make(map[int64]bool)

	for i, p := range targets {
		target := jsBySKU[normalizeSKU(p.SKU)]

		if target.VariantID > 0 {
			if touchedParents[target.ProductID] {
				skipped++
				uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, 0, updated, fails.count())
				continue
			}
			touchedParents[target.ProductID] = true
		}

		price := p.Price
		input := domain.UpdateProductInput{
			Name:        p.Name,
			Price:       &price,
			Description: p.Description,
			Weight:      uc.probabilityWeightForStore(ctx, p, data.storeWeightUnit),
			Height:      p.Height,
			Width:       p.Width,
			Length:      p.Length,
		}

		if target.VariantID > 0 {
			input.Price = nil
			input.Name = ""
		}

		if uerr := uc.client.UpdateProduct(ctx, data.cred, target.ProductID, input); uerr != nil {
			uc.logger.Error(ctx).Err(uerr).
				Str("sku", p.SKU).
				Int64("product_id", target.ProductID).
				Msg("Error al actualizar producto en Jumpseller")
			fails.add(p.SKU)
		} else {
			updated++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToJumpseller, i+1, total, 0, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToJumpseller,
		"mode":           ModeUpdate,
		"total":          total,
		"created":        0,
		"updated":        updated,
		"skipped":        skipped,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	jsBySKU := indexJumpsellerBySKU(data.jumpseller)
	probBySKU := indexProbabilityBySKU(data.probability)

	associated, err := uc.associatedSKUs(ctx, uint(integIDUint))
	if err != nil {
		return err
	}

	targets := make([]string, 0)
	if len(skus) > 0 {
		for _, sku := range skus {
			if key := normalizeSKU(sku); key != "" {
				targets = append(targets, key)
			}
		}
	} else {
		for key := range probBySKU {
			if jsBySKU[key].ExternalID != "" && !associated[key] {
				targets = append(targets, key)
			}
		}
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.started", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
	})

	fails := &failedSKUs{}
	updated := 0
	for i, key := range targets {
		p, okProb := probBySKU[key]
		j, okJS := jsBySKU[key]

		if !okProb || !okJS || j.ExternalID == "" || associated[key] {
			uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, "associate", i+1, total, 0, updated, fails.count())
			continue
		}

		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), j.ExternalID); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto con Jumpseller")
			fails.add(p.SKU)
		} else {
			associated[key] = true
			updated++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, "associate", i+1, total, 0, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), "jumpseller.product.sync.completed", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
		"created":        0,
		"updated":        updated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *jumpsellerUseCase) SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	if err := uc.ApplyProductsToJumpseller(ctx, integrationID, businessID, correlationID); err != nil {
		return err
	}
	return uc.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID)
}

func (uc *jumpsellerUseCase) maybeProductProgress(ctx context.Context, businessID, integrationID uint, correlationID, direction string, processed, total, created, updated, failed int) {
	if processed%syncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, "jumpseller.product.sync.progress", map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      direction,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
}
