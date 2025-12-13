package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

// CreateIntegration crea una nueva integración
func (uc *IntegrationUseCase) CreateIntegration(ctx context.Context, dto domain.CreateIntegrationDTO) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "CreateIntegration")

	// Validaciones
	if dto.Name == "" {
		return nil, domain.ErrIntegrationNameRequired
	}
	if dto.Code == "" {
		return nil, domain.ErrIntegrationCodeRequired
	}
	if dto.IntegrationTypeID == 0 {
		return nil, domain.ErrIntegrationTypeRequired
	}
	if !domain.IsValidCategory(dto.Category) {
		return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationCategoryInvalid, dto.Category)
	}

	// Validar que el tipo de integración exista (necesitamos el repositorio de tipos)
	// TODO: Inyectar IIntegrationTypeRepository en el use case

	// Validar que no exista otra integración con el mismo código
	exists, err := uc.repo.ExistsIntegrationByCode(ctx, dto.Code, dto.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al verificar existencia de código")
		return nil, fmt.Errorf("error al verificar código: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationCodeExists, dto.Code)
	}

	// TODO: Validar reglas específicas del tipo de integración
	// Por ejemplo, si el tipo es WhatsApp, debe ser global (BusinessID = NULL)
	// Esto se puede hacer consultando el IntegrationType y sus reglas

	// Obtener el tipo de integración para validar la conexión
	integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, dto.IntegrationTypeID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_type_id", dto.IntegrationTypeID).Msg("Error al obtener tipo de integración")
		return nil, fmt.Errorf("error al obtener tipo de integración: %w", err)
	}

	// VALIDAR CONEXIÓN ANTES DE GUARDAR
	// Obtener tester registrado para este tipo
	tester, err := uc.testerReg.GetTester(integrationType.Code)
	if err != nil {
		uc.log.Warn(ctx).
			Str("type_code", integrationType.Code).
			Msg("No hay tester registrado, solo validando credenciales básicas")
		// Fallback: validación básica si no hay tester
		if err := uc.validateBasicCredentials(ctx, integrationType.Code, dto.Credentials); err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTestFailed, err)
		}
	} else {
		// Deserializar Config a map para el tester
		var configMap map[string]interface{}
		if len(dto.Config) > 0 {
			if err := json.Unmarshal(dto.Config, &configMap); err != nil {
				return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationConfigDeserialize, err)
			}
		}

		// Testear conexión con el tester específico
		if err := tester.TestConnection(ctx, configMap, dto.Credentials); err != nil {
			uc.log.Error(ctx).
				Err(err).
				Str("type_code", integrationType.Code).
				Str("integration_code", dto.Code).
				Msg("Test de conexión falló al crear integración")
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTestFailed, err)
		}
		uc.log.Info(ctx).
			Str("type_code", integrationType.Code).
			Str("integration_code", dto.Code).
			Msg("Test de conexión exitoso antes de crear integración")
	}

	// Convertir Config a datatypes.JSON
	var configJSON datatypes.JSON
	if len(dto.Config) > 0 {
		configBytes, err := json.Marshal(dto.Config)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationConfigSerialize, err)
		}
		configJSON = configBytes
	}

	// Convertir Credentials a []byte (se encriptará en el repository)
	var credentialsJSON datatypes.JSON
	if len(dto.Credentials) > 0 {
		credentialsBytes, err := json.Marshal(dto.Credentials)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsSerialize, err)
		}
		credentialsJSON = credentialsBytes
	}

	// Crear entidad de dominio
	integration := &domain.Integration{
		Name:              dto.Name,
		Code:              dto.Code,
		IntegrationTypeID: dto.IntegrationTypeID,
		Category:          dto.Category,
		BusinessID:        dto.BusinessID,
		IsActive:          dto.IsActive,
		IsDefault:         dto.IsDefault,
		Config:            configJSON,
		Credentials:       credentialsJSON,
		Description:       dto.Description,
		CreatedByID:       dto.CreatedByID,
	}

	// Guardar en repository (encriptará credenciales automáticamente)
	if err := uc.repo.CreateIntegration(ctx, integration); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al crear integración")
		return nil, fmt.Errorf("error al crear integración: %w", err)
	}

	uc.log.Info(ctx).
		Uint("id", integration.ID).
		Uint("integration_type_id", integration.IntegrationTypeID).
		Str("code", integration.Code).
		Msg("Integración creada exitosamente")

	// Notificar observadores (e.g., para auto-sync)
	// Hacemos esto de forma asíncrona para no bloquear la respuesta HTTP
	go func() {
		for _, observer := range uc.observers {
			// Crear un nuevo contexto desconectado del request HTTP cancelar
			bgCtx := context.Background()
			// Tratar pánicos en observadores para no romper nada
			defer func() {
				if r := recover(); r != nil {
					uc.log.Error(bgCtx).Interface("recover", r).Msg("Pánico en observador de creación de integración")
				}
			}()
			observer(bgCtx, integration)
		}
	}()

	return integration, nil
}
