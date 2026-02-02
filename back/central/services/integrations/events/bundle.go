package events

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/primary"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/primary/handlers"
	eventsevents "github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/secondary/events"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa el módulo de eventos de integraciones
func New(router *gin.RouterGroup, logger log.ILogger) (domain.IIntegrationEventService, domain.IIntegrationEventPublisher) {
	// 1. Init Event Manager (para SSE y eventos en tiempo real)
	eventManager := eventsevents.New(logger)

	// 2. Init Event Service
	eventService := app.NewIntegrationEventService(eventManager)

	// 3. Init SSE Handler
	sseHandler := handlers.New(eventManager, logger)

	// 4. Init Routes
	routes := primary.New(sseHandler)

	// 5. Register Routes
	routes.RegisterRoutes(router)

	return eventService, eventManager
}
