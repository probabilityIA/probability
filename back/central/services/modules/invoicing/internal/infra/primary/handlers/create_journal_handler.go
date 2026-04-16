package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// CreateJournal crea un comprobante contable (journal) manualmente
func (h *handler) CreateJournal(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreateJournal
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body for journal")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")
	dto := &dtos.CreateJournalDTO{
		OrderID:         req.OrderID,
		IsManual:        true,
		CreatedByUserID: &userID,
	}

	h.log.Info(ctx).
		Str("order_id", req.OrderID).
		Msg("Creating journal manually")

	invoice, err := h.useCase.CreateJournal(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Str("order_id", req.OrderID).Msg("Failed to create journal")
		handleDomainError(c, err, "journal_creation_failed")
		return
	}

	baseURL, bucket := h.getS3Config()
	resp := mappers.InvoiceToResponse(invoice, true, baseURL, bucket)

	h.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Msg("Journal created successfully")

	c.JSON(http.StatusCreated, resp)
}
