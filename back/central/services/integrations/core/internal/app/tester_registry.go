package app

import (
	"context"
	"fmt"
	"sync"
)

// ITestIntegration define la interfaz que cada integración debe implementar para testear su conexión
// Esta es la versión interna de la interfaz
type ITestIntegration interface {
	// TestConnection prueba la conexión con las credenciales y configuración dadas
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

// IntegrationTesterRegistry mantiene un registro de testers por tipo de integración
type IntegrationTesterRegistry struct {
	testers map[string]ITestIntegration
	mu      sync.RWMutex
}

// NewIntegrationTesterRegistry crea una nueva instancia del registry
func NewIntegrationTesterRegistry() *IntegrationTesterRegistry {
	return &IntegrationTesterRegistry{
		testers: make(map[string]ITestIntegration),
	}
}

// Register registra un tester para un tipo de integración
func (r *IntegrationTesterRegistry) Register(integrationType string, tester ITestIntegration) error {
	if integrationType == "" {
		return fmt.Errorf("tipo de integración no puede estar vacío")
	}
	if tester == nil {
		return fmt.Errorf("tester no puede ser nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.testers[integrationType] = tester
	return nil
}

// GetTester obtiene el tester registrado para un tipo de integración
func (r *IntegrationTesterRegistry) GetTester(integrationType string) (ITestIntegration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tester, exists := r.testers[integrationType]
	if !exists {
		return nil, fmt.Errorf("tester no registrado para tipo: %s", integrationType)
	}

	return tester, nil
}

// IsRegistered verifica si hay un tester registrado para un tipo
func (r *IntegrationTesterRegistry) IsRegistered(integrationType string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.testers[integrationType]
	return exists
}

// ListRegisteredTypes retorna la lista de tipos registrados
func (r *IntegrationTesterRegistry) ListRegisteredTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.testers))
	for t := range r.testers {
		types = append(types, t)
	}
	return types
}
