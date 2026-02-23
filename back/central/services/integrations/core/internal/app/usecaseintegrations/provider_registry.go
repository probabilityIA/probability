package usecaseintegrations

import (
	"sync"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

// providerRegistry mantiene un registro unificado de providers por tipo de integración.
// Reemplaza el antiguo IntegrationTesterRegistry y el map de integrationCore.
type providerRegistry struct {
	providers map[int]domain.IIntegrationContract
	mu        sync.RWMutex
}

func newProviderRegistry() *providerRegistry {
	return &providerRegistry{
		providers: make(map[int]domain.IIntegrationContract),
	}
}

// Register registra un provider para un tipo de integración
func (r *providerRegistry) Register(integrationType int, provider domain.IIntegrationContract) {
	if integrationType == 0 || provider == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.providers[integrationType] = provider
}

// Get obtiene el provider registrado para un tipo de integración
func (r *providerRegistry) Get(integrationType int) (domain.IIntegrationContract, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[integrationType]
	return provider, exists
}

// ListRegisteredTypes retorna la lista de tipos registrados
func (r *providerRegistry) ListRegisteredTypes() []int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]int, 0, len(r.providers))
	for t := range r.providers {
		types = append(types, t)
	}
	return types
}
