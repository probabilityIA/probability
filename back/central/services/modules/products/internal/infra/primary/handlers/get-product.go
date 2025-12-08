package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// GetProductByID godoc
// @Summary      Obtener producto por ID
// @Description  Obtiene un producto específico por su ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del producto (hash alfanumérico)"
// @Security     BearerAuth
// @Success      200  {object}  domain.ProductResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func (h *Handlers) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto inválido",
			"error":   "El ID es requerido",
		})
		return
	}

	// Llamar al caso de uso
	product, err := h.uc.GetProductByID(c.Request.Context(), id)
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
			"message": "Error al obtener producto",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Producto obtenido exitosamente",
		"data":    product,
	})
}
