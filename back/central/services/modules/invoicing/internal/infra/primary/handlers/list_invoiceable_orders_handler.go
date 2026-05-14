package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

func (h *handler) ListInvoiceableOrders(c *gin.Context) {
	ctx := c.Request.Context()

	businessIDValue, exists := c.Get("business_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Business ID not found in context"})
		return
	}
	businessID, ok := businessIDValue.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid business ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filterBusinessID := businessID
	if businessID == 0 {
		if businessIDParam := c.Query("business_id"); businessIDParam != "" {
			if paramID, err := strconv.ParseUint(businessIDParam, 10, 32); err == nil {
				filterBusinessID = uint(paramID)
			}
		}
	}

	filter := dtos.InvoiceableOrdersFilter{
		BusinessID:    filterBusinessID,
		Page:          page,
		PageSize:      pageSize,
		OrderNumber:   c.Query("order_number"),
		CustomerName:  c.Query("customer_name"),
		CustomerEmail: c.Query("customer_email"),
		SortBy:        c.Query("sort_by"),
		SortOrder:     c.Query("sort_order"),
	}
	if v := c.Query("payment_status_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filter.PaymentStatusID = uint(id)
		}
	}
	if v := c.Query("fulfillment_status_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filter.FulfillmentStatusID = uint(id)
		}
	}
	if v := c.Query("start_date"); v != "" {
		if t, err := parseFilterDate(v, false); err == nil {
			filter.StartDate = &t
		}
	}
	if v := c.Query("end_date"); v != "" {
		if t, err := parseFilterDate(v, true); err == nil {
			filter.EndDate = &t
		}
	}

	orders, total, err := h.repo.GetInvoiceableOrders(ctx, filter)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to get invoiceable orders")
		c.JSON(500, gin.H{"error": "Failed to get invoiceable orders"})
		return
	}

	orderResponses := make([]response.InvoiceableOrder, len(orders))
	for i, order := range orders {
		orderResponses[i] = mappers.ToInvoiceableOrderResponse(order)
	}

	c.JSON(200, response.PaginatedInvoiceableOrders{
		Data:     orderResponses,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	})
}

func parseFilterDate(s string, endOfDay bool) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			if layout == "2006-01-02" && endOfDay {
				t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			}
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid date: " + s)
}
