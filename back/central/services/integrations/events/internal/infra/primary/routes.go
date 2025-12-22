package primary

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/primary/handlers"
)

type routes struct {
	sseHandler handlers.IntegrationSSEHandlerInterface
}

type IRoutes interface {
	RegisterRoutes(router *gin.RouterGroup)
}

func New(sseHandler handlers.IntegrationSSEHandlerInterface) IRoutes {
	return &routes{
		sseHandler: sseHandler,
	}
}

func (r *routes) RegisterRoutes(router *gin.RouterGroup) {
	eventsGroup := router.Group("/events")
	{
		// SSE endpoint para eventos de integraciones por business_id
		// Ejemplo: /events/sse/:businessID?integration_id=13&event_types=integration.sync.order.created
		eventsGroup.GET("/sse/:businessID", r.sseHandler.HandleSSE)

		// SSE endpoint para super usuario
		// Ejemplo: /events/sse?integration_id=13&event_types=integration.sync.order.created
		eventsGroup.GET("/sse", r.sseHandler.HandleSSE)

		// Endpoint para consultar estado de sincronización
		// GET /events/sync-status/:integrationID
		eventsGroup.GET("/sync-status/:integrationID", r.sseHandler.GetSyncStatus)
	}
}
