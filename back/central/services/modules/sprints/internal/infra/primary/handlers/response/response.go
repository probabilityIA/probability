package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
)

type SprintResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Goal          string    `json:"goal"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        string    `json:"status"`
	CreatedByID   uint      `json:"created_by_id"`
	CreatedByName string    `json:"created_by_name"`
	TicketCount   int64     `json:"ticket_count"`
	DoneCount     int64     `json:"done_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func FromSprint(s *entities.Sprint) SprintResponse {
	return SprintResponse{
		ID:            s.ID,
		Name:          s.Name,
		Goal:          s.Goal,
		StartDate:     s.StartDate,
		EndDate:       s.EndDate,
		Status:        s.Status,
		CreatedByID:   s.CreatedByID,
		CreatedByName: s.CreatedByName,
		TicketCount:   s.TicketCount,
		DoneCount:     s.DoneCount,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
