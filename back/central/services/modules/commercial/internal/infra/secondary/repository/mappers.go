package repository

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func modelToProspect(m *models.CommercialProspect) entities.Prospect {
	var breakdown map[string]interface{}
	if len(m.ScoreBreakdown) > 0 {
		_ = json.Unmarshal(m.ScoreBreakdown, &breakdown)
	}

	return entities.Prospect{
		ID:                     m.ID,
		Domain:                 m.Domain,
		MeliSellerID:           m.MeliSellerID,
		BusinessName:           m.BusinessName,
		Platform:               m.Platform,
		City:                   m.City,
		Phone:                  m.Phone,
		WhatsAppNumber:         m.WhatsAppNumber,
		InstagramHandle:        m.InstagramHandle,
		InstagramFollowers:     m.InstagramFollowers,
		ProductCount:           m.ProductCount,
		HasCOD:                 m.HasCOD,
		NumCouriersMentioned:   m.NumCouriersMentioned,
		UsesProbabilityChannel: m.UsesProbabilityChannel,
		RobotsOK:               m.RobotsOK,
		DiscoveryQuery:         m.DiscoveryQuery,
		Score:                  m.Score,
		ScoreBreakdown:         breakdown,
		Status:                 m.Status,
		Seen:                   m.Seen,
		SourceDiscoveredAt:     m.SourceDiscoveredAt,
		SourceLastScannedAt:    m.SourceLastScannedAt,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}
