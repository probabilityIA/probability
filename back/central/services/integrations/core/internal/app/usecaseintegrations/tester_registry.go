package usecaseintegrations

import (
	"context"
	"fmt"
	"sync"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

// ITestIntegration define la interfaz que cada integración debe implementar para testear su conexión
// Esta es la versión interna de la interfaz
type ITestIntegration interface {
	// TestConnection prueba la conexión con las credenciales y configuración dadas
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

// IntegrationTesterRegistry mantiene un registro de testers por tipo de integración
type IntegrationTesterRegistry struct {
	testers map[int]ITestIntegration
	mu      sync.RWMutex
}

// NewIntegrationTesterRegistry crea una nueva instancia del registry
func NewIntegrationTesterRegistry() *IntegrationTesterRegistry {
	return &IntegrationTesterRegistry{
		testers: make(map[int]ITestIntegration),
	}
}

// Register registra un tester para un tipo de integración
func (r *IntegrationTesterRegistry) Register(integrationType int, tester ITestIntegration) error {
	if integrationType == 0 {
		return domain.ErrTesterTypeEmpty
	}
	if tester == nil {
		return domain.ErrTesterNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.testers[integrationType] = tester
	return nil
}

// GetTester obtiene el tester registrado para un tipo de integración
func (r *IntegrationTesterRegistry) GetTester(integrationType int) (ITestIntegration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tester, exists := r.testers[integrationType]
	if !exists {
		return nil, fmt.Errorf("%w: %d", domain.ErrTesterNotRegistered, integrationType)
	}

	return tester, nil
}

// IsRegistered verifica si hay un tester registrado para un tipo
func (r *IntegrationTesterRegistry) IsRegistered(integrationType int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.testers[integrationType]
	return exists
}

// ListRegisteredTypes retorna la lista de tipos registrados
func (r *IntegrationTesterRegistry) ListRegisteredTypes() []int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]int, 0, len(r.testers))
	for t := range r.testers {
		types = append(types, t)
	}
	return types
}
