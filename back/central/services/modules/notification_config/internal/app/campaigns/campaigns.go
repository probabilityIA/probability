package campaigns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

const (
	defaultTimezone        = "America/Bogota"
	defaultSendWindowStart = "09:00"
	defaultSendWindowEnd   = "19:00"
	defaultDailySendCap    = uint(250)
	defaultBatchSize       = uint(50)
	audienceSampleSize     = 5
)

func (uc *useCase) Create(ctx context.Context, dto dtos.CreateCampaignDTO) (*entities.Campaign, error) {
	campaign, err := uc.build(ctx, dto)
	if err != nil {
		return nil, err
	}

	campaign.Status = entities.CampaignStatusDraft

	if err := uc.campaigns.CreateCampaign(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

func (uc *useCase) Update(ctx context.Context, dto dtos.UpdateCampaignDTO) (*entities.Campaign, error) {
	current, err := uc.GetByID(ctx, dto.ID, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("campana no encontrada")
	}
	if !current.IsEditable() {
		return nil, fmt.Errorf("una campana en estado %s ya no se puede editar", current.Status)
	}

	campaign, err := uc.build(ctx, dto.CreateCampaignDTO)
	if err != nil {
		return nil, err
	}

	campaign.ID = current.ID
	campaign.Status = current.Status
	campaign.CreatedByID = current.CreatedByID
	campaign.StartedAt = current.StartedAt
	campaign.FinishedAt = current.FinishedAt

	if err := uc.campaigns.UpdateCampaign(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

func (uc *useCase) GetByID(ctx context.Context, id, businessID uint) (*entities.Campaign, error) {
	campaign, err := uc.campaigns.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, nil
	}
	if campaign.BusinessID != businessID {
		return nil, fmt.Errorf("la campana no pertenece al negocio")
	}
	return campaign, nil
}

func (uc *useCase) List(ctx context.Context, businessID uint, status string, page, pageSize int) ([]entities.Campaign, int64, error) {
	page, pageSize = normalizePaging(page, pageSize)
	return uc.campaigns.ListCampaigns(ctx, businessID, strings.TrimSpace(status), page, pageSize)
}

func (uc *useCase) Delete(ctx context.Context, id, businessID uint) error {
	campaign, err := uc.GetByID(ctx, id, businessID)
	if err != nil {
		return err
	}
	if campaign == nil {
		return fmt.Errorf("campana no encontrada")
	}
	if campaign.Status == entities.CampaignStatusRunning {
		return fmt.Errorf("pausa la campana antes de eliminarla")
	}
	return uc.campaigns.DeleteCampaign(ctx, id)
}

func (uc *useCase) ListSends(ctx context.Context, campaignID, businessID uint, status string, page, pageSize int) ([]entities.CampaignSend, int64, error) {
	campaign, err := uc.GetByID(ctx, campaignID, businessID)
	if err != nil {
		return nil, 0, err
	}
	if campaign == nil {
		return nil, 0, fmt.Errorf("campana no encontrada")
	}

	page, pageSize = normalizePaging(page, pageSize)
	return uc.sends.ListSends(ctx, campaignID, strings.TrimSpace(status), page, pageSize)
}

func (uc *useCase) PreviewAudience(ctx context.Context, dto dtos.CreateCampaignDTO) (*dtos.CampaignAudiencePreviewDTO, error) {
	if dto.BusinessID == 0 {
		return nil, fmt.Errorf("business_id es obligatorio")
	}

	audienceType := normalizeAudienceType(dto.AudienceType)
	if !entities.IsSupportedAudience(audienceType) {
		return nil, fmt.Errorf("audiencia no soportada: %s", audienceType)
	}

	params := buildAudienceParams(dto)

	total, optedOut, noPhone, reachable, err := uc.audience.CountCampaignAudience(ctx, dto.BusinessID, params, audienceType)
	if err != nil {
		return nil, err
	}

	candidates, err := uc.audience.FindCampaignCandidates(ctx, dto.BusinessID, params, audienceType, audienceSampleSize)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		names = append(names, candidate.Name)
	}

	return &dtos.CampaignAudiencePreviewDTO{
		Total:      total,
		OptedOut:   optedOut,
		NoPhone:    noPhone,
		Reachable:  reachable,
		SampleName: names,
	}, nil
}

func (uc *useCase) Launch(ctx context.Context, id, businessID uint) (*entities.Campaign, error) {
	campaign, err := uc.GetByID(ctx, id, businessID)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, fmt.Errorf("campana no encontrada")
	}
	if campaign.Status != entities.CampaignStatusDraft && campaign.Status != entities.CampaignStatusScheduled {
		return nil, fmt.Errorf("la campana ya fue lanzada")
	}

	if _, err := uc.resolveTemplate(ctx, campaign); err != nil {
		return nil, err
	}
	if err := uc.requireOwnNumber(ctx, businessID); err != nil {
		return nil, err
	}

	queued, err := uc.materializeAudience(ctx, campaign)
	if err != nil {
		return nil, err
	}
	if queued == 0 {
		return nil, fmt.Errorf("la audiencia quedo vacia: nadie cumple los filtros o ninguno acepta marketing")
	}

	now := time.Now()
	campaign.Status = entities.CampaignStatusScheduled
	campaign.AudienceCount = uint(queued)
	campaign.StartedAt = &now
	if campaign.ScheduledAt == nil {
		campaign.ScheduledAt = &now
	}

	if err := uc.campaigns.UpdateCampaign(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

func (uc *useCase) Pause(ctx context.Context, id, businessID uint) (*entities.Campaign, error) {
	return uc.transition(ctx, id, businessID, entities.CampaignStatusPaused, func(campaign *entities.Campaign) error {
		if campaign.Status != entities.CampaignStatusRunning && campaign.Status != entities.CampaignStatusScheduled {
			return fmt.Errorf("solo se puede pausar una campana programada o en curso")
		}
		return nil
	})
}

func (uc *useCase) Resume(ctx context.Context, id, businessID uint) (*entities.Campaign, error) {
	return uc.transition(ctx, id, businessID, entities.CampaignStatusScheduled, func(campaign *entities.Campaign) error {
		if campaign.Status != entities.CampaignStatusPaused {
			return fmt.Errorf("solo se puede reanudar una campana pausada")
		}
		return nil
	})
}

func (uc *useCase) Cancel(ctx context.Context, id, businessID uint) (*entities.Campaign, error) {
	return uc.transition(ctx, id, businessID, entities.CampaignStatusCancelled, func(campaign *entities.Campaign) error {
		if campaign.IsFinished() {
			return fmt.Errorf("la campana ya termino")
		}
		return nil
	})
}

func (uc *useCase) transition(ctx context.Context, id, businessID uint, status string, guard func(*entities.Campaign) error) (*entities.Campaign, error) {
	campaign, err := uc.GetByID(ctx, id, businessID)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, fmt.Errorf("campana no encontrada")
	}
	if err := guard(campaign); err != nil {
		return nil, err
	}

	campaign.Status = status
	if status == entities.CampaignStatusCancelled {
		now := time.Now()
		campaign.FinishedAt = &now
	}

	if err := uc.campaigns.UpdateCampaign(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

func (uc *useCase) MarkSendResult(ctx context.Context, result dtos.CampaignSendResult) error {
	if err := uc.sends.MarkSendResult(ctx, result.SendID, result.Status, result.MessageID, result.ErrorMessage); err != nil {
		return err
	}
	return nil
}

func (uc *useCase) requireOwnNumber(ctx context.Context, businessID uint) error {
	integrationID, _, err := uc.sender.GetOwnWhatsappSender(ctx, businessID)
	if err != nil {
		return err
	}
	if integrationID == 0 {
		return fmt.Errorf("las campanas de marketing solo salen desde el numero propio del negocio: conecta un numero de WhatsApp antes de lanzar")
	}
	return nil
}

func (uc *useCase) resolveTemplate(ctx context.Context, campaign *entities.Campaign) (*entities.WhatsappTemplate, error) {
	if campaign.WhatsappTemplateID == nil || *campaign.WhatsappTemplateID == 0 {
		return nil, fmt.Errorf("la campana no tiene plantilla asociada")
	}

	template, err := uc.templates.GetTemplateByID(ctx, *campaign.WhatsappTemplateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("la plantilla de la campana ya no existe")
	}
	if !template.BelongsTo(campaign.BusinessID) {
		return nil, fmt.Errorf("la plantilla no pertenece al negocio")
	}
	if !template.IsUsable() {
		return nil, fmt.Errorf("la plantilla esta en estado %s: solo se envia si esta aprobada por Meta", template.Status)
	}

	return template, nil
}

func (uc *useCase) build(ctx context.Context, dto dtos.CreateCampaignDTO) (*entities.Campaign, error) {
	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, fmt.Errorf("el nombre de la campana es obligatorio")
	}

	if dto.BusinessID == 0 {
		return nil, fmt.Errorf("business_id es obligatorio")
	}

	audienceType := normalizeAudienceType(dto.AudienceType)
	if !entities.IsSupportedAudience(audienceType) {
		return nil, fmt.Errorf("audiencia no soportada: %s", audienceType)
	}

	if dto.WhatsappTemplateID == nil || *dto.WhatsappTemplateID == 0 {
		return nil, fmt.Errorf("la campana necesita una plantilla de WhatsApp")
	}

	template, err := uc.templates.GetTemplateByID(ctx, *dto.WhatsappTemplateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("la plantilla no existe")
	}
	if !template.BelongsTo(dto.BusinessID) {
		return nil, fmt.Errorf("la plantilla no pertenece al negocio")
	}

	timezone := strings.TrimSpace(dto.Timezone)
	if timezone == "" {
		timezone = defaultTimezone
	}
	if _, err := time.LoadLocation(timezone); err != nil {
		return nil, fmt.Errorf("zona horaria invalida: %s", timezone)
	}

	windowStart := strings.TrimSpace(dto.SendWindowStart)
	if windowStart == "" {
		windowStart = defaultSendWindowStart
	}
	windowEnd := strings.TrimSpace(dto.SendWindowEnd)
	if windowEnd == "" {
		windowEnd = defaultSendWindowEnd
	}
	if err := validateWindow(windowStart, windowEnd); err != nil {
		return nil, err
	}

	dailyCap := dto.DailySendCap
	if dailyCap == 0 {
		dailyCap = defaultDailySendCap
	}
	if dailyCap > entities.CampaignMaxDailySendCap {
		return nil, fmt.Errorf("el tope diario no puede superar %d mensajes", entities.CampaignMaxDailySendCap)
	}

	batchSize := dto.BatchSize
	if batchSize == 0 {
		batchSize = defaultBatchSize
	}
	if batchSize > entities.CampaignMaxBatchSize {
		batchSize = entities.CampaignMaxBatchSize
	}
	if batchSize > dailyCap {
		batchSize = dailyCap
	}

	return &entities.Campaign{
		BusinessID:         dto.BusinessID,
		IntegrationID:      dto.IntegrationID,
		WhatsappTemplateID: dto.WhatsappTemplateID,
		Name:               name,
		Description:        strings.TrimSpace(dto.Description),
		SenderName:         strings.TrimSpace(dto.SenderName),
		AudienceType:       audienceType,
		AudienceParams:     buildAudienceParams(dto),
		VariableValues:     dto.VariableValues,
		Timezone:           timezone,
		SendWindowStart:    windowStart,
		SendWindowEnd:      windowEnd,
		ScheduledAt:        dto.ScheduledAt,
		DailySendCap:       dailyCap,
		BatchSize:          batchSize,
		CreatedByID:        dto.CreatedBy,
		Template:           template,
	}, nil
}

func buildAudienceParams(dto dtos.CreateCampaignDTO) entities.CampaignAudienceParams {
	return entities.CampaignAudienceParams{
		City:             strings.TrimSpace(dto.City),
		CreatedFromDays:  dto.CreatedFromDays,
		OnlyWithoutOrder: dto.OnlyWithoutOrder,
		ClientIDs:        dto.ClientIDs,
	}
}

func normalizeAudienceType(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return entities.CampaignAudienceFiltered
	}
	return trimmed
}

func validateWindow(start, end string) error {
	startMinutes, err := parseClock(start)
	if err != nil {
		return err
	}
	endMinutes, err := parseClock(end)
	if err != nil {
		return err
	}
	if startMinutes >= endMinutes {
		return fmt.Errorf("la ventana de envio debe empezar antes de terminar")
	}
	return nil
}

func parseClock(value string) (int, error) {
	parsed, err := time.Parse("15:04", strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("hora invalida (formato HH:MM): %s", value)
	}
	return parsed.Hour()*60 + parsed.Minute(), nil
}

func normalizePaging(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
