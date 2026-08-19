package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	DirectionToTiendanube  = "to_tiendanube"
	DirectionToProbability = "to_probability"

	ModeCreate = "create"
	ModeUpdate = "update"

	eventProductSyncStarted   = "tiendanube.product.sync.started"
	eventProductSyncProgress  = "tiendanube.product.sync.progress"
	eventProductSyncCompleted = "tiendanube.product.sync.completed"
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

type tiendanubeSKU struct {
	SKU        string
	Barcode    string
	Name       string
	ImageURL   string
	Price      float64
	Stock      int
	ExternalID string
	ProductID  int64
	VariantID  int64
	Weight     float64
	Height     float64
	Width      float64
	Depth      float64
}

type reconcileData struct {
	cred        domain.Credential
	probability []domain.ProductForSync
	tiendanube  []tiendanubeSKU
	rules       []productmatch.Rule
	outcome     productmatch.Outcome
}

func (t tiendanubeSKU) MatchItem() productmatch.Item {
	variantID := ""
	if t.VariantID != 0 {
		variantID = strconv.FormatInt(t.VariantID, 10)
	}
	return productmatch.Item{
		SKU:        t.SKU,
		Barcode:    t.Barcode,
		ExternalID: strconv.FormatInt(t.ProductID, 10),
		VariantID:  variantID,
		Name:       t.Name,
	}
}

func tiendanubeRefs(t tiendanubeSKU) productmatch.ExternalRefs {
	variantID := ""
	if t.VariantID != 0 {
		variantID = strconv.FormatInt(t.VariantID, 10)
	}
	return productmatch.ExternalRefs{
		ProductID: t.ExternalID,
		VariantID: variantID,
		SKU:       t.SKU,
		Barcode:   t.Barcode,
	}
}

func probabilityItems(products []domain.ProductForSync) []productmatch.Item {
	items := make([]productmatch.Item, len(products))
	for i, p := range products {
		items[i] = p.MatchItem()
	}
	return items
}

func tiendanubeItems(skus []tiendanubeSKU) []productmatch.Item {
	items := make([]productmatch.Item, len(skus))
	for i, t := range skus {
		items[i] = t.MatchItem()
	}
	return items
}

func flattenProductSKUs(products []domain.TiendanubeProduct) []tiendanubeSKU {
	flat := make([]tiendanubeSKU, 0, len(products))
	for _, product := range products {
		for _, variant := range product.Variants {
			flat = append(flat, tiendanubeSKU{
				SKU:        variant.SKU,
				Barcode:    variant.Barcode,
				Name:       product.Name,
				ImageURL:   product.ImageURL,
				Price:      variant.Price,
				Stock:      variant.Stock,
				ExternalID: strconv.FormatInt(product.ID, 10) + ":" + strconv.FormatInt(variant.ID, 10),
				ProductID:  product.ID,
				VariantID:  variant.ID,
				Weight:     variant.Weight,
				Height:     variant.Height,
				Width:      variant.Width,
				Depth:      variant.Depth,
			})
		}
	}
	return flat
}

func matchedByProbabilityIndex(data *reconcileData) map[int]tiendanubeSKU {
	matched := make(map[int]tiendanubeSKU, len(data.outcome.Pairs))
	for _, pair := range data.outcome.Pairs {
		matched[pair.ProbabilityIndex] = data.tiendanube[pair.ChannelIndex]
	}
	return matched
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

func (uc *tiendanubeUseCase) loadReconcileData(ctx context.Context, integrationID string, businessID uint) (*reconcileData, error) {
	integration, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	probProducts, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}

	tnProducts, err := uc.client.GetProducts(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("listing tiendanube products: %w", err)
	}

	flat := flattenProductSKUs(tnProducts)
	rules := productmatch.Sanitize(integration.ProductMatchRules)

	return &reconcileData{
		cred:        cred,
		probability: probProducts,
		tiendanube:  flat,
		rules:       rules,
		outcome:     productmatch.Reconcile(rules, probabilityItems(probProducts), tiendanubeItems(flat)),
	}, nil
}

func (uc *tiendanubeUseCase) associatedSKUs(ctx context.Context, integrationID uint) (map[string]bool, error) {
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

func (uc *tiendanubeUseCase) ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error) {
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
		MatchedItems:         []domain.ProductBrief{},
		MatchedNotAssociated: []domain.ProductBrief{},
		OnlyInProbability:    []domain.ProductBrief{},
		OnlyInTiendanube:     []domain.ProductBrief{},
		ProbabilityNoSKU:     data.outcome.ProbabilityUnmatchable,
		TiendanubeNoSKU:      data.outcome.ChannelUnmatchable,
		MatchRules:           data.rules,
	}

	for _, pair := range data.outcome.Pairs {
		prob := data.probability[pair.ProbabilityIndex]
		channel := data.tiendanube[pair.ChannelIndex]
		brief := domain.ProductBrief{
			SKU:          prob.SKU,
			Name:         channel.Name,
			MatchedBy:    pair.Rule.Key(),
			MatchedValue: channel.MatchItem().Values()[pair.Rule.Channel],
		}
		result.MatchedItems = append(result.MatchedItems, brief)
		if associated[normalizeSKU(prob.SKU)] {
			result.Matched++
		} else {
			result.MatchedNotAssociated = append(result.MatchedNotAssociated, brief)
		}
	}

	for _, idx := range data.outcome.OnlyInChannel {
		t := data.tiendanube[idx]
		result.OnlyInTiendanube = append(result.OnlyInTiendanube, domain.ProductBrief{SKU: t.SKU, Name: t.Name})
	}
	for _, idx := range data.outcome.OnlyInProbability {
		p := data.probability[idx]
		result.OnlyInProbability = append(result.OnlyInProbability, domain.ProductBrief{SKU: p.SKU, Name: p.Name})
	}

	return result, nil
}

func (uc *tiendanubeUseCase) probabilityWeightKg(ctx context.Context, p domain.ProductForSync) *float64 {
	if p.Weight == nil || *p.Weight <= 0 {
		return nil
	}
	weightKg := *p.Weight
	if unit := normalizeSKU(p.WeightUnit); unit != "" && unit != probabilityWeightUnit {
		factor, known := weightFactor(unit)
		if !known {
			uc.logger.Warn(ctx).
				Str("sku", p.SKU).
				Str("weight_unit", p.WeightUnit).
				Msg("Unidad de peso desconocida en el producto de Probability, no se envia el peso a Tiendanube")
			return nil
		}
		weightKg = *p.Weight * factor
	}
	return positive(weightKg)
}

func (uc *tiendanubeUseCase) upsertMsgFromTiendanube(businessID, integrationID uint, t tiendanubeSKU) providerUpsertMsg {
	msg := providerUpsertMsg{
		BusinessID:     businessID,
		IntegrationID:  integrationID,
		ImageURL:       t.ImageURL,
		SKU:            t.SKU,
		Name:           t.Name,
		TrackInventory: true,
		Price:          t.Price,
		ExternalID:     t.ExternalID,
		Length:         positive(t.Depth),
		Width:          positive(t.Width),
		Height:         positive(t.Height),
	}

	if msg.Length != nil || msg.Width != nil || msg.Height != nil {
		msg.DimensionUnit = probabilityDimensionUnit
	}

	if weight := positive(t.Weight); weight != nil {
		msg.Weight = weight
		msg.WeightUnit = probabilityWeightUnit
	}

	return msg
}

func (uc *tiendanubeUseCase) ApplyProductsToTiendanube(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	matched := matchedByProbabilityIndex(data)

	only := selectedSKUs(skus)
	targets := make([]int, 0, len(data.probability))
	for i, p := range data.probability {
		key := normalizeSKU(p.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		targets = append(targets, i)
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), eventProductSyncStarted, map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToTiendanube,
		"mode":           ModeCreate,
		"total":          total,
	})

	fails := &failedSKUs{}
	created, updated := 0, 0
	for i, idx := range targets {
		p := data.probability[idx]

		if existing, ok := matched[idx]; ok {
			if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), tiendanubeRefs(existing)); merr != nil {
				uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al mapear producto existente de Tiendanube")
				fails.add(p.SKU)
			} else {
				updated++
			}
			uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToTiendanube, i+1, total, created, updated, fails.count())
			continue
		}

		productID, variantID, cerr := uc.client.CreateProduct(ctx, data.cred, domain.CreateProductInput{
			Name:          p.Name,
			SKU:           p.SKU,
			Barcode:       p.Barcode,
			Price:         p.Price,
			Description:   p.Description,
			StockQuantity: p.StockQuantity,
			ManageStock:   p.TrackInventory,
			ImageURL:      p.ImageURL,
			Weight:        uc.probabilityWeightKg(ctx, p),
			Height:        p.Height,
			Width:         p.Width,
			Length:        p.Length,
		})
		if cerr != nil {
			uc.logger.Error(ctx).Err(cerr).Str("sku", p.SKU).Msg("Error al crear producto en Tiendanube")
			fails.add(p.SKU)
			uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToTiendanube, i+1, total, created, updated, fails.count())
			continue
		}

		externalID := strconv.FormatInt(productID, 10)
		variantRef := ""
		if variantID > 0 {
			externalID = externalID + ":" + strconv.FormatInt(variantID, 10)
			variantRef = strconv.FormatInt(variantID, 10)
		}

		newRefs := productmatch.ExternalRefs{ProductID: externalID, VariantID: variantRef, SKU: p.SKU, Barcode: p.Barcode}
		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), newRefs); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Producto creado en Tiendanube pero fallo el mapeo")
			fails.add(p.SKU)
		} else {
			created++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToTiendanube, i+1, total, created, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), eventProductSyncCompleted, map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToTiendanube,
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

func (uc *tiendanubeUseCase) UpdateProductsToTiendanube(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	matched := matchedByProbabilityIndex(data)

	only := selectedSKUs(skus)
	targets := make([]int, 0)
	for i, p := range data.probability {
		key := normalizeSKU(p.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		if _, ok := matched[i]; ok {
			targets = append(targets, i)
		}
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), eventProductSyncStarted, map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToTiendanube,
		"mode":           ModeUpdate,
		"total":          total,
	})

	fails := &failedSKUs{}
	updated := 0
	touchedParents := make(map[int64]bool)

	for i, idx := range targets {
		p := data.probability[idx]
		target := matched[idx]
		failed := false

		if !touchedParents[target.ProductID] {
			touchedParents[target.ProductID] = true
			if uerr := uc.client.UpdateProduct(ctx, data.cred, target.ProductID, domain.UpdateProductInput{
				Name:        p.Name,
				Description: p.Description,
			}); uerr != nil {
				uc.logger.Error(ctx).Err(uerr).
					Str("sku", p.SKU).
					Int64("product_id", target.ProductID).
					Msg("Error al actualizar producto en Tiendanube")
				fails.add(p.SKU)
				failed = true
			}
		}

		if !failed && target.VariantID > 0 {
			price := p.Price
			if verr := uc.client.UpdateVariant(ctx, data.cred, target.ProductID, target.VariantID, domain.UpdateVariantInput{
				Price:   &price,
				Barcode: p.Barcode,
				Weight:  uc.probabilityWeightKg(ctx, p),
				Depth:   p.Length,
				Width:   p.Width,
				Height:  p.Height,
			}); verr != nil {
				uc.logger.Error(ctx).Err(verr).
					Str("sku", p.SKU).
					Int64("variant_id", target.VariantID).
					Msg("Error al actualizar la variante en Tiendanube")
				fails.add(p.SKU)
				failed = true
			}
		}

		if !failed {
			updated++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, DirectionToTiendanube, i+1, total, 0, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), eventProductSyncCompleted, map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToTiendanube,
		"mode":           ModeUpdate,
		"total":          total,
		"created":        0,
		"updated":        updated,
		"failed":         fails.count(),
		"failed_skus":    fails.list(),
		"failed_hidden":  fails.truncated(),
	})

	return nil
}

func (uc *tiendanubeUseCase) ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)
	missing := make([]tiendanubeSKU, 0)
	for _, idx := range data.outcome.OnlyInChannel {
		t := data.tiendanube[idx]
		key := normalizeSKU(t.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		missing = append(missing, t)
	}

	return uc.publishUpserts(ctx, businessID, uint(integIDUint), correlationID, ModeCreate, missing)
}

func (uc *tiendanubeUseCase) UpdateProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)
	existing := make([]tiendanubeSKU, 0)
	for _, pair := range data.outcome.Pairs {
		t := data.tiendanube[pair.ChannelIndex]
		key := normalizeSKU(t.SKU)
		if key == "" || (only != nil && !only[key]) {
			continue
		}
		existing = append(existing, t)
	}

	return uc.publishUpserts(ctx, businessID, uint(integIDUint), correlationID, ModeUpdate, existing)
}

func (uc *tiendanubeUseCase) publishUpserts(ctx context.Context, businessID, integrationID uint, correlationID, mode string, items []tiendanubeSKU) error {
	total := len(items)
	uc.emitSyncEvent(ctx, businessID, integrationID, eventProductSyncStarted, map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      DirectionToProbability,
		"mode":           mode,
		"total":          total,
	})

	if uc.rabbit == nil {
		uc.emitSyncEvent(ctx, businessID, integrationID, eventProductSyncCompleted, map[string]interface{}{
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
	for i, t := range items {
		msg := uc.upsertMsgFromTiendanube(businessID, integrationID, t)

		payload, merr := json.Marshal(msg)
		if merr != nil {
			fails.add(t.SKU)
		} else if perr := uc.rabbit.Publish(ctx, rabbitmq.QueueProductsProviderUpsert, payload); perr != nil {
			uc.logger.Error(ctx).Err(perr).Str("sku", t.SKU).Msg("Error al publicar producto hacia Probability")
			fails.add(t.SKU)
		} else {
			applied++
		}
		uc.maybeProductProgress(ctx, businessID, integrationID, correlationID, DirectionToProbability, i+1, total, applied, 0, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, integrationID, eventProductSyncCompleted, map[string]interface{}{
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

func (uc *tiendanubeUseCase) AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	data, err := uc.loadReconcileData(ctx, integrationID, businessID)
	if err != nil {
		return err
	}

	associated, err := uc.associatedSKUs(ctx, uint(integIDUint))
	if err != nil {
		return err
	}

	only := selectedSKUs(skus)

	targets := make([]productmatch.Pair, 0, len(data.outcome.Pairs))
	for _, pair := range data.outcome.Pairs {
		key := normalizeSKU(data.probability[pair.ProbabilityIndex].SKU)
		if only != nil {
			if key == "" || !only[key] {
				continue
			}
		} else if associated[key] {
			continue
		}
		if data.tiendanube[pair.ChannelIndex].ExternalID == "" {
			continue
		}
		targets = append(targets, pair)
	}

	total := len(targets)
	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), eventProductSyncStarted, map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      "associate",
		"total":          total,
	})

	fails := &failedSKUs{}
	updated := 0
	for i, pair := range targets {
		p := data.probability[pair.ProbabilityIndex]
		t := data.tiendanube[pair.ChannelIndex]

		if merr := uc.productRepo.UpsertProductIntegrationMapping(ctx, p.ID, businessID, uint(integIDUint), tiendanubeRefs(t)); merr != nil {
			uc.logger.Error(ctx).Err(merr).Str("sku", p.SKU).Msg("Error al asociar producto con Tiendanube")
			fails.add(p.SKU)
		} else {
			updated++
		}
		uc.maybeProductProgress(ctx, businessID, uint(integIDUint), correlationID, "associate", i+1, total, 0, updated, fails.count())
	}

	uc.emitSyncEvent(ctx, businessID, uint(integIDUint), eventProductSyncCompleted, map[string]interface{}{
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

func (uc *tiendanubeUseCase) SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error {
	if err := uc.ApplyProductsToTiendanube(ctx, integrationID, businessID, correlationID); err != nil {
		return err
	}
	return uc.ApplyProductsToProbability(ctx, integrationID, businessID, correlationID)
}

func (uc *tiendanubeUseCase) maybeProductProgress(ctx context.Context, businessID, integrationID uint, correlationID, direction string, processed, total, created, updated, failed int) {
	if processed%syncProgressBatch != 0 && processed != total {
		return
	}
	uc.emitSyncEvent(ctx, businessID, integrationID, eventProductSyncProgress, map[string]interface{}{
		"correlation_id": correlationID,
		"direction":      direction,
		"processed":      processed,
		"total":          total,
		"created":        created,
		"updated":        updated,
		"failed":         failed,
	})
}
