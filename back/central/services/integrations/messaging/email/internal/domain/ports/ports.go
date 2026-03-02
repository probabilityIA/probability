package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/entities"
)

// IEmailClient envía emails HTML via el servicio compartido (shared/email)
type IEmailClient interface {
	SendHTML(ctx context.Context, to, subject, html string) error
}

// IResultPublisher publica resultados de entrega a RabbitMQ
// para que notification_config los persista en email_logs
type IResultPublisher interface {
	PublishResult(ctx context.Context, result *entities.DeliveryResult) error
}

// IEmailUseCase caso de uso para enviar notificaciones por email
type IEmailUseCase interface {
	SendNotificationEmail(ctx context.Context, dto dtos.SendEmailDTO) error
}
