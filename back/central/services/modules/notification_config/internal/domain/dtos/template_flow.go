package dtos

type TemplateFlowDTO struct {
	ButtonText       string `json:"button_text"`
	TargetTemplateID uint   `json:"target_template_id"`
	Enabled          *bool  `json:"enabled"`
}

type ReplaceTemplateFlowsDTO struct {
	BusinessID       uint              `json:"-"`
	SourceTemplateID uint              `json:"-"`
	Flows            []TemplateFlowDTO `json:"flows"`
}

type FlowSendMessage struct {
	BusinessID     uint   `json:"business_id"`
	Phone          string `json:"phone"`
	TemplateName   string `json:"template_name"`
	Language       string `json:"language"`
	HeaderImageURL string `json:"header_image_url"`
}

type ButtonReplyEvent struct {
	BusinessID       uint   `json:"business_id"`
	PhoneNumber      string `json:"phone_number"`
	ButtonText       string `json:"button_text"`
	ContextMessageID string `json:"context_message_id"`
	MessageID        string `json:"message_id"`
}
