package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (h *Handlers) invoiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrInvoiceNotFound),
		errors.Is(err, domainerrors.ErrConceptNotFound),
		errors.Is(err, domainerrors.ErrBusinessNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
	case errors.Is(err, domainerrors.ErrInvoiceNotDraft),
		errors.Is(err, domainerrors.ErrDianNotConfigured),
		errors.Is(err, domainerrors.ErrDianAlreadyEmitted),
		errors.Is(err, domainerrors.ErrDianDocumentRequired),
		errors.Is(err, domainerrors.ErrDianCredentials),
		errors.Is(err, domainerrors.ErrInvoiceNoItems),
		errors.Is(err, domainerrors.ErrInvoiceInvalidItem),
		errors.Is(err, domainerrors.ErrInvoiceAlreadyPaid),
		errors.Is(err, domainerrors.ErrInvoiceCancelled),
		errors.Is(err, domainerrors.ErrInvoiceNoEmail),
		errors.Is(err, domainerrors.ErrInvoiceConceptKind),
		errors.Is(err, domainerrors.ErrInvoiceNotPayable):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
	}
}
