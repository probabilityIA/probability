package dtos

import "time"

type CreateCampaignDTO struct {
	BusinessID         uint   `json:"business_id"`
	IntegrationID      *uint  `json:"integration_id"`
	WhatsappTemplateID *uint  `json:"whatsapp_template_id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	SenderName         string `json:"sender_name"`

	AudienceType     string `json:"audience_type"`
	City             string `json:"city"`
	CreatedFromDays  int    `json:"created_from_days"`
	OnlyWithoutOrder bool   `json:"only_without_order"`
	ClientIDs        []uint `json:"client_ids"`

	VariableValues map[string]string `json:"variable_values"`

	Timezone        string     `json:"timezone"`
	SendWindowStart string     `json:"send_window_start"`
	SendWindowEnd   string     `json:"send_window_end"`
	ScheduledAt     *time.Time `json:"scheduled_at"`

	DailySendCap uint `json:"daily_send_cap"`
	BatchSize    uint `json:"batch_size"`

	CreatedBy *uint `json:"-"`
}

type UpdateCampaignDTO struct {
	ID uint `json:"-"`
	CreateCampaignDTO
}

type CampaignAudiencePreviewDTO struct {
	Total      uint     `json:"total"`
	OptedOut   uint     `json:"opted_out"`
	NoPhone    uint     `json:"no_phone"`
	Reachable  uint     `json:"reachable"`
	SampleName []string `json:"sample_names"`
}

type CampaignSendMessage struct {
	SendID       uint              `json:"send_id"`
	CampaignID   uint              `json:"campaign_id"`
	BusinessID   uint              `json:"business_id"`
	ClientID     uint              `json:"client_id"`
	Phone        string            `json:"phone"`
	TemplateName string            `json:"template_name"`
	Language     string            `json:"language"`
	Parameters   []string          `json:"parameters"`
	Context      map[string]string `json:"context"`
}

type CampaignSendResult struct {
	SendID       uint   `json:"send_id"`
	Status       string `json:"status"`
	MessageID    string `json:"message_id"`
	ErrorMessage string `json:"error_message"`
}
