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

	var shipmentID uint
	if shipmentReq.OrderID != nil && *shipmentReq.OrderID != "" {
		existing, _ := h.uc.Repo().GetShipmentsByOrderID(c.Request.Context(), *shipmentReq.OrderID)
		for i := range existing {
			if shipmentHasActiveGuide(&existing[i]) {
				tracking := ""
				if existing[i].TrackingNumber != nil {
					tracking = *existing[i].TrackingNumber
				}
				c.JSON(http.StatusConflict, gin.H{
					"error":           "La orden ya tiene una guia activa. Cancela la guia existente antes de generar una nueva.",
					"shipment_id":     existing[i].ID,
					"tracking_number": tracking,
				})
				return
			}
		}
		for i := range existing {
			s := &existing[i]
			if s.Status == "pending" && (s.TrackingNumber == nil || *s.TrackingNumber == "") && (s.GuideURL == nil || *s.GuideURL == "") {
				shipmentID = s.ID
				break
			}
		}
	}

	if shipmentID == 0 {
		shipmentResp, err := h.uc.CreateShipment(c.Request.Context(), shipmentReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear registro de envio: " + err.Error()})
			return
		}
		shipmentID = shipmentResp.ID
	} else {
		updateReq := &domain.UpdateShipmentRequest{
			TotalCost:     shipmentReq.TotalCost,
			CodCarrierFee: shipmentReq.CodCarrierFee,
			Carrier:       shipmentReq.Carrier,
			CarrierCode:   shipmentReq.CarrierCode,
			Weight:        shipmentReq.Weight,
			Height:        shipmentReq.Height,
			Width:         shipmentReq.Width,
			Length:        shipmentReq.Length,
		}
		if _, err := h.uc.UpdateShipment(c.Request.Context(), shipmentID, updateReq); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar envio: " + err.Error()})
			return
		}
	}

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
	}

	if err := h.transportPub.PublishTransportRequest(c.Request.Context(), msg); err != nil {
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

func shipmentHasActiveGuide(s *domain.Shipment) bool {
	if s == nil {
		return false
	}
	if s.Status == "cancelled" || s.Status == "failed" {
		return false
	}
	if s.TrackingNumber != nil && *s.TrackingNumber != "" {
		return true
	}
	if s.GuideURL != nil && *s.GuideURL != "" {
		return true
	}
	return false
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
