package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

// ITestIntegration define la interfaz que cada integración debe implementar para testear su conexión
// Esta interfaz es pública para que las integraciones puedan implementarla
type ITestIntegration interface {
	// TestConnection prueba la conexión con las credenciales y configuración dadas
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

// Constantes públicas para tipos de integración
const (
	IntegrationTypeWhatsApp     = "whatsapp"
	IntegrationTypeShopify      = "shopify"
	IntegrationTypeMercadoLibre = "mercado_libre"
)

// IntegrationWithCredentials representa una integración con credenciales desencriptadas
// Este es un tipo público que envuelve el tipo interno
type IntegrationWithCredentials = domain.IntegrationWithCredentials

// IIntegrationCore es la interfaz pública que expone Core para que otras integraciones lo consuman
type IIntegrationCore interface {
	// GetIntegrationByType obtiene una integración con credenciales desencriptadas (para uso interno)
	GetIntegrationByType(ctx context.Context, integrationType string, businessID *uint) (*IntegrationWithCredentials, error)

	// GetIntegrationConfig obtiene solo la configuración de una integración (sin credenciales)
	GetIntegrationConfig(ctx context.Context, integrationType string, businessID *uint) (map[string]interface{}, error)

	// TestIntegration testea la conexión de una integración usando su tester registrado
	TestIntegration(ctx context.Context, integrationType string, config map[string]interface{}, credentials map[string]interface{}) error

	// RegisterTester registra un tester para un tipo de integración
	RegisterTester(integrationType string, tester ITestIntegration) error
}

// integrationCore implementa IIntegrationCore
type integrationCore struct {
	useCase usecaseintegrations.IIntegrationUseCase
}

// NewIntegrationCore crea una nueva instancia de IIntegrationCore
func NewIntegrationCore(useCase usecaseintegrations.IIntegrationUseCase) IIntegrationCore {
	return &integrationCore{
		useCase: useCase,
	}
}

// GetIntegrationByType obtiene una integración con credenciales desencriptadas
func (ic *integrationCore) GetIntegrationByType(ctx context.Context, integrationType string, businessID *uint) (*domain.IntegrationWithCredentials, error) {
	return ic.useCase.GetIntegrationByType(ctx, integrationType, businessID)
}

// GetIntegrationConfig obtiene solo la configuración (sin credenciales)
func (ic *integrationCore) GetIntegrationConfig(ctx context.Context, integrationType string, businessID *uint) (map[string]interface{}, error) {
	integration, err := ic.useCase.GetIntegrationByType(ctx, integrationType, businessID)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if len(integration.Config) > 0 {
		// Convertir datatypes.JSON a map
		if err := json.Unmarshal(integration.Config, &config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// TestIntegration testea la conexión usando el tester registrado
func (ic *integrationCore) TestIntegration(ctx context.Context, integrationType string, config map[string]interface{}, credentials map[string]interface{}) error {
	// Obtener el usecase interno para acceder al registry
	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		return fmt.Errorf("error interno: no se puede acceder al registry de testers")
	}

	tester, err := useCaseImpl.GetTesterRegistry().GetTester(integrationType)
	if err != nil {
		return fmt.Errorf("no hay tester registrado para tipo %s: %w", integrationType, err)
	}

	// Convertir ITestIntegration pública a interna usando un adapter
	adapter := &testIntegrationAdapter{tester: tester}
	return adapter.TestConnection(ctx, config, credentials)
}

// RegisterTester registra un tester para un tipo de integración
func (ic *integrationCore) RegisterTester(integrationType string, tester ITestIntegration) error {
	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		return fmt.Errorf("error interno: no se puede acceder al registry de testers")
	}

	// Convertir ITestIntegration pública a interna usando un adapter
	adapter := &testIntegrationAdapter{tester: tester}
	return useCaseImpl.GetTesterRegistry().Register(integrationType, adapter)
}

// testIntegrationAdapter adapta ITestIntegration pública a la interfaz interna
type testIntegrationAdapter struct {
	tester ITestIntegration
}

func (a *testIntegrationAdapter) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return a.tester.TestConnection(ctx, config, credentials)
}
