package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) GenerateGuide(c *gin.Context) {
	businessID, err := h.resolveBusinessIDFromOrder(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.overageChecker != nil {
		blocked, reason, fee, err := h.overageChecker(c.Request.Context(), businessID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar el plan de suscripcion: " + err.Error()})
			return
		}
		if blocked {
			message := fmt.Sprintf("Superaste los envios incluidos en tu plan gratuito. Desde ahora cada guia adicional tiene un cargo de $%.0f que se facturara al cierre del periodo. Acepta el cargo para seguir generando guias.", fee)
			if reason == "overage_payment_due" {
				message = fmt.Sprintf("Tienes un cargo pendiente de $%.0f por excedente de envios del ciclo anterior. Debes pagarlo para seguir generando guias.", fee)
			}
			c.JSON(http.StatusPaymentRequired, gin.H{
				"error":       message,
				"code":        reason,
				"overage_fee": fee,
			})
			return
		}
	}

	carrier, err := h.resolveCarrier(c, businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := raw["origin"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin es requerido"})
		return
	}
	if _, ok := raw["destination"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "destination es requerido"})
		return
	}
	if _, ok := raw["packages"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "packages es requerido"})
		return
	}

	shipmentReq := buildShipmentRequest(raw, carrier)

	actor := h.resolveActor(c)
	shipmentReq.CreatedBy = actor.ID
	shipmentReq.CreatedByName = actor.Name

	shipmentID, conflict, err := h.reserveShipmentForGuide(c.Request.Context(), shipmentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar el envio: " + err.Error()})
		return
	}
	if conflict != nil {
		c.JSON(conflict.status, conflict.body)
		return
	}

	h.overrideCodValue(c, raw, shipmentReq)

	correlationID := uuid.New().String()

	effectiveBaseURL := carrier.BaseURL
	if carrier.IsTesting && carrier.BaseURLTest != "" {
		effectiveBaseURL = carrier.BaseURLTest
	}

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
		Payload:           raw,
		UserID:            actor.ID,
		TriggeredBy:       "user",
	}

	if err := h.transportPub.PublishTransportRequest(c.Request.Context(), msg); err != nil {
		_ = h.uc.Repo().ReleaseShipmentGenerating(c.Request.Context(), shipmentID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al enviar solicitud de generacion de guia: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Solicitud de generacion de guia enviada. Sera procesada en breve.",
		"correlation_id": correlationID,
		"shipment_id":    shipmentID,
	})
}

func (h *Handlers) overrideCodValue(c *gin.Context, raw map[string]interface{}, req *domain.CreateShipmentRequest) {
	if v, ok := raw["codValue"].(float64); !ok || v <= 0 {
		return
	}
	if req.OrderID == nil || *req.OrderID == "" {
		return
	}

	basis, err := h.uc.Repo().GetOrderCodBasis(c.Request.Context(), *req.OrderID)
	if err != nil || basis == nil {
		return
	}

	totalCost := 0.0
	if req.TotalCost != nil {
		totalCost = *req.TotalCost
	}
	carrierFee := 0.0
	if req.CodCarrierFee != nil {
		carrierFee = *req.CodCarrierFee
	}

	amount := basis.AmountToCollect(totalCost, carrierFee)
	if amount > 0 {
		raw["codValue"] = amount
	}

	if netTarget := basis.NetTarget(totalCost); netTarget > 0 {
		raw["codNetTarget"] = netTarget
	}
}

func shipmentHasActiveGuide(s *domain.Shipment) bool {
	return s.HasActiveGuide()
}

func buildShipmentRequest(raw map[string]interface{}, carrier *domain.CarrierInfo) *domain.CreateShipmentRequest {
	req := &domain.CreateShipmentRequest{
		Status:      "pending",
		CarrierCode: strPtr(carrier.ProviderCode),
	}

	if v, ok := raw["carrier"].(string); ok && v != "" {
		req.Carrier = strPtr(v)
	}

	if v, ok := raw["order_uuid"].(string); ok && v != "" {
		req.OrderID = strPtr(v)
	}

	if v, ok := raw["totalCost"].(float64); ok {
		req.TotalCost = float64Ptr(v)
	}

	if v, ok := raw["codCarrierFee"].(float64); ok && v > 0 {
		req.CodCarrierFee = float64Ptr(v)
	}

	if dest, ok := raw["destination"].(map[string]interface{}); ok {
		firstName, _ := dest["firstName"].(string)
		lastName, _ := dest["lastName"].(string)
		address, _ := dest["address"].(string)
		city, _ := dest["city"].(string)
		state, _ := dest["state"].(string)
		suburb, _ := dest["suburb"].(string)
		req.ClientName = fmt.Sprintf("%s %s", firstName, lastName)
		req.DestinationAddress = address
		req.DestinationCity = city
		req.DestinationState = state
		req.DestinationSuburb = suburb
	}

	if pkgs, ok := raw["packages"].([]interface{}); ok && len(pkgs) > 0 {
		if pkg, ok := pkgs[0].(map[string]interface{}); ok {
			if v, ok := pkg["weight"].(float64); ok {
				req.Weight = float64Ptr(v)
			}
			if v, ok := pkg["height"].(float64); ok {
				req.Height = float64Ptr(v)
			}
			if v, ok := pkg["width"].(float64); ok {
				req.Width = float64Ptr(v)
			}
			if v, ok := pkg["length"].(float64); ok {
				req.Length = float64Ptr(v)
			}
		}
	}

	return req
}

func strPtr(s string) *string       { return &s }
func float64Ptr(f float64) *float64 { return &f }
