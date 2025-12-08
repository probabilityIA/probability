package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ListShipments godoc
// @Summary      Listar envíos
// @Description  Obtiene una lista paginada de envíos con filtros opcionales y búsquedas optimizadas
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        page              query    int     false  "Número de página (default: 1, min: 1)"
// @Param        page_size         query    int     false  "Tamaño de página (default: 10, min: 1, max: 100)"
// @Param        order_id          query    string  false  "Filtrar por ID de orden"
// @Param        order_ids         query    string  false  "Filtrar por múltiples IDs de orden separados por coma"
// @Param        tracking_number   query    string  false  "Filtrar por número de tracking (búsqueda parcial)"
// @Param        tracking_numbers  query    string  false  "Filtrar por múltiples números de tracking separados por coma"
// @Param        carrier           query    string  false  "Filtrar por transportista (búsqueda parcial)"
// @Param        carrier_code      query    string  false  "Filtrar por código de transportista"
// @Param        status            query    string  false  "Filtrar por estado (pending, in_transit, delivered, failed)"
// @Param        statuses          query    string  false  "Filtrar por múltiples estados separados por coma"
// @Param        guide_id          query    string  false  "Filtrar por ID de guía"
// @Param        warehouse_id      query    int     false  "Filtrar por ID de almacén"
// @Param        driver_id         query    int     false  "Filtrar por ID de conductor"
// @Param        is_last_mile      query    bool    false  "Filtrar por última milla"
// @Param        business_id       query    int     false  "Filtrar por ID de negocio (a través de orden)"
// @Param        integration_id   query    int     false  "Filtrar por ID de integración (a través de orden)"
// @Param        integration_type query    string  false  "Filtrar por tipo de integración (ej: shopify, whatsapp)"
// @Param        shipped_after    query    string  false  "Envíos enviados después de esta fecha (YYYY-MM-DD)"
// @Param        shipped_before   query    string  false  "Envíos enviados antes de esta fecha (YYYY-MM-DD)"
// @Param        delivered_after  query    string  false  "Envíos entregados después de esta fecha (YYYY-MM-DD)"
// @Param        delivered_before query    string  false  "Envíos entregados antes de esta fecha (YYYY-MM-DD)"
// @Param        start_date       query    string  false  "Fecha de inicio creación (YYYY-MM-DD)"
// @Param        end_date         query    string  false  "Fecha de fin creación (YYYY-MM-DD)"
// @Param        created_after    query    string  false  "Envíos creados después de esta fecha (YYYY-MM-DD)"
// @Param        created_before   query    string  false  "Envíos creados antes de esta fecha (YYYY-MM-DD)"
// @Param        updated_after    query    string  false  "Envíos actualizados después de esta fecha (YYYY-MM-DD)"
// @Param        updated_before  query    string  false  "Envíos actualizados antes de esta fecha (YYYY-MM-DD)"
// @Param        sort_by         query    string  false  "Campo para ordenar (id, order_id, tracking_number, status, carrier, shipped_at, delivered_at, created_at, updated_at, warehouse_id, driver_id) (default: created_at)"
// @Param        sort_order      query    string  false  "Orden (asc, desc) (default: desc)"
// @Security     BearerAuth
// @Success      200  {object}  domain.ShipmentsListResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments [get]
func (h *Handlers) ListShipments(c *gin.Context) {
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

	// Construir filtros
	filters := make(map[string]interface{})

	// Filtro por order_id
	if orderID := c.Query("order_id"); orderID != "" {
		filters["order_id"] = strings.TrimSpace(orderID)
	}

	// Filtro por múltiples order_ids
	if orderIDs := c.Query("order_ids"); orderIDs != "" {
		idList := strings.Split(orderIDs, ",")
		trimmedIDs := make([]string, 0, len(idList))
		for _, id := range idList {
			trimmed := strings.TrimSpace(id)
			if trimmed != "" {
				trimmedIDs = append(trimmedIDs, trimmed)
			}
		}
		if len(trimmedIDs) > 0 {
			filters["order_ids"] = trimmedIDs
		}
	}

	// Filtro por tracking_number
	if trackingNumber := c.Query("tracking_number"); trackingNumber != "" {
		filters["tracking_number"] = strings.TrimSpace(trackingNumber)
	}

	// Filtro por múltiples tracking_numbers
	if trackingNumbers := c.Query("tracking_numbers"); trackingNumbers != "" {
		numList := strings.Split(trackingNumbers, ",")
		trimmedNums := make([]string, 0, len(numList))
		for _, num := range numList {
			trimmed := strings.TrimSpace(num)
			if trimmed != "" {
				trimmedNums = append(trimmedNums, trimmed)
			}
		}
		if len(trimmedNums) > 0 {
			filters["tracking_numbers"] = trimmedNums
		}
	}

	// Filtro por carrier
	if carrier := c.Query("carrier"); carrier != "" {
		filters["carrier"] = strings.TrimSpace(carrier)
	}

	// Filtro por carrier_code
	if carrierCode := c.Query("carrier_code"); carrierCode != "" {
		filters["carrier_code"] = strings.TrimSpace(carrierCode)
	}

	// Filtro por status
	if status := c.Query("status"); status != "" {
		filters["status"] = strings.TrimSpace(status)
	}

	// Filtro por múltiples statuses
	if statuses := c.Query("statuses"); statuses != "" {
		statusList := strings.Split(statuses, ",")
		trimmedStatuses := make([]string, 0, len(statusList))
		for _, s := range statusList {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				trimmedStatuses = append(trimmedStatuses, trimmed)
			}
		}
		if len(trimmedStatuses) > 0 {
			filters["statuses"] = trimmedStatuses
		}
	}

	// Filtro por guide_id
	if guideID := c.Query("guide_id"); guideID != "" {
		filters["guide_id"] = strings.TrimSpace(guideID)
	}

	// Filtro por warehouse_id
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		if id, err := strconv.ParseUint(warehouseID, 10, 32); err == nil && id > 0 {
			filters["warehouse_id"] = uint(id)
		}
	}

	// Filtro por driver_id
	if driverID := c.Query("driver_id"); driverID != "" {
		if id, err := strconv.ParseUint(driverID, 10, 32); err == nil && id > 0 {
			filters["driver_id"] = uint(id)
		}
	}

	// Filtro por is_last_mile
	if isLastMile := c.Query("is_last_mile"); isLastMile != "" {
		if lastMile, err := strconv.ParseBool(isLastMile); err == nil {
			filters["is_last_mile"] = lastMile
		}
	}

	// Filtro por business_id
	if businessID := c.Query("business_id"); businessID != "" {
		if id, err := strconv.ParseUint(businessID, 10, 32); err == nil && id > 0 {
			filters["business_id"] = uint(id)
		}
	}

	// Filtro por integration_id
	if integrationID := c.Query("integration_id"); integrationID != "" {
		if id, err := strconv.ParseUint(integrationID, 10, 32); err == nil && id > 0 {
			filters["integration_id"] = uint(id)
		}
	}

	// Filtro por integration_type
	if integrationType := c.Query("integration_type"); integrationType != "" {
		filters["integration_type"] = strings.TrimSpace(integrationType)
	}

	// Filtros de fecha - shipped_at
	if shippedAfter := c.Query("shipped_after"); shippedAfter != "" {
		filters["shipped_after"] = strings.TrimSpace(shippedAfter)
	}

	if shippedBefore := c.Query("shipped_before"); shippedBefore != "" {
		filters["shipped_before"] = strings.TrimSpace(shippedBefore)
	}

	// Filtros de fecha - delivered_at
	if deliveredAfter := c.Query("delivered_after"); deliveredAfter != "" {
		filters["delivered_after"] = strings.TrimSpace(deliveredAfter)
	}

	if deliveredBefore := c.Query("delivered_before"); deliveredBefore != "" {
		filters["delivered_before"] = strings.TrimSpace(deliveredBefore)
	}

	// Filtros de fecha - created_at (compatibilidad)
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = strings.TrimSpace(startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = strings.TrimSpace(endDate)
	}

	if createdAfter := c.Query("created_after"); createdAfter != "" {
		filters["created_after"] = strings.TrimSpace(createdAfter)
	}

	if createdBefore := c.Query("created_before"); createdBefore != "" {
		filters["created_before"] = strings.TrimSpace(createdBefore)
	}

	// Filtros de fecha - updated_at
	if updatedAfter := c.Query("updated_after"); updatedAfter != "" {
		filters["updated_after"] = strings.TrimSpace(updatedAfter)
	}

	if updatedBefore := c.Query("updated_before"); updatedBefore != "" {
		filters["updated_before"] = strings.TrimSpace(updatedBefore)
	}

	// Ordenamiento
	if sortBy := c.Query("sort_by"); sortBy != "" {
		// Validar campos permitidos para ordenar
		allowedSortFields := map[string]bool{
			"id":              true,
			"order_id":        true,
			"tracking_number": true,
			"status":          true,
			"carrier":         true,
			"shipped_at":      true,
			"delivered_at":    true,
			"created_at":      true,
			"updated_at":      true,
			"warehouse_id":    true,
			"driver_id":       true,
		}
		if allowedSortFields[strings.ToLower(sortBy)] {
			filters["sort_by"] = strings.ToLower(sortBy)
		}
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		sortOrderLower := strings.ToLower(sortOrder)
		if sortOrderLower == "asc" || sortOrderLower == "desc" {
			filters["sort_order"] = sortOrderLower
		}
	}

	// Llamar al caso de uso
	response, err := h.uc.ListShipments(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener envíos",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Envíos obtenidos exitosamente",
		"data":        response.Data,
		"total":       response.Total,
		"page":        response.Page,
		"page_size":   response.PageSize,
		"total_pages": response.TotalPages,
	})
}

