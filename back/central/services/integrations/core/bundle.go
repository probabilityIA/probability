package core

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrationtype"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/encryption"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
)

// ============================================
// Re-exports para backward compatibility
// ============================================

// Type aliases — consumidores externos siguen usando core.IIntegrationContract etc.
type IIntegrationContract = domain.IIntegrationContract
type BaseIntegration = domain.BaseIntegration
type WebhookInfo = domain.WebhookInfo
type IntegrationWithCredentials = domain.IntegrationWithCredentials
type PublicIntegration = domain.PublicIntegration

// Sentinel error
var ErrNotSupported = domain.ErrNotSupported

// Integration type constants
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
)

// ============================================
// Interfaces públicas
// ============================================

// IIntegrationService expone las operaciones de consulta y configuración que los módulos
// consumidores (facturación, ecommerce, etc.) necesitan del core de integraciones.
type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*domain.PublicIntegration, error)
	GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.PublicIntegration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error
	GetIntegrationConfig(ctx context.Context, integrationID string) (map[string]interface{}, error)
	// GetPlatformCredential decrypts a field from the integration type's platform credentials.
	// Use when an integration has use_platform_token=true in its config.
	GetPlatformCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
}

// IIntegrationCore es la interfaz completa del core de integraciones.
// Embeds IIntegrationService + métodos de registro y operaciones internas.
// Solo debe usarse en integrations/bundle.go y shopify/bundle.go.
type IIntegrationCore interface {
	IIntegrationService
	RegisterIntegration(integrationType int, integration IIntegrationContract)
	OnIntegrationCreated(integrationType int, observer func(context.Context, *domain.PublicIntegration))
	GetRegisteredIntegration(integrationType int) (IIntegrationContract, bool)
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
	SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error
	SyncOrdersByBusiness(ctx context.Context, businessID uint) error
	GetWebhookURL(ctx context.Context, integrationID uint) (*WebhookInfo, error)
	ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
}

// ============================================
// Thin facade wrapping the use case
// ============================================

type integrationCore struct {
	useCase domain.IIntegrationUseCase
}

// ============================================
// Constructor
// ============================================

func New(router *gin.RouterGroup, db db.IDatabase, redisClient redis.IRedis, logger log.ILogger, config env.IConfig, s3 storage.IS3Service) IIntegrationCore {
	// 1. Inicializar Servicio de Encriptación
	encryptionService := encryption.New(config, logger)

	// 2. Inicializar Cache Service
	integrationCache := cache.New(redisClient, logger)
	logger.Info(context.Background()).Msg("Integration cache initialized")

	// Registrar prefijos de caché para startup logs
	redisClient.RegisterCachePrefix("integration:meta:*")
	redisClient.RegisterCachePrefix("integration:creds:*")
	redisClient.RegisterCachePrefix("integration:code:*")
	redisClient.RegisterCachePrefix("integration:idx:*")

	// 3. Inicializar Repositorio (con cache)
	repo := repository.New(db, logger, encryptionService, integrationCache)

	// 4. Inicializar Casos de Uso
	integrationUseCase := usecaseintegrations.New(repo, encryptionService, integrationCache, logger, config)
	integrationTypeUseCase := usecaseintegrationtype.New(repo, s3, integrationCache, logger, config, encryptionService)

	// 5. Inicializar Handlers (solo dependen del use case)
	handlerIntegrations := handlerintegrations.New(integrationUseCase, logger, config)
	handlerIntegrationType := handlerintegrationtype.New(integrationTypeUseCase, logger, config)

	// 6. Registrar Rutas
	handlerIntegrations.RegisterRoutes(router, logger)
	handlerIntegrationType.RegisterRoutes(router, logger)

	// 7. Cache Warming en background
	go func() {
		bgCtx := context.Background()
		logger.Info(bgCtx).Msg("Starting cache warming in background...")
		if err := integrationUseCase.WarmCache(bgCtx); err != nil {
			logger.Error(bgCtx).Err(err).Msg("Cache warming failed")
		} else {
			logger.Info(bgCtx).Msg("Cache warming completed successfully")
		}
	}()

	return &integrationCore{useCase: integrationUseCase}
}

// IIntegrationService pass-throughs
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

// IIntegrationCore pass-throughs
func (ic *integrationCore) RegisterIntegration(integrationType int, integration IIntegrationContract) {
	ic.useCase.RegisterProvider(integrationType, integration)
}

func (ic *integrationCore) OnIntegrationCreated(integrationType int, observer func(context.Context, *domain.PublicIntegration)) {
	ic.useCase.OnIntegrationCreated(integrationType, observer)
}

func (ic *integrationCore) GetRegisteredIntegration(integrationType int) (IIntegrationContract, bool) {
	return ic.useCase.GetProvider(integrationType)
}

func (ic *integrationCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return ic.useCase.TestConnectionFromConfig(ctx, config, credentials)
}

func (ic *integrationCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return ic.useCase.SyncOrdersByIntegrationID(ctx, integrationID)
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
