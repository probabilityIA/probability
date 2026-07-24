package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func invoiceItemsToDTO(items []request.InvoiceItemRequest) []dtos.InvoiceItemDTO {
	out := make([]dtos.InvoiceItemDTO, len(items))
	for i, item := range items {
		out[i] = dtos.InvoiceItemDTO{Description: item.Description, Quantity: item.Quantity, UnitPrice: item.UnitPrice}
	}
	return out
}

func parseInvoiceDates(issue, due string) (time.Time, *time.Time, bool) {
	issueDate, err := time.Parse("2006-01-02", issue)
	if err != nil {
		return time.Time{}, nil, false
	}
	var dueDate *time.Time
	if due != "" {
		d, err := time.Parse("2006-01-02", due)
		if err != nil {
			return time.Time{}, nil, false
		}
		dueDate = &d
	}
	return issueDate, dueDate, true
}

func (h *Handlers) CreateInvoice(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	var req request.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	issueDate, dueDate, ok := parseInvoiceDates(req.IssueDate, req.DueDate)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "fechas deben tener formato YYYY-MM-DD"})
		return
	}
	inv, err := h.uc.CreateInvoice(c.Request.Context(), dtos.CreateInvoiceDTO{
		BusinessID: req.BusinessID,
		ConceptID:  req.ConceptID,
		IssueDate:  issueDate,
		DueDate:    dueDate,
		Notes:      req.Notes,
		EmailTo:    req.EmailTo,
		TaxIDs:     req.TaxIDs,
		Items:      invoiceItemsToDTO(req.Items),
	})
	if err != nil {
		h.invoiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": response.FromInvoice(inv)})
}
