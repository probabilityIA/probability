package dtos

type TemplateVariableDTO struct {
	Position int    `json:"position"`
	Source   string `json:"source"`
	Label    string `json:"label"`
	Fallback string `json:"fallback"`
}

type TemplateButtonDTO struct {
	Type string `json:"type"`
	Text string `json:"text"`
	URL  string `json:"url"`
}

type CreateTemplateDTO struct {
	BusinessID uint                  `json:"business_id"`
	Scope      string                `json:"scope"`
	Name       string                `json:"name"`
	Language   string                `json:"language"`
	Category   string                `json:"category"`
	HeaderText string                `json:"header_text"`
	BodyText   string                `json:"body_text"`
	FooterText string                `json:"footer_text"`
	Variables  []TemplateVariableDTO `json:"variables"`
	Buttons    []TemplateButtonDTO   `json:"buttons"`
	CreatedBy  *uint                 `json:"-"`
}

type UpdateTemplateDTO struct {
	ID         uint                  `json:"-"`
	BusinessID uint                  `json:"business_id"`
	HeaderText string                `json:"header_text"`
	BodyText   string                `json:"body_text"`
	FooterText string                `json:"footer_text"`
	Category   string                `json:"category"`
	Variables  []TemplateVariableDTO `json:"variables"`
	Buttons    []TemplateButtonDTO   `json:"buttons"`
}

type TemplateSubmissionMessage struct {
	Action         string           `json:"action"`
	TemplateID     uint             `json:"template_id"`
	BusinessID     uint             `json:"business_id"`
	Name           string           `json:"name"`
	Language       string           `json:"language"`
	Category       string           `json:"category"`
	MetaTemplateID string           `json:"meta_template_id"`
	Components     []map[string]any `json:"components"`
}

type TemplateSubmissionResult struct {
	TemplateID     uint   `json:"template_id"`
	BusinessID     uint   `json:"business_id"`
	MetaTemplateID string `json:"meta_template_id"`
	WABAID         string `json:"waba_id"`
	Name           string `json:"name"`
	Language       string `json:"language"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	ErrorMessage   string `json:"error_message"`
}
