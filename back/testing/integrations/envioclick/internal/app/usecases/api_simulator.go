package usecases

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/domain"
)

// HandleQuote simulates the POST /api/v2/quotation endpoint
func (s *APISimulator) HandleQuote(req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	s.logger.Info().
		Str("origin", req.Origin.DaneCode).
		Str("destination", req.Destination.DaneCode).
		Int("packages", len(req.Packages)).
		Float64("content_value", req.ContentValue).
		Msg("Simulando cotizacion de envio")

	// Validate DANE codes
	if req.Origin.DaneCode == "" {
		return nil, fmt.Errorf("error: el codigo dane de origen no existe o no es valido")
	}
	if !IsValidDaneCode(req.Origin.DaneCode) {
		return nil, fmt.Errorf("error: el codigo dane de origen no existe o no es valido")
	}
	if req.Destination.DaneCode == "" {
		return nil, fmt.Errorf("error: el codigo dane del destino no existe o no es valido")
	}
	if !IsValidDaneCode(req.Destination.DaneCode) {
		return nil, fmt.Errorf("error: el codigo dane del destino no existe o no es valido")
	}

	// Validate packages
	if len(req.Packages) == 0 {
		return nil, fmt.Errorf("Error de validacion: Faltan datos obligatorios o hay datos invalidos en la solicitud")
	}
	for _, pkg := range req.Packages {
		if pkg.Weight <= 0 {
			return nil, fmt.Errorf("El peso del paquete es invalido")
		}
		if pkg.Height <= 0 || pkg.Width <= 0 || pkg.Length <= 0 {
			return nil, fmt.Errorf("Las dimensiones del paquete son invalidas")
		}
	}

	// Validate content value
	if req.ContentValue <= 0 {
		return nil, fmt.Errorf("El valor declarado es invalido o esta fuera de rango")
	}

	rates := GenerateRates(req)

	s.logger.Info().
		Int("rates_count", len(rates)).
		Str("origin_city", GetCityName(req.Origin.DaneCode)).
		Str("dest_city", GetCityName(req.Destination.DaneCode)).
		Msg("Cotizacion generada exitosamente")

	return &domain.QuoteResponse{
		Status: "success",
		Data: domain.QuoteData{
			Rates: rates,
		},
	}, nil
}

// HandleGenerate simulates the POST /api/v2/shipment endpoint
func (s *APISimulator) HandleGenerate(req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	s.logger.Info().
		Int64("id_rate", req.IDRate).
		Str("origin", req.Origin.DaneCode).
		Str("destination", req.Destination.DaneCode).
		Msg("Simulando generacion de guia")

	// Same validations as quote
	if req.Origin.DaneCode == "" || !IsValidDaneCode(req.Origin.DaneCode) {
		return nil, fmt.Errorf("error: el codigo dane de origen no existe o no es valido")
	}
	if req.Destination.DaneCode == "" || !IsValidDaneCode(req.Destination.DaneCode) {
		return nil, fmt.Errorf("error: el codigo dane del destino no existe o no es valido")
	}
	if len(req.Packages) == 0 {
		return nil, fmt.Errorf("Error de validacion: Faltan datos obligatorios o hay datos invalidos en la solicitud")
	}
	for _, pkg := range req.Packages {
		if pkg.Weight <= 0 {
			return nil, fmt.Errorf("El peso del paquete es invalido")
		}
		if pkg.Height <= 0 || pkg.Width <= 0 || pkg.Length <= 0 {
			return nil, fmt.Errorf("Las dimensiones del paquete son invalidas")
		}
	}
	if req.ContentValue <= 0 {
		return nil, fmt.Errorf("El valor declarado es invalido o esta fuera de rango")
	}

	// Pick a carrier based on IDRate or random
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var carrier *carrierInfo
	if req.IDRate > 0 {
		// Try to find carrier by the rate's carrier ID - for mock just pick a random one
		carrier = &carriers[rng.Intn(len(carriers))]
	} else {
		carrier = &carriers[rng.Intn(len(carriers))]
	}

	shipmentID := s.Repository.GenerateShipmentID()
	trackingNumber := GenerateTrackingNumber(*carrier, rng)
	labelURL := fmt.Sprintf("https://envioclick-mock.local/labels/%s.pdf", shipmentID)

	// Calculate flete
	totalWeight := 0.0
	for _, pkg := range req.Packages {
		totalWeight += pkg.Weight
	}
	flete := calculateFlete(totalWeight, req.ContentValue, rng)

	// Store the shipment
	now := time.Now()
	shipment := &domain.StoredShipment{
		ID:             shipmentID,
		TrackingNumber: trackingNumber,
		Carrier:        carrier.Name,
		CarrierID:      carrier.ID,
		Product:        carrier.Products[rng.Intn(len(carrier.Products))].Name,
		Origin:         req.Origin,
		Destination:    req.Destination,
		Packages:       req.Packages,
		ContentValue:   req.ContentValue,
		Flete:          flete,
		Status:         "created",
		LabelURL:       labelURL,
		CreatedAt:      now,
	}
	s.Repository.SaveShipment(shipment)

	s.logger.Info().
		Str("shipment_id", shipmentID).
		Str("tracking", trackingNumber).
		Str("carrier", carrier.Name).
		Float64("flete", flete).
		Msg("Guia generada exitosamente")

	return &domain.GenerateResponse{
		Status: "success",
		Data: domain.GenerateData{
			TrackingNumber:   trackingNumber,
			LabelURL:         labelURL,
			MyGuideReference: req.MyShipmentReference,
		},
	}, nil
}

// HandleTrack simulates the POST /api/v2/track endpoint
func (s *APISimulator) HandleTrack(trackingNumber string) (*domain.TrackingResponse, error) {
	s.logger.Info().
		Str("tracking", trackingNumber).
		Msg("Simulando rastreo de envio")

	if trackingNumber == "" {
		return nil, fmt.Errorf("tracking code is required")
	}

	shipment, exists := s.Repository.GetByTracking(trackingNumber)
	if !exists {
		return nil, fmt.Errorf("shipment not found for tracking: %s", trackingNumber)
	}

	status := "in_transit"
	if shipment.Status == "cancelled" {
		status = "cancelled"
	}

	history := GenerateTrackingHistory(shipment.Carrier, shipment.CreatedAt)

	s.logger.Info().
		Str("tracking", trackingNumber).
		Str("carrier", shipment.Carrier).
		Str("status", status).
		Int("events", len(history)).
		Msg("Rastreo generado exitosamente")

	return &domain.TrackingResponse{
		Status: "success",
		Data: domain.TrackingData{
			TrackingNumber: trackingNumber,
			Carrier:        shipment.Carrier,
			Status:         status,
			Events:         history,
		},
	}, nil
}

// HandleCancel simulates the DELETE /api/v2/shipment/:id endpoint
func (s *APISimulator) HandleCancel(shipmentID string) (*domain.CancelResponse, error) {
	s.logger.Info().
		Str("shipment_id", shipmentID).
		Msg("Simulando cancelacion de envio")

	if shipmentID == "" {
		return nil, fmt.Errorf("shipment ID is required")
	}

	shipment, exists := s.Repository.GetByID(shipmentID)
	if !exists {
		return nil, fmt.Errorf("shipment not found: %s", shipmentID)
	}

	if shipment.Status == "cancelled" {
		return nil, fmt.Errorf("shipment already cancelled: %s", shipmentID)
	}

	now := time.Now()
	shipment.CancelledAt = &now
	s.Repository.MarkCancelled(shipmentID)

	s.logger.Info().
		Str("shipment_id", shipmentID).
		Str("tracking", shipment.TrackingNumber).
		Msg("Envio cancelado exitosamente")

	return &domain.CancelResponse{
		Status:  "success",
		Message: "Cancelacion exitosa",
	}, nil
}
