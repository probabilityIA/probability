package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// ListInvoiceableOrders lista órdenes facturables
// GET /api/v1/invoicing/invoices/invoiceable-orders
//
// Para super admin (business_id = 0):
//   - Sin query param: lista órdenes de TODOS los businesses
//   - Con ?business_id=X: filtra solo ese business específico
// Para usuarios normales:
//   - Siempre filtra por su business_id (ignora query param)
func (h *handler) ListInvoiceableOrders(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener business_id del JWT (inyectado por middleware)
	businessIDValue, exists := c.Get("business_id")
	if !exists {
		h.log.Error(ctx).Msg("Business ID not found in context")
		c.JSON(401, gin.H{"error": "Business ID not found in context"})
		return
	}

	businessID, ok := businessIDValue.(uint)
	if !ok {
		h.log.Error(ctx).Msg("Invalid business ID type in context")
		c.JSON(500, gin.H{"error": "Invalid business ID"})
		return
	}

	// Parse query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Validar límites
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Determinar business_id a usar para filtrar
	filterBusinessID := businessID

	// Si es super admin (business_id = 0), permitir filtrar por business específico
	if businessID == 0 {
		if businessIDParam := c.Query("business_id"); businessIDParam != "" {
			if paramID, err := strconv.ParseUint(businessIDParam, 10, 32); err == nil {
				filterBusinessID = uint(paramID)
				h.log.Debug(ctx).
					Uint("filter_business_id", filterBusinessID).
					Msg("Super admin filtering by specific business")
			}
		} else {
			h.log.Debug(ctx).Msg("Super admin querying all businesses")
		}
	}

	h.log.Info(ctx).
		Uint("jwt_business_id", businessID).
		Uint("filter_business_id", filterBusinessID).
		Bool("is_super_admin", businessID == 0).
		Int("page", page).
		Int("page_size", pageSize).
		Msg("Listing invoiceable orders")

	// Obtener órdenes facturables del repositorio
	orders, total, err := h.orderRepo.GetInvoiceableOrders(ctx, filterBusinessID, page, pageSize)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to get invoiceable orders")
		c.JSON(500, gin.H{"error": "Failed to get invoiceable orders"})
		return
	}

	// Mapear a response
	orderResponses := make([]response.InvoiceableOrder, len(orders))
	for i, order := range orders {
		orderResponses[i] = mappers.ToInvoiceableOrderResponse(order)
	}

	// Retornar respuesta paginada
	c.JSON(200, response.PaginatedInvoiceableOrders{
		Data:     orderResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
