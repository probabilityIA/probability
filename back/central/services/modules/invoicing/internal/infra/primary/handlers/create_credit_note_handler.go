package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// CreateCreditNote crea una nota de crédito para una factura
func (h *handler) CreateCreditNote(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener invoice ID del path
	idStr := c.Param("id")
	invoiceID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid invoice ID",
		})
		return
	}

	// Parsear request
	var req request.CreateCreditNote
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Asegurar que invoice_id del path coincide con el del body
	req.InvoiceID = uint(invoiceID)

	h.log.Info(ctx).
		Uint("invoice_id", req.InvoiceID).
		Str("note_type", req.NoteType).
		Float64("amount", req.Amount).
		Msg("Creating credit note")

	// Convertir a DTO de dominio
	dto := mappers.CreateCreditNoteRequestToDTO(&req)

	// Llamar caso de uso
	creditNote, err := h.useCase.CreateCreditNote(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", req.InvoiceID).Msg("Failed to create credit note")
		handleDomainError(c, err, "credit_note_creation_failed")
		return
	}

	// Convertir a response
	resp := mappers.CreditNoteToResponse(creditNote)

	h.log.Info(ctx).
		Uint("credit_note_id", creditNote.ID).
		Str("credit_note_number", creditNote.CreditNoteNumber).
		Msg("Credit note created successfully")

	c.JSON(http.StatusCreated, resp)
}
