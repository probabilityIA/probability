package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) RevokeOverride(ctx context.Context, businessID uint, moduleCode string, actorUserID uint) error {
	if err := uc.repo.DeleteOverride(ctx, businessID, moduleCode); err != nil {
		return err
	}

	uc.recordAudit(ctx, businessID, actorUserID, entities.AuditActionOverrideRevoked,
		fmt.Sprintf("revoco el acceso al modulo %s", moduleCode))
	return nil
}
