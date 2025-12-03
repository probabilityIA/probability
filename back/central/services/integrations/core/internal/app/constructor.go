package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type integrationUseCase struct {
	repo       domain.IIntegrationRepository
	encryption domain.IEncryptionService
	testerReg  *IntegrationTesterRegistry
	log        log.ILogger
}

// IntegrationUseCase expone el tipo para uso en public.go
type IntegrationUseCase = integrationUseCase

// New crea una nueva instancia del caso de uso de integraciones
func New(repo domain.IIntegrationRepository, encryption domain.IEncryptionService, logger log.ILogger) domain.IIntegrationUseCase {
	return &integrationUseCase{
		repo:       repo,
		encryption: encryption,
		testerReg:  NewIntegrationTesterRegistry(),
		log:        logger,
	}
}

// GetTesterRegistry retorna el registry de testers (para uso interno del core)
func (uc *integrationUseCase) GetTesterRegistry() *IntegrationTesterRegistry {
	return uc.testerReg
}
