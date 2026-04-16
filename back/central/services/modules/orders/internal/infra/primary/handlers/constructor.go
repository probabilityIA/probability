package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Handlers contiene todos los handlers del módulo orders
type Handlers struct {
	orderCRUD                ports.IOrderUseCase
	createUC                 ports.IOrderCreateUseCase
	requestConfirmationUC    ports.IRequestConfirmationUseCase
	sendGuideNotificationUC  ports.ISendGuideNotificationUseCase
	statusUC                 ports.IOrderStatusUseCase
	logger                   log.ILogger
}

// resolveBusinessID obtiene el business_id efectivo.
// Para usuarios normales usa el del JWT.
// Para super admins (business_id=0 en JWT) lee el query param ?business_id=X.
func (h *Handlers) resolveBusinessID(c *gin.Context) (uint, bool) {
	if id, ok := c.Get("business_id"); ok {
		if bID, ok := id.(uint); ok && bID > 0 {
			return bID, true
		}
	}
	// Super admin: leer de query param
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}

// New crea una nueva instancia de Handlers
func New(
	orderCRUD ports.IOrderUseCase,
	createUC ports.IOrderCreateUseCase,
	requestConfirmationUC ports.IRequestConfirmationUseCase,
	sendGuideNotificationUC ports.ISendGuideNotificationUseCase,
	statusUC ports.IOrderStatusUseCase,
	logger log.ILogger,
) *Handlers {
	return &Handlers{
		orderCRUD:                orderCRUD,
		createUC:                 createUC,
		requestConfirmationUC:    requestConfirmationUC,
		sendGuideNotificationUC:  sendGuideNotificationUC,
		statusUC:                 statusUC,
		logger:                   logger,
	}
}
