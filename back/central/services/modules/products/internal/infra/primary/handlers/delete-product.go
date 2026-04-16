package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// DeleteProduct godoc
// @Summary      Eliminar producto
// @Description  Elimina (soft delete) un producto del sistema
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id           path      string  true   "ID del producto (hash alfanumérico)"
// @Param        business_id  query     int     false  "ID del negocio (requerido para super admin)"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [delete]
func (h *Handlers) DeleteProduct(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto inválido",
			"error":   "El ID es requerido",
		})
		return
	}

	// Llamar al caso de uso (valida que el producto pertenezca al negocio)
	err := h.uc.DeleteProduct(c.Request.Context(), businessID, id)
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
			"message": "Error al eliminar producto",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Producto eliminado exitosamente",
	})
}
