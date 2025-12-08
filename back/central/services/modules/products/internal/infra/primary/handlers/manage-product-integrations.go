package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// AddProductIntegration godoc
// @Summary      Asociar producto con integración
// @Description  Asocia un producto con una integración dentro del mismo negocio
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      string                                true   "ID del producto (hash alfanumérico)"
// @Param        body body      domain.AddProductIntegrationRequest  true   "Datos de la integración"
// @Security     BearerAuth
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id}/integrations [post]
func (h *Handlers) AddProductIntegration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto inválido",
			"error":   "El ID es requerido",
		})
		return
	}

	var req domain.AddProductIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	integration, err := h.uc.AddProductIntegration(c.Request.Context(), id, &req)
	if err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Producto no encontrado",
				"error":   err.Error(),
			})
			return
		}

		// Check for conflict errors (duplicate association, business mismatch)
		switch err.Error() {
		case "product is already associated with this integration":
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "El producto ya está asociado con esta integración",
				"error":   err.Error(),
			})
			return
		case "integration does not belong to the same business as the product":
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "La integración no pertenece al mismo negocio que el producto",
				"error":   err.Error(),
			})
			return
		case "integration not found":
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Integración no encontrada",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al asociar producto con integración",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Integración asociada exitosamente",
		"data":    integration,
	})
}

// RemoveProductIntegration godoc
// @Summary      Remover integración de producto
// @Description  Remueve la asociación entre un producto y una integración
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id              path      string  true  "ID del producto (hash alfanumérico)"
// @Param        integration_id  path      int  true  "ID de la integración"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id}/integrations/{integration_id} [delete]
func (h *Handlers) RemoveProductIntegration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto inválido",
			"error":   "El ID es requerido",
		})
		return
	}

	integrationIDStr := c.Param("integration_id")
	integrationID, err := strconv.ParseUint(integrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de integración inválido",
			"error":   "El ID debe ser un número válido",
		})
		return
	}

	err = h.uc.RemoveProductIntegration(c.Request.Context(), id, uint(integrationID))
	if err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Producto no encontrado",
				"error":   err.Error(),
			})
			return
		}

		if err.Error() == "product integration association not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Asociación no encontrada",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al remover integración",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Integración removida exitosamente",
	})
}

// GetProductIntegrations godoc
// @Summary      Listar integraciones de producto
// @Description  Obtiene todas las integraciones asociadas a un producto
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del producto (hash alfanumérico)"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id}/integrations [get]
func (h *Handlers) GetProductIntegrations(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto inválido",
			"error":   "El ID es requerido",
		})
		return
	}

	integrations, err := h.uc.GetProductIntegrations(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Producto no encontrado",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener integraciones",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Integraciones obtenidas exitosamente",
		"data":    integrations,
		"total":   len(integrations),
	})
}
