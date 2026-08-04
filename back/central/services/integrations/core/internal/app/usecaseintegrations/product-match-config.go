package usecaseintegrations

import (
	"context"
	"fmt"

	"gorm.io/datatypes"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func (uc *IntegrationUseCase) loadOwnedIntegration(ctx context.Context, id, businessID uint) (*domain.Integration, error) {
	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if integration == nil {
		return nil, fmt.Errorf("integracion no encontrada")
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, fmt.Errorf("la integracion no pertenece al negocio")
	}
	return integration, nil
}

func buildProductMatchConfig(integration *domain.Integration) *domain.ProductMatchConfig {
	var typeDefaultRaw, typeOptionsRaw []byte
	if integration.IntegrationType != nil {
		typeDefaultRaw = integration.IntegrationType.DefaultProductMatchRules
		typeOptionsRaw = integration.IntegrationType.ProductMatchOptions
	}

	options := productmatch.ParseOptions(typeOptionsRaw)
	defaults := productmatch.Resolve(nil, typeDefaultRaw)
	own := productmatch.Parse(integration.ProductMatchRules)

	effective := defaults
	if own != nil {
		effective = own
	}

	return &domain.ProductMatchConfig{
		IntegrationID: integration.ID,
		Rules:         productmatch.First(productmatch.AllowedByOptions(effective, options)),
		DefaultRules:  productmatch.First(productmatch.AllowedByOptions(defaults, options)),
		Options:       options,
		IsOverride:    own != nil,
	}
}

func (uc *IntegrationUseCase) GetProductMatchConfig(ctx context.Context, id, businessID uint) (*domain.ProductMatchConfig, error) {
	ctx = log.WithFunctionCtx(ctx, "GetProductMatchConfig")

	integration, err := uc.loadOwnedIntegration(ctx, id, businessID)
	if err != nil {
		return nil, err
	}
	return buildProductMatchConfig(integration), nil
}

func (uc *IntegrationUseCase) UpdateProductMatchConfig(ctx context.Context, id, businessID uint, rules []productmatch.Rule) (*domain.ProductMatchConfig, error) {
	ctx = log.WithFunctionCtx(ctx, "UpdateProductMatchConfig")

	integration, err := uc.loadOwnedIntegration(ctx, id, businessID)
	if err != nil {
		return nil, err
	}

	var typeOptionsRaw []byte
	if integration.IntegrationType != nil {
		typeOptionsRaw = integration.IntegrationType.ProductMatchOptions
	}
	options := productmatch.ParseOptions(typeOptionsRaw)

	if len(rules) != 1 {
		return nil, fmt.Errorf("debe seleccionar exactamente una regla de match")
	}
	for _, r := range rules {
		if !r.Valid() {
			return nil, fmt.Errorf("regla de match invalida: %s", r.Key())
		}
	}

	sanitized := productmatch.First(productmatch.Sanitize(rules))
	if len(productmatch.AllowedByOptions(sanitized, options)) != len(sanitized) {
		return nil, fmt.Errorf("alguna regla usa campos no soportados por este tipo de integracion")
	}

	raw, err := productmatch.Marshal(sanitized)
	if err != nil {
		return nil, fmt.Errorf("error al serializar las reglas de match: %w", err)
	}

	if err := uc.repo.UpdateProductMatchRules(ctx, id, datatypes.JSON(raw)); err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_id", id).Msg("Error al guardar las reglas de match de producto")
		return nil, err
	}

	if cerr := uc.cache.InvalidateIntegration(ctx, id); cerr != nil {
		uc.log.Warn(ctx).Err(cerr).Uint("integration_id", id).Msg("No se pudo invalidar el cache de la integracion")
	}

	integration.ProductMatchRules = datatypes.JSON(raw)
	return buildProductMatchConfig(integration), nil
}
