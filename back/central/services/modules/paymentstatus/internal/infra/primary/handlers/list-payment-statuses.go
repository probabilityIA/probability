package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/infra/primary/handlers/response"
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
func (h *PaymentStatusHandlers) ListPaymentStatuses(c *gin.Context) {
	var isActive *bool

	// Filtro opcional por is_active
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActiveValue, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &isActiveValue
		}
	}

	// ✅ Obtener DTOs de dominio desde usecase
	result, err := h.uc.ListPaymentStatuses(c.Request.Context(), isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener estados de pago",
			"error":   err.Error(),
		})
		return
	}

	// ✅ Convertir a response HTTP usando mapper
	c.JSON(http.StatusOK, response.PaymentStatusListResponse{
		Success: true,
		Message: "Estados de pago obtenidos exitosamente",
		Data:    mappers.ToPaymentStatusListResponse(result),
	})
}
