package notification_config

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa y registra el módulo de configuración de notificaciones
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger) {
	logger = logger.WithModule("notification_config")

	// 1. Infraestructura secundaria (adaptadores de salida)
	repo := repository.New(database, logger)
	notificationTypeRepo := repository.NewNotificationTypeRepository(database, logger)
	notificationEventTypeRepo := repository.NewNotificationEventTypeRepository(database, logger)

	// 2. Capa de aplicación (casos de uso)
	useCase := app.New(repo, notificationTypeRepo, notificationEventTypeRepo, logger)

	// 3. Infraestructura primaria (adaptadores de entrada)
	configHandler := notification_config.New(useCase, logger)
	typeHandler := notification_type.New(useCase, logger)
	eventTypeHandler := notification_event_type.New(useCase, logger)

	// 4. Registrar rutas HTTP
	configHandler.RegisterRoutes(router)
	typeHandler.RegisterRoutes(router)
	eventTypeHandler.RegisterRoutes(router)

	logger.Info().Msg("Notification config module initialized")
}
