package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) RetrySavedQuoteGuide(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}

	businessID, ok := h.resolveBusinessIDParam(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo identificar la empresa"})
		return
	}

	repo := h.uc.Repo()

	quote, err := repo.GetSavedQuoteByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if quote == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cotizacion no encontrada"})
		return
	}
	if quote.BusinessID != businessID {
		c.JSON(http.StatusConflict, gin.H{"error": "la cotizacion no pertenece a este negocio"})
		return
	}
	if quote.OrderUUID == nil || *quote.OrderUUID == "" {
		c.JSON(http.StatusConflict, gin.H{"error": "la cotizacion no tiene una orden asociada"})
		return
	}
	orderUUID := *quote.OrderUUID

	existing, _ := repo.GetShipmentsByOrderID(ctx, orderUUID)
	for i := range existing {
		if shipmentHasActiveGuide(&existing[i]) {
			c.JSON(http.StatusConflict, gin.H{
				"error":       "La orden ya tiene una guia activa. Cancelala antes de generar otra.",
				"shipment_id": existing[i].ID,
			})
			return
		}
	}

	if quote.SelectedIDRate == nil || *quote.SelectedIDRate == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "la cotizacion no tiene una tarifa elegida, vuelve a cotizar"})
		return
	}
	if quote.ExpiresAt == nil || !quote.ExpiresAt.After(time.Now()) {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "La tarifa de esta cotizacion ya expiro. Genera la guia desde la orden para cotizar de nuevo.",
			"expired": true,
		})
		return
	}

	recipient, err := repo.GetOrderRecipient(ctx, orderUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if recipient == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "orden no encontrada"})
		return
	}

	if msg := missingRecipientData(recipient); msg != "" {
		h.failQuoteWithMessage(ctx, quote, msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "needs_data": true})
		return
	}

	carrier, err := repo.GetActiveShippingCarrier(ctx, businessID)
	if err != nil || carrier == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "el negocio no tiene una transportadora activa"})
		return
	}

	payload := buildRetryPayload(quote, recipient, orderUUID)

	shipmentID, err := h.createOrReuseShipmentForRetry(ctx, payload, carrier, orderUUID, existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar el envio: " + err.Error()})
		return
	}

	effectiveBaseURL := carrier.BaseURL
	if carrier.IsTesting && carrier.BaseURLTest != "" {
		effectiveBaseURL = carrier.BaseURLTest
	}

	correlationID := uuid.New().String()
	msg := &domain.TransportRequestMessage{
		ShipmentID:        &shipmentID,
		Provider:          carrier.ProviderCode,
		IntegrationTypeID: carrier.IntegrationTypeID,
		Operation:         "generate",
		CorrelationID:     correlationID,
		BusinessID:        businessID,
		IntegrationID:     carrier.IntegrationID,
		BaseURL:           effectiveBaseURL,
		IsTest:            carrier.IsTesting,
		Timestamp:         time.Now(),
		Payload:           payload,
	}

	if err := h.transportPub.PublishTransportRequest(ctx, msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar la solicitud: " + err.Error()})
		return
	}

	quote.Status = domain.QuoteStatusGenerating
	quote.ErrorMessage = ""
	_ = repo.UpdateSavedQuote(ctx, quote)

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Reintento enviado. La guia se generara en breve.",
		"correlation_id": correlationID,
		"shipment_id":    shipmentID,
	})
}

func (h *Handlers) failQuoteWithMessage(ctx context.Context, quote *domain.SavedQuote, message string) {
	quote.Status = domain.QuoteStatusFailed
	quote.ErrorMessage = domain.SafeQuoteMessage(message)
	_ = h.uc.Repo().UpdateSavedQuote(ctx, quote)
}

func (h *Handlers) createOrReuseShipmentForRetry(
	ctx context.Context,
	payload map[string]interface{},
	carrier *domain.CarrierInfo,
	orderUUID string,
	existing []domain.Shipment,
) (uint, error) {
	for i := range existing {
		s := existing[i]
		if s.Status == "pending" && (s.TrackingNumber == nil || *s.TrackingNumber == "") && (s.GuideURL == nil || *s.GuideURL == "") {
			return s.ID, nil
		}
	}

	req := buildShipmentRequest(payload, carrier)
	req.OrderID = &orderUUID
	resp, err := h.uc.CreateShipment(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func missingRecipientData(r *domain.OrderRecipient) string {
	if !looksLikeEmail(r.Email) {
		return "La orden no tiene un correo valido del destinatario. Completalo en la orden y vuelve a reintentar."
	}
	digits := onlyDigits(r.Phone)
	if len(digits) < 7 || len(digits) > 10 {
		return "El telefono del destinatario no tiene un formato valido (entre 7 y 10 digitos). Corrigelo en la orden y vuelve a reintentar."
	}
	if strings.TrimSpace(r.Address) == "" {
		return "La orden no tiene direccion de destino. Completala en la orden y vuelve a reintentar."
	}
	return ""
}

func looksLikeEmail(raw string) bool {
	email := strings.TrimSpace(raw)
	at := strings.Index(email, "@")
	if at <= 0 || at == len(email)-1 {
		return false
	}
	domain := email[at+1:]
	dot := strings.LastIndex(domain, ".")
	return dot > 0 && dot < len(domain)-1 && !strings.Contains(email, " ")
}

func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	digits := b.String()
	if len(digits) == 12 && strings.HasPrefix(digits, "57") {
		return digits[2:]
	}
	return digits
}

func buildRetryPayload(quote *domain.SavedQuote, r *domain.OrderRecipient, orderUUID string) map[string]interface{} {
	payload := make(map[string]interface{}, len(quote.RequestPayload)+4)
	for k, v := range quote.RequestPayload {
		payload[k] = v
	}

	dest := make(map[string]interface{})
	if existing, ok := payload["destination"].(map[string]interface{}); ok {
		for k, v := range existing {
			dest[k] = v
		}
	}

	dest["firstName"] = r.FirstName
	dest["lastName"] = r.LastName
	dest["email"] = r.Email
	dest["phone"] = onlyDigits(r.Phone)
	dest["address"] = r.Address
	if r.Neighborhood != "" {
		dest["suburb"] = r.Neighborhood
	}
	if strings.TrimSpace(asString(dest["suburb"])) == "" {
		dest["suburb"] = r.City
	}
	if strings.TrimSpace(asString(dest["company"])) == "" {
		dest["company"] = r.FullName
	}

	payload["destination"] = dest
	payload["idRate"] = *quote.SelectedIDRate
	payload["carrier"] = quote.SelectedCarrier
	payload["order_uuid"] = orderUUID
	payload["external_order_id"] = orderUUID
	payload["myShipmentReference"] = orderUUID

	return payload
}

func asString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
