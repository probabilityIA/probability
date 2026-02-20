package usecaseintegrations

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IntegrationCreatedObserver func(ctx context.Context, integration *domain.Integration)

type IIntegrationUseCase interface {
	CreateIntegration(ctx context.Context, dto domain.CreateIntegrationDTO) (*domain.Integration, error)
	UpdateIntegration(ctx context.Context, id uint, dto domain.UpdateIntegrationDTO) (*domain.Integration, error)
	GetIntegrationByID(ctx context.Context, id uint) (*domain.Integration, error)
	GetIntegrationByIDWithCredentials(ctx context.Context, id uint) (*domain.IntegrationWithCredentials, error)
	GetIntegrationByType(ctx context.Context, integrationTypeCode string, businessID *uint) (*domain.IntegrationWithCredentials, error)
	GetPublicIntegrationByID(ctx context.Context, integrationID string) (*PublicIntegration, error)
	GetIntegrationConfig(ctx context.Context, integrationType string, businessID *uint) (map[string]interface{}, error)
	DecryptCredentialField(ctx context.Context, integrationID string, fieldName string) (string, error)
	DeleteIntegration(ctx context.Context, id uint) error
	ListIntegrations(ctx context.Context, filters domain.IntegrationFilters) ([]*domain.Integration, int64, error)
	TestIntegration(ctx context.Context, id uint) error
	TestConnectionRaw(ctx context.Context, integrationTypeCode string, config map[string]interface{}, credentials map[string]interface{}) error
	ActivateIntegration(ctx context.Context, id uint) error
	DeactivateIntegration(ctx context.Context, id uint) error
	SetAsDefault(ctx context.Context, id uint) error
	UpdateLastSync(ctx context.Context, integrationID string) error
	RegisterObserver(observer IntegrationCreatedObserver)
	WarmCache(ctx context.Context) error // ✅ NUEVO - Pre-carga cache al iniciar
	SetWebhookCreator(creator IWebhookCreator)
}

type IntegrationUseCase struct {
	repo           domain.IRepository
	encryption     domain.IEncryptionService
	cache          domain.IIntegrationCache
	testerReg      *IntegrationTesterRegistry
	log            log.ILogger
	observers      []IntegrationCreatedObserver
	webhookCreator IWebhookCreator
}

// New crea una nueva instancia del caso de uso de integraciones
func New(repo domain.IRepository, encryption domain.IEncryptionService, cache domain.IIntegrationCache, logger log.ILogger) IIntegrationUseCase {
	return &IntegrationUseCase{
		repo:       repo,
		encryption: encryption,
		cache:      cache,
		testerReg:  NewIntegrationTesterRegistry(),
		log:        logger,
		observers:  make([]IntegrationCreatedObserver, 0),
	}
}

func (uc *IntegrationUseCase) RegisterObserver(observer IntegrationCreatedObserver) {
	uc.observers = append(uc.observers, observer)
}

// IWebhookCreator define la capacidad de crear webhooks (implementado por integrationCore)
type IWebhookCreator interface {
	CreateWebhook(ctx context.Context, integrationID string) (interface{}, error)
}

// GetTesterRegistry retorna el registry de testers (para uso interno del core)
func (uc *IntegrationUseCase) GetTesterRegistry() *IntegrationTesterRegistry {
	return uc.testerReg
}

// SetWebhookCreator inyecta la dependencia de creación de webhooks (para romper ciclo con core)
func (uc *IntegrationUseCase) SetWebhookCreator(creator IWebhookCreator) {
	uc.webhookCreator = creator
}
