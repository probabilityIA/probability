package campaigns

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

const (
	dueCampaignsBatch    = 20
	audienceMaterialize  = 5000
	senderNameVariable   = "sender.name"
	campaignNameVariable = "campaign.name"
)

func (uc *useCase) materializeAudience(ctx context.Context, campaign *entities.Campaign) (int, error) {
	candidates, err := uc.audience.FindCampaignCandidates(
		ctx,
		campaign.BusinessID,
		campaign.AudienceParams,
		campaign.AudienceType,
		audienceMaterialize,
	)
	if err != nil {
		return 0, err
	}

	if len(candidates) == 0 {
		return 0, nil
	}

	sends := make([]entities.CampaignSend, 0, len(candidates))
	for _, candidate := range candidates {
		sends = append(sends, entities.CampaignSend{
			CampaignID: campaign.ID,
			ClientID:   candidate.ClientID,
			BusinessID: campaign.BusinessID,
			Phone:      candidate.Phone,
		})
	}

	return uc.sends.BulkCreateSends(ctx, sends)
}

func (uc *useCase) RunDueCampaigns(ctx context.Context) error {
	campaigns, err := uc.campaigns.ListDueCampaigns(ctx, time.Now(), dueCampaignsBatch)
	if err != nil {
		return err
	}

	for i := range campaigns {
		campaign := campaigns[i]

		if !uc.isInsideWindow(&campaign, time.Now()) {
			continue
		}

		if err := uc.dispatchBatch(ctx, &campaign); err != nil {
			uc.logger.Error().Err(err).Uint("campaign_id", campaign.ID).
				Msg("Error despachando la tanda de la campana")
		}
	}

	return nil
}

func (uc *useCase) dispatchBatch(ctx context.Context, campaign *entities.Campaign) error {
	template, err := uc.resolveTemplate(ctx, campaign)
	if err != nil {
		return uc.failCampaign(ctx, campaign, err)
	}

	sentToday, err := uc.sends.CountSentSince(ctx, campaign.ID, startOfDay(campaign, time.Now()))
	if err != nil {
		return err
	}

	remaining := int(campaign.DailySendCap) - int(sentToday)
	if remaining <= 0 {
		return nil
	}

	limit := min(int(campaign.BatchSize), remaining)

	pending, err := uc.sends.ListPendingSends(ctx, campaign.ID, limit)
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		return uc.completeCampaign(ctx, campaign)
	}

	now := time.Now()
	campaign.Status = entities.CampaignStatusRunning

	for _, send := range pending {
		message := dtos.CampaignSendMessage{
			SendID:       send.ID,
			CampaignID:   campaign.ID,
			BusinessID:   campaign.BusinessID,
			ClientID:     send.ClientID,
			Phone:        send.Phone,
			TemplateName: template.Name,
			Language:     template.Language,
			Parameters:   buildParameters(template, campaign, send),
		}

		if uc.publisher == nil {
			continue
		}

		if err := uc.sends.MarkSendQueued(ctx, send.ID, now); err != nil {
			continue
		}

		if err := uc.publisher.PublishCampaignSend(ctx, message); err != nil {
			if markErr := uc.sends.MarkSendResult(ctx, send.ID, entities.CampaignSendStatusFailed, "", err.Error()); markErr != nil {
				uc.logger.Error().Err(markErr).Uint("send_id", send.ID).
					Msg("Error marcando el envio de campana como fallido")
			}
		}
	}

	if err := uc.campaigns.MarkCampaignBatch(ctx, campaign.ID, now); err != nil {
		return err
	}
	if err := uc.campaigns.UpdateCampaign(ctx, campaign); err != nil {
		return err
	}

	return uc.campaigns.UpdateCampaignCounters(ctx, campaign.ID)
}

func (uc *useCase) completeCampaign(ctx context.Context, campaign *entities.Campaign) error {
	now := time.Now()
	campaign.Status = entities.CampaignStatusCompleted
	campaign.FinishedAt = &now

	if err := uc.campaigns.UpdateCampaign(ctx, campaign); err != nil {
		return err
	}

	return uc.campaigns.UpdateCampaignCounters(ctx, campaign.ID)
}

func (uc *useCase) failCampaign(ctx context.Context, campaign *entities.Campaign, cause error) error {
	campaign.Status = entities.CampaignStatusPaused
	campaign.ErrorMessage = cause.Error()

	if err := uc.campaigns.UpdateCampaign(ctx, campaign); err != nil {
		return err
	}

	uc.logger.Warn().
		Uint("campaign_id", campaign.ID).
		Str("motivo", cause.Error()).
		Msg("Campana pausada por un problema de configuracion")

	return nil
}

func (uc *useCase) isInsideWindow(campaign *entities.Campaign, now time.Time) bool {
	location, err := time.LoadLocation(campaign.Timezone)
	if err != nil {
		location = time.UTC
	}

	local := now.In(location)
	current := local.Hour()*60 + local.Minute()

	start, err := parseClock(campaign.SendWindowStart)
	if err != nil {
		return true
	}
	end, err := parseClock(campaign.SendWindowEnd)
	if err != nil {
		return true
	}

	return current >= start && current < end
}

func startOfDay(campaign *entities.Campaign, now time.Time) time.Time {
	location, err := time.LoadLocation(campaign.Timezone)
	if err != nil {
		location = time.UTC
	}

	local := now.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
}

func buildParameters(template *entities.WhatsappTemplate, campaign *entities.Campaign, send entities.CampaignSend) []string {
	parameters := make([]string, 0, len(template.Variables))

	for _, variable := range template.Variables {
		value := resolveVariable(variable, campaign, send)
		if value == "" {
			value = variable.Fallback
		}
		if value == "" {
			value = "-"
		}
		parameters = append(parameters, value)
	}

	return parameters
}

func resolveVariable(variable entities.TemplateVariable, campaign *entities.Campaign, send entities.CampaignSend) string {
	if campaign.VariableValues != nil {
		if fixed, ok := campaign.VariableValues[variable.Source]; ok && strings.TrimSpace(fixed) != "" {
			return strings.TrimSpace(fixed)
		}
	}

	switch variable.Source {
	case "customer.first_name":
		return firstName(send.ClientName)
	case "customer.full_name":
		return strings.TrimSpace(send.ClientName)
	case senderNameVariable:
		return strings.TrimSpace(campaign.SenderName)
	case campaignNameVariable:
		return strings.TrimSpace(campaign.Name)
	default:
		return ""
	}
}

func firstName(name string) string {
	fields := strings.Fields(strings.TrimSpace(name))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
