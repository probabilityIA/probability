package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/mappers"
)

// ListPaymentMethods godoc
// @Summary      Listar métodos de pago
// @Description  Obtiene una lista paginada de métodos de pago
// @Tags         Payment Methods
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Número de página"      default(1)
// @Param        pageSize  query     int     false  "Tamaño de página"      default(10)
// @Param        category  query     string  false  "Filtrar por categoría"
// @Param        is_active query     bool    false  "Filtrar por estado activo"
// @Param        search    query     string  false  "Buscar por nombre o código"
// @Success      200       {object}  response.PaymentMethodsList
// @Failure      400       {object}  response.Error
// @Router       /payments/methods [get]
func (h *PaymentHandlers) ListPaymentMethods(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"], _ = strconv.ParseBool(isActive)
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	domainResponse, err := h.uc.ListPaymentMethods(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir a response HTTP
	httpResponse := mappers.DomainToPaymentMethodsListResponse(domainResponse)

	c.JSON(http.StatusOK, httpResponse)
}
