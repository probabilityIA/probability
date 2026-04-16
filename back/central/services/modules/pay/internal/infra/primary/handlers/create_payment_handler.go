package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/request"
)

// CreatePayment maneja POST /pay/transactions
func (h *handler) CreatePayment(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener business_id desde el JWT
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parsear body
	var req request.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Construir DTO
	dto := &dtos.CreatePaymentDTO{
		BusinessID:    businessID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		GatewayCode:   req.GatewayCode,
		PaymentMethod: req.PaymentMethod,
		Description:   req.Description,
		CallbackURL:   req.CallbackURL,
		Metadata:      req.Metadata,
	}

	// Ejecutar caso de uso
	tx, err := h.useCase.CreatePayment(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to create payment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mappers.ToResponse(tx))
}
