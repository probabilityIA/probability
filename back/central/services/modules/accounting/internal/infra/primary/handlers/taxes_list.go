package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListTaxes(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	items, err := h.uc.ListTaxes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	data := make([]response.TaxResponse, len(items))
	for i := range items {
		data[i] = response.FromTax(&items[i])
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
