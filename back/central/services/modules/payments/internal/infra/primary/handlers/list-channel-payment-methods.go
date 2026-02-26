package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/response"
)

func (h *PaymentHandlers) ListChannelPaymentMethods(c *gin.Context) {
	var integrationType *string
	if it := c.Query("integration_type"); it != "" {
		integrationType = &it
	}

	var isActive *bool
	if ia := c.Query("is_active"); ia != "" {
		val := ia == "true"
		isActive = &val
	}

	methods, err := h.uc.ListChannelPaymentMethods(c.Request.Context(), integrationType, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	data := make([]response.ChannelPaymentMethodResponse, len(methods))
	for i, m := range methods {
		data[i] = response.ChannelPaymentMethodResponse{
			ID:              m.ID,
			IntegrationType: m.IntegrationType,
			Code:            m.Code,
			Name:            m.Name,
			Description:     m.Description,
			IsActive:        m.IsActive,
			DisplayOrder:    m.DisplayOrder,
		}
	}

	c.JSON(http.StatusOK, response.ChannelPaymentMethodListResponse{
		Success: true,
		Message: "Channel payment methods retrieved successfully",
		Data:    data,
	})
}
