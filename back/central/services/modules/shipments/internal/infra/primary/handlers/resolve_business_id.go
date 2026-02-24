package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// resolveBusinessIDFromOrder returns the business ID for the current request.
// For normal users, it returns the business_id from the JWT.
// For super admins (business_id == 0), it peeks at the JSON body to extract
// "order_uuid", then looks up the order's business_id from the database.
// The body is restored after peeking so downstream ShouldBindJSON still works.
func (h *Handlers) resolveBusinessIDFromOrder(c *gin.Context) (uint, error) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		return 0, errors.New("no se pudo identificar la empresa")
	}

	// Normal user: return JWT business_id directly
	if !middleware.IsSuperAdmin(c) {
		return businessID, nil
	}

	// Super admin: peek at body to find order_uuid, then DB lookup
	if c.Request.Body == nil {
		return 0, errors.New("super admin: body vacío, no se puede determinar la empresa")
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil || len(bodyBytes) == 0 {
		return 0, errors.New("super admin: no se pudo leer el body de la solicitud")
	}

	// Restore body so downstream ShouldBindJSON works
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		return 0, errors.New("super admin: body no es JSON válido")
	}

	orderUUID, _ := bodyMap["order_uuid"].(string)
	if orderUUID == "" {
		return 0, errors.New("super admin: order_uuid es requerido en el body para determinar la empresa")
	}

	// Look up the order's business_id from the database
	bid, err := h.uc.Repo().GetOrderBusinessID(c.Request.Context(), orderUUID)
	if err != nil {
		return 0, err
	}

	return bid, nil
}

// resolveBusinessIDFromShipment returns the business ID for track/cancel operations.
// For normal users, it returns the business_id from the JWT.
// For super admins, it looks up the shipment by tracking number or numeric ID
// and resolves the business_id via shipment → order → business_id.
func (h *Handlers) resolveBusinessIDFromShipment(c *gin.Context, identifier string) (uint, error) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		return 0, errors.New("no se pudo identificar la empresa")
	}

	// Normal user: return JWT business_id directly
	if !middleware.IsSuperAdmin(c) {
		return businessID, nil
	}

	// Super admin: try tracking number first
	bid, err := h.uc.Repo().GetShipmentBusinessIDByTracking(c.Request.Context(), identifier)
	if err == nil {
		return bid, nil
	}

	// Fallback: try as numeric shipment ID
	if numID, parseErr := strconv.ParseUint(identifier, 10, 64); parseErr == nil {
		bid, err = h.uc.Repo().GetShipmentBusinessIDByID(c.Request.Context(), uint(numID))
		if err == nil {
			return bid, nil
		}
	}

	return 0, fmt.Errorf("no se pudo resolver la empresa para el envío '%s'", identifier)
}

// resolveCarrier finds the active shipping carrier for the given business.
// If no carrier is configured, the error includes the business name for clarity.
func (h *Handlers) resolveCarrier(c *gin.Context, businessID uint) (*domain.CarrierInfo, error) {
	carrier, err := h.carrierResolver.GetActiveShippingCarrier(c.Request.Context(), businessID)
	if err != nil {
		return nil, fmt.Errorf("error al resolver transportadora: %w", err)
	}
	if carrier == nil {
		name, _ := h.uc.Repo().GetBusinessName(c.Request.Context(), businessID)
		return nil, fmt.Errorf("el negocio '%s' no tiene una transportadora activa configurada", name)
	}
	return carrier, nil
}
