package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
)

type IRepository interface {
	CreateLead(ctx context.Context, lead *entities.MarketingLead) error
	SetWhatsAppMessageID(ctx context.Context, leadID uint, messageID string) error
}

// IWhatsAppSender es satisfecha estructuralmente por *whatsapp.Bundle
// (services/integrations/messaging/whatsapp), sin importar su paquete interno —
// mismo patron ya usado por publicsite.IBoldGateway con *pay.Bundle.
type IWhatsAppSender interface {
	SendDiagnosticResultTemplate(ctx context.Context, phone, name, surveyTitle, level string) (string, error)
}
