package request

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

	Event                   string `json:"event,omitempty"`
	MessageTemplateID       int64  `json:"message_template_id,omitempty"`
	MessageTemplateName     string `json:"message_template_name,omitempty"`
	MessageTemplateLanguage string `json:"message_template_language,omitempty"`
	Reason                  string `json:"reason,omitempty"`
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
	From        string               `json:"from"`
	ID          string               `json:"id"`
	Timestamp   string               `json:"timestamp"`
	Type        string               `json:"type"`
	Text        *TextContent         `json:"text,omitempty"`
	Button      *ButtonResponse      `json:"button,omitempty"`
	Interactive *InteractiveResponse `json:"interactive,omitempty"`
	Context     *MessageContext      `json:"context,omitempty"`
}

type TextContent struct {
	Body string `json:"body"`
}

type ButtonResponse struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

type InteractiveResponse struct {
	Type        string           `json:"type"`
	ButtonReply *ButtonReplyData `json:"button_reply,omitempty"`
	ListReply   *ListReplyData   `json:"list_reply,omitempty"`
}

type ButtonReplyData struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ListReplyData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MessageContext struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

type WebhookStatus struct {
	ID                 string            `json:"id"`
	Status             string            `json:"status"`
	Timestamp          string            `json:"timestamp"`
	RecipientID        string            `json:"recipient_id"`
	RecipientLogicalID string            `json:"recipient_logical_id"`
	Conversation       *ConversationInfo `json:"conversation,omitempty"`
	Pricing            *PricingInfo      `json:"pricing,omitempty"`
	Errors             []WebhookError    `json:"errors,omitempty"`
}

type ConversationInfo struct {
	ID                  string             `json:"id"`
	Origin              ConversationOrigin `json:"origin"`
	ExpirationTimestamp string             `json:"expiration_timestamp,omitempty"`
}

type ConversationOrigin struct {
	Type string `json:"type"`
}

type PricingInfo struct {
	Billable     bool   `json:"billable"`
	PricingModel string `json:"pricing_model"`
	Category     string `json:"category"`
	Type         string `json:"type"`
}

type WebhookError struct {
	Code    int    `json:"code"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Details string `json:"details"`
}

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

func (m *WebhookMessage) IsButtonResponse() bool {
	return m.Type == "button" || (m.Type == "interactive" && m.Interactive != nil && m.Interactive.ButtonReply != nil)
}
