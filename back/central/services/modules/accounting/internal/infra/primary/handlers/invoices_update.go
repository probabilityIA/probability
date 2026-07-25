package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateInvoice(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	var req request.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	issueDate, dueDate, ok := parseInvoiceDates(req.IssueDate, req.DueDate)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "fechas deben tener formato YYYY-MM-DD"})
		return
	}
	inv, err := h.uc.UpdateInvoice(c.Request.Context(), dtos.UpdateInvoiceDTO{
		ID:               uint(id),
		ConceptID:        req.ConceptID,
		IssueDate:        issueDate,
		DueDate:          dueDate,
		Notes:            req.Notes,
		EmailTo:          req.EmailTo,
		CustomerDocument: req.CustomerDocument,
		CustomerPhone:    req.CustomerPhone,
		CustomerAddress:  req.CustomerAddress,
		IsTest:           req.IsTest,
		TaxIDs:           req.TaxIDs,
		Items:            invoiceItemsToDTO(req.Items),
	})
	if err != nil {
		h.invoiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": response.FromInvoice(inv)})
}
