package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// ListInvoices lista facturas con filtros
func (h *handler) ListInvoices(c *gin.Context) {
	ctx := c.Request.Context()

	// Parsear parámetros de query
	filters := make(map[string]interface{})

	// Business ID (obligatorio)
	if businessID := c.Query("business_id"); businessID != "" {
		id, _ := strconv.ParseUint(businessID, 10, 32)
		filters["business_id"] = uint(id)
	}

	// Order ID
	if orderID := c.Query("order_id"); orderID != "" {
		filters["order_id"] = orderID
	}

	// Integration ID
	if integrationID := c.Query("integration_id"); integrationID != "" {
		id, _ := strconv.ParseUint(integrationID, 10, 32)
		filters["integration_id"] = uint(id)
	}

	// Status
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Provider ID
	if providerID := c.Query("provider_id"); providerID != "" {
		id, _ := strconv.ParseUint(providerID, 10, 32)
		filters["invoicing_provider_id"] = uint(id)
	}

	// Paginación
	page := 1
	if p := c.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}

	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	}

	if pageSize > 100 {
		pageSize = 100
	}

	filters["page"] = page
	filters["page_size"] = pageSize

	h.log.Debug(ctx).
		Interface("filters", filters).
		Msg("Listing invoices")

	// Llamar caso de uso
	invoices, err := h.useCase.ListInvoices(ctx, filters)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to list invoices")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "list_invoices_failed",
			Message: err.Error(),
		})
		return
	}

	// Convertir a response
	resp := mappers.InvoicesToResponse(invoices, int64(len(invoices)), page, pageSize)

	c.JSON(http.StatusOK, resp)
}
