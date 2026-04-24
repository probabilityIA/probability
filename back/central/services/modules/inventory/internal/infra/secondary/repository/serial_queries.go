package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateSerial(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
	model := mappers.SerialEntityToModel(serial)
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return mappers.SerialModelToEntity(model), nil
}

func (r *Repository) GetSerialByID(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
	var m models.InventorySerial
	err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", serialID, businessID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrSerialNotFound
		}
		return nil, err
	}
	return mappers.SerialModelToEntity(&m), nil
}

func (r *Repository) ListSerials(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error) {
	var modelsList []models.InventorySerial
	var total int64

	query := r.db.Conn(ctx).Model(&models.InventorySerial{}).Where("business_id = ?", params.BusinessID)
	if params.ProductID != "" {
		query = query.Where("product_id = ?", params.ProductID)
	}
	if params.LotID != nil {
		query = query.Where("lot_id = ?", *params.LotID)
	}
	if params.StateID != nil {
		query = query.Where("current_state_id = ?", *params.StateID)
	}
	if params.LocationID != nil {
		query = query.Where("current_location_id = ?", *params.LocationID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	serials := make([]entities.InventorySerial, len(modelsList))
	for i := range modelsList {
		serials[i] = *mappers.SerialModelToEntity(&modelsList[i])
	}
	return serials, total, nil
}

func (r *Repository) UpdateSerial(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error) {
	updates := map[string]any{
		"lot_id":              serial.LotID,
		"current_location_id": serial.CurrentLocationID,
		"current_state_id":    serial.CurrentStateID,
		"received_at":         serial.ReceivedAt,
		"sold_at":             serial.SoldAt,
	}
	res := r.db.Conn(ctx).Model(&models.InventorySerial{}).
		Where("id = ? AND business_id = ?", serial.ID, serial.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrSerialNotFound
	}
	return r.GetSerialByID(ctx, serial.BusinessID, serial.ID)
}

func (r *Repository) SerialExists(ctx context.Context, businessID uint, productID, serial string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.InventorySerial{}).
		Where("business_id = ? AND product_id = ? AND serial_number = ?", businessID, productID, serial)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
