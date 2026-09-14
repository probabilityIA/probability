package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/store"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/webhook"
)

type templateLanguage struct {
	Code string `json:"code"`
}

type templateParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type templateComponent struct {
	Type       string              `json:"type"`
	SubType    string              `json:"sub_type"`
	Index      string              `json:"index"`
	Parameters []templateParameter `json:"parameters"`
}

type templatePayload struct {
	Name       string              `json:"name"`
	Language   templateLanguage    `json:"language"`
	Components []templateComponent `json:"components"`
}

type sendMessageRequest struct {
	MessagingProduct string           `json:"messaging_product"`
	To               string           `json:"to"`
	Type             string           `json:"type"`
	Template         *templatePayload `json:"template"`
	Text             *textPayload     `json:"text"`
	Image            *mediaPayload    `json:"image"`
	Document         *mediaPayload    `json:"document"`
	Audio            *mediaPayload    `json:"audio"`
	Video            *mediaPayload    `json:"video"`
}

func newMessageID() string {
	return "wamid.MOCK" + strings.ReplaceAll(uuid.New().String(), "-", "")[:24]
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mock": "whatsapp-graph"})
}

func (h *Handler) SendMessage(c *gin.Context) {
	phoneNumberID := c.Param("id")

	var body sendMessageRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error(), "code": 100}})
		return
	}

	if body.Template == nil {
		h.sendNonTemplate(c, phoneNumberID, body)
		return
	}

	if body.To == "" || body.Template.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "falta el destinatario o la plantilla", "code": 100}})
		return
	}

	messageID := newMessageID()

	parameters := []string{}
	for _, component := range body.Template.Components {
		if strings.EqualFold(component.Type, "body") {
			for _, parameter := range component.Parameters {
				parameters = append(parameters, parameter.Text)
			}
		}
	}

	send := &store.Send{
		MessageID:     messageID,
		PhoneNumberID: phoneNumberID,
		To:            body.To,
		TemplateName:  body.Template.Name,
		Language:      body.Template.Language.Code,
		Parameters:    parameters,
		SentAt:        time.Now(),
	}
	h.store.AddSend(send)

	h.logger.Info().
		Str("template", send.TemplateName).
		Str("to", send.To).
		Str("message_id", messageID).
		Msg("mock: plantilla enviada")

	c.JSON(http.StatusOK, gin.H{
		"messaging_product": "whatsapp",
		"contacts":          []gin.H{{"input": body.To, "wa_id": strings.TrimPrefix(body.To, "+")}},
		"messages":          []gin.H{{"id": messageID}},
	})

	go h.autoReply(send)
}

func (h *Handler) autoReply(send *store.Send) {
	reply, ok := h.store.ReplyFor(send.TemplateName)
	if !ok {
		h.logger.Info().
			Str("template", send.TemplateName).
			Msg("mock: la plantilla no tiene respuesta programada, la rama termina aca")
		return
	}

	time.Sleep(h.replyDelay)

	h.sendStatus(send, "delivered")
	time.Sleep(150 * time.Millisecond)
	h.sendStatus(send, "read")
	time.Sleep(150 * time.Millisecond)

	h.sendButtonReply(send, reply)
}

func (h *Handler) sendStatus(send *store.Send, status string) {
	payload := webhook.Payload{
		Object: "whatsapp_business_account",
		Entry: []webhook.Entry{{
			ID: send.PhoneNumberID,
			Changes: []webhook.Change{{
				Field: "messages",
				Value: webhook.Value{
					MessagingProduct: "whatsapp",
					Metadata: &webhook.Metadata{
						DisplayPhoneNumber: h.displayPhone,
						PhoneNumberID:      send.PhoneNumberID,
					},
					Statuses: []webhook.Status{{
						ID:          send.MessageID,
						Status:      status,
						Timestamp:   webhook.Now(),
						RecipientID: send.To,
					}},
				},
			}},
		}},
	}

	if err := h.webhook.Send(payload); err != nil {
		h.logger.Error().Err(err).Str("status", status).Msg("mock: error enviando el estado")
	}
}

func (h *Handler) sendButtonReply(send *store.Send, buttonText string) {
	payload := webhook.Payload{
		Object: "whatsapp_business_account",
		Entry: []webhook.Entry{{
			ID: send.PhoneNumberID,
			Changes: []webhook.Change{{
				Field: "messages",
				Value: webhook.Value{
					MessagingProduct: "whatsapp",
					Metadata: &webhook.Metadata{
						DisplayPhoneNumber: h.displayPhone,
						PhoneNumberID:      send.PhoneNumberID,
					},
					Contacts: []webhook.Contact{{
						Profile: webhook.Profile{Name: "Cliente Mock"},
						WaID:    send.To,
					}},
					Messages: []webhook.Message{{
						From:      send.To,
						ID:        newMessageID(),
						Timestamp: webhook.Now(),
						Type:      "button",
						Button: &webhook.Button{
							Payload: buttonText,
							Text:    buttonText,
						},
						Context: &webhook.Context{
							From: h.displayPhone,
							ID:   send.MessageID,
						},
					}},
				},
			}},
		}},
	}

	if err := h.webhook.Send(payload); err != nil {
		h.logger.Error().Err(err).Str("button", buttonText).Msg("mock: error enviando la respuesta de boton")
		return
	}

	h.store.MarkReplied(send.MessageID, buttonText)

	h.logger.Info().
		Str("template", send.TemplateName).
		Str("button", buttonText).
		Str("context_message_id", send.MessageID).
		Msg("mock: el cliente respondio")
}

type createTemplateRequest struct {
	Name       string           `json:"name"`
	Language   string           `json:"language"`
	Category   string           `json:"category"`
	Components []map[string]any `json:"components"`
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	wabaID := c.Param("id")

	var body createTemplateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error(), "code": 100}})
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "la plantilla necesita nombre", "code": 100}})
		return
	}

	if !hasBodyComponent(body.Components) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message":        "Invalid parameter",
			"error_user_msg": "components: la plantilla necesita un componente BODY",
			"code":           100,
			"error_subcode":  2388043,
		}})
		return
	}

	metaID := fmt.Sprintf("mock_tpl_%d", time.Now().UnixNano())

	template := &store.Template{
		MetaID:   metaID,
		WABAID:   wabaID,
		Name:     body.Name,
		Language: body.Language,
		Status:   "PENDING",
		SentAt:   time.Now(),
	}
	h.store.AddTemplate(template)

	h.logger.Info().
		Str("name", body.Name).
		Str("waba_id", wabaID).
		Str("meta_id", metaID).
		Msg("mock: plantilla creada, se aprueba en un momento")

	c.JSON(http.StatusOK, gin.H{"id": metaID, "status": "PENDING", "category": body.Category})

	go h.approveTemplate(template)
}

func (h *Handler) approveTemplate(template *store.Template) {
	time.Sleep(h.approveDelay)
	h.emitTemplateStatus(template, "APPROVED", "")
}

func (h *Handler) emitTemplateStatus(template *store.Template, event, reason string) {
	payload := webhook.Payload{
		Object: "whatsapp_business_account",
		Entry: []webhook.Entry{{
			ID: template.WABAID,
			Changes: []webhook.Change{{
				Field: "message_template_status_update",
				Value: webhook.Value{
					Event:                   event,
					MessageTemplateName:     template.Name,
					MessageTemplateLanguage: template.Language,
					Reason:                  reason,
				},
			}},
		}},
	}

	if err := h.webhook.Send(payload); err != nil {
		h.logger.Error().Err(err).Str("name", template.Name).Msg("mock: error notificando el estado de la plantilla")
		return
	}

	template.Status = event

	h.logger.Info().Str("name", template.Name).Str("event", event).Msg("mock: estado de plantilla notificado")
}

func (h *Handler) ListMetaTemplates(c *gin.Context) {
	data := []gin.H{}
	for _, template := range h.store.Templates() {
		data = append(data, gin.H{
			"id":       template.MetaID,
			"name":     template.Name,
			"language": template.Language,
			"status":   template.Status,
			"category": "MARKETING",
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) DeleteMetaTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) StartUpload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": "upload:mock" + strings.ReplaceAll(uuid.New().String(), "-", "")[:12]})
}

func (h *Handler) FinishUpload(c *gin.Context) {
	handle := "mock_handle_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
	c.JSON(http.StatusOK, gin.H{"h": handle, "id": handle})
}

func hasBodyComponent(components []map[string]any) bool {
	for _, component := range components {
		componentType, _ := component["type"].(string)
		text, _ := component["text"].(string)
		if strings.EqualFold(componentType, "BODY") && strings.TrimSpace(text) != "" {
			return true
		}
	}
	return false
}
