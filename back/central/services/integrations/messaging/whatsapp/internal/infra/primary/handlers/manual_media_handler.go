package handlers

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

func (h *handler) SendManualMedia(c *gin.Context) {
	ctx := c.Request.Context()
	conversationID := c.Param("id")

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "contexto de negocio no encontrado"})
		return
	}
	if businessID == 0 {
		parsed, err := strconv.ParseUint(c.PostForm("business_id"), 10, 64)
		if err != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "business_id es requerido para super admin"})
			return
		}
		businessID = uint(parsed)
	}

	phoneNumber := strings.TrimSpace(c.PostForm("phone_number"))
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "phone_number es obligatorio"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "adjunta un archivo"})
		return
	}
	if fileHeader.Size > entities.MaxAttachmentBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file_too_large", "message": "el archivo supera 10 MB"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "no se pudo leer el archivo"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, entities.MaxAttachmentBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "no se pudo leer el archivo"})
		return
	}

	sentBy := ""
	if userID, exists := middleware.GetUserID(c); exists {
		sentBy = strconv.FormatUint(uint64(userID), 10)
	}

	media := entities.OutboundMedia{
		Filename: fileHeader.Filename,
		MimeType: detectUploadMime(fileHeader.Header.Get("Content-Type"), fileHeader.Filename, data),
		Data:     data,
	}

	messageID, err := h.useCase.SendManualMedia(ctx, conversationID, phoneNumber, businessID, media, c.PostForm("caption"), sentBy)
	if err != nil {
		status := http.StatusInternalServerError
		code := "send_failed"
		if errors.Is(err, entities.ErrMediaNotAllowed) {
			status = http.StatusBadRequest
			code = "invalid_media"
		}
		h.log.Error(ctx).Err(err).Str("conversation_id", conversationID).Msg("[ManualMedia Handler] - error enviando archivo")
		c.JSON(status, gin.H{"error": code, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message_id": messageID, "status": "sent"})
}

func detectUploadMime(declared, filename string, data []byte) string {
	normalized := entities.NormalizeMime(declared)
	if normalized != "" && normalized != "application/octet-stream" {
		return normalized
	}
	if byExtension := entities.NormalizeMime(mime.TypeByExtension(strings.ToLower(filepath.Ext(filename)))); byExtension != "" {
		return byExtension
	}
	return entities.NormalizeMime(http.DetectContentType(data))
}
