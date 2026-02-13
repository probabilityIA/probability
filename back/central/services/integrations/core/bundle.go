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

type IIntegrationContract interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
	SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error
	SyncOrdersByBusiness(ctx context.Context, businessID uint) error
	// GetWebhookURL construye la URL del webhook para esta integración
	// baseURL es la URL base del servidor (ej: URL_BASE_SWAGGER)
	// integrationID es el ID de la integración específica
	GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*WebhookInfo, error)
}

// IWebhookOperations es una interfaz opcional que las integraciones pueden implementar
// para soportar operaciones de webhooks (listar, eliminar, verificar, crear)
type IWebhookOperations interface {
	ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
	VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error)
	CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error)
}

// WebhookInfo es un alias del tipo de domain
type WebhookInfo = domain.WebhookInfo

type IIntegrationCore interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	GetIntegrationByStoreID(ctx context.Context, storeID string, integrationType int) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	RegisterIntegration(integrationType int, integration IIntegrationContract)
	// GetRegisteredIntegration obtiene el bundle registrado para un tipo de integración
	GetRegisteredIntegration(integrationType int) (IIntegrationContract, bool)
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
	SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error
	SyncOrdersByBusiness(ctx context.Context, businessID uint) error
	RegisterObserverForType(integrationType int, observer func(context.Context, *Integration))
	// GetWebhookURL obtiene la URL del webhook para una integración específica
	GetWebhookURL(ctx context.Context, integrationID uint) (*WebhookInfo, error)
	// UpdateIntegrationConfig actualiza el config de una integración haciendo merge con el config existente
	UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error
	// ListWebhooks lista todos los webhooks de una integración (solo para integraciones que lo soporten)
	ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error)
	// DeleteWebhook elimina un webhook de una integración (solo para integraciones que lo soporten)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
}

type integrationCore struct {
	useCase      usecaseintegrations.IIntegrationUseCase
	integrations map[int]IIntegrationContract
	logger       log.ILogger
	config       env.IConfig
}

func New(router *gin.RouterGroup, db db.IDatabase, redisClient redis.IRedis, logger log.ILogger, config env.IConfig, s3 storage.IS3Service) IIntegrationCore {
	// 1. Inicializar Servicio de Encriptación
	encryptionService := encryption.New(config, logger)

	// 2. ✅ NUEVO - Inicializar Cache Service
	integrationCache := cache.New(redisClient, logger)
	logger.Info(context.Background()).Msg("✅ Integration cache initialized")

	// Registrar prefijos de caché para startup logs
	redisClient.RegisterCachePrefix("integration:meta:*")
	redisClient.RegisterCachePrefix("integration:creds:*")
	redisClient.RegisterCachePrefix("integration:code:*")
	redisClient.RegisterCachePrefix("integration:idx:*")

	// 3. Inicializar Repositorio (con cache)
	repo := repository.New(db, logger, encryptionService, integrationCache)

	// 4. Inicializar Casos de Uso (con cache)
	IntegrationUseCase := usecaseintegrations.New(repo, encryptionService, integrationCache, logger)
	integrationTypeUseCase := usecaseintegrationtype.New(repo, s3, logger, config)

	// 5. Inicializar Handlers
	coreIntegration := &integrationCore{
		useCase:      IntegrationUseCase,
		integrations: make(map[int]IIntegrationContract),
		logger:       logger.WithModule("integrations-core"),
		config:       config,
	}

	handlerIntegrations := handlerintegrations.New(IntegrationUseCase, logger, coreIntegration, config)
	handlerIntegrationType := handlerintegrationtype.New(integrationTypeUseCase, logger, config)

	// 6. Registrar Rutas
	handlerIntegrations.RegisterRoutes(router, logger)
	handlerIntegrationType.RegisterRoutes(router, logger)

	// 7. ✅ NUEVO - Cache Warming en background (pre-carga integraciones activas)
	go func() {
		bgCtx := context.Background()
		logger.Info(bgCtx).Msg("🔥 Starting cache warming in background...")
		if err := IntegrationUseCase.WarmCache(bgCtx); err != nil {
			logger.Error(bgCtx).Err(err).Msg("❌ Cache warming failed")
		} else {
			logger.Info(bgCtx).Msg("✅ Cache warming completed successfully")
		}
	}()

	return coreIntegration
}
