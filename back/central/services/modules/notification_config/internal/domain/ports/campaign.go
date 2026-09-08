package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

type ICampaignRepository interface {
	CreateCampaign(ctx context.Context, campaign *entities.Campaign) error
	UpdateCampaign(ctx context.Context, campaign *entities.Campaign) error
	GetCampaignByID(ctx context.Context, id uint) (*entities.Campaign, error)
	ListCampaigns(ctx context.Context, businessID uint, status string, page, pageSize int) ([]entities.Campaign, int64, error)
	DeleteCampaign(ctx context.Context, id uint) error
	ListDueCampaigns(ctx context.Context, now time.Time, limit int) ([]entities.Campaign, error)
	UpdateCampaignCounters(ctx context.Context, campaignID uint) error
	MarkCampaignBatch(ctx context.Context, campaignID uint, lastBatchAt time.Time) error
}

type ICampaignSendRepository interface {
	BulkCreateSends(ctx context.Context, sends []entities.CampaignSend) (int, error)
	ListPendingSends(ctx context.Context, campaignID uint, limit int) ([]entities.CampaignSend, error)
	ListSends(ctx context.Context, campaignID uint, status string, page, pageSize int) ([]entities.CampaignSend, int64, error)
	MarkSendQueued(ctx context.Context, sendID uint, queuedAt time.Time) error
	MarkSendResult(ctx context.Context, sendID uint, status, messageID, errorMessage string) error
	CountSentSince(ctx context.Context, campaignID uint, since time.Time) (int64, error)
}

type ICampaignAudienceQuerier interface {
	FindCampaignCandidates(ctx context.Context, businessID uint, params entities.CampaignAudienceParams, audienceType string, limit int) ([]entities.CampaignCandidate, error)
	CountCampaignAudience(ctx context.Context, businessID uint, params entities.CampaignAudienceParams, audienceType string) (total, optedOut, noPhone, reachable uint, err error)
}

type ICampaignSenderQuerier interface {
	GetOwnWhatsappSender(ctx context.Context, businessID uint) (integrationID uint, phoneNumber string, err error)
}
