package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/log"
)

var communes = []domain.Commune{
	{ID: 308, Name: "LAS CONDES"},
	{ID: 309, Name: "SANTIAGO"},
	{ID: 310, Name: "PROVIDENCIA"},
	{ID: 311, Name: "NUNOA"},
	{ID: 312, Name: "MAIPU"},
	{ID: 313, Name: "VALPARAISO"},
	{ID: 314, Name: "VINA DEL MAR"},
	{ID: 315, Name: "CONCEPCION"},
	{ID: 316, Name: "ANTOFAGASTA"},
	{ID: 317, Name: "TEMUCO"},
	{ID: 318, Name: "PUERTO MONTT"},
	{ID: 319, Name: "LA FLORIDA"},
	{ID: 320, Name: "PUENTE ALTO"},
	{ID: 321, Name: "BOGOTA"},
	{ID: 322, Name: "MEDELLIN"},
}

var couriers = []struct {
	Name      string
	BasePrice float64
	PerKg     float64
	Days      int
}{
	{"starken", 3200, 450, 2},
	{"chilexpress", 3800, 520, 1},
	{"bluexpress", 2900, 400, 3},
}

type APISimulator struct {
	Repository    *domain.Repository
	log           log.ILogger
	webhookTarget string
	baseURL       string
}

func NewAPISimulator(logger log.ILogger, webhookTarget, baseURL string) *APISimulator {
	return &APISimulator{
		Repository:    domain.NewRepository(),
		log:           logger,
		webhookTarget: webhookTarget,
		baseURL:       baseURL,
	}
}

func (s *APISimulator) GetCommunes() []domain.Commune {
	return communes
}

func (s *APISimulator) HandleRates(req domain.RatesRequest) (*domain.RatesResponse, error) {
	if req.Parcel.OriginID == 0 || req.Parcel.DestinyID == 0 {
		return nil, fmt.Errorf("origin_id y destiny_id son requeridos")
	}

	weight := req.Parcel.Weight
	if weight <= 0 {
		weight = 1
	}

	prices := make([]domain.Price, 0, len(couriers))
	for _, courier := range couriers {
		prices = append(prices, domain.Price{
			Courier:             domain.PriceCourier{Name: courier.Name},
			Name:                "ENVIO ESTANDAR",
			Price:               courier.BasePrice + courier.PerKg*weight,
			Days:                courier.Days,
			AvailableToShipping: true,
		})
	}

	return &domain.RatesResponse{Algorithm: "1", Prices: prices}, nil
}

func (s *APISimulator) HandleCreateShipment(req domain.ShipmentRequest) (*domain.ShipmentResponse, error) {
	body := req.Shipment
	if body.Reference == "" {
		return nil, fmt.Errorf("reference es requerido")
	}
	if body.Destiny.CommuneName == "" && body.Destiny.CommuneID == 0 {
		return nil, fmt.Errorf("destiny.commune_id o destiny.commune_name es requerido")
	}

	id := s.Repository.NextID()
	now := time.Now().Format(time.RFC3339)

	courier := body.Courier.Client
	if courier == "" {
		courier = couriers[int(id)%len(couriers)].Name
	}

	trackingNumber := fmt.Sprintf("SHPT%06d", id)
	ticketURL := fmt.Sprintf("%s/labels/%d.pdf", s.baseURL, id)

	stored := &domain.StoredShipment{
		ID:             id,
		Reference:      body.Reference,
		TrackingNumber: trackingNumber,
		Courier:        courier,
		Status:         "in_preparation",
		Sandbox:        body.Sandbox,
		Destiny:        body.Destiny,
		CreatedAt:      now,
		UpdatedAt:      now,
		Statuses: []domain.TrackingStatus{
			{
				ID:               1,
				TrackingID:       id,
				GenericStatus:    0,
				CourierStatus:    "in_preparation",
				CourierUpdatedAt: now,
				CreatedAt:        now,
			},
		},
	}
	s.Repository.Save(stored)

	s.log.Info().
		Str("reference", body.Reference).
		Str("tracking", trackingNumber).
		Str("courier", courier).
		Bool("sandbox", body.Sandbox).
		Msg("Shipit mock: shipment created")

	return &domain.ShipmentResponse{
		ID:               id,
		Reference:        body.Reference,
		Status:           "created",
		SubStatus:        "waiting_courier",
		TrackingNumber:   &trackingNumber,
		ShipitCode:       fmt.Sprintf("Shipit / %s", body.Reference),
		TrackingLink:     fmt.Sprintf("%s/tracking/%s", s.baseURL, trackingNumber),
		CourierURL:       "",
		CourierForClient: courier,
		TicketPDFURL:     &ticketURL,
		TicketURL:        &ticketURL,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (s *APISimulator) HandleTrack(trackingNumber string) (map[string]any, error) {
	stored := s.Repository.GetByTracking(trackingNumber)
	if stored == nil {
		return nil, fmt.Errorf("tracking number no encontrado")
	}

	return map[string]any{
		"data": map[string]any{
			"id":           stored.ID,
			"courier":      stored.Courier,
			"number":       stored.TrackingNumber,
			"package_id":   stored.ID,
			"is_delivered": strings.EqualFold(stored.Status, "delivered"),
			"created_at":   stored.CreatedAt,
			"statuses":     stored.Statuses,
		},
	}, nil
}

func (s *APISimulator) SimulateStatusChange(trackingNumber, status string) (*domain.StoredShipment, error) {
	stored := s.Repository.GetByTracking(trackingNumber)
	if stored == nil {
		return nil, fmt.Errorf("tracking number no encontrado: %s", trackingNumber)
	}

	now := time.Now().Format(time.RFC3339)
	stored.Status = status
	stored.UpdatedAt = now
	stored.Statuses = append(stored.Statuses, domain.TrackingStatus{
		ID:               int64(len(stored.Statuses) + 1),
		TrackingID:       stored.ID,
		GenericStatus:    len(stored.Statuses),
		CourierStatus:    status,
		CourierUpdatedAt: now,
		CreatedAt:        now,
	})
	s.Repository.Save(stored)

	if s.webhookTarget != "" {
		go s.sendWebhook(stored)
	}

	return stored, nil
}

func (s *APISimulator) sendWebhook(stored *domain.StoredShipment) {
	payload := map[string]any{
		"id":                 stored.ID,
		"reference":          stored.Reference,
		"status":             stored.Status,
		"courier_status":     stored.Status,
		"tracking_number":    stored.TrackingNumber,
		"courier_for_client": stored.Courier,
		"full_name":          stored.Destiny.FullName,
		"created_at":         stored.CreatedAt,
		"updated_at":         stored.UpdatedAt,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	resp, err := http.Post(s.webhookTarget, "application/json", bytes.NewReader(data))
	if err != nil {
		s.log.Error().Err(err).Str("target", s.webhookTarget).Msg("Shipit mock: webhook delivery failed")
		return
	}
	defer resp.Body.Close()

	s.log.Info().
		Str("tracking", stored.TrackingNumber).
		Str("status", stored.Status).
		Int("response", resp.StatusCode).
		Msg("Shipit mock: webhook sent")
}
