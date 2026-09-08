package scheduled

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

const dueRulesBatch = 20

func (uc *useCase) RunDueRules(ctx context.Context) error {
	rules, err := uc.rules.ListDueRules(ctx, time.Now(), dueRulesBatch)
	if err != nil {
		return err
	}

	for i := range rules {
		rule := rules[i]

		if !uc.isInsideWindow(&rule, time.Now()) {
			uc.logger.Debug().Uint("rule_id", rule.ID).
				Msg("Regla fuera de la ventana de envio, se reprograma")
			uc.reschedule(ctx, &rule)
			continue
		}

		if _, err := uc.execute(ctx, &rule); err != nil {
			uc.logger.Error().Err(err).Uint("rule_id", rule.ID).
				Msg("Error ejecutando la regla programada")
		}
	}

	return nil
}

func (uc *useCase) RunRuleNow(ctx context.Context, id, businessID uint) (*entities.ScheduledRun, error) {
	rule, err := uc.GetRule(ctx, id, businessID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, fmt.Errorf("regla no encontrada")
	}

	return uc.execute(ctx, rule)
}

func (uc *useCase) execute(ctx context.Context, rule *entities.ScheduledRule) (*entities.ScheduledRun, error) {
	run := &entities.ScheduledRun{
		RuleID:     rule.ID,
		BusinessID: rule.BusinessID,
		StartedAt:  time.Now(),
		Status:     entities.RunStatusRunning,
	}

	if err := uc.runs.CreateRun(ctx, run); err != nil {
		return nil, err
	}

	finish := func(status, message string) {
		now := time.Now()
		run.FinishedAt = &now
		run.Status = status
		run.ErrorMessage = message
		if err := uc.runs.FinishRun(ctx, run); err != nil {
			uc.logger.Error().Err(err).Uint("run_id", run.ID).Msg("Error cerrando la corrida")
		}
		uc.reschedule(ctx, rule)
	}

	template, err := uc.resolveTemplate(ctx, rule)
	if err != nil {
		finish(entities.RunStatusFailed, err.Error())
		return run, err
	}

	sendDate := time.Now().Format("2006-01-02")

	alreadySent, err := uc.sends.CountSendsToday(ctx, rule.ID, sendDate)
	if err != nil {
		finish(entities.RunStatusFailed, err.Error())
		return run, err
	}

	remaining := int(rule.DailySendCap) - int(alreadySent)
	if remaining <= 0 {
		run.CapReachedFlag = true
		finish(entities.RunStatusCompleted, "")
		return run, nil
	}

	limit := min(int(rule.BatchSizeCap), remaining)

	candidates, err := uc.segments.FindInactiveCustomers(
		ctx,
		rule.BusinessID,
		rule.SegmentParams,
		rule.RequiresOptIn,
		rule.CooldownDays,
		rule.ID,
		limit,
	)
	if err != nil {
		finish(entities.RunStatusFailed, err.Error())
		return run, err
	}

	run.MatchedCount = uint(len(candidates))

	for _, candidate := range candidates {
		if candidate.Phone == "" {
			run.SkippedCount++
			continue
		}

		send := &entities.ScheduledSend{
			RuleID:     rule.ID,
			ClientID:   candidate.ClientID,
			SendDate:   sendDate,
			RunID:      &run.ID,
			BusinessID: rule.BusinessID,
			Phone:      candidate.Phone,
		}

		reserved, err := uc.sends.ReserveSend(ctx, send)
		if err != nil {
			run.FailedCount++
			continue
		}
		if !reserved {
			run.SkippedCount++
			continue
		}

		message := dtos.ScheduledSendMessage{
			SendID:       send.ID,
			RuleID:       rule.ID,
			BusinessID:   rule.BusinessID,
			ClientID:     candidate.ClientID,
			Phone:        candidate.Phone,
			TemplateName: template.Name,
			Language:     template.Language,
			Parameters:   buildParameters(template, candidate),
		}

		if uc.publisher == nil {
			run.FailedCount++
			continue
		}

		if err := uc.publisher.PublishScheduledSend(ctx, message); err != nil {
			run.FailedCount++
			if markErr := uc.sends.MarkSendResult(ctx, send.ID, entities.SendStatusFailed, "", err.Error()); markErr != nil {
				uc.logger.Error().Err(markErr).Uint("send_id", send.ID).
					Msg("Error marcando el envio como fallido")
			}
			continue
		}

		run.QueuedCount++
	}

	if run.QueuedCount >= uint(remaining) {
		run.CapReachedFlag = true
	}

	uc.logger.Info().
		Uint("rule_id", rule.ID).
		Uint("run_id", run.ID).
		Uint("matched", run.MatchedCount).
		Uint("queued", run.QueuedCount).
		Uint("skipped", run.SkippedCount).
		Uint("failed", run.FailedCount).
		Msg("Corrida de regla programada terminada")

	finish(entities.RunStatusCompleted, "")

	return run, nil
}

func (uc *useCase) resolveTemplate(ctx context.Context, rule *entities.ScheduledRule) (*entities.WhatsappTemplate, error) {
	if rule.WhatsappTemplateID == nil || *rule.WhatsappTemplateID == 0 {
		return nil, fmt.Errorf("la regla no tiene plantilla asociada")
	}

	template, err := uc.templates.GetTemplateByID(ctx, *rule.WhatsappTemplateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("la plantilla de la regla ya no existe")
	}
	if !template.IsUsable() {
		return nil, fmt.Errorf("la plantilla esta en estado %s: solo se envia si esta aprobada", template.Status)
	}

	return template, nil
}

func (uc *useCase) reschedule(ctx context.Context, rule *entities.ScheduledRule) {
	now := time.Now()
	next := now.Add(time.Duration(rule.FrequencyMinutes) * time.Minute)

	if err := uc.rules.MarkRuleRan(ctx, rule.ID, now, next); err != nil {
		uc.logger.Error().Err(err).Uint("rule_id", rule.ID).
			Msg("Error reprogramando la regla")
	}
}

func (uc *useCase) isInsideWindow(rule *entities.ScheduledRule, now time.Time) bool {
	location, err := time.LoadLocation(rule.Timezone)
	if err != nil {
		location = time.UTC
	}

	local := now.In(location)
	current := local.Hour()*60 + local.Minute()

	start, err := parseClock(rule.SendWindowStart)
	if err != nil {
		return true
	}
	end, err := parseClock(rule.SendWindowEnd)
	if err != nil {
		return true
	}

	return current >= start && current < end
}

func buildParameters(template *entities.WhatsappTemplate, candidate entities.SegmentCandidate) []string {
	parameters := make([]string, 0, len(template.Variables))

	for _, variable := range template.Variables {
		value := resolveVariable(variable, candidate)
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

func resolveVariable(variable entities.TemplateVariable, candidate entities.SegmentCandidate) string {
	switch variable.Source {
	case "customer.first_name":
		return firstName(candidate.Name)
	case "customer.full_name":
		return strings.TrimSpace(candidate.Name)
	case "customer.days_inactive":
		return strconv.Itoa(candidate.DaysInactive)
	case "customer.total_orders":
		return strconv.Itoa(candidate.TotalOrders)
	case "customer.last_product":
		return strings.TrimSpace(candidate.LastProduct)
	case "business.name":
		return strings.TrimSpace(candidate.BusinessName)
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
