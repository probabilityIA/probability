package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List godoc
// @Summary      Listar mapeos de estado de orden
// @Description  Obtiene una lista paginada de mapeos de estado de orden con filtros opcionales. Filtra por ID del tipo de integración.
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        page                 query     int     false  "Número de página (default: 1)"
// @Param        page_size            query     int     false  "Tamaño de página (default: 10, max: 100)"
// @Param        integration_type_id  query     int     false  "Filtrar por ID del tipo de integración (1=shopify, 2=whatsapp, 3=mercado_libre, 4=woocommerce)"
// @Param        is_active            query     bool    false  "Filtrar por estado activo/inactivo"
// @Success      200                  {object}  response.OrderStatusMappingsListResponse
// @Failure      400                  {object}  map[string]string
// @Failure      500                  {object}  map[string]string
// @Router       /order-status-mappings [get]
func (h *OrderStatusMappingHandlers) List(c *gin.Context) {
	// Obtener y validar parámetros de paginación
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parámetro 'page' inválido. Debe ser un número entero mayor a 0",
			"error":   "invalid page parameter",
		})
		return
	}

	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parámetro 'page_size' inválido. Debe ser un número entero entre 1 y 100",
			"error":   "invalid page_size parameter",
		})
		return
	}

	// Limitar el tamaño máximo de página
	if pageSize > 100 {
		pageSize = 100
	}

	filters := make(map[string]interface{})
	filters["page"] = page
	filters["page_size"] = pageSize

	// Filtrar por integration_type_id
	if integrationTypeIDStr := c.Query("integration_type_id"); integrationTypeIDStr != "" {
		if integrationTypeID, err := strconv.ParseUint(integrationTypeIDStr, 10, 32); err == nil {
			filters["integration_type_id"] = uint(integrationTypeID)
		}
	}

	// Filtro opcional por is_active
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filters["is_active"] = isActive
		}
	}

	result, total, err := h.uc.ListOrderStatusMappings(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toListResponse(result, total, page, pageSize, h.getImageURLBase()))
}
