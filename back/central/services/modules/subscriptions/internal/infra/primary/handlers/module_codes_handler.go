package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetModuleCodes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": h.uc.GetModuleCodes()})
}

func (h *Handlers) GetModuleCatalog(c *gin.Context) {
	catalog := h.uc.GetModuleCatalog()
	data := make([]gin.H, 0, len(catalog))
	for _, m := range catalog {
		data = append(data, gin.H{"code": m.Code, "name": m.Name})
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handlers) GetMyModules(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	modules, err := h.uc.GetAccessibleModules(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve accessible modules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": modules})
}
