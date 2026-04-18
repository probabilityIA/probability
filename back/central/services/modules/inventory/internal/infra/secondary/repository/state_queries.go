package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) ListInventoryStates(ctx context.Context) ([]entities.InventoryState, error) {
	var modelsList []models.InventoryState
	if err := r.db.Conn(ctx).Where("is_active = ?", true).Order("id ASC").Find(&modelsList).Error; err != nil {
		return nil, err
	}
	states := make([]entities.InventoryState, len(modelsList))
	for i := range modelsList {
		states[i] = *mappers.StateModelToEntity(&modelsList[i])
	}
	return states, nil
}

func (r *Repository) GetInventoryStateByCode(ctx context.Context, code string) (*entities.InventoryState, error) {
	var m models.InventoryState
	err := r.db.Conn(ctx).Where("code = ?", code).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrStateNotFound
		}
		return nil, err
	}
	return mappers.StateModelToEntity(&m), nil
}
