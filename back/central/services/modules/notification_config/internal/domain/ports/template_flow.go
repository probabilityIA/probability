package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

type ITemplateFlowRepository interface {
	ListBySource(ctx context.Context, businessID, sourceTemplateID uint) ([]entities.TemplateFlow, error)
	ListByBusiness(ctx context.Context, businessID uint) ([]entities.TemplateFlow, error)
	ReplaceForSource(ctx context.Context, businessID, sourceTemplateID uint, flows []entities.TemplateFlow) error
	Resolve(ctx context.Context, businessID, sourceTemplateID uint, buttonText string) (*entities.TemplateFlow, error)
	TemplateNameByOutboundMessageID(ctx context.Context, messageID string) (string, error)
}
