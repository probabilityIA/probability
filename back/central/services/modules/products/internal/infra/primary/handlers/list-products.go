package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) ListProducts(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

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

	if pageSize > 100 {
		pageSize = 100
	}

	filters := make(map[string]interface{})

	if integrationID := c.Query("integration_id"); integrationID != "" {
		if id, err := strconv.ParseUint(integrationID, 10, 32); err == nil && id > 0 {
			filters["integration_id"] = uint(id)
		}
	}

	if integrationType := c.Query("integration_type"); integrationType != "" {
		filters["integration_type"] = strings.TrimSpace(integrationType)
	}

	if sku := c.Query("sku"); sku != "" {
		filters["sku"] = strings.TrimSpace(sku)
	}

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

	if familyID := c.Query("family_id"); familyID != "" {
		if id, err := strconv.ParseUint(familyID, 10, 32); err == nil && id > 0 {
			filters["family_id"] = uint(id)
		}
	}

	if barcode := c.Query("barcode"); barcode != "" {
		filters["barcode"] = strings.TrimSpace(barcode)
	}

	if name := c.Query("name"); name != "" {
		filters["name"] = strings.TrimSpace(name)
	}

	if search := c.Query("search"); search != "" {
		filters["search"] = strings.TrimSpace(search)
	}

	if externalID := c.Query("external_id"); externalID != "" {
		filters["external_id"] = strings.TrimSpace(externalID)
	}

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

	if updatedAfter := c.Query("updated_after"); updatedAfter != "" {
		filters["updated_after"] = strings.TrimSpace(updatedAfter)
	}

	if updatedBefore := c.Query("updated_before"); updatedBefore != "" {
		filters["updated_before"] = strings.TrimSpace(updatedBefore)
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = strings.TrimSpace(status)
	}

	if category := c.Query("category"); category != "" {
		filters["category"] = strings.TrimSpace(category)
	}

	if brand := c.Query("brand"); brand != "" {
		filters["brand"] = strings.TrimSpace(brand)
	}

	if hasFamily := c.Query("has_family"); hasFamily != "" {
		if hasFamily == "true" {
			filters["has_family"] = true
		} else if hasFamily == "false" {
			filters["has_family"] = false
		}
	}

	if priceMin := c.Query("price_min"); priceMin != "" {
		if v, err := strconv.ParseFloat(priceMin, 64); err == nil {
			filters["price_min"] = v
		}
	}

	if priceMax := c.Query("price_max"); priceMax != "" {
		if v, err := strconv.ParseFloat(priceMax, 64); err == nil {
			filters["price_max"] = v
		}
	}

	if stockMin := c.Query("stock_min"); stockMin != "" {
		if v, err := strconv.Atoi(stockMin); err == nil {
			filters["stock_min"] = v
		}
	}

	if stockMax := c.Query("stock_max"); stockMax != "" {
		if v, err := strconv.Atoi(stockMax); err == nil {
			filters["stock_max"] = v
		}
	}

	if weightMin := c.Query("weight_min"); weightMin != "" {
		if v, err := strconv.ParseFloat(weightMin, 64); err == nil {
			filters["weight_min"] = v
		}
	}

	if weightMax := c.Query("weight_max"); weightMax != "" {
		if v, err := strconv.ParseFloat(weightMax, 64); err == nil {
			filters["weight_max"] = v
		}
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		allowedSortFields := map[string]bool{
			"id":         true,
			"sku":        true,
			"name":       true,
			"created_at": true,
			"updated_at": true,
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

	response, err := h.uc.ListProducts(c.Request.Context(), businessID, page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener productos",
			"error":   err.Error(),
		})
		return
	}

	for i := range response.Data {
		response.Data[i].ImageURL = h.buildImageURL(response.Data[i].ImageURL)
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
