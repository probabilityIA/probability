package core

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrationtype"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/encryption"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
)

type IIntegrationContract = domain.IIntegrationContract
type BaseIntegration = domain.BaseIntegration
type WebhookInfo = domain.WebhookInfo
type IntegrationWithCredentials = domain.IntegrationWithCredentials
type PublicIntegration = domain.PublicIntegration

var ErrNotSupported = domain.ErrNotSupported

const (
	IntegrationTypeShopify      = domain.IntegrationTypeShopify
	IntegrationTypeWhatsApp     = domain.IntegrationTypeWhatsApp
	IntegrationTypeMercadoLibre = domain.IntegrationTypeMercadoLibre
	IntegrationTypeWoocommerce  = domain.IntegrationTypeWoocommerce
	IntegrationTypeInvoicing    = domain.IntegrationTypeInvoicing
	IntegrationTypePlatform     = domain.IntegrationTypePlatform
	IntegrationTypeFactus       = domain.IntegrationTypeFactus
	IntegrationTypeSiigo        = domain.IntegrationTypeSiigo
	IntegrationTypeAlegra       = domain.IntegrationTypeAlegra
	IntegrationTypeWorldOffice  = domain.IntegrationTypeWorldOffice
	IntegrationTypeHelisa       = domain.IntegrationTypeHelisa
	IntegrationTypeEnvioClick   = domain.IntegrationTypeEnvioClick
	IntegrationTypeEnviame      = domain.IntegrationTypeEnviame
	IntegrationTypeTu           = domain.IntegrationTypeTu
	IntegrationTypeMiPaquete    = domain.IntegrationTypeMiPaquete
	IntegrationTypeVTEX         = domain.IntegrationTypeVTEX
	IntegrationTypeTiendanube   = domain.IntegrationTypeTiendanube
	IntegrationTypeMagento      = domain.IntegrationTypeMagento
	IntegrationTypeAmazon       = domain.IntegrationTypeAmazon
	IntegrationTypeFalabella    = domain.IntegrationTypeFalabella
	IntegrationTypeExito        = domain.IntegrationTypeExito
	IntegrationTypeEmail        = domain.IntegrationTypeEmail
	IntegrationTypeTienda       = domain.IntegrationTypeTienda
	IntegrationTypeTiendaWeb    = domain.IntegrationTypeTiendaWeb
	IntegrationTypeJumpseller   = domain.IntegrationTypeJumpseller
)

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*domain.PublicIntegration, error)
	GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.PublicIntegration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error
	UpdateIntegrationCredentials(ctx context.Context, integrationID string, credentials map[string]interface{}) error
	GetIntegrationConfig(ctx context.Context, integrationID string) (map[string]interface{}, error)
	GetPlatformCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
}

type IIntegrationCore interface {
	IIntegrationService
	RegisterIntegration(integrationType int, integration IIntegrationContract)
	OnIntegrationCreated(integrationType int, observer func(context.Context, *domain.PublicIntegration))
	GetRegisteredIntegration(integrationType int) (IIntegrationContract, bool)
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
	SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error
	SyncOrdersByIntegrationIDWithBatches(ctx context.Context, integrationID string, params *domain.SyncBatchParams) error
	SyncOrdersByBusiness(ctx context.Context, businessID uint) error
	GetWebhookURL(ctx context.Context, integrationID uint) (*WebhookInfo, error)
	ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
	GetCachedPlatformCredentials(ctx context.Context, integrationTypeID uint) (map[string]any, error)
	GetIntegrationIDByBusinessAndType(ctx context.Context, businessID, integrationTypeID uint) (uint, error)
	GetIntegrationTypeByCode(ctx context.Context, code string) (*domain.IntegrationType, error)
	SetEcommerceLimitChecker(checker domain.EcommerceLimitChecker)
}

type integrationCore struct {
	useCase    domain.IIntegrationUseCase
	cache      domain.IIntegrationCache
	repo       domain.IRepository
	encryption domain.IEncryptionService
}

func New(router *gin.RouterGroup, db db.IDatabase, redisClient redis.IRedis, logger log.ILogger, config env.IConfig, s3 storage.IS3Service, rabbitMQ rabbitmq.IQueue) IIntegrationCore {
	encryptionService := encryption.New(config, logger)

	integrationCache := cache.New(redisClient, logger)
	logger.Info(context.Background()).Msg("Integration cache initialized")

	redisClient.RegisterCachePrefix("integration:meta:*")
	redisClient.RegisterCachePrefix("integration:creds:*")
	redisClient.RegisterCachePrefix("integration:code:*")
	redisClient.RegisterCachePrefix("integration:idx:*")
	redisClient.RegisterCachePrefix("integration:platform_creds:*")

	repo := repository.New(db, logger, encryptionService, integrationCache)

	integrationUseCase := usecaseintegrations.New(repo, encryptionService, integrationCache, logger, config, rabbitMQ)
	integrationTypeUseCase := usecaseintegrationtype.New(repo, s3, integrationCache, logger, config, encryptionService)

	handlerIntegrations := handlerintegrations.New(integrationUseCase, logger, config)
	handlerIntegrationType := handlerintegrationtype.New(integrationTypeUseCase, logger, config)

	handlerIntegrations.RegisterRoutes(router, logger)
	handlerIntegrationType.RegisterRoutes(router, logger)

	go func() {
		bgCtx := context.Background()
		logger.Info(bgCtx).Msg("Starting cache warming in background...")
		if err := integrationUseCase.WarmCache(bgCtx); err != nil {
			logger.Error(bgCtx).Err(err).Msg("Cache warming failed")
		} else {
			logger.Info(bgCtx).Msg("Cache warming completed successfully")
		}
	}()

	syncBatchConsumer := consumer.NewSyncBatchConsumer(
		rabbitMQ,
		func(integrationTypeID int) (domain.IIntegrationContract, bool) {
			return integrationUseCase.GetProvider(integrationTypeID)
		},
		logger,
	)
	go func() {
		bgCtx := context.Background()
		if err := syncBatchConsumer.Start(bgCtx); err != nil {
			logger.Error(bgCtx).Err(err).Msg("Error al iniciar SyncBatchConsumer")
		}
	}()

	return &integrationCore{useCase: integrationUseCase, cache: integrationCache, repo: repo, encryption: encryptionService}
}

func (ic *integrationCore) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.PublicIntegration, error) {
	return ic.useCase.GetPublicIntegrationByID(ctx, integrationID)
}

func (ic *integrationCore) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.PublicIntegration, error) {
	return ic.useCase.GetIntegrationByExternalID(ctx, externalID, integrationType)
}

func (ic *integrationCore) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return ic.useCase.DecryptCredentialField(ctx, integrationID, fieldName)
}

func (ic *integrationCore) UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error {
	return ic.useCase.UpdateIntegrationConfig(ctx, integrationID, newConfig)
}

func (ic *integrationCore) UpdateIntegrationCredentials(ctx context.Context, integrationID string, credentials map[string]interface{}) error {
	id, err := strconv.ParseUint(integrationID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integration id %q: %w", integrationID, err)
	}
	creds := credentials
	_, err = ic.useCase.UpdateIntegration(ctx, uint(id), domain.UpdateIntegrationDTO{Credentials: &creds})
	return err
}

func (ic *integrationCore) GetIntegrationConfig(ctx context.Context, integrationID string) (map[string]interface{}, error) {
	integration, err := ic.useCase.GetPublicIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	if integration == nil {
		return nil, nil
	}
	return integration.Config, nil
}

func (ic *integrationCore) GetPlatformCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return ic.useCase.GetPlatformCredentialByIntegrationID(ctx, integrationID, fieldName)
}

func (ic *integrationCore) RegisterIntegration(integrationType int, integration IIntegrationContract) {
	ic.useCase.RegisterProvider(integrationType, integration)
}

func (ic *integrationCore) OnIntegrationCreated(integrationType int, observer func(context.Context, *domain.PublicIntegration)) {
	ic.useCase.OnIntegrationCreated(integrationType, observer)
}

func (ic *integrationCore) GetRegisteredIntegration(integrationType int) (IIntegrationContract, bool) {
	return ic.useCase.GetProvider(integrationType)
}

func (ic *integrationCore) SetEcommerceLimitChecker(checker domain.EcommerceLimitChecker) {
	ic.useCase.SetEcommerceLimitChecker(checker)
}

func (ic *integrationCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return ic.useCase.TestConnectionFromConfig(ctx, config, credentials)
}

func (ic *integrationCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return ic.useCase.SyncOrdersByIntegrationID(ctx, integrationID)
}

func (ic *integrationCore) SyncOrdersByIntegrationIDWithBatches(ctx context.Context, integrationID string, params *domain.SyncBatchParams) error {
	return ic.useCase.SyncOrdersByIntegrationIDWithBatches(ctx, integrationID, params)
}

func (ic *integrationCore) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	return ic.useCase.SyncOrdersByBusiness(ctx, businessID)
}

func (ic *integrationCore) GetWebhookURL(ctx context.Context, integrationID uint) (*WebhookInfo, error) {
	return ic.useCase.GetWebhookURL(ctx, integrationID)
}

func (ic *integrationCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	return ic.useCase.ListWebhooks(ctx, integrationID)
}

func (ic *integrationCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return ic.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}

func (ic *integrationCore) GetCachedPlatformCredentials(ctx context.Context, integrationTypeID uint) (map[string]any, error) {
	creds, err := ic.cache.GetPlatformCredentials(ctx, integrationTypeID)
	if err == nil {
		return creds, nil
	}

	intType, err := ic.repo.GetIntegrationTypeByID(ctx, integrationTypeID)
	if err != nil {
		return nil, err
	}
	if len(intType.PlatformCredentialsEncrypted) == 0 {
		return nil, fmt.Errorf("no platform credentials for integration type %d", integrationTypeID)
	}

	creds, err = ic.encryption.DecryptCredentials(ctx, intType.PlatformCredentialsEncrypted)
	if err != nil {
		return nil, err
	}

	_ = ic.cache.SetPlatformCredentials(ctx, integrationTypeID, creds)

	return creds, nil
}

func (ic *integrationCore) GetIntegrationTypeByCode(ctx context.Context, code string) (*domain.IntegrationType, error) {
	return ic.repo.GetIntegrationTypeByCode(ctx, code)
}

func (ic *integrationCore) GetIntegrationIDByBusinessAndType(ctx context.Context, businessID, integrationTypeID uint) (uint, error) {
	if cached, err := ic.cache.GetByBusinessAndType(ctx, businessID, integrationTypeID); err == nil && cached != nil {
		return cached.ID, nil
	}

	bizID := businessID
	integration, err := ic.repo.GetActiveIntegrationByIntegrationTypeID(ctx, integrationTypeID, &bizID)
	if err != nil {
		return 0, fmt.Errorf("integration not found for business %d type %d: %w", businessID, integrationTypeID, err)
	}
	if integration == nil || integration.ID == 0 {
		return 0, fmt.Errorf("no active integration for business %d type %d", businessID, integrationTypeID)
	}

	_ = ic.cache.SetBusinessTypeIndex(ctx, businessID, integrationTypeID, integration.ID)

	return integration.ID, nil
}
