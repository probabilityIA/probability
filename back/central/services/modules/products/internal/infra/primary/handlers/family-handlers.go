package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (h *Handlers) CreateProductFamily(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	var req domain.CreateProductFamilyStandaloneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos de entrada inválidos", "error": err.Error()})
		return
	}

	req.BusinessID = businessID
	family, err := h.uc.CreateProductFamily(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al crear familia", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Familia creada exitosamente", "data": family})
}

func (h *Handlers) GetProductFamilyByID(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	familyID, ok := parseFamilyID(c)
	if !ok {
		return
	}

	family, err := h.uc.GetProductFamilyByID(c.Request.Context(), businessID, familyID)
	if err != nil {
		if errors.Is(err, domain.ErrProductFamilyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Familia no encontrada", "error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al obtener familia", "error": err.Error()})
		return
	}

	for i := range family.Variants {
		family.Variants[i].ImageURL = h.buildImageURL(family.Variants[i].ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Familia obtenida exitosamente", "data": family})
}

func (h *Handlers) ListProductFamilies(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	page := 1
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		parsed, err := strconv.Atoi(pageStr)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Parámetro 'page' inválido", "error": "invalid page parameter"})
			return
		}
		page = parsed
	}

	pageSize := 10
	if pageSizeStr := c.DefaultQuery("page_size", "10"); pageSizeStr != "" {
		parsed, err := strconv.Atoi(pageSizeStr)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Parámetro 'page_size' inválido", "error": "invalid page_size parameter"})
			return
		}
		if parsed > 100 {
			parsed = 100
		}
		pageSize = parsed
	}

	filters := make(map[string]interface{})
	for _, field := range []string{"name", "category", "brand", "status", "sort_by", "sort_order"} {
		if value := strings.TrimSpace(c.Query(field)); value != "" {
			filters[field] = strings.ToLower(value)
			if field == "name" || field == "category" || field == "brand" {
				filters[field] = value
			}
		}
	}

	response, err := h.uc.ListProductFamilies(c.Request.Context(), businessID, page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al listar familias", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Familias obtenidas exitosamente",
		"data":        response.Data,
		"total":       response.Total,
		"page":        response.Page,
		"page_size":   response.PageSize,
		"total_pages": response.TotalPages,
	})
}

func (h *Handlers) UpdateProductFamily(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	familyID, ok := parseFamilyID(c)
	if !ok {
		return
	}

	var req domain.UpdateProductFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos de entrada inválidos", "error": err.Error()})
		return
	}

	family, err := h.uc.UpdateProductFamily(c.Request.Context(), businessID, familyID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrProductFamilyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Familia no encontrada", "error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al actualizar familia", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Familia actualizada exitosamente", "data": family})
}

func (h *Handlers) ListProductFamilyVariants(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	familyID, ok := parseFamilyID(c)
	if !ok {
		return
	}

	if _, err := h.uc.GetProductFamilyByID(c.Request.Context(), businessID, familyID); err != nil {
		if errors.Is(err, domain.ErrProductFamilyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Familia no encontrada", "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al validar familia", "error": err.Error()})
		return
	}

	variants, err := h.uc.ListProductsByFamilyID(c.Request.Context(), businessID, familyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al listar variantes", "error": err.Error()})
		return
	}

	response := make([]domain.ProductResponse, len(variants))
	for i := range variants {
		product := variants[i]
		response[i] = domain.ProductResponse{
			ID:                product.ID,
			CreatedAt:         product.CreatedAt,
			UpdatedAt:         product.UpdatedAt,
			DeletedAt:         product.DeletedAt,
			BusinessID:        product.BusinessID,
			SKU:               product.SKU,
			ExternalID:        product.ExternalID,
			Barcode:           product.Barcode,
			FamilyID:          product.FamilyID,
			Name:              product.Name,
			Title:             product.Title,
			Description:       product.Description,
			ShortDescription:  product.ShortDescription,
			Slug:              product.Slug,
			VariantLabel:      product.VariantLabel,
			VariantAttributes: product.VariantAttributes,
			Price:             product.Price,
			CompareAtPrice:    product.CompareAtPrice,
			CostPrice:         product.CostPrice,
			Currency:          product.Currency,
			StockQuantity:     product.StockQuantity,
			TrackInventory:    product.TrackInventory,
			AllowBackorder:    product.AllowBackorder,
			LowStockThreshold: product.LowStockThreshold,
			ImageURL:          h.buildImageURL(product.ImageURL),
			Images:            product.Images,
			VideoURL:          product.VideoURL,
			Weight:            product.Weight,
			WeightUnit:        product.WeightUnit,
			Length:            product.Length,
			Width:             product.Width,
			Height:            product.Height,
			DimensionUnit:     product.DimensionUnit,
			Category:          product.Category,
			Tags:              product.Tags,
			Brand:             product.Brand,
			Status:            product.Status,
			IsActive:          product.IsActive,
			IsFeatured:        product.IsFeatured,
			Metadata:          product.Metadata,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Variantes obtenidas exitosamente",
		"data":    response,
		"total":   len(response),
	})
}

func (h *Handlers) DeleteProductFamily(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	familyID, ok := parseFamilyID(c)
	if !ok {
		return
	}

	if err := h.uc.DeleteProductFamily(c.Request.Context(), businessID, familyID); err != nil {
		switch {
		case errors.Is(err, domain.ErrProductFamilyNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Familia no encontrada", "error": err.Error()})
		case errors.Is(err, domain.ErrFamilyHasActiveVariants):
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "No se puede eliminar una familia con variantes activas", "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al eliminar familia", "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Familia eliminada exitosamente"})
}

func parseFamilyID(c *gin.Context) (uint, bool) {
	familyIDStr := c.Param("family_id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil || familyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de familia inválido", "error": "El ID debe ser un número válido"})
		return 0, false
	}

	return uint(familyID), true
}
