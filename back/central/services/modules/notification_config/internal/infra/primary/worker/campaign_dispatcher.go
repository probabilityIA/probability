package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/campaigns"
	"github.com/secamc93/probability/back/central/shared/log"
)

const campaignTickInterval = 2 * time.Minute

type CampaignDispatcher struct {
	useCase campaigns.IUseCase
	logger  log.ILogger
}

func NewCampaignDispatcher(useCase campaigns.IUseCase, logger log.ILogger) *CampaignDispatcher {
	return &CampaignDispatcher{
		useCase: useCase,
		logger:  logger.WithModule("whatsapp_campaigns_worker"),
	}
}

func (d *CampaignDispatcher) Start(ctx context.Context) {
	ticker := time.NewTicker(campaignTickInterval)
	defer ticker.Stop()

	d.logger.Info(ctx).Dur("interval", campaignTickInterval).
		Msg("Worker de campanas de WhatsApp iniciado")

	for {
		select {
		case <-ctx.Done():
			d.logger.Info(ctx).Msg("Worker de campanas de WhatsApp detenido")
			return
		case <-ticker.C:
			if err := d.useCase.RunDueCampaigns(ctx); err != nil {
				d.logger.Error(ctx).Err(err).
					Msg("Error despachando campanas de WhatsApp")
			}
		}
	}
}
