package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// DeleteInvoice elimina una factura pendiente con 3+ intentos de consulta sin respuesta
func (h *handler) DeleteInvoice(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid invoice ID",
		})
		return
	}

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Deleting pending invoice")

	err = h.useCase.DeletePendingInvoice(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to delete invoice")
		handleDomainError(c, err, "delete_invoice_failed")
		return
	}

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Invoice deleted successfully")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Factura eliminada exitosamente",
		"invoice_id": uint(id),
	})
}
