package primary

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/primary/handlers"
)

type routes struct {
	sseHandler handlers.SSEHandlerInterface
}

type IRoutes interface {
	RegisterRoutes(router *gin.RouterGroup)
}

func New(sseHandler handlers.SSEHandlerInterface) IRoutes {
	return &routes{
		sseHandler: sseHandler,
	}
}

func (r *routes) RegisterRoutes(router *gin.RouterGroup) {
	notifyGroup := router.Group("/notify")
	{
		// SSE endpoint para notificaciones de órdenes por business_id
		// Ejemplo: /notify/sse/order-notify/:businessID?integration_id=123&event_types=order.created,order.updated
		notifyGroup.GET("/sse/order-notify/:businessID", r.sseHandler.HandleSSE)

		// SSE endpoint para super usuario (business_id = 0 o sin parámetro)
		// Ejemplo: /notify/sse/order-notify?integration_id=123&event_types=order.created
		notifyGroup.GET("/sse/order-notify", r.sseHandler.HandleSSE)
	}
}
