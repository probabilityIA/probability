package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOrders godoc
// @Summary      Listar órdenes
// @Description  Obtiene una lista paginada de órdenes con filtros opcionales
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        page              query    int     false  "Número de página (default: 1)"
// @Param        page_size         query    int     false  "Tamaño de página (default: 10, max: 100)"
// @Param        business_id       query    int     false  "Filtrar por ID de negocio"
// @Param        integration_id    query    int     false  "Filtrar por ID de integración"
// @Param        integration_type  query    string  false  "Filtrar por tipo de integración"
// @Param        status            query    string  false  "Filtrar por estado"
// @Param        customer_email    query    string  false  "Filtrar por email del cliente"
// @Param        customer_phone    query    string  false  "Filtrar por teléfono del cliente"
// @Param        platform          query    string  false  "Filtrar por plataforma"
// @Param        is_paid           query    bool    false  "Filtrar por estado de pago"
// @Param        warehouse_id      query    int     false  "Filtrar por ID de almacén"
// @Param        driver_id         query    int     false  "Filtrar por ID de conductor"
// @Param        start_date        query    string  false  "Fecha de inicio (YYYY-MM-DD)"
// @Param        end_date          query    string  false  "Fecha de fin (YYYY-MM-DD)"
// @Param        sort_by           query    string  false  "Campo para ordenar (default: created_at)"
// @Param        sort_order        query    string  false  "Orden (asc, desc) (default: desc)"
// @Security     BearerAuth
// @Success      200  {object}  domain.OrdersListResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders [get]
func (h *Handlers) ListOrders(c *gin.Context) {
	// Obtener parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Construir filtros
	filters := make(map[string]interface{})

	if businessID := c.Query("business_id"); businessID != "" {
		if id, err := strconv.ParseUint(businessID, 10, 32); err == nil {
			filters["business_id"] = uint(id)
		}
	}

	if integrationID := c.Query("integration_id"); integrationID != "" {
		if id, err := strconv.ParseUint(integrationID, 10, 32); err == nil {
			filters["integration_id"] = uint(id)
		}
	}

	if integrationType := c.Query("integration_type"); integrationType != "" {
		filters["integration_type"] = integrationType
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if customerEmail := c.Query("customer_email"); customerEmail != "" {
		filters["customer_email"] = customerEmail
	}

	if customerPhone := c.Query("customer_phone"); customerPhone != "" {
		filters["customer_phone"] = customerPhone
	}

	if platform := c.Query("platform"); platform != "" {
		filters["platform"] = platform
	}

	if isPaid := c.Query("is_paid"); isPaid != "" {
		if paid, err := strconv.ParseBool(isPaid); err == nil {
			filters["is_paid"] = paid
		}
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		if id, err := strconv.ParseUint(warehouseID, 10, 32); err == nil {
			filters["warehouse_id"] = uint(id)
		}
	}

	if driverID := c.Query("driver_id"); driverID != "" {
		if id, err := strconv.ParseUint(driverID, 10, 32); err == nil {
			filters["driver_id"] = uint(id)
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters["sort_by"] = sortBy
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filters["sort_order"] = sortOrder
	}

	// Llamar al caso de uso
	response, err := h.uc.ListOrders(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener órdenes",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Órdenes obtenidas exitosamente",
		"data":        response.Data,
		"total":       response.Total,
		"page":        response.Page,
		"page_size":   response.PageSize,
		"total_pages": response.TotalPages,
	})
}
