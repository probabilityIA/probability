package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) EmitInvoiceDian(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	var req request.EmitInvoiceDianRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
	}
	inv, err := h.uc.EmitInvoiceDian(c.Request.Context(), dtos.EmitInvoiceDianDTO{
		ID:               uint(id),
		CustomerDocument: req.CustomerDocument,
		CustomerPhone:    req.CustomerPhone,
		CustomerAddress:  req.CustomerAddress,
	})
	if err != nil {
		h.invoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": response.FromInvoice(inv)})
}
