package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

// ITestIntegration define la interfaz que cada integración debe implementar para testear su conexión
type ITestIntegration interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

// Constantes públicas para tipos de integración
const (
	IntegrationTypeWhatsApp     = "whatsapp"
	IntegrationTypeShopify      = "shopify"
	IntegrationTypeMercadoLibre = "mercado_libre"
)

// IntegrationWithCredentials representa una integración con credenciales desencriptadas
type IntegrationWithCredentials = domain.IntegrationWithCredentials

// Integration es el tipo público para representar una integración
type Integration struct {
	ID         uint
	BusinessID *uint
	Name       string
	Config     interface{}
}

// IIntegrationCore es la interfaz pública que expone Core para que otras integraciones lo consuman
type IIntegrationCore interface {
	GetIntegrationByType(ctx context.Context, integrationType string, businessID *uint) (*IntegrationWithCredentials, error)
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	GetIntegrationConfig(ctx context.Context, integrationType string, businessID *uint) (map[string]interface{}, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateLastSync(ctx context.Context, integrationID string) error
	TestIntegration(ctx context.Context, integrationType string, config map[string]interface{}, credentials map[string]interface{}) error
	RegisterTester(integrationType string, tester ITestIntegration) error
	RegisterObserver(observer func(context.Context, *Integration))
}

// integrationCore implementa IIntegrationCore
type integrationCore struct {
	useCase usecaseintegrations.IIntegrationUseCase
}

// NewIntegrationCore crea una nueva instancia de IIntegrationCore
func NewIntegrationCore(useCase usecaseintegrations.IIntegrationUseCase) IIntegrationCore {
	return &integrationCore{useCase: useCase}
}

// GetIntegrationByType obtiene una integración con credenciales desencriptadas
func (ic *integrationCore) GetIntegrationByType(ctx context.Context, integrationType string, businessID *uint) (*domain.IntegrationWithCredentials, error) {
	return ic.useCase.GetIntegrationByType(ctx, integrationType, businessID)
}

// GetIntegrationByID obtiene una integración por su ID
func (ic *integrationCore) GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error) {
	var id uint
	if _, err := fmt.Sscanf(integrationID, "%d", &id); err != nil {
		return nil, fmt.Errorf("invalid integration ID: %w", err)
	}

	integration, err := ic.useCase.GetIntegrationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &config); err != nil {
			return nil, err
		}
	}

	return &Integration{
		ID:         integration.ID,
		BusinessID: integration.BusinessID,
		Name:       integration.Name,
		Config:     config,
	}, nil
}

// GetIntegrationConfig obtiene solo la configuración (sin credenciales)
func (ic *integrationCore) GetIntegrationConfig(ctx context.Context, integrationType string, businessID *uint) (map[string]interface{}, error) {
	integration, err := ic.useCase.GetIntegrationByType(ctx, integrationType, businessID)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// DecryptCredential desencripta un campo específico de las credenciales
func (ic *integrationCore) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	var id uint
	if _, err := fmt.Sscanf(integrationID, "%d", &id); err != nil {
		return "", fmt.Errorf("invalid integration ID: %w", err)
	}

	integration, err := ic.useCase.GetIntegrationByIDWithCredentials(ctx, id)
	if err != nil {
		return "", err
	}

	if integration.DecryptedCredentials == nil {
		return "", fmt.Errorf("no credentials found for integration")
	}

	value, ok := integration.DecryptedCredentials[fieldName]
	if !ok {
		return "", fmt.Errorf("field %s not found in credentials", fieldName)
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("field %s is not a string", fieldName)
	}

	return strValue, nil
}

// UpdateLastSync actualiza el timestamp de última sincronización
func (ic *integrationCore) UpdateLastSync(ctx context.Context, integrationID string) error {
	return ic.useCase.UpdateLastSync(ctx, integrationID)
}

// TestIntegration testea la conexión usando el tester registrado
func (ic *integrationCore) TestIntegration(ctx context.Context, integrationType string, config map[string]interface{}, credentials map[string]interface{}) error {
	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		return fmt.Errorf("error interno: no se puede acceder al registry de testers")
	}

	tester, err := useCaseImpl.GetTesterRegistry().GetTester(integrationType)
	if err != nil {
		return fmt.Errorf("no hay tester registrado para tipo %s: %w", integrationType, err)
	}

	adapter := &testIntegrationAdapter{tester: tester}
	return adapter.TestConnection(ctx, config, credentials)
}

// RegisterTester registra un tester para un tipo de integración
func (ic *integrationCore) RegisterTester(integrationType string, tester ITestIntegration) error {
	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		return fmt.Errorf("error interno: no se puede acceder al registry de testers")
	}

	adapter := &testIntegrationAdapter{tester: tester}
	return useCaseImpl.GetTesterRegistry().Register(integrationType, adapter)
}

// RegisterObserver registra un observador para eventos de creación
func (ic *integrationCore) RegisterObserver(observer func(context.Context, *Integration)) {
	ic.useCase.RegisterObserver(func(ctx context.Context, integration *domain.Integration) {
		// Map domain.Integration to public core.Integration
		var config map[string]interface{}
		if len(integration.Config) > 0 {
			_ = json.Unmarshal(integration.Config, &config)
		}

		publicIntegration := &Integration{
			ID:         integration.ID,
			BusinessID: integration.BusinessID,
			Name:       integration.Name,
			Config:     config,
		}
		observer(ctx, publicIntegration)
	})
}

// testIntegrationAdapter adapta ITestIntegration pública a la interfaz interna
type testIntegrationAdapter struct {
	tester ITestIntegration
}

func (a *testIntegrationAdapter) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return a.tester.TestConnection(ctx, config, credentials)
}
