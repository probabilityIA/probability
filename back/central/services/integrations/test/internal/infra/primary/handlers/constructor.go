package handlers

import (
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/app/usecases"
)

// Handlers contiene todos los handlers del módulo test
type Handlers struct {
	uc *usecases.UseCases
}

// New crea una nueva instancia de Handlers
func New(uc *usecases.UseCases) *Handlers {
	return &Handlers{
		uc: uc,
	}
}
