package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) MarkInvoicePaid(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	inv, err := h.uc.MarkInvoicePaid(c.Request.Context(), uint(id))
	if err != nil {
		h.invoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": response.FromInvoice(inv)})
}
