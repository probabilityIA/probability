package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) ListCategories(c *gin.Context) {
	categories, err := h.uc.ListCategories(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	if categories == nil {
		categories = []string{}
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}
