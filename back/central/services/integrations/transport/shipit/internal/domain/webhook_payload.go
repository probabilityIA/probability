package domain

type WebhookPayload struct {
	ID               int64  `json:"id"`
	Reference        string `json:"reference"`
	Status           string `json:"status"`
	CourierStatus    string `json:"courier_status"`
	TrackingNumber   string `json:"tracking_number"`
	CourierForClient string `json:"courier_for_client"`
	CourierURL       string `json:"courier_url"`
	FullName         string `json:"full_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type NormalizedWebhookUpdate struct {
	TrackingNumber    string
	Reference         string
	Carrier           string
	ProbabilityStatus ProbabilityShipmentStatus
	RawStatus         string
	RawStatusDetail   string
	IsUnknownStatus   bool
	EventTimestamp    string
	ShippedAt         *string
	DeliveredAt       *string
}

func (p *WebhookPayload) ToNormalizedUpdate() *NormalizedWebhookUpdate {
	if p == nil || p.TrackingNumber == "" {
		return nil
	}

	rawStatus := p.Status
	if rawStatus == "" {
		rawStatus = p.CourierStatus
	}
	if rawStatus == "" {
		return nil
	}

	status, known := MapShipitStatus(rawStatus)
	if !known {
		if courierMapped, courierKnown := MapShipitStatus(p.CourierStatus); courierKnown {
			status = courierMapped
			known = true
		}
	}

	update := &NormalizedWebhookUpdate{
		TrackingNumber:    p.TrackingNumber,
		Reference:         p.Reference,
		Carrier:           p.CourierForClient,
		ProbabilityStatus: status,
		RawStatus:         rawStatus,
		RawStatusDetail:   p.CourierStatus,
		IsUnknownStatus:   !known,
		EventTimestamp:    p.UpdatedAt,
	}

	switch status {
	case StatusInTransit, StatusPickedUp, StatusOutForDelivery:
		if p.UpdatedAt != "" {
			ts := p.UpdatedAt
			update.ShippedAt = &ts
		}
	case StatusDelivered:
		if p.UpdatedAt != "" {
			ts := p.UpdatedAt
			update.DeliveredAt = &ts
		}
	}

	return update
}
