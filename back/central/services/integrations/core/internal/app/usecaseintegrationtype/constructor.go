package usecaseintegrationtype

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IIntegrationTypeUseCase interface {
	CreateIntegrationType(ctx context.Context, dto domain.CreateIntegrationTypeDTO) (*domain.IntegrationType, error)
	UpdateIntegrationType(ctx context.Context, id uint, dto domain.UpdateIntegrationTypeDTO) (*domain.IntegrationType, error)
	GetIntegrationTypeByID(ctx context.Context, id uint) (*domain.IntegrationType, error)
	GetIntegrationTypeByCode(ctx context.Context, code string) (*domain.IntegrationType, error)
	DeleteIntegrationType(ctx context.Context, id uint) error
	ListIntegrationTypes(ctx context.Context, categoryID *uint) ([]*domain.IntegrationType, error)
	ListActiveIntegrationTypes(ctx context.Context) ([]*domain.IntegrationType, error)

	// Platform credentials (admin only — returns decrypted map)
	GetPlatformCredentials(ctx context.Context, id uint) (map[string]interface{}, error)

	// Integration Categories
	ListIntegrationCategories(ctx context.Context) ([]*domain.IntegrationCategory, error)
}

type integrationTypeUseCase struct {
	repo       domain.IRepository
	s3         domain.IS3Service
	cache      domain.IIntegrationCache
	log        log.ILogger
	env        env.IConfig
	encryption domain.IEncryptionService
}

// New crea una nueva instancia del caso de uso de tipos de integración
func New(
	repo domain.IRepository,
	s3 domain.IS3Service,
	cache domain.IIntegrationCache,
	logger log.ILogger,
	env env.IConfig,
	encryption domain.IEncryptionService,
) IIntegrationTypeUseCase {
	return &integrationTypeUseCase{
		repo:       repo,
		s3:         s3,
		cache:      cache,
		log:        logger,
		env:        env,
		encryption: encryption,
	}
}
