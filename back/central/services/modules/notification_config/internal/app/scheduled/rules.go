package scheduled

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
	defaultFrequency       = uint(1440)
	defaultCooldownDays    = uint(30)
	defaultDailySendCap    = uint(500)
	defaultBatchSizeCap    = uint(200)
	maxBatchSizeCap        = uint(1000)
	defaultSendWindowStart = "09:00"
	defaultSendWindowEnd   = "19:00"
)

func (uc *useCase) CreateRule(ctx context.Context, dto dtos.CreateScheduledRuleDTO) (*entities.ScheduledRule, error) {
	rule, err := uc.buildRule(ctx, dto)
	if err != nil {
		return nil, err
	}

	next := time.Now()
	rule.NextRunAt = &next

	if err := uc.rules.CreateRule(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (uc *useCase) UpdateRule(ctx context.Context, dto dtos.UpdateScheduledRuleDTO) (*entities.ScheduledRule, error) {
	current, err := uc.GetRule(ctx, dto.ID, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("regla no encontrada")
	}

	rule, err := uc.buildRule(ctx, dto.CreateScheduledRuleDTO)
	if err != nil {
		return nil, err
	}

	rule.ID = current.ID
	rule.CreatedByID = current.CreatedByID
	rule.LastRunAt = current.LastRunAt
	rule.NextRunAt = current.NextRunAt

	if err := uc.rules.UpdateRule(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (uc *useCase) GetRule(ctx context.Context, id, businessID uint) (*entities.ScheduledRule, error) {
	rule, err := uc.rules.GetRuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, nil
	}
	if rule.BusinessID != businessID {
		return nil, fmt.Errorf("la regla no pertenece al negocio")
	}
	return rule, nil
}

func (uc *useCase) ListRules(ctx context.Context, businessID uint, page, pageSize int) ([]entities.ScheduledRule, int64, error) {
	page, pageSize = normalizePaging(page, pageSize)
	return uc.rules.ListRules(ctx, businessID, page, pageSize)
}

func (uc *useCase) DeleteRule(ctx context.Context, id, businessID uint) error {
	rule, err := uc.GetRule(ctx, id, businessID)
	if err != nil {
		return err
	}
	if rule == nil {
		return fmt.Errorf("regla no encontrada")
	}
	return uc.rules.DeleteRule(ctx, id)
}

func (uc *useCase) ListRuns(ctx context.Context, ruleID, businessID uint, page, pageSize int) ([]entities.ScheduledRun, int64, error) {
	rule, err := uc.GetRule(ctx, ruleID, businessID)
	if err != nil {
		return nil, 0, err
	}
	if rule == nil {
		return nil, 0, fmt.Errorf("regla no encontrada")
	}

	page, pageSize = normalizePaging(page, pageSize)
	return uc.runs.ListRuns(ctx, ruleID, page, pageSize)
}

func (uc *useCase) MarkSendResult(ctx context.Context, sendID uint, status, messageID, errorMessage string) error {
	return uc.sends.MarkSendResult(ctx, sendID, status, messageID, errorMessage)
}

func (uc *useCase) buildRule(ctx context.Context, dto dtos.CreateScheduledRuleDTO) (*entities.ScheduledRule, error) {
	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, fmt.Errorf("el nombre de la regla es obligatorio")
	}

	segmentType := strings.TrimSpace(dto.SegmentType)
	if segmentType == "" {
		segmentType = entities.SegmentCustomersInactive
	}
	if !entities.IsSupportedSegment(segmentType) {
		return nil, fmt.Errorf("segmento no soportado: %s", segmentType)
	}

	if dto.DaysWithoutPurchase <= 0 {
		return nil, fmt.Errorf("days_without_purchase debe ser mayor a cero")
	}
	if dto.MaxOrders > 0 && dto.MinOrders > dto.MaxOrders {
		return nil, fmt.Errorf("min_orders no puede ser mayor que max_orders")
	}

	if dto.NotificationTypeID == 0 {
		return nil, fmt.Errorf("notification_type_id es obligatorio")
	}

	if dto.WhatsappTemplateID == nil || *dto.WhatsappTemplateID == 0 {
		return nil, fmt.Errorf("la regla necesita una plantilla de WhatsApp")
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

	frequency := dto.FrequencyMinutes
	if frequency == 0 {
		frequency = defaultFrequency
	}
	if frequency < 60 {
		return nil, fmt.Errorf("frequency_minutes no puede ser menor a 60")
	}

	cooldown := dto.CooldownDays
	if cooldown == 0 {
		cooldown = defaultCooldownDays
	}

	dailyCap := dto.DailySendCap
	if dailyCap == 0 {
		dailyCap = defaultDailySendCap
	}

	batchCap := dto.BatchSizeCap
	if batchCap == 0 {
		batchCap = defaultBatchSizeCap
	}
	if batchCap > maxBatchSizeCap {
		batchCap = maxBatchSizeCap
	}

	requiresOptIn := true
	if dto.RequiresOptIn != nil {
		requiresOptIn = *dto.RequiresOptIn
	}

	enabled := true
	if dto.Enabled != nil {
		enabled = *dto.Enabled
	}

	return &entities.ScheduledRule{
		BusinessID:         dto.BusinessID,
		IntegrationID:      dto.IntegrationID,
		NotificationTypeID: dto.NotificationTypeID,
		WhatsappTemplateID: dto.WhatsappTemplateID,
		Name:               name,
		Description:        strings.TrimSpace(dto.Description),
		SegmentType:        segmentType,
		SegmentParams: entities.SegmentParams{
			DaysWithoutPurchase: dto.DaysWithoutPurchase,
			MinOrders:           dto.MinOrders,
			MaxOrders:           dto.MaxOrders,
		},
		Timezone:         timezone,
		FrequencyMinutes: frequency,
		SendWindowStart:  windowStart,
		SendWindowEnd:    windowEnd,
		CooldownDays:     cooldown,
		DailySendCap:     dailyCap,
		BatchSizeCap:     batchCap,
		RequiresOptIn:    requiresOptIn,
		Enabled:          enabled,
		CreatedByID:      dto.CreatedBy,
		Template:         template,
	}, nil
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
