package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/store"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/webhook"
)

const samplePDF = "%PDF-1.4\n1 0 obj<< /Type /Catalog >>endobj\ntrailer<< /Root 1 0 R >>\n%%EOF\n"

type textPayload struct {
	Body string `json:"body"`
}

type mediaPayload struct {
	ID       string `json:"id"`
	Link     string `json:"link"`
	Caption  string `json:"caption"`
	Filename string `json:"filename"`
}

type inboundMediaRequest struct {
	Phone    string `json:"phone"`
	Type     string `json:"type"`
	Caption  string `json:"caption"`
	Filename string `json:"filename"`
}

func (h *Handler) sendNonTemplate(c *gin.Context, phoneNumberID string, body sendMessageRequest) {
	if body.To == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "falta el destinatario", "code": 100}})
		return
	}

	send := &store.Send{
		MessageID:     newMessageID(),
		PhoneNumberID: phoneNumberID,
		To:            body.To,
		Type:          body.Type,
		SentAt:        time.Now(),
	}

	var media *mediaPayload
	switch body.Type {
	case "text":
		if body.Text == nil || body.Text.Body == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "falta el texto", "code": 100}})
			return
		}
		send.Text = body.Text.Body
	case "image":
		media = body.Image
	case "document":
		media = body.Document
	case "audio":
		media = body.Audio
	case "video":
		media = body.Video
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "tipo de mensaje no soportado por el mock: " + body.Type, "code": 100}})
		return
	}

	if body.Type != "text" {
		if media == nil || (media.ID == "" && media.Link == "") {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "falta el archivo del mensaje", "code": 100}})
			return
		}
		if media.ID != "" && h.store.Media(media.ID) == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "el media id no existe en el mock", "code": 100}})
			return
		}
		send.MediaType = body.Type
		send.MediaID = media.ID
		send.Caption = media.Caption
		send.Filename = media.Filename
	}

	h.store.AddSend(send)
	h.logger.Info().Str("type", send.Type).Str("to", send.To).Str("message_id", send.MessageID).Msg("mock: mensaje libre enviado")

	c.JSON(http.StatusOK, gin.H{
		"messaging_product": "whatsapp",
		"contacts":          []gin.H{{"input": body.To, "wa_id": strings.TrimPrefix(body.To, "+")}},
		"messages":          []gin.H{{"id": send.MessageID}},
	})
}

func (h *Handler) UploadMedia(c *gin.Context) {
	if c.PostForm("messaging_product") != "whatsapp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "messaging_product debe ser whatsapp", "code": 100}})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "falta el archivo", "code": 100}})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error(), "code": 100}})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error(), "code": 100}})
		return
	}

	mimeType := c.PostForm("type")
	if mimeType == "" {
		mimeType = fileHeader.Header.Get("Content-Type")
	}

	media := &store.Media{
		ID:       newMediaID(),
		MimeType: mimeType,
		Filename: fileHeader.Filename,
		Size:     len(data),
		Data:     data,
	}
	h.store.AddMedia(media)
	h.logger.Info().Str("media_id", media.ID).Str("mime", mimeType).Int("size", len(data)).Msg("mock: archivo subido")

	c.JSON(http.StatusOK, gin.H{"id": media.ID})
}

func (h *Handler) GetMediaInfo(c *gin.Context) {
	media := h.store.Media(c.Param("id"))
	if media == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "objeto no encontrado", "code": 100}})
		return
	}

	sum := sha256.Sum256(media.Data)
	c.JSON(http.StatusOK, gin.H{
		"messaging_product": "whatsapp",
		"url":               "http://" + c.Request.Host + "/_mock/media/" + media.ID,
		"mime_type":         media.MimeType,
		"sha256":            hex.EncodeToString(sum[:]),
		"file_size":         media.Size,
		"id":                media.ID,
	})
}

func (h *Handler) DownloadMedia(c *gin.Context) {
	if !strings.HasPrefix(c.GetHeader("Authorization"), "Bearer ") {
		c.Status(http.StatusUnauthorized)
		return
	}
	media := h.store.Media(c.Param("id"))
	if media == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Data(http.StatusOK, media.MimeType, media.Data)
}

func (h *Handler) ListMedia(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": h.store.MediaList()})
}

func (h *Handler) InboundMedia(c *gin.Context) {
	var body inboundMediaRequest
	if err := c.ShouldBindJSON(&body); err != nil || body.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone es obligatorio"})
		return
	}

	send := h.store.LastSendTo(body.Phone)
	if send == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no hay ningun mensaje enviado a ese numero"})
		return
	}

	media := &store.Media{ID: newMediaID(), Filename: body.Filename}
	switch body.Type {
	case "document":
		media.MimeType = "application/pdf"
		media.Data = []byte(samplePDF)
		if media.Filename == "" {
			media.Filename = "comprobante.pdf"
		}
	default:
		body.Type = "image"
		media.MimeType = "image/png"
		media.Data = samplePNG()
		media.Filename = ""
	}
	media.Size = len(media.Data)
	h.store.AddMedia(media)

	sum := sha256.Sum256(media.Data)
	content := &webhook.Media{
		ID:       media.ID,
		MimeType: media.MimeType,
		SHA256:   hex.EncodeToString(sum[:]),
		Caption:  body.Caption,
		Filename: media.Filename,
	}

	message := webhook.Message{
		From:      body.Phone,
		ID:        newMessageID(),
		Timestamp: webhook.Now(),
		Type:      body.Type,
	}
	if body.Type == "document" {
		message.Document = content
	} else {
		message.Image = content
	}

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
						WaID:    body.Phone,
					}},
					Messages: []webhook.Message{message},
				},
			}},
		}},
	}

	if err := h.webhook.Send(payload); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"media_id": media.ID, "type": body.Type, "message_id": message.ID})
}

func newMediaID() string {
	return "mock_media_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:16]
}

func samplePNG() []byte {
	const width, height = 320, 200
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			canvas.Set(x, y, color.RGBA{R: uint8(40 + x*180/width), G: uint8(120 + y*100/height), B: 200, A: 255})
		}
	}
	for y := 70; y < 130; y++ {
		for x := 130; x < 190; x++ {
			canvas.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, canvas)
	return buf.Bytes()
}
