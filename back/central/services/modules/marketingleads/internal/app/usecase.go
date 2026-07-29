package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/errors"
)

var surveyTitles = map[string]string{
	"diagnostico-logistico": "Diagnóstico Logístico",
	"escanea-tu-logistica":  "Diagnóstico de Operación E-commerce",
}

func surveyTitle(slug string) string {
	if title, ok := surveyTitles[slug]; ok {
		return title
	}
	return "Diagnóstico"
}

func (uc *UseCase) CreateLead(ctx context.Context, dto dtos.CreateLeadDTO) error {
	name := strings.TrimSpace(dto.Name)
	email := strings.TrimSpace(dto.Email)
	phone := strings.TrimSpace(dto.Phone)
	if name == "" || email == "" || phone == "" {
		return domainerrors.ErrInvalidLead
	}

	lead := &entities.MarketingLead{
		Name:       name,
		Email:      email,
		Phone:      phone,
		SurveySlug: dto.SurveySlug,
		ScoreTotal: dto.ScoreTotal,
		ScoreMax:   dto.ScoreMax,
		Level:      dto.Level,
	}

	if err := uc.repo.CreateLead(ctx, lead); err != nil {
		return err
	}

	if uc.wa != nil {
		messageID, err := uc.wa.SendDiagnosticResultTemplate(ctx, phone, name, surveyTitle(dto.SurveySlug), dto.Level)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).
				Str("phone", phone).
				Str("survey_slug", dto.SurveySlug).
				Msg("marketingleads: no se pudo enviar el WhatsApp de resultado (el lead ya quedo guardado)")
		} else if messageID != "" {
			if err := uc.repo.SetWhatsAppMessageID(ctx, lead.ID, messageID); err != nil {
				uc.logger.Warn(ctx).Err(err).Msg("marketingleads: no se pudo guardar el whatsapp_message_id")
			}
		}
	}

	return nil
}

func (uc *UseCase) ListLeads(ctx context.Context, page, pageSize int) ([]entities.MarketingLead, dtos.ListLeadsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	leads, total, err := uc.repo.ListLeads(ctx, page, pageSize)
	if err != nil {
		return nil, dtos.ListLeadsResult{}, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	return leads, dtos.ListLeadsResult{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
