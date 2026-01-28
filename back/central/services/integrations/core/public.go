package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"gorm.io/datatypes"
)

const (
	IntegrationTypeShopify      = 1
	IntegrationTypeWhatsApp     = 2 // Whastap en la BD
	IntegrationTypeMercadoLibre = 3
	IntegrationTypeWoocommerce  = 4
)

// SyncOrdersParams contiene los parámetros opcionales para sincronizar órdenes
type SyncOrdersParams struct {
	CreatedAtMin      *time.Time
	CreatedAtMax      *time.Time
	Status            string
	FinancialStatus   string
	FulfillmentStatus string
}

type IntegrationWithCredentials = domain.IntegrationWithCredentials

type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          interface{}
}

func (ic *integrationCore) GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error) {
	publicIntegration, err := ic.useCase.GetPublicIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	return &Integration{
		ID:              publicIntegration.ID,
		BusinessID:      publicIntegration.BusinessID,
		Name:            publicIntegration.Name,
		StoreID:         publicIntegration.StoreID,
		IntegrationType: publicIntegration.IntegrationType,
		Config:          publicIntegration.Config,
	}, nil
}

// GetIntegrationByStoreID busca una integración por StoreID (ej: shop domain) y tipo
func (ic *integrationCore) GetIntegrationByStoreID(ctx context.Context, storeID string, integrationType int) (*Integration, error) {
	var typeID *uint
	if integrationType > 0 {
		tid := uint(integrationType)
		typeID = &tid
	}

	filters := domain.IntegrationFilters{
		Page:              1,
		PageSize:          1,
		IntegrationTypeID: typeID,
		StoreID:           &storeID,
	}

	integrations, _, err := ic.useCase.ListIntegrations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing integrations: %w", err)
	}
	if len(integrations) == 0 {
		return nil, fmt.Errorf("integration not found for store_id %s", storeID)
	}

	integration := integrations[0]

	var config map[string]interface{}
	if len(integration.Config) > 0 {
		_ = json.Unmarshal(integration.Config, &config)
	}

	integrationTypeCode := integrationType
	if integrationTypeCode == 0 && integration.IntegrationType != nil {
		integrationTypeCode = getIntegrationTypeCodeAsInt(integration.IntegrationType.Code)
	} else if integrationTypeCode == 0 {
		integrationTypeCode = int(integration.IntegrationTypeID)
	}

	return &Integration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integrationTypeCode,
		Config:          config,
	}, nil
}

func (ic *integrationCore) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return ic.useCase.DecryptCredentialField(ctx, integrationID, fieldName)
}

func (ic *integrationCore) RegisterIntegration(integrationType int, integration IIntegrationContract) {
	if integrationType == 0 {
		ic.logger.Error().Msg("RegisterIntegration: integration type cannot be zero")
		return
	}
	if integration == nil {
		ic.logger.Error().Msg("RegisterIntegration: integration cannot be nil")
		return
	}

	ic.integrations[integrationType] = integration

	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		ic.logger.Error().Msg("RegisterIntegration: error interno: no se puede acceder al registry de testers")
		return
	}

	adapter := &integrationAdapter{integration: integration}
	if err := useCaseImpl.GetTesterRegistry().Register(integrationType, adapter); err != nil {
		ic.logger.Error().Err(err).Int("integration_type", integrationType).Msg("RegisterIntegration: error al registrar tester")
		return
	}

	ic.logger.Info().Int("integration_type", integrationType).Msg("Integration registered successfully")
}

func (ic *integrationCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	// Intentar obtener como int primero, luego como string para compatibilidad
	var integrationType int
	if intVal, ok := config["integration_type"].(int); ok {
		integrationType = intVal
	} else if floatVal, ok := config["integration_type"].(float64); ok {
		// JSON numbers se deserializan como float64
		integrationType = int(floatVal)
	} else if strVal, ok := config["integration_type"].(string); ok {
		// Compatibilidad con strings antiguos
		integrationType = getIntegrationTypeCodeAsInt(strVal)
	} else {
		return fmt.Errorf("integration_type is required in config and must be int or string")
	}

	if integrationType == 0 {
		return fmt.Errorf("integration_type cannot be zero")
	}

	integration, ok := ic.integrations[integrationType]
	if !ok {
		return fmt.Errorf("integration no registrada para tipo %d", integrationType)
	}
	return integration.TestConnection(ctx, config, credentials)
}

func (ic *integrationCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	integration, err := ic.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return err
	}

	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return fmt.Errorf("integration type %d not registered", integration.IntegrationType)
	}

	// Intentar usar el método con parámetros si está disponible
	if syncSvc, ok := integrationImpl.(domain.IOrderSyncService); ok {
		return syncSvc.SyncOrdersByIntegrationIDWithParams(ctx, integrationID, params)
	}

	// Fallback al método sin parámetros
	return integrationImpl.SyncOrdersByIntegrationID(ctx, integrationID)
}

func (ic *integrationCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	integration, err := ic.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return err
	}

	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return fmt.Errorf("integration no registrada para tipo %d", integration.IntegrationType)
	}

	return integrationImpl.SyncOrdersByIntegrationID(ctx, integrationID)
}

func (ic *integrationCore) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		return fmt.Errorf("error interno: no se puede acceder al use case")
	}

	businessIDPtr := &businessID
	filters := domain.IntegrationFilters{
		BusinessID: businessIDPtr,
		IsActive:   &[]bool{true}[0],
	}

	integrations, _, err := useCaseImpl.ListIntegrations(ctx, filters)
	if err != nil {
		return fmt.Errorf("error al obtener integraciones: %w", err)
	}

	for _, integration := range integrations {
		if integration.IntegrationType == nil {
			continue
		}

		integrationID := fmt.Sprintf("%d", integration.ID)
		if err := ic.SyncOrdersByIntegrationID(ctx, integrationID); err != nil {
			continue
		}
	}

	return nil
}

func (ic *integrationCore) RegisterObserverForType(integrationType int, observer func(context.Context, *Integration)) {
	ic.useCase.RegisterObserver(func(ctx context.Context, integration *domain.Integration) {
		var integrationTypeCode int
		if integration.IntegrationType != nil {
			// Convertir el código del tipo de integración a int si es necesario
			// Por ahora asumimos que el código puede ser convertido o comparado
			integrationTypeCode = getIntegrationTypeCodeAsInt(integration.IntegrationType.Code)
		}

		if integrationTypeCode == integrationType {
			publicIntegration := mapDomainToPublicIntegration(ic.useCase, integration)
			observer(ctx, publicIntegration)
		}
	})
}

// getIntegrationTypeCodeAsInt convierte el código de tipo de integración a int
// Esta función mapea los códigos (strings) a los IDs de la tabla integration_types
func getIntegrationTypeCodeAsInt(code string) int {
	// Convertir a minúsculas para comparación case-insensitive
	lowerCode := strings.ToLower(code)
	switch lowerCode {
	case "shopify":
		return IntegrationTypeShopify // 1
	case "whatsapp", "whatsap", "whastap":
		return IntegrationTypeWhatsApp // 2
	case "mercado_libre", "mercado libre":
		return IntegrationTypeMercadoLibre // 3
	case "woocommerce", "woocormerce":
		return IntegrationTypeWoocommerce // 4
	default:
		return 0
	}
}

type integrationAdapter struct {
	integration IIntegrationContract
}

func (a *integrationAdapter) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return a.integration.TestConnection(ctx, config, credentials)
}

// GetWebhookURL obtiene la URL del webhook para una integración específica
func (ic *integrationCore) GetWebhookURL(ctx context.Context, integrationID uint) (*WebhookInfo, error) {
	// Obtener la integración
	integration, err := ic.GetIntegrationByID(ctx, fmt.Sprintf("%d", integrationID))
	if err != nil {
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener la implementación de la integración
	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return nil, fmt.Errorf("integración no registrada para tipo %d", integration.IntegrationType)
	}

	// Obtener la URL base del servidor para webhooks
	// Prioridad: WEBHOOK_BASE_URL > URL_BASE_SWAGGER
	baseURL := ic.config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = ic.config.Get("URL_BASE_SWAGGER")
	}
	if baseURL == "" {
		return nil, fmt.Errorf("WEBHOOK_BASE_URL o URL_BASE_SWAGGER no está configurada")
	}

	// Delegar a la integración específica
	return integrationImpl.GetWebhookURL(ctx, baseURL, integrationID)
}

// UpdateIntegrationConfig actualiza el config de una integración haciendo merge con el config existente
func (ic *integrationCore) UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error {
	// Obtener integración existente
	id, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		return fmt.Errorf("ID de integración inválido: %w", err)
	}

	existing, err := ic.useCase.GetPublicIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener config existente (ya viene como map desde GetPublicIntegrationByID)
	existingConfig := existing.Config
	if existingConfig == nil {
		existingConfig = make(map[string]interface{})
	}

	// Hacer merge del config existente con el nuevo
	for k, v := range newConfig {
		existingConfig[k] = v
	}

	// Convertir a JSON
	configBytes, err := json.Marshal(existingConfig)
	if err != nil {
		return fmt.Errorf("error al serializar config: %w", err)
	}

	// Actualizar usando el useCase
	dto := domain.UpdateIntegrationDTO{
		Config: func() *datatypes.JSON {
			configJSON := datatypes.JSON(configBytes)
			return &configJSON
		}(),
	}

	_, err = ic.useCase.UpdateIntegration(ctx, uint(id), dto)
	if err != nil {
		return fmt.Errorf("error al actualizar integración: %w", err)
	}

	return nil
}

// ListWebhooks lista todos los webhooks de una integración (solo soportado para integraciones que implementan IWebhookOperations)
func (ic *integrationCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	// Obtener la integración
	integration, err := ic.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener la implementación de la integración
	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return nil, fmt.Errorf("integración no registrada para tipo %d", integration.IntegrationType)
	}

	// Verificar si implementa IWebhookOperations
	webhookOps, ok := integrationImpl.(IWebhookOperations)
	if !ok {
		return nil, fmt.Errorf("esta integración no soporta operaciones de webhooks")
	}

	// Listar webhooks
	return webhookOps.ListWebhooks(ctx, integrationID)
}

// DeleteWebhook elimina un webhook de una integración (solo soportado para integraciones que implementan IWebhookOperations)
func (ic *integrationCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	// Obtener la integración
	integration, err := ic.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener la implementación de la integración
	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return fmt.Errorf("integración no registrada para tipo %d", integration.IntegrationType)
	}

	// Verificar si implementa IWebhookOperations
	webhookOps, ok := integrationImpl.(IWebhookOperations)
	if !ok {
		return fmt.Errorf("esta integración no soporta operaciones de webhooks")
	}

	// Eliminar webhook
	return webhookOps.DeleteWebhook(ctx, integrationID, webhookID)
}

// VerifyWebhooksByURL verifica webhooks existentes que coincidan con nuestra URL
func (ic *integrationCore) VerifyWebhooksByURL(ctx context.Context, integrationID string) ([]interface{}, error) {
	// Obtener la integración
	integration, err := ic.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener la implementación de la integración
	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return nil, fmt.Errorf("integración no registrada para tipo %d", integration.IntegrationType)
	}

	// Verificar si implementa IWebhookOperations
	webhookOps, ok := integrationImpl.(IWebhookOperations)
	if !ok {
		return nil, fmt.Errorf("esta integración no soporta operaciones de webhooks")
	}

	// Obtener baseURL desde configuración
	baseURL := ic.config.Get("URL_BASE_SWAGGER")
	if baseURL == "" {
		return nil, fmt.Errorf("URL_BASE_SWAGGER no está configurada")
	}

	// Verificar webhooks
	return webhookOps.VerifyWebhooksByURL(ctx, integrationID, baseURL)
}

// CreateWebhook crea webhooks en la plataforma externa después de verificar y eliminar duplicados
func (ic *integrationCore) CreateWebhook(ctx context.Context, integrationID string) (interface{}, error) {
	// Obtener la integración
	integration, err := ic.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener la implementación de la integración
	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return nil, fmt.Errorf("integración no registrada para tipo %d", integration.IntegrationType)
	}

	// Verificar si implementa IWebhookOperations
	webhookOps, ok := integrationImpl.(IWebhookOperations)
	if !ok {
		return nil, fmt.Errorf("esta integración no soporta operaciones de webhooks")
	}

	// Obtener baseURL desde configuración
	baseURL := ic.config.Get("URL_BASE_SWAGGER")
	if baseURL == "" {
		return nil, fmt.Errorf("URL_BASE_SWAGGER no está configurada")
	}

	// Crear webhooks
	return webhookOps.CreateWebhook(ctx, integrationID, baseURL)
}
