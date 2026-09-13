package dtos

import "time"

type CreateCampaignDTO struct {
	BusinessID         uint   `json:"business_id"`
	IntegrationID      *uint  `json:"integration_id"`
	WhatsappTemplateID *uint  `json:"whatsapp_template_id"`
	FlowID             *uint  `json:"flow_id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	SenderName         string `json:"sender_name"`

	AudienceType     string `json:"audience_type"`
	City             string `json:"city"`
	State            string `json:"state"`
	CreatedFromDays  int    `json:"created_from_days"`
	OnlyWithoutOrder bool   `json:"only_without_order"`
	ClientIDs        []uint `json:"client_ids"`

	ExcludeRecentDays int `json:"exclude_recent_days"`
	ExcludeRecentMax  int `json:"exclude_recent_max"`

	RegisteredBeforeDays   int     `json:"registered_before_days"`
	MinOrders              int     `json:"min_orders"`
	MinSpent               float64 `json:"min_spent"`
	LastPurchaseBeforeDays int     `json:"last_purchase_before_days"`

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

type CampaignAudienceClientDTO struct {
	ClientID uint   `json:"client_id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	City     string `json:"city"`
}

type AudienceLocationDTO struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Clients uint   `json:"clients"`
}

type CampaignAudiencePreviewDTO struct {
	Total      uint                        `json:"total"`
	OptedOut   uint                        `json:"opted_out"`
	NoPhone    uint                        `json:"no_phone"`
	Reachable  uint                        `json:"reachable"`
	SampleName []string                    `json:"sample_names"`
	Clients    []CampaignAudienceClientDTO `json:"clients"`
}

type CampaignSendMessage struct {
	SendID         uint              `json:"send_id"`
	CampaignID     uint              `json:"campaign_id"`
	BusinessID     uint              `json:"business_id"`
	ClientID       uint              `json:"client_id"`
	Phone          string            `json:"phone"`
	TemplateName   string            `json:"template_name"`
	Language       string            `json:"language"`
	Parameters     []string          `json:"parameters"`
	HeaderImageURL string            `json:"header_image_url"`
	Context        map[string]string `json:"context"`
}

type CampaignSendResult struct {
	SendID       uint   `json:"send_id"`
	Status       string `json:"status"`
	MessageID    string `json:"message_id"`
	ErrorMessage string `json:"error_message"`
}
