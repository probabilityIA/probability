package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (h *Handlers) AddProductIntegration(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de producto inválido", "error": "El ID es requerido"})
		return
	}

	var req domain.AddProductIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos de entrada inválidos", "error": err.Error()})
		return
	}

	integration, err := h.uc.AddProductIntegration(c.Request.Context(), businessID, id, &req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Producto no encontrado", "error": err.Error()})
		case err.Error() == "product is already associated with this integration":
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "El producto ya está asociado con esta integración", "error": err.Error()})
		case err.Error() == "integration does not belong to the same business as the product":
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "La integración no pertenece al mismo negocio que el producto", "error": err.Error()})
		case err.Error() == "integration not found":
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Integración no encontrada", "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al asociar producto con integración", "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Integración asociada exitosamente", "data": integration})
}

func (h *Handlers) UpdateProductIntegration(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de producto inválido", "error": "El ID es requerido"})
		return
	}

	integrationIDStr := c.Param("integration_id")
	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de integración inválido", "error": "El ID debe ser un número válido"})
		return
	}

	var req domain.UpdateProductIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Datos de entrada inválidos", "error": err.Error()})
		return
	}

	integration, err := h.uc.UpdateProductIntegration(c.Request.Context(), businessID, id, uint(integrationID), &req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Producto no encontrado", "error": err.Error()})
		case errors.Is(err, domain.ErrProductIntegrationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Mapping de integración no encontrado", "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al actualizar mapping de integración", "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mapping actualizado exitosamente", "data": integration})
}

func (h *Handlers) RemoveProductIntegration(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de producto inválido", "error": "El ID es requerido"})
		return
	}

	integrationIDStr := c.Param("integration_id")
	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de integración inválido", "error": "El ID debe ser un número válido"})
		return
	}

	err = h.uc.RemoveProductIntegration(c.Request.Context(), businessID, id, uint(integrationID))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Producto no encontrado", "error": err.Error()})
		case err.Error() == "product integration association not found":
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Asociación no encontrada", "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al remover integración", "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Integración removida exitosamente"})
}

func (h *Handlers) GetProductIntegrations(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID de producto inválido", "error": "El ID es requerido"})
		return
	}

	integrations, err := h.uc.GetProductIntegrations(c.Request.Context(), businessID, id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Producto no encontrado", "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al obtener integraciones", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Integraciones obtenidas exitosamente", "data": integrations, "total": len(integrations)})
}

func (h *Handlers) LookupProductByExternalRef(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "integration_id es requerido", "error": "query param faltante"})
		return
	}
	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "integration_id inválido", "error": err.Error()})
		return
	}

	strPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	externalVariantID := strPtr(c.Query("external_variant_id"))
	externalSKU := strPtr(c.Query("external_sku"))
	externalProductID := strPtr(c.Query("external_product_id"))
	externalBarcode := strPtr(c.Query("external_barcode"))

	if externalVariantID == nil && externalSKU == nil && externalProductID == nil && externalBarcode == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Debe proveer al menos un identificador externo", "error": "sin filtros de búsqueda"})
		return
	}

	product, err := h.uc.LookupProductByExternalRef(c.Request.Context(), businessID, uint(integrationID), externalVariantID, externalSKU, externalProductID, externalBarcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al buscar producto", "error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Producto no encontrado para las referencias externas proporcionadas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": product})
}
