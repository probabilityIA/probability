package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

// CreateIntegration crea una nueva integración
func (uc *integrationUseCase) CreateIntegration(ctx context.Context, dto domain.CreateIntegrationDTO) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "CreateIntegration")

	// Validaciones
	if dto.Name == "" {
		return nil, fmt.Errorf("el nombre de la integración es obligatorio")
	}
	if dto.Code == "" {
		return nil, fmt.Errorf("el código de la integración es obligatorio")
	}
	if !domain.IsValidType(dto.Type) {
		return nil, fmt.Errorf("tipo de integración inválido: %s", dto.Type)
	}
	if !domain.IsValidCategory(dto.Category) {
		return nil, fmt.Errorf("categoría de integración inválida: %s", dto.Category)
	}

	// Validar que no exista otra integración con el mismo código
	exists, err := uc.repo.ExistsByCode(ctx, dto.Code, dto.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al verificar existencia de código")
		return nil, fmt.Errorf("error al verificar código: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("ya existe una integración con el código '%s'", dto.Code)
	}

	// Validación especial: WhatsApp debe ser global (BusinessID = NULL)
	if dto.Type == domain.IntegrationTypeWhatsApp && dto.BusinessID != nil {
		return nil, fmt.Errorf("WhatsApp debe ser una integración global (business_id debe ser NULL)")
	}

	// Validar que no exista otra integración de WhatsApp si estamos creando WhatsApp
	if dto.Type == domain.IntegrationTypeWhatsApp {
		existing, err := uc.repo.GetByType(ctx, domain.IntegrationTypeWhatsApp, nil)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("ya existe una integración de WhatsApp. Solo puede haber una integración de WhatsApp global")
		}
	}

	// Convertir Config a datatypes.JSON
	var configJSON datatypes.JSON
	if len(dto.Config) > 0 {
		configBytes, err := json.Marshal(dto.Config)
		if err != nil {
			return nil, fmt.Errorf("error al serializar configuración: %w", err)
		}
		configJSON = configBytes
	}

	// Convertir Credentials a []byte (se encriptará en el repository)
	var credentialsJSON datatypes.JSON
	if len(dto.Credentials) > 0 {
		credentialsBytes, err := json.Marshal(dto.Credentials)
		if err != nil {
			return nil, fmt.Errorf("error al serializar credenciales: %w", err)
		}
		credentialsJSON = credentialsBytes
	}

	// Crear entidad de dominio
	integration := &domain.Integration{
		Name:        dto.Name,
		Code:        dto.Code,
		Type:        dto.Type,
		Category:    dto.Category,
		BusinessID:  dto.BusinessID,
		IsActive:    dto.IsActive,
		IsDefault:   dto.IsDefault,
		Config:      configJSON,
		Credentials: credentialsJSON,
		Description: dto.Description,
		CreatedByID: dto.CreatedByID,
	}

	// Guardar en repository (encriptará credenciales automáticamente)
	if err := uc.repo.Create(ctx, integration); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al crear integración")
		return nil, fmt.Errorf("error al crear integración: %w", err)
	}

	uc.log.Info(ctx).
		Uint("id", integration.ID).
		Str("type", integration.Type).
		Str("code", integration.Code).
		Msg("Integración creada exitosamente")

	return integration, nil
}
