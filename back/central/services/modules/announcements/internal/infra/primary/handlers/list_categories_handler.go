package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/mappers"
)

func (h *handler) ListCategories(c *gin.Context) {
	cats, err := h.uc.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	resp := mappers.CategoriesToResponse(cats)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}
