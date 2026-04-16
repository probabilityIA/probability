package dtos

// WebhookPayloadDTO representa el payload del webhook en el dominio (sin tags JSON)
type WebhookPayloadDTO struct {
	Object string
	Entry  []WebhookEntryDTO
}

// WebhookEntryDTO representa una entrada en el webhook
type WebhookEntryDTO struct {
	ID      string
	Changes []WebhookChangeDTO
}

// WebhookChangeDTO representa un cambio en el webhook
type WebhookChangeDTO struct {
	Field string // "messages" o "message_template_status_update"
	Value WebhookValueDTO
}

// WebhookValueDTO contiene los datos del webhook
type WebhookValueDTO struct {
	MessagingProduct string
	Metadata         WebhookMetadataDTO
	Contacts         []WebhookContactDTO
	Messages         []WebhookMessageDTO
	Statuses         []WebhookStatusDTO
}

// WebhookMetadataDTO contiene metadata del webhook
type WebhookMetadataDTO struct {
	DisplayPhoneNumber string
	PhoneNumberID      string
}

// WebhookContactDTO representa información de contacto
type WebhookContactDTO struct {
	Profile WebhookProfileDTO
	WaID    string
}

// WebhookProfileDTO contiene el perfil del contacto
type WebhookProfileDTO struct {
	Name string
}

// WebhookMessageDTO representa un mensaje recibido (dominio puro)
type WebhookMessageDTO struct {
	From        string
	ID          string
	Timestamp   string
	Type        string // "text", "button", "interactive"
	Text        *TextContentDTO
	Button      *ButtonResponseDTO
	Interactive *InteractiveResponseDTO
	Context     *MessageContextDTO
}

// TextContentDTO representa contenido de texto
type TextContentDTO struct {
	Body string
}

// ButtonResponseDTO representa la respuesta a un botón quick reply
type ButtonResponseDTO struct {
	Payload string
	Text    string
}

// InteractiveResponseDTO representa respuesta a mensaje interactivo
type InteractiveResponseDTO struct {
	Type        string
	ButtonReply *ButtonReplyDataDTO
	ListReply   *ListReplyDataDTO
}

// ButtonReplyDataDTO contiene datos de respuesta de botón
type ButtonReplyDataDTO struct {
	ID    string
	Title string
}

// ListReplyDataDTO contiene datos de respuesta de lista
type ListReplyDataDTO struct {
	ID          string
	Title       string
	Description string
}

// MessageContextDTO representa el contexto de un mensaje (reply)
type MessageContextDTO struct {
	From string
	ID   string
}

// WebhookStatusDTO representa un cambio de estado de mensaje
type WebhookStatusDTO struct {
	ID           string
	Status       string // "sent", "delivered", "read", "failed"
	Timestamp    string
	RecipientID  string
	Conversation *ConversationInfoDTO
	Pricing      *PricingInfoDTO
	Errors       []WebhookErrorDTO
}

// ConversationInfoDTO contiene información de la conversación
type ConversationInfoDTO struct {
	ID                  string
	Origin              string
	ExpirationTimestamp string
}

// PricingInfoDTO contiene información de precio del mensaje
type PricingInfoDTO struct {
	Billable     bool
	PricingModel string
	Category     string
}

// WebhookErrorDTO representa un error en el webhook
type WebhookErrorDTO struct {
	Code    int
	Title   string
	Message string
	Details string
}

// GetMessageText extrae el texto del mensaje independientemente del tipo
func (m *WebhookMessageDTO) GetMessageText() string {
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
func (m *WebhookMessageDTO) IsButtonResponse() bool {
	return m.Type == "button" || (m.Type == "interactive" && m.Interactive != nil && m.Interactive.ButtonReply != nil)
}
