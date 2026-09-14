package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
)

const chatMediaURLTTL = 30 * time.Minute

func (uc *useCase) SetChatMediaSigner(signer ports.IChatMediaSigner) {
	uc.chatMediaSigner = signer
}

func (uc *useCase) messageMediaDTO(ctx context.Context, media *entities.MessageMedia) *dtos.MessageMediaResponseDTO {
	if media == nil {
		return nil
	}

	dto := &dtos.MessageMediaResponseDTO{
		Type:     media.Type,
		MimeType: media.Mime,
		Filename: media.Filename,
		Size:     media.Size,
	}

	if media.Key == "" || uc.chatMediaSigner == nil {
		return dto
	}

	url, err := uc.chatMediaSigner.PresignGet(ctx, media.Key, media.Filename, chatMediaURLTTL)
	if err != nil {
		uc.logger.Warn().Err(err).Str("key", media.Key).Msg("No se pudo firmar el enlace del adjunto")
		return dto
	}

	dto.URL = url
	dto.Available = true
	return dto
}
