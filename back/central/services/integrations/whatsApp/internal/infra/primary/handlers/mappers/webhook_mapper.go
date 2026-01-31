package mappers

import (
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/handlers/request"
)

// WebhookPayloadToDomain convierte request.WebhookPayload → dtos.WebhookPayloadDTO
func WebhookPayloadToDomain(req request.WebhookPayload) dtos.WebhookPayloadDTO {
	entries := make([]dtos.WebhookEntryDTO, len(req.Entry))
	for i, entry := range req.Entry {
		entries[i] = mapWebhookEntry(entry)
	}

	return dtos.WebhookPayloadDTO{
		Object: req.Object,
		Entry:  entries,
	}
}

func mapWebhookEntry(req request.WebhookEntry) dtos.WebhookEntryDTO {
	changes := make([]dtos.WebhookChangeDTO, len(req.Changes))
	for i, change := range req.Changes {
		changes[i] = mapWebhookChange(change)
	}

	return dtos.WebhookEntryDTO{
		ID:      req.ID,
		Changes: changes,
	}
}

func mapWebhookChange(req request.WebhookChange) dtos.WebhookChangeDTO {
	return dtos.WebhookChangeDTO{
		Field: req.Field,
		Value: mapWebhookValue(req.Value),
	}
}

func mapWebhookValue(req request.WebhookValue) dtos.WebhookValueDTO {
	contacts := make([]dtos.WebhookContactDTO, len(req.Contacts))
	for i, contact := range req.Contacts {
		contacts[i] = mapWebhookContact(contact)
	}

	messages := make([]dtos.WebhookMessageDTO, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = mapWebhookMessage(msg)
	}

	statuses := make([]dtos.WebhookStatusDTO, len(req.Statuses))
	for i, status := range req.Statuses {
		statuses[i] = mapWebhookStatus(status)
	}

	return dtos.WebhookValueDTO{
		MessagingProduct: req.MessagingProduct,
		Metadata: dtos.WebhookMetadataDTO{
			DisplayPhoneNumber: req.Metadata.DisplayPhoneNumber,
			PhoneNumberID:      req.Metadata.PhoneNumberID,
		},
		Contacts: contacts,
		Messages: messages,
		Statuses: statuses,
	}
}

func mapWebhookContact(req request.WebhookContact) dtos.WebhookContactDTO {
	return dtos.WebhookContactDTO{
		Profile: dtos.WebhookProfileDTO{
			Name: req.Profile.Name,
		},
		WaID: req.WaID,
	}
}

func mapWebhookMessage(req request.WebhookMessage) dtos.WebhookMessageDTO {
	var text *dtos.TextContentDTO
	if req.Text != nil {
		text = &dtos.TextContentDTO{Body: req.Text.Body}
	}

	var button *dtos.ButtonResponseDTO
	if req.Button != nil {
		button = &dtos.ButtonResponseDTO{
			Payload: req.Button.Payload,
			Text:    req.Button.Text,
		}
	}

	var interactive *dtos.InteractiveResponseDTO
	if req.Interactive != nil {
		interactive = &dtos.InteractiveResponseDTO{
			Type: req.Interactive.Type,
		}
		if req.Interactive.ButtonReply != nil {
			interactive.ButtonReply = &dtos.ButtonReplyDataDTO{
				ID:    req.Interactive.ButtonReply.ID,
				Title: req.Interactive.ButtonReply.Title,
			}
		}
		if req.Interactive.ListReply != nil {
			interactive.ListReply = &dtos.ListReplyDataDTO{
				ID:          req.Interactive.ListReply.ID,
				Title:       req.Interactive.ListReply.Title,
				Description: req.Interactive.ListReply.Description,
			}
		}
	}

	var context *dtos.MessageContextDTO
	if req.Context != nil {
		context = &dtos.MessageContextDTO{
			From: req.Context.From,
			ID:   req.Context.ID,
		}
	}

	return dtos.WebhookMessageDTO{
		From:        req.From,
		ID:          req.ID,
		Timestamp:   req.Timestamp,
		Type:        req.Type,
		Text:        text,
		Button:      button,
		Interactive: interactive,
		Context:     context,
	}
}

func mapWebhookStatus(req request.WebhookStatus) dtos.WebhookStatusDTO {
	errors := make([]dtos.WebhookErrorDTO, len(req.Errors))
	for i, err := range req.Errors {
		errors[i] = dtos.WebhookErrorDTO{
			Code:    err.Code,
			Title:   err.Title,
			Message: err.Message,
			Details: err.Details,
		}
	}

	var conversation *dtos.ConversationInfoDTO
	if req.Conversation != nil {
		conversation = &dtos.ConversationInfoDTO{
			ID:                  req.Conversation.ID,
			Origin:              req.Conversation.Origin,
			ExpirationTimestamp: req.Conversation.ExpirationTimestamp,
		}
	}

	var pricing *dtos.PricingInfoDTO
	if req.Pricing != nil {
		pricing = &dtos.PricingInfoDTO{
			Billable:     req.Pricing.Billable,
			PricingModel: req.Pricing.PricingModel,
			Category:     req.Pricing.Category,
		}
	}

	return dtos.WebhookStatusDTO{
		ID:           req.ID,
		Status:       req.Status,
		Timestamp:    req.Timestamp,
		RecipientID:  req.RecipientID,
		Conversation: conversation,
		Pricing:      pricing,
		Errors:       errors,
	}
}
