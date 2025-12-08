package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// UpdateProduct godoc
// @Summary      Actualizar producto
// @Description  Actualiza un producto existente
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "ID del producto (hash alfanumérico)"
// @Param        product  body      domain.UpdateProductRequest true  "Datos a actualizar"
// @Security     BearerAuth
// @Success      200  {object}  domain.ProductResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [put]
func (h *Handlers) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto inválido",
			"error":   "El ID es requerido",
		})
		return
	}

	var req domain.UpdateProductRequest

	// Validar el request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Llamar al caso de uso
	product, err := h.uc.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Producto no encontrado",
				"error":   err.Error(),
			})
			return
		}

		if err == domain.ErrProductAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Producto con este SKU ya existe",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al actualizar producto",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Producto actualizado exitosamente",
		"data":    product,
	})
}
