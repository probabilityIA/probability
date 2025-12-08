package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// CreateProduct godoc
// @Summary      Crear producto
// @Description  Crea un nuevo producto en el sistema
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      domain.CreateProductRequest  true  "Datos del producto"
// @Security     BearerAuth
// @Success      201  {object}  domain.ProductResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products [post]
func (h *Handlers) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest

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
	product, err := h.uc.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		// Verificar si es un error de duplicado
		if err == domain.ErrProductAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Producto ya existe",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear producto",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Producto creado exitosamente",
		"data":    product,
	})
}

