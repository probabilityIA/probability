package notification_config

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa y registra el módulo de configuración de notificaciones
func New(router *gin.RouterGroup, database db.IDatabase, redisClient redisclient.IRedis, logger log.ILogger) {
	logger = logger.WithModule("notification_config")

	// 1. Infraestructura secundaria (adaptadores de salida)
	repo := repository.New(database, logger)
	notificationTypeRepo := repository.NewNotificationTypeRepository(database, logger)
	notificationEventTypeRepo := repository.NewNotificationEventTypeRepository(database, logger)
	orderStatusQuerier := repository.NewOrderStatusQuerier(database, logger)
	messageAuditQuerier := repository.NewMessageAuditQuerier(database, logger)

	// Cache Manager
	cacheManager := cache.New(redisClient, repo, orderStatusQuerier, logger)

	// Warmup inicial del cache
	ctx := context.Background()
	if err := cacheManager.WarmupCache(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("❌ Error en warmup de cache de notification configs - sistema continuará sin cache")
	}

	// 2. Capa de aplicación (casos de uso) - inyectar cache manager
	useCase := app.New(repo, notificationTypeRepo, notificationEventTypeRepo, cacheManager, messageAuditQuerier, logger)

	// 3. Infraestructura primaria (adaptadores de entrada)
	configHandler := notification_config.New(useCase, logger)
	typeHandler := notification_type.New(useCase, logger)
	eventTypeHandler := notification_event_type.New(useCase, logger)
	auditHandler := message_audit.New(useCase, logger)

	// 4. Registrar rutas HTTP
	configHandler.RegisterRoutes(router)
	typeHandler.RegisterRoutes(router)
	eventTypeHandler.RegisterRoutes(router)
	auditHandler.RegisterRoutes(router)
}
