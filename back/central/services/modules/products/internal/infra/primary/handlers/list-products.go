package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ListProducts godoc
// @Summary      Listar productos
// @Description  Obtiene una lista paginada de productos con filtros opcionales y búsquedas optimizadas
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page            query    int     false  "Número de página (default: 1, min: 1)"
// @Param        page_size       query    int     false  "Tamaño de página (default: 10, min: 1, max: 100)"
// @Param        business_id     query    int     false  "Filtrar por ID de negocio"
// @Param        integration_id  query    int     false  "Filtrar por ID de integración (productos de negocios con esta integración)"
// @Param        integration_type query    string  false  "Filtrar por tipo de integración (ej: shopify, whatsapp)"
// @Param        sku             query    string  false  "Filtrar por SKU (búsqueda parcial, case-insensitive)"
// @Param        skus            query    string  false  "Filtrar por múltiples SKUs separados por coma (búsqueda exacta)"
// @Param        name            query    string  false  "Filtrar por nombre (búsqueda parcial, case-insensitive)"
// @Param        external_id     query    string  false  "Filtrar por ID externo (búsqueda exacta)"
// @Param        external_ids    query    string  false  "Filtrar por múltiples IDs externos separados por coma"
// @Param        start_date      query    string  false  "Fecha de inicio (YYYY-MM-DD o YYYY-MM-DD HH:MM:SS)"
// @Param        end_date        query    string  false  "Fecha de fin (YYYY-MM-DD o YYYY-MM-DD HH:MM:SS)"
// @Param        created_after   query    string  false  "Productos creados después de esta fecha (YYYY-MM-DD)"
// @Param        created_before  query    string  false  "Productos creados antes de esta fecha (YYYY-MM-DD)"
// @Param        updated_after   query    string  false  "Productos actualizados después de esta fecha (YYYY-MM-DD)"
// @Param        updated_before  query    string  false  "Productos actualizados antes de esta fecha (YYYY-MM-DD)"
// @Param        sort_by        query    string  false  "Campo para ordenar (id, sku, name, created_at, updated_at, business_id) (default: created_at)"
// @Param        sort_order      query    string  false  "Orden (asc, desc) (default: desc)"
// @Security     BearerAuth
// @Success      200  {object}  domain.ProductsListResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products [get]
func (h *Handlers) ListProducts(c *gin.Context) {
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

	// Filtro por business_id
	if businessID := c.Query("business_id"); businessID != "" {
		if id, err := strconv.ParseUint(businessID, 10, 32); err == nil && id > 0 {
			filters["business_id"] = uint(id)
		}
	}

	// Filtro por integration_id (a través de JOIN con Business -> Integrations)
	if integrationID := c.Query("integration_id"); integrationID != "" {
		if id, err := strconv.ParseUint(integrationID, 10, 32); err == nil && id > 0 {
			filters["integration_id"] = uint(id)
		}
	}

	// Filtro por integration_type
	if integrationType := c.Query("integration_type"); integrationType != "" {
		filters["integration_type"] = strings.TrimSpace(integrationType)
	}

	// Filtro por SKU (búsqueda parcial)
	if sku := c.Query("sku"); sku != "" {
		filters["sku"] = strings.TrimSpace(sku)
	}

	// Filtro por múltiples SKUs (búsqueda exacta)
	if skus := c.Query("skus"); skus != "" {
		skuList := strings.Split(skus, ",")
		trimmedSKUs := make([]string, 0, len(skuList))
		for _, s := range skuList {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				trimmedSKUs = append(trimmedSKUs, trimmed)
			}
		}
		if len(trimmedSKUs) > 0 {
			filters["skus"] = trimmedSKUs
		}
	}

	// Filtro por nombre (búsqueda parcial)
	if name := c.Query("name"); name != "" {
		filters["name"] = strings.TrimSpace(name)
	}

	// Filtro por external_id (búsqueda exacta)
	if externalID := c.Query("external_id"); externalID != "" {
		filters["external_id"] = strings.TrimSpace(externalID)
	}

	// Filtro por múltiples external_ids
	if externalIDs := c.Query("external_ids"); externalIDs != "" {
		idList := strings.Split(externalIDs, ",")
		trimmedIDs := make([]string, 0, len(idList))
		for _, id := range idList {
			trimmed := strings.TrimSpace(id)
			if trimmed != "" {
				trimmedIDs = append(trimmedIDs, trimmed)
			}
		}
		if len(trimmedIDs) > 0 {
			filters["external_ids"] = trimmedIDs
		}
	}

	// Filtros de fecha (compatibilidad con formato anterior)
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = strings.TrimSpace(startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = strings.TrimSpace(endDate)
	}

	// Filtros de fecha mejorados
	if createdAfter := c.Query("created_after"); createdAfter != "" {
		filters["created_after"] = strings.TrimSpace(createdAfter)
	}

	if createdBefore := c.Query("created_before"); createdBefore != "" {
		filters["created_before"] = strings.TrimSpace(createdBefore)
	}

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
			"id":          true,
			"sku":         true,
			"name":        true,
			"created_at":  true,
			"updated_at":  true,
			"business_id": true,
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
	response, err := h.uc.ListProducts(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener productos",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Productos obtenidos exitosamente",
		"data":        response.Data,
		"total":       response.Total,
		"page":        response.Page,
		"page_size":   response.PageSize,
		"total_pages": response.TotalPages,
	})
}

