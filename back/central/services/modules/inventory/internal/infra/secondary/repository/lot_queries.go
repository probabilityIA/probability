package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateLot(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
	model := mappers.LotEntityToModel(lot)
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return mappers.LotModelToEntity(model), nil
}

func (r *Repository) GetLotByID(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
	var m models.InventoryLot
	err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", lotID, businessID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrLotNotFound
		}
		return nil, err
	}
	return mappers.LotModelToEntity(&m), nil
}

func (r *Repository) ListLots(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error) {
	var modelsList []models.InventoryLot
	var total int64

	query := r.db.Conn(ctx).Model(&models.InventoryLot{}).Where("business_id = ?", params.BusinessID)
	if params.ProductID != "" {
		query = query.Where("product_id = ?", params.ProductID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.ExpiringInDays > 0 {
		threshold := time.Now().AddDate(0, 0, params.ExpiringInDays)
		query = query.Where("expiration_date IS NOT NULL AND expiration_date <= ?", threshold)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("expiration_date ASC, id DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	lots := make([]entities.InventoryLot, len(modelsList))
	for i := range modelsList {
		lots[i] = *mappers.LotModelToEntity(&modelsList[i])
	}
	return lots, total, nil
}

func (r *Repository) UpdateLot(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
	updates := map[string]any{
		"lot_code":         lot.LotCode,
		"manufacture_date": lot.ManufactureDate,
		"expiration_date":  lot.ExpirationDate,
		"received_at":      lot.ReceivedAt,
		"supplier_id":      lot.SupplierID,
		"status":           lot.Status,
	}
	res := r.db.Conn(ctx).Model(&models.InventoryLot{}).
		Where("id = ? AND business_id = ?", lot.ID, lot.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrLotNotFound
	}
	return r.GetLotByID(ctx, lot.BusinessID, lot.ID)
}

func (r *Repository) DeleteLot(ctx context.Context, businessID, lotID uint) error {
	res := r.db.Conn(ctx).Where("id = ? AND business_id = ?", lotID, businessID).Delete(&models.InventoryLot{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrLotNotFound
	}
	return nil
}

func (r *Repository) LotExistsByCode(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.InventoryLot{}).
		Where("business_id = ? AND product_id = ? AND lot_code = ?", businessID, productID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *Repository) ListLotsForReserve(ctx context.Context, productID string, warehouseID, businessID uint, strategy string) ([]entities.InventoryLot, error) {
	var modelsList []models.InventoryLot

	query := r.db.Conn(ctx).Model(&models.InventoryLot{}).
		Where("business_id = ? AND product_id = ? AND status = ?", businessID, productID, "active")

	switch strategy {
	case "fefo":
		query = query.Where("expiration_date IS NOT NULL").Order("expiration_date ASC, id ASC")
	case "fifo":
		query = query.Order("received_at ASC NULLS LAST, id ASC")
	default:
		query = query.Order("id ASC")
	}

	if err := query.Find(&modelsList).Error; err != nil {
		return nil, err
	}

	lots := make([]entities.InventoryLot, len(modelsList))
	for i := range modelsList {
		lots[i] = *mappers.LotModelToEntity(&modelsList[i])
	}
	return lots, nil
}
