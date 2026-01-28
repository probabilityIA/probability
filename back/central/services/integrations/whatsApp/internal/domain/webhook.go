package domain

// WebhookPayload representa el payload completo de un webhook de WhatsApp
type WebhookPayload struct {
	Object  string         `json:"object"`  // Siempre "whatsapp_business_account"
	Entry   []WebhookEntry `json:"entry"`
}

// WebhookEntry representa una entrada en el webhook
type WebhookEntry struct {
	ID      string          `json:"id"`      // ID de la cuenta de WhatsApp Business
	Changes []WebhookChange `json:"changes"`
}

// WebhookChange representa un cambio en el webhook
type WebhookChange struct {
	Value WebhookValue `json:"value"`
	Field string       `json:"field"` // "messages" o "message_template_status_update"
}

// WebhookValue contiene los datos del webhook
type WebhookValue struct {
	MessagingProduct string            `json:"messaging_product"` // "whatsapp"
	Metadata         WebhookMetadata   `json:"metadata"`
	Contacts         []WebhookContact  `json:"contacts,omitempty"`
	Messages         []WebhookMessage  `json:"messages,omitempty"`
	Statuses         []WebhookStatus   `json:"statuses,omitempty"`
}

// WebhookMetadata contiene metadata del webhook
type WebhookMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"` // Número de teléfono del negocio
	PhoneNumberID      string `json:"phone_number_id"`      // ID del número de teléfono
}

// WebhookContact representa información de contacto
type WebhookContact struct {
	Profile WebhookProfile `json:"profile"`
	WaID    string         `json:"wa_id"` // WhatsApp ID (número de teléfono)
}

// WebhookProfile contiene el perfil del contacto
type WebhookProfile struct {
	Name string `json:"name"` // Nombre del contacto
}

// WebhookMessage representa un mensaje recibido
type WebhookMessage struct {
	From        string                `json:"from"`               // Número de teléfono del remitente
	ID          string                `json:"id"`                 // ID del mensaje
	Timestamp   string                `json:"timestamp"`          // Unix timestamp
	Type        string                `json:"type"`               // "text", "button", "interactive"
	Text        *TextContent          `json:"text,omitempty"`     // Solo si type="text"
	Button      *ButtonResponse       `json:"button,omitempty"`   // Solo si type="button"
	Interactive *InteractiveResponse  `json:"interactive,omitempty"` // Solo si type="interactive"
	Context     *MessageContext       `json:"context,omitempty"`  // Contexto del mensaje (reply)
}

// TextContent representa contenido de texto
type TextContent struct {
	Body string `json:"body"` // Texto del mensaje
}

// ButtonResponse representa la respuesta a un botón quick reply
type ButtonResponse struct {
	Payload string `json:"payload"` // Payload del botón (opcional)
	Text    string `json:"text"`    // Texto del botón presionado
}

// InteractiveResponse representa respuesta a mensaje interactivo
type InteractiveResponse struct {
	Type           string                `json:"type"` // "button_reply", "list_reply"
	ButtonReply    *ButtonReplyData      `json:"button_reply,omitempty"`
	ListReply      *ListReplyData        `json:"list_reply,omitempty"`
}

// ButtonReplyData contiene datos de respuesta de botón
type ButtonReplyData struct {
	ID    string `json:"id"`    // ID del botón
	Title string `json:"title"` // Título del botón
}

// ListReplyData contiene datos de respuesta de lista
type ListReplyData struct {
	ID          string `json:"id"`          // ID de la opción seleccionada
	Title       string `json:"title"`       // Título de la opción
	Description string `json:"description"` // Descripción de la opción
}

// MessageContext representa el contexto de un mensaje (reply)
type MessageContext struct {
	From string `json:"from"` // ID del mensaje al que responde
	ID   string `json:"id"`   // ID del mensaje original
}

// WebhookStatus representa un cambio de estado de mensaje
type WebhookStatus struct {
	ID           string              `json:"id"`           // ID del mensaje
	Status       string              `json:"status"`       // "sent", "delivered", "read", "failed"
	Timestamp    string              `json:"timestamp"`    // Unix timestamp
	RecipientID  string              `json:"recipient_id"` // Número de teléfono del destinatario
	Conversation *ConversationInfo   `json:"conversation,omitempty"`
	Pricing      *PricingInfo        `json:"pricing,omitempty"`
	Errors       []WebhookError      `json:"errors,omitempty"`
}

// ConversationInfo contiene información de la conversación
type ConversationInfo struct {
	ID               string `json:"id"`                 // ID de la conversación
	Origin           string `json:"origin"`             // "business_initiated", "user_initiated"
	ExpirationTimestamp string `json:"expiration_timestamp,omitempty"` // Timestamp de expiración
}

// PricingInfo contiene información de precio del mensaje
type PricingInfo struct {
	Billable     bool   `json:"billable"`      // Si es facturable
	PricingModel string `json:"pricing_model"` // "CBP" (Conversation-Based Pricing)
	Category     string `json:"category"`      // "business_initiated", "user_initiated"
}

// WebhookError representa un error en el webhook
type WebhookError struct {
	Code    int    `json:"code"`    // Código de error
	Title   string `json:"title"`   // Título del error
	Message string `json:"message"` // Mensaje de error
	Details string `json:"details"` // Detalles adicionales
}

// GetMessageText extrae el texto del mensaje independientemente del tipo
func (m *WebhookMessage) GetMessageText() string {
	switch m.Type {
	case "text":
		if m.Text != nil {
			return m.Text.Body
		}
	case "button":
		if m.Button != nil {
			return m.Button.Text
		}
	case "interactive":
		if m.Interactive != nil {
			if m.Interactive.ButtonReply != nil {
				return m.Interactive.ButtonReply.Title
			}
			if m.Interactive.ListReply != nil {
				return m.Interactive.ListReply.Title
			}
		}
	}
	return ""
}

// IsButtonResponse verifica si el mensaje es una respuesta de botón
func (m *WebhookMessage) IsButtonResponse() bool {
	return m.Type == "button" || (m.Type == "interactive" && m.Interactive != nil && m.Interactive.ButtonReply != nil)
}
