package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaymentMethod godoc
// @Summary      Obtener método de pago
// @Description  Obtiene un método de pago por su ID
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del método de pago"
// @Success      200  {object}  domain.PaymentMethodResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /payments/methods/{id} [get]
func (h *PaymentHandlers) GetPaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	response, err := h.uc.GetPaymentMethodByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
