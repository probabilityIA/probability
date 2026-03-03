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
// @Param        business_id  query     int                         false  "ID del negocio (requerido para super admin)"
// @Param        product      body      domain.CreateProductRequest true   "Datos del producto"
// @Security     BearerAuth
// @Success      201  {object}  domain.ProductResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products [post]
func (h *Handlers) CreateProduct(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

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

	// El business_id siempre viene del JWT (o query param para super admin), nunca del body
	req.BusinessID = businessID

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
