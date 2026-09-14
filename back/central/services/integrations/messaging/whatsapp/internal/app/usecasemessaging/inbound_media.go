package usecasemessaging

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

func (u *usecases) attachInboundMedia(ctx context.Context, messageLog *entities.MessageLog, message dtos.WebhookMessageDTO, businessID uint) {
	media := message.GetMedia()
	if media == nil {
		return
	}

	messageLog.MediaType = message.Type
	messageLog.MediaMime = entities.NormalizeMime(media.MimeType)
	messageLog.MediaFilename = entities.SafeMediaFilename(media.Filename)

	if u.mediaFactory == nil || u.chatMedia == nil || media.ID == "" || businessID == 0 {
		return
	}

	whatsappConfig, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		u.log.Warn(ctx).Err(err).Uint("business_id", businessID).Msg("[WhatsApp Webhook] - sin credenciales para descargar el archivo del cliente")
		return
	}

	api := u.mediaFactory(whatsappConfig.WhatsAppURL)

	info, err := api.GetMediaInfo(ctx, media.ID, whatsappConfig.AccessToken)
	if err != nil {
		u.log.Warn(ctx).Err(err).Str("media_id", media.ID).Msg("[WhatsApp Webhook] - no se pudo consultar el archivo del cliente")
		return
	}

	data, err := api.DownloadMedia(ctx, info.URL, whatsappConfig.AccessToken)
	if err != nil {
		u.log.Warn(ctx).Err(err).Str("media_id", media.ID).Msg("[WhatsApp Webhook] - no se pudo descargar el archivo del cliente")
		return
	}

	mimeType := messageLog.MediaMime
	if normalized := entities.NormalizeMime(info.MimeType); normalized != "" {
		mimeType = normalized
	}

	key := entities.ChatMediaKey(businessID, mimeType, messageLog.MediaFilename, uuid.NewString(), time.Now())
	if err := u.chatMedia.Put(ctx, key, data, mimeType); err != nil {
		u.log.Warn(ctx).Err(err).Str("media_id", media.ID).Msg("[WhatsApp Webhook] - no se pudo guardar el archivo del cliente")
		return
	}

	messageLog.MediaKey = key
	messageLog.MediaMime = mimeType
	messageLog.MediaSize = int64(len(data))
}

func messagePreview(messageText string, message dtos.WebhookMessageDTO) string {
	if messageText != "" {
		return messageText
	}
	return entities.MediaPreviewLabel(message.Type)
}
