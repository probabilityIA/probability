package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateTicketArea(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Ticket{}, &models.TicketStatusHistory{}); err != nil {
		return fmt.Errorf("failed to auto-migrate ticket area columns: %w", err)
	}
	return nil
}
