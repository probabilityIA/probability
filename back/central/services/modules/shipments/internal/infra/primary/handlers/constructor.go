package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecaseenvioclick"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
)

// Handlers contiene todos los handlers del módulo shipments
type Handlers struct {
	uc           *usecases.UseCases
	envioClickUC *usecaseenvioclick.UseCaseEnvioClick
}

// New crea una nueva instancia de Handlers
func New(uc *usecases.UseCases, envioClickUC *usecaseenvioclick.UseCaseEnvioClick) *Handlers {
	return &Handlers{
		uc:           uc,
		envioClickUC: envioClickUC,
	}
}
