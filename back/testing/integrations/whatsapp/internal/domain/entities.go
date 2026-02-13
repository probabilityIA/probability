package domain

import "time"

// Conversation representa una conversación de WhatsApp simplificada para pruebas
type Conversation struct {
	ID           string
	PhoneNumber  string
	CurrentState string
	OrderNumber  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MessageLog representa un mensaje enviado o recibido
type MessageLog struct {
	ID             string
	ConversationID string
	Direction      string // "OUTBOUND" o "INBOUND"
	MessageType    string // "template", "button", "text"
	Content        string
	Status         string
	CreatedAt      time.Time
}

// TemplateRequest representa una petición para enviar un template
type TemplateRequest struct {
	TemplateName string
	PhoneNumber  string
	OrderNumber  string
	BusinessID   uint
	Variables    map[string]string
}

// WebhookPayload representa el payload de un webhook de WhatsApp (igual que Meta)
type WebhookPayload struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}

type WebhookEntry struct {
	ID      string          `json:"id"`
	Changes []WebhookChange `json:"changes"`
}

type WebhookChange struct {
	Value WebhookValue `json:"value"`
	Field string       `json:"field"`
}

type WebhookValue struct {
	MessagingProduct string           `json:"messaging_product"`
	Metadata         WebhookMetadata  `json:"metadata"`
	Contacts         []WebhookContact `json:"contacts,omitempty"`
	Messages         []WebhookMessage `json:"messages,omitempty"`
	Statuses         []WebhookStatus  `json:"statuses,omitempty"`
}

type WebhookMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type WebhookContact struct {
	Profile WebhookProfile `json:"profile"`
	WaID    string         `json:"wa_id"`
}

type WebhookProfile struct {
	Name string `json:"name"`
}

type WebhookMessage struct {
	From      string          `json:"from"`
	ID        string          `json:"id"`
	Timestamp string          `json:"timestamp"`
	Type      string          `json:"type"`
	Button    *ButtonResponse `json:"button,omitempty"`
	Text      *TextContent    `json:"text,omitempty"`
}

type ButtonResponse struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

type TextContent struct {
	Body string `json:"body"`
}

type WebhookStatus struct {
	ID           string            `json:"id"`
	Status       string            `json:"status"`
	Timestamp    string            `json:"timestamp"`
	RecipientID  string            `json:"recipient_id"`
	Conversation *ConversationInfo `json:"conversation,omitempty"`
}

type ConversationInfo struct {
	ID     string             `json:"id"`
	Origin ConversationOrigin `json:"origin"`
}

type ConversationOrigin struct {
	Type string `json:"type"`
}
