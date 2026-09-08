package templates

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ISubmissionPublisher interface {
	PublishTemplateSubmission(ctx context.Context, message dtos.TemplateSubmissionMessage) error
}

type IUseCase interface {
	Create(ctx context.Context, dto dtos.CreateTemplateDTO) (*entities.WhatsappTemplate, error)
	Update(ctx context.Context, dto dtos.UpdateTemplateDTO) (*entities.WhatsappTemplate, error)
	GetByID(ctx context.Context, id, businessID uint) (*entities.WhatsappTemplate, error)
	List(ctx context.Context, businessID uint, scope, status string, page, pageSize int) ([]entities.WhatsappTemplate, int64, error)
	Delete(ctx context.Context, id, businessID uint) error
	Resubmit(ctx context.Context, id, businessID uint) (*entities.WhatsappTemplate, error)
	ApplySubmissionResult(ctx context.Context, result dtos.TemplateSubmissionResult) error
	ApplyMetaStatus(ctx context.Context, wabaID, name, language, event, reason string) error
	VariableCatalog() map[string]string
}

type useCase struct {
	repository ports.ITemplateRepository
	publisher  ISubmissionPublisher
	logger     log.ILogger
}

func New(repository ports.ITemplateRepository, publisher ISubmissionPublisher, logger log.ILogger) IUseCase {
	return &useCase{
		repository: repository,
		publisher:  publisher,
		logger:     logger.WithModule("whatsapp_templates_usecase"),
	}
}
