package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/request"
)

// ChangeStatus maneja la petición PUT /orders/:id/status
func (h *Handlers) ChangeStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	var req request.ChangeStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	// Extraer usuario del JWT context
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}
	var userName string
	if name, exists := c.Get("user_name"); exists {
		if n, ok := name.(string); ok {
			userName = n
		}
	}

	// Mapear a DTO de dominio
	domainReq := &dtos.ChangeStatusRequest{
		Status:   req.Status,
		Metadata: req.Metadata,
		UserID:   userID,
		UserName: userName,
	}

	// Ejecutar caso de uso
	result, err := h.statusUC.ChangeStatus(c.Request.Context(), id, domainReq)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidStatusTransition):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrOrderInTerminalState):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, mappers.OrderToResponse(result))
}
