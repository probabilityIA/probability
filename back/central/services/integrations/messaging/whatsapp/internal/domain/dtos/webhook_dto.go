package dtos

type WebhookPayloadDTO struct {
	Object string
	Entry  []WebhookEntryDTO
}

type WebhookEntryDTO struct {
	ID      string
	Changes []WebhookChangeDTO
}

type WebhookChangeDTO struct {
	Field string
	Value WebhookValueDTO
}

type WebhookValueDTO struct {
	MessagingProduct string
	Metadata         WebhookMetadataDTO
	Contacts         []WebhookContactDTO
	Messages         []WebhookMessageDTO
	Statuses         []WebhookStatusDTO

	TemplateEvent    string
	TemplateID       string
	TemplateName     string
	TemplateLanguage string
	TemplateReason   string
}

type WebhookMetadataDTO struct {
	DisplayPhoneNumber string
	PhoneNumberID      string
}

type WebhookContactDTO struct {
	Profile WebhookProfileDTO
	WaID    string
}

type WebhookProfileDTO struct {
	Name string
}

type WebhookMessageDTO struct {
	From        string
	ID          string
	Timestamp   string
	Type        string
	Text        *TextContentDTO
	Button      *ButtonResponseDTO
	Interactive *InteractiveResponseDTO
	Context     *MessageContextDTO
}

type TextContentDTO struct {
	Body string
}

type ButtonResponseDTO struct {
	Payload string
	Text    string
}

type InteractiveResponseDTO struct {
	Type        string
	ButtonReply *ButtonReplyDataDTO
	ListReply   *ListReplyDataDTO
}

type ButtonReplyDataDTO struct {
	ID    string
	Title string
}

type ListReplyDataDTO struct {
	ID          string
	Title       string
	Description string
}

type MessageContextDTO struct {
	From string
	ID   string
}

type WebhookStatusDTO struct {
	ID           string
	Status       string
	Timestamp    string
	RecipientID  string
	Conversation *ConversationInfoDTO
	Pricing      *PricingInfoDTO
	Errors       []WebhookErrorDTO
}

type ConversationInfoDTO struct {
	ID                  string
	Origin              string
	ExpirationTimestamp string
}

type PricingInfoDTO struct {
	Billable     bool
	PricingModel string
	Category     string
}

type WebhookErrorDTO struct {
	Code    int
	Title   string
	Message string
	Details string
}

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

func (m *WebhookMessageDTO) IsButtonResponse() bool {
	return m.Type == "button" || (m.Type == "interactive" && m.Interactive != nil && m.Interactive.ButtonReply != nil)
}
