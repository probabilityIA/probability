package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func EntityToModel(s *entities.Sprint) *models.Sprint {
	return &models.Sprint{
		Name:        s.Name,
		Goal:        s.Goal,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		Status:      s.Status,
		CreatedByID: s.CreatedByID,
	}
}

func ModelToEntity(m *models.Sprint) *entities.Sprint {
	return &entities.Sprint{
		ID:            m.ID,
		Name:          m.Name,
		Goal:          m.Goal,
		StartDate:     m.StartDate,
		EndDate:       m.EndDate,
		Status:        m.Status,
		CreatedByID:   m.CreatedByID,
		CreatedByName: m.CreatedBy.Name,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
