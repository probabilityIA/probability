package usecaseintegrations

import (
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IntegrationUseCase struct {
	repo        domain.IRepository
	encryption  domain.IEncryptionService
	cache       domain.IIntegrationCache
	providerReg *providerRegistry
	log         log.ILogger
	observers   []domain.IntegrationCreatedObserver
	config      env.IConfig
}

// New crea una nueva instancia del caso de uso de integraciones
func New(repo domain.IRepository, encryption domain.IEncryptionService, cache domain.IIntegrationCache, logger log.ILogger, config env.IConfig) *IntegrationUseCase {
	return &IntegrationUseCase{
		repo:        repo,
		encryption:  encryption,
		cache:       cache,
		providerReg: newProviderRegistry(),
		log:         logger,
		observers:   make([]domain.IntegrationCreatedObserver, 0),
		config:      config,
	}
}

func (uc *IntegrationUseCase) RegisterObserver(observer domain.IntegrationCreatedObserver) {
	uc.observers = append(uc.observers, observer)
}

// RegisterProvider registra un provider para un tipo de integración
func (uc *IntegrationUseCase) RegisterProvider(integrationType int, provider domain.IIntegrationContract) {
	uc.providerReg.Register(integrationType, provider)
}

// GetProvider obtiene el provider registrado para un tipo de integración
func (uc *IntegrationUseCase) GetProvider(integrationType int) (domain.IIntegrationContract, bool) {
	return uc.providerReg.Get(integrationType)
}
