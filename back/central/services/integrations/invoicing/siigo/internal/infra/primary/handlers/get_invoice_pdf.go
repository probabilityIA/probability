package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type IInvoiceReadRepository interface {
	GetInvoiceRef(ctx context.Context, invoiceID uint) (*dtos.InvoiceRef, error)
}

func (h *Handler) GetInvoicePDF(c *gin.Context) {
	ctx := c.Request.Context()

	if h.invoiceRead == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "error": "consulta de facturas no disponible"})
		return
	}

	invoiceID, err := strconv.ParseUint(c.Param("invoiceID"), 10, 64)
	if err != nil || invoiceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id de factura invalido"})
		return
	}

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return
	}
	if businessID == 0 {
		param := c.Query("business_id")
		parsed, parseErr := strconv.ParseUint(param, 10, 64)
		if parseErr != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
			return
		}
		businessID = uint(parsed)
	}

	ref, err := h.invoiceRead.GetInvoiceRef(ctx, uint(invoiceID))
	if err != nil || ref == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "factura no encontrada"})
		return
	}

	if ref.BusinessID != businessID {
		h.log.Warn(ctx).
			Uint("invoice_id", ref.ID).
			Uint("invoice_business_id", ref.BusinessID).
			Uint("request_business_id", businessID).
			Msg("Intento de descargar el PDF de una factura de otro negocio")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "factura no encontrada"})
		return
	}

	if ref.ExternalID == "" {
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": "la factura no tiene identificador de Siigo: no se puede descargar el PDF"})
		return
	}
	if ref.IntegrationID == 0 {
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": "la factura no tiene integracion de facturacion asociada"})
		return
	}

	contenido, nombre, err := h.useCase.GetInvoicePDF(ctx, ref.IntegrationID, ref.ExternalID)
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", ref.ID).Msg("No se pudo descargar el PDF de la factura en Siigo")
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", nombre))
	c.Data(http.StatusOK, "application/pdf", contenido)
}
