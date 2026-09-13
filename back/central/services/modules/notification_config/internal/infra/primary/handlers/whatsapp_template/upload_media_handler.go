package whatsapp_template

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxHeaderImageBytes = 5 * 1024 * 1024

var allowedHeaderImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}

func (h *handler) UploadMedia(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	if h.s3 == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "error": "el almacenamiento de archivos no esta disponible"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "adjunta la imagen en el campo file"})
		return
	}

	contentType := strings.ToLower(strings.TrimSpace(file.Header.Get("Content-Type")))
	if !allowedHeaderImageTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Meta solo acepta JPG o PNG en el encabezado"})
		return
	}

	if file.Size > maxHeaderImageBytes {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "la imagen supera los 5 MB que acepta Meta"})
		return
	}

	path, err := h.s3.UploadImage(c.Request.Context(), file, fmt.Sprintf("whatsapp-templates/%d", businessID))
	if err != nil {
		h.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error subiendo la imagen del encabezado")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo guardar la imagen"})
		return
	}

	url := path
	if !strings.HasPrefix(strings.ToLower(url), "http") {
		url = h.s3.GetImageURL(path)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"url": url, "path": path}})
}
