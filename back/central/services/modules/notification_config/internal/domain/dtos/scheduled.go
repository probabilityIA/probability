package dtos

type CreateScheduledRuleDTO struct {
	BusinessID         uint   `json:"business_id"`
	IntegrationID      *uint  `json:"integration_id"`
	NotificationTypeID uint   `json:"notification_type_id"`
	WhatsappTemplateID *uint  `json:"whatsapp_template_id"`
	Name               string `json:"name"`
	Description        string `json:"description"`

	SegmentType         string `json:"segment_type"`
	DaysWithoutPurchase int    `json:"days_without_purchase"`
	MinOrders           int    `json:"min_orders"`
	MaxOrders           int    `json:"max_orders"`

	Timezone         string `json:"timezone"`
	FrequencyMinutes uint   `json:"frequency_minutes"`
	SendWindowStart  string `json:"send_window_start"`
	SendWindowEnd    string `json:"send_window_end"`

	CooldownDays  uint  `json:"cooldown_days"`
	DailySendCap  uint  `json:"daily_send_cap"`
	BatchSizeCap  uint  `json:"batch_size_cap"`
	RequiresOptIn *bool `json:"requires_opt_in"`

	Enabled   *bool `json:"enabled"`
	CreatedBy *uint `json:"-"`
}

type UpdateScheduledRuleDTO struct {
	ID uint `json:"-"`
	CreateScheduledRuleDTO
}

type ScheduledSendMessage struct {
	SendID       uint              `json:"send_id"`
	RuleID       uint              `json:"rule_id"`
	BusinessID   uint              `json:"business_id"`
	ClientID     uint              `json:"client_id"`
	Phone        string            `json:"phone"`
	TemplateName string            `json:"template_name"`
	Language     string            `json:"language"`
	Parameters   []string          `json:"parameters"`
	Context      map[string]string `json:"context"`
}

type ScheduledSendResult struct {
	SendID       uint   `json:"send_id"`
	Status       string `json:"status"`
	MessageID    string `json:"message_id"`
	ErrorMessage string `json:"error_message"`
}
