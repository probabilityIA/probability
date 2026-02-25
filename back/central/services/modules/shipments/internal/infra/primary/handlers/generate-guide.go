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
	// 1. Resolve business_id (JWT for normal users, order DB lookup for super admin)
	businessID, err := h.resolveBusinessIDFromOrder(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Resolve active shipping carrier
	carrier, err := h.resolveCarrier(c, businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Parse request body
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

	// 4. Pre-create shipment record so the response consumer can update it
	shipmentReq := buildShipmentRequest(raw, carrier)
	shipmentResp, err := h.uc.CreateShipment(c.Request.Context(), shipmentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear registro de envío: " + err.Error()})
		return
	}
	shipmentID := shipmentResp.ID

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
			"error": "Error al enviar solicitud de generación de guía: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":        true,
		"message":        "Solicitud de generación de guía enviada. Será procesada en breve.",
		"correlation_id": correlationID,
		"shipment_id":    shipmentID,
	})
}

// buildShipmentRequest extracts fields from the raw generate payload to pre-create the DB record.
func buildShipmentRequest(raw map[string]interface{}, carrier *domain.CarrierInfo) *domain.CreateShipmentRequest {
	req := &domain.CreateShipmentRequest{
		Status:      "pending",
		CarrierCode: strPtr(carrier.ProviderCode),
	}

	// order_uuid
	if v, ok := raw["order_uuid"].(string); ok && v != "" {
		req.OrderID = strPtr(v)
	}

	// totalCost
	if v, ok := raw["totalCost"].(float64); ok {
		req.TotalCost = float64Ptr(v)
	}

	// destination → ClientName, DestinationAddress
	if dest, ok := raw["destination"].(map[string]interface{}); ok {
		firstName, _ := dest["firstName"].(string)
		lastName, _ := dest["lastName"].(string)
		address, _ := dest["address"].(string)
		req.ClientName = fmt.Sprintf("%s %s", firstName, lastName)
		req.DestinationAddress = address
	}

	// packages[0] → dimensions
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

func strPtr(s string) *string    { return &s }
func float64Ptr(f float64) *float64 { return &f }
