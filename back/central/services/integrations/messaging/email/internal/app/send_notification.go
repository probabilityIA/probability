package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/errors"
)

// SendNotificationEmail envía un email de notificación y publica el resultado
func (uc *useCase) SendNotificationEmail(ctx context.Context, dto dtos.SendEmailDTO) error {
	if dto.CustomerEmail == "" {
		uc.logger.Warn(ctx).
			Str("event_type", dto.EventType).
			Uint("business_id", dto.BusinessID).
			Msg("Email de destinatario vacío, saltando notificación")
		return domainerrors.ErrMissingRecipient
	}

	subject := buildSubject(dto.EventType)
	html := buildHTML(dto.EventType, dto.EventData)

	result := &entities.DeliveryResult{
		Channel:       "email",
		BusinessID:    dto.BusinessID,
		IntegrationID: dto.IntegrationID,
		ConfigID:      dto.ConfigID,
		To:            dto.CustomerEmail,
		Subject:       subject,
		EventType:     dto.EventType,
		SentAt:        time.Now(),
	}

	err := uc.emailClient.SendHTML(ctx, dto.CustomerEmail, subject, html)
	if err != nil {
		result.Status = "failed"
		result.ErrorMessage = err.Error()

		uc.logger.Error(ctx).
			Err(err).
			Str("to", dto.CustomerEmail).
			Str("event_type", dto.EventType).
			Uint("config_id", dto.ConfigID).
			Msg("Error enviando email de notificación")
	} else {
		result.Status = "sent"

		uc.logger.Info(ctx).
			Str("to", dto.CustomerEmail).
			Str("event_type", dto.EventType).
			Uint("config_id", dto.ConfigID).
			Msg("Email de notificación enviado")
	}

	// Publicar resultado a notification_config (best-effort)
	if pubErr := uc.resultPub.PublishResult(ctx, result); pubErr != nil {
		uc.logger.Error(ctx).
			Err(pubErr).
			Msg("Error publicando resultado de entrega de email")
	}

	return err
}
