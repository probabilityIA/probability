package repository

import (
	"fmt"
	"sync"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	db                 db.IDatabase
	cache              IInventoryCache
	availableStateID   uint
	availableStateOnce sync.Once
	availableStateErr  error
}

func New(database db.IDatabase, cache IInventoryCache) ports.IRepository {
	return &Repository{db: database, cache: cache}
}

func (r *Repository) resolveAvailableStateID(tx *gorm.DB) (uint, error) {
	r.availableStateOnce.Do(func() {
		var s models.InventoryState
		if err := tx.Where("code = ?", "available").First(&s).Error; err != nil {
			r.availableStateErr = fmt.Errorf("inventory state 'available' not found: %w", err)
			return
		}
		r.availableStateID = s.ID
	})
	if r.availableStateErr != nil {
		return 0, r.availableStateErr
	}
	return r.availableStateID, nil
}
