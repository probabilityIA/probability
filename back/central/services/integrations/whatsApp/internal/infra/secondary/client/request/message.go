package request

// Language representa el idioma del template
type Language struct {
	Code string `json:"code"`
}

// Parameter representa un parámetro del template
type Parameter struct {
	Type          string `json:"type"`
	ParameterName string `json:"parameter_name"`
	Text          string `json:"text"`
}

// Component representa un componente del template
type Component struct {
	Type       string      `json:"type"`
	SubType    string      `json:"sub_type,omitempty"` // Para botones: "quick_reply", "url", "phone_number"
	Index      int         `json:"index,omitempty"`    // Para botones: índice del botón (0-based)
	Parameters []Parameter `json:"parameters,omitempty"`
}

// Template representa la configuración del template a enviar
type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components,omitempty"`
}

// TemplateMessageRequest representa el payload para enviar un mensaje tipo template
// a la API de WhatsApp Cloud.
type TemplateMessageRequest struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

// TextBody representa el cuerpo de un mensaje de texto
type TextBody struct {
	Body string `json:"body"`
}

// TextMessageRequest representa el payload para enviar un mensaje de texto
// a la API de WhatsApp Cloud.
type TextMessageRequest struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Text             TextBody `json:"text"`
}
