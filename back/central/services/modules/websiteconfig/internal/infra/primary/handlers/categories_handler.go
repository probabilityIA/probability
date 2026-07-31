package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetCategories(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "business_id es requerido"})
		return
	}

	categories, err := h.uc.GetProductCategories(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "error consultando categorias"})
		return
	}

	items := make([]gin.H, 0, len(categories))
	for _, cat := range categories {
		items = append(items, gin.H{"name": cat.Name, "product_count": cat.ProductCount})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}
