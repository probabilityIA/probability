package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListPaymentMethods godoc
// @Summary      Listar métodos de pago
// @Description  Obtiene una lista paginada de métodos de pago con filtros opcionales
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        page       query    int     false  "Número de página" default(1)
// @Param        page_size  query    int     false  "Tamaño de página" default(10)
// @Param        category   query    string  false  "Filtrar por categoría (card, digital_wallet, bank_transfer, cash)"
// @Param        search     query    string  false  "Buscar por nombre o código"
// @Success      200  {object}  domain.PaymentMethodsListResponse
// @Failure      500  {object}  map[string]string
// @Router       /payments/methods [get]
func (h *PaymentHandlers) ListPaymentMethods(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	response, err := h.uc.ListPaymentMethods(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
