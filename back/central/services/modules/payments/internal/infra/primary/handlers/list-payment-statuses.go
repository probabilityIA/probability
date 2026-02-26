package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/response"
)

// ListPaymentStatuses godoc
// @Summary      Listar estados de pago de Probability
// @Description  Obtiene una lista de todos los estados de pago de Probability. Opcionalmente puede filtrar por estado activo/inactivo.
// @Tags         Payment Statuses
// @Accept       json
// @Produce      json
// @Param        is_active  query     bool    false  "Filtrar por estado activo/inactivo (true=activos, false=inactivos, omitir=todos)"
// @Success      200        {object}  response.PaymentStatusListResponse
// @Failure      500        {object}  map[string]string
// @Router       /payment-statuses [get]
func (h *PaymentHandlers) ListPaymentStatuses(c *gin.Context) {
	var isActive *bool

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActiveValue, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &isActiveValue
		}
	}

	result, err := h.uc.ListPaymentStatuses(c.Request.Context(), isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener estados de pago",
			"error":   err.Error(),
		})
		return
	}

	data := make([]response.PaymentStatusResponse, len(result))
	for i, s := range result {
		data[i] = response.PaymentStatusResponse{
			ID:          s.ID,
			Code:        s.Code,
			Name:        s.Name,
			Description: s.Description,
			Category:    s.Category,
			Color:       s.Color,
			IsActive:    s.IsActive,
		}
	}

	c.JSON(http.StatusOK, response.PaymentStatusListResponse{
		Success: true,
		Message: "Estados de pago obtenidos exitosamente",
		Data:    data,
	})
}
