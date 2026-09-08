package usecasetemplates

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type TemplateAPIFactory func(baseURL string) ports.ITemplateAPI

type ProvisionResult struct {
	IntegrationID uint
	BusinessID    uint
	WABAID        string
	Created       []string
	AlreadyExists []string
	Skipped       []string
	Failed        map[string]string
	Templates     []ports.TemplateStatus
}

type IUseCase interface {
	Provision(ctx context.Context, businessID uint) (*ProvisionResult, error)
	GetStatus(ctx context.Context, businessID uint, refresh bool) (*ports.WABATemplatesSnapshot, error)
	HandleStatusUpdate(ctx context.Context, wabaID, name, language, event, reason string) error
	SubmitCustom(ctx context.Context, submission CustomTemplateSubmission) CustomTemplateResult
	DeleteCustom(ctx context.Context, submission CustomTemplateSubmission) error
	PreviewByEvent(ctx context.Context, businessID uint, eventCode string) ([]TemplatePreview, error)
}

type IStatusPublisher interface {
	PublishTemplateStatus(ctx context.Context, result CustomTemplateResult) error
}

type usecase struct {
	credentialsCache ports.ICredentialsCache
	templatesCache   ports.ITemplateStatusCache
	apiFactory       TemplateAPIFactory
	statusPublisher  IStatusPublisher
	log              log.ILogger
}

func New(
	credentialsCache ports.ICredentialsCache,
	templatesCache ports.ITemplateStatusCache,
	apiFactory TemplateAPIFactory,
	statusPublisher IStatusPublisher,
	logger log.ILogger,
) IUseCase {
	return &usecase{
		credentialsCache: credentialsCache,
		templatesCache:   templatesCache,
		apiFactory:       apiFactory,
		statusPublisher:  statusPublisher,
		log:              logger.WithModule("whatsapp-templates"),
	}
}
