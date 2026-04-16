package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// UploadProductImage godoc
// @Summary      Subir imagen de producto
// @Description  Sube una imagen para un producto existente a S3
// @Tags         Products
// @Accept       mpfd
// @Produce      json
// @Param        id           path      string  true   "ID del producto"
// @Param        business_id  query     int     false  "ID del negocio (requerido para super admin)"
// @Param        image        formData  file    true   "Imagen del producto"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id}/image [post]
func (h *Handlers) UploadProductImage(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto es requerido",
		})
		return
	}

	// Verificar que el producto existe y pertenece al negocio
	_, err := h.uc.GetProductByID(c.Request.Context(), businessID, productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Producto no encontrado",
			"error":   err.Error(),
		})
		return
	}

	// Obtener archivo del form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Imagen es requerida",
			"error":   err.Error(),
		})
		return
	}

	// Subir a S3 con folder products/{business_id}
	folder := fmt.Sprintf("products/%d", businessID)
	relativePath, err := h.s3.UploadImage(c.Request.Context(), file, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al subir imagen",
			"error":   err.Error(),
		})
		return
	}

	// Actualizar image_url del producto con el path relativo
	imageURL := relativePath
	updateReq := &domain.UpdateProductRequest{ImageURL: &imageURL}
	updatedProduct, err := h.uc.UpdateProduct(c.Request.Context(), businessID, productID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al actualizar producto con la imagen",
			"error":   err.Error(),
		})
		return
	}

	// Construir URL completa para la respuesta
	fullURL := h.buildImageURL(relativePath)
	updatedProduct.ImageURL = fullURL

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Imagen subida exitosamente",
		"image_url": fullURL,
		"data":      updatedProduct,
	})
}

// buildImageURL construye la URL completa desde el path relativo
// Si ya es una URL completa (http/https), la retorna sin modificar
func (h *Handlers) buildImageURL(relativePath string) string {
	if relativePath == "" {
		return relativePath
	}
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}
	base := h.env.Get("URL_BASE_DOMAIN_S3")
	if base == "" {
		return relativePath
	}
	return base + relativePath
}
