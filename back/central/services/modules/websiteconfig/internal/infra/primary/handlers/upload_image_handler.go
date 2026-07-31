package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) UploadImage(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "business_id es requerido",
		})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Imagen es requerida",
			"error":   err.Error(),
		})
		return
	}

	folder := fmt.Sprintf("website/%d", businessID)
	relativePath, err := h.s3.UploadImage(c.Request.Context(), file, folder)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).Msg("error uploading website image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al subir imagen",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Imagen subida exitosamente",
		"image_url": h.buildImageURL(relativePath),
	})
}

func (h *Handlers) DeleteImage(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "business_id es requerido",
		})
		return
	}

	imageURL := c.Query("image_url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "image_url es requerido",
		})
		return
	}

	key := imageURL
	if base := h.env.Get("URL_BASE_DOMAIN_S3"); base != "" && strings.HasPrefix(key, base) {
		key = strings.TrimPrefix(key, base)
	} else if parsed, err := url.Parse(imageURL); err == nil && parsed.Host != "" {
		key = strings.TrimPrefix(parsed.Path, "/")
	}

	expectedPrefix := fmt.Sprintf("website/%d/", businessID)
	if !strings.HasPrefix(key, expectedPrefix) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "la imagen no pertenece al negocio",
		})
		return
	}

	if err := h.s3.DeleteImage(c.Request.Context(), key); err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Str("key", key).Uint("business_id", businessID).Msg("error deleting website image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar imagen",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Imagen eliminada",
	})
}

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
