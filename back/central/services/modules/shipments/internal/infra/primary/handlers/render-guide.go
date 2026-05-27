package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) RenderGuide(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	formatCode := c.Query("format")

	rendered, err := h.uc.RenderGuide(c.Request.Context(), uint(id64), formatCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	disposition := "inline"
	if c.Query("download") == "1" {
		disposition = "attachment"
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", disposition+"; filename=\""+rendered.Filename+"\"")
	c.Header("Cache-Control", "private, max-age=300")
	c.Data(http.StatusOK, "application/pdf", rendered.PDF)
}

func (h *Handlers) ListGuideFormats(c *gin.Context) {
	carrier := c.Query("carrier")
	formats, err := h.uc.ListGuideFormats(c.Request.Context(), carrier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": formats})
}
