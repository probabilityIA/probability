package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateTickets(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.Ticket{},
		&models.TicketComment{},
		&models.TicketAttachment{},
		&models.TicketStatusHistory{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate tickets: %w", err)
	}
	return nil
}
