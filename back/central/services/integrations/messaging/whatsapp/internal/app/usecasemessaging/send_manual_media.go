package usecasemessaging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

type MediaAPIFactory func(baseURL string) ports.IMediaAPI

func (u *usecases) SetMediaDependencies(factory MediaAPIFactory, storage ports.IChatMediaStorage) {
	u.mediaFactory = factory
	u.chatMedia = storage
}

func (u *usecases) SendManualMedia(
	ctx context.Context,
	conversationID string,
	phoneNumber string,
	businessID uint,
	media entities.OutboundMedia,
	caption string,
	sentBy string,
) (string, error) {
	if u.mediaFactory == nil || u.chatMedia == nil {
		return "", fmt.Errorf("el envio de archivos no esta disponible")
	}

	mimeType := entities.NormalizeMime(media.MimeType)
	mediaType, err := entities.ClassifyOutboundMedia(mimeType, len(media.Data))
	if err != nil {
		return "", err
	}

	caption = strings.TrimSpace(caption)
	if len([]rune(caption)) > entities.MaxCaptionLength {
		return "", fmt.Errorf("%w: el texto del archivo supera %d caracteres", entities.ErrMediaNotAllowed, entities.MaxCaptionLength)
	}

	phoneNumber = NormalizePhoneNumber(phoneNumber)
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		return "", fmt.Errorf("numero de telefono invalido: %w", err)
	}

	whatsappConfig, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		return "", fmt.Errorf("error obteniendo configuracion de WhatsApp: %w", err)
	}

	filename := entities.SafeMediaFilename(media.Filename)
	key := entities.ChatMediaKey(businessID, mimeType, filename, uuid.NewString(), time.Now())

	if err := u.chatMedia.Put(ctx, key, media.Data, mimeType); err != nil {
		return "", fmt.Errorf("no se pudo guardar el archivo: %w", err)
	}

	api := u.mediaFactory(whatsappConfig.WhatsAppURL)

	mediaID, err := api.UploadMedia(ctx, whatsappConfig.PhoneNumberID, whatsappConfig.AccessToken, filename, mimeType, media.Data)
	if err != nil {
		u.discardChatMedia(ctx, key)
		return "", fmt.Errorf("error subiendo el archivo a WhatsApp: %w", err)
	}

	messageID, err := api.SendMediaMessage(ctx, whatsappConfig.PhoneNumberID, whatsappConfig.AccessToken, phoneNumber, mediaType, mediaID, caption, filename)
	if err != nil {
		u.discardChatMedia(ctx, key)
		return "", fmt.Errorf("error al enviar el archivo: %w", err)
	}

	if hsErr := u.conversationCache.ActivateHumanSession(ctx, phoneNumber, conversationID, businessID); hsErr != nil {
		u.log.Error(ctx).Err(hsErr).Str("conversation_id", conversationID).Msg("[WhatsApp UseCase] - error activando human session")
	}

	messageLog := &entities.MessageLog{
		ConversationID: conversationID,
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      messageID,
		Content:        caption,
		Status:         entities.MessageStatusSent,
		CreatedAt:      time.Now(),
		MediaType:      mediaType,
		MediaKey:       key,
		MediaMime:      mimeType,
		MediaFilename:  filename,
		MediaSize:      int64(len(media.Data)),
	}
	if err := u.persistPublisher.PublishMessageLogCreated(ctx, messageLog); err != nil {
		u.log.Error(ctx).Err(err).Str("message_id", messageID).Msg("[WhatsApp UseCase] - error publicando el archivo enviado en el log")
	}

	u.log.Info(ctx).
		Str("message_id", messageID).
		Str("conversation_id", conversationID).
		Str("media_type", mediaType).
		Str("sent_by", sentBy).
		Msg("[WhatsApp UseCase] - archivo enviado desde el dashboard")

	return messageID, nil
}

func (u *usecases) discardChatMedia(ctx context.Context, key string) {
	if err := u.chatMedia.Delete(ctx, key); err != nil {
		u.log.Warn(ctx).Err(err).Str("key", key).Msg("[WhatsApp UseCase] - no se pudo borrar el adjunto que no salio")
	}
}
