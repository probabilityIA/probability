package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

var closedTicketStatuses = []string{"resolved", "closed", "wont_fix"}

func (r *Repository) countTicketsBySprint(ctx context.Context, sprintIDs []uint) (map[uint]entities.SprintCounts, error) {
	out := make(map[uint]entities.SprintCounts, len(sprintIDs))
	if len(sprintIDs) == 0 {
		return out, nil
	}

	var rows []struct {
		SprintID    uint
		TicketCount int64
		DoneCount   int64
	}

	err := r.db.Conn(ctx).Model(&models.Ticket{}).
		Select("sprint_id, COUNT(*) AS ticket_count, COUNT(*) FILTER (WHERE status IN ?) AS done_count", closedTicketStatuses).
		Where("sprint_id IN ?", sprintIDs).
		Group("sprint_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		out[row.SprintID] = entities.SprintCounts{TicketCount: row.TicketCount, DoneCount: row.DoneCount}
	}
	return out, nil
}

func detachTicketsFromSprint(tx *gorm.DB, sprintID uint) error {
	return tx.Model(&models.Ticket{}).
		Where("sprint_id = ?", sprintID).
		Update("sprint_id", nil).Error
}
