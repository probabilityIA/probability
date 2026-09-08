package campaigns

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ICampaignPublisher interface {
	PublishCampaignSend(ctx context.Context, message dtos.CampaignSendMessage) error
}

type IUseCase interface {
	Create(ctx context.Context, dto dtos.CreateCampaignDTO) (*entities.Campaign, error)
	Update(ctx context.Context, dto dtos.UpdateCampaignDTO) (*entities.Campaign, error)
	GetByID(ctx context.Context, id, businessID uint) (*entities.Campaign, error)
	List(ctx context.Context, businessID uint, status string, page, pageSize int) ([]entities.Campaign, int64, error)
	Delete(ctx context.Context, id, businessID uint) error
	PreviewAudience(ctx context.Context, dto dtos.CreateCampaignDTO) (*dtos.CampaignAudiencePreviewDTO, error)
	Launch(ctx context.Context, id, businessID uint) (*entities.Campaign, error)
	Pause(ctx context.Context, id, businessID uint) (*entities.Campaign, error)
	Resume(ctx context.Context, id, businessID uint) (*entities.Campaign, error)
	Cancel(ctx context.Context, id, businessID uint) (*entities.Campaign, error)
	ListSends(ctx context.Context, campaignID, businessID uint, status string, page, pageSize int) ([]entities.CampaignSend, int64, error)
	RunDueCampaigns(ctx context.Context) error
	MarkSendResult(ctx context.Context, result dtos.CampaignSendResult) error
}

type useCase struct {
	campaigns ports.ICampaignRepository
	sends     ports.ICampaignSendRepository
	audience  ports.ICampaignAudienceQuerier
	sender    ports.ICampaignSenderQuerier
	templates ports.ITemplateRepository
	publisher ICampaignPublisher
	logger    log.ILogger
}

func New(
	campaigns ports.ICampaignRepository,
	sends ports.ICampaignSendRepository,
	audience ports.ICampaignAudienceQuerier,
	sender ports.ICampaignSenderQuerier,
	templates ports.ITemplateRepository,
	publisher ICampaignPublisher,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		campaigns: campaigns,
		sends:     sends,
		audience:  audience,
		sender:    sender,
		templates: templates,
		publisher: publisher,
		logger:    logger.WithModule("whatsapp_campaigns_usecase"),
	}
}
