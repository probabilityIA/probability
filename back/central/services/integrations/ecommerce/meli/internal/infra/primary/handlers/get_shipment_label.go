package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (h *meliHandler) GetShipmentLabel(c *gin.Context) {
	idStr := c.Param("shipment_id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "shipment_id invalido"})
		return
	}

	var queryBusinessID *uint
	if raw := c.Query("business_id"); raw != "" {
		parsed, parseErr := strconv.ParseUint(raw, 10, 64)
		if parseErr != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id invalido"})
			return
		}
		value := uint(parsed)
		queryBusinessID = &value
	}

	businessID, ok := h.resolveBusinessID(c, queryBusinessID)
	if !ok {
		return
	}

	label, err := h.useCase.GetShipmentLabel(c.Request.Context(), uint(id64), businessID, c.Query("response_type"))
	if err != nil {
		status, message := labelErrorResponse(err)
		c.JSON(status, gin.H{"success": false, "error": message})
		return
	}

	disposition := "inline"
	if c.Query("download") == "1" {
		disposition = "attachment"
	}

	c.Header("Content-Type", label.ContentType)
	c.Header("Content-Disposition", disposition+"; filename=\""+label.Filename+"\"")
	c.Header("Cache-Control", "private, max-age=300")
	c.Data(http.StatusOK, label.ContentType, label.Content)
}

func labelErrorResponse(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrShipmentNotFound):
		return http.StatusNotFound, "envio no encontrado"
	case errors.Is(err, domain.ErrShipmentNotOwned):
		return http.StatusForbidden, "el envio no pertenece al negocio"
	case errors.Is(err, domain.ErrLabelAlreadyShipped):
		return http.StatusConflict, "la guia ya no se puede descargar: la transportadora recogio el envio y el canal no permite reimprimirla"
	case errors.Is(err, domain.ErrLabelNotPDF):
		return http.StatusBadGateway, "el canal devolvio un archivo ilegible en vez de la guia, intenta de nuevo en unos minutos"
	case errors.Is(err, domain.ErrLabelNotAvailable):
		return http.StatusNotFound, "la guia todavia no esta disponible para este envio"
	case errors.Is(err, domain.ErrIntegrationNotFound):
		return http.StatusNotFound, "integracion de MercadoLibre no encontrada"
	case errors.Is(err, domain.ErrTokenExpired), errors.Is(err, domain.ErrTokenRefreshFailed):
		return http.StatusBadGateway, "no fue posible autenticar con MercadoLibre"
	case errors.Is(err, domain.ErrRateLimited):
		return http.StatusTooManyRequests, "MercadoLibre esta limitando las solicitudes, intenta de nuevo"
	default:
		return http.StatusInternalServerError, "no fue posible descargar la guia, intenta de nuevo"
	}
}
