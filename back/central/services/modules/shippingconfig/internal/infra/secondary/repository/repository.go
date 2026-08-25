package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) domain.IRepository {
	return &Repository{db: database}
}

func (r *Repository) GetConfig(ctx context.Context, businessID uint, warehouseID *uint) (*domain.ShippingConfig, error) {
	query := r.db.Conn(ctx).Model(&models.ShippingConfig{}).Where("business_id = ?", businessID)
	if warehouseID == nil {
		query = query.Where("warehouse_id IS NULL")
	} else {
		query = query.Where("warehouse_id = ?", *warehouseID)
	}

	var row models.ShippingConfig
	if err := query.First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	cfg := toDomain(&row)
	return &cfg, nil
}

func (r *Repository) ListConfigs(ctx context.Context, businessID uint) ([]domain.ShippingConfig, error) {
	var rows []models.ShippingConfig
	if err := r.db.Conn(ctx).Where("business_id = ?", businessID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.ShippingConfig, 0, len(rows))
	for i := range rows {
		out = append(out, toDomain(&rows[i]))
	}
	return out, nil
}

func (r *Repository) UpsertConfig(ctx context.Context, cfg *domain.ShippingConfig) error {
	boxes, err := json.Marshal(cfg.Boxes)
	if err != nil {
		return err
	}
	carriers, err := json.Marshal(cfg.Carriers)
	if err != nil {
		return err
	}

	query := r.db.Conn(ctx).Model(&models.ShippingConfig{}).Where("business_id = ?", cfg.BusinessID)
	if cfg.WarehouseID == nil {
		query = query.Where("warehouse_id IS NULL")
	} else {
		query = query.Where("warehouse_id = ?", *cfg.WarehouseID)
	}

	var existing models.ShippingConfig
	err = query.First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		row := models.ShippingConfig{
			BusinessID:      cfg.BusinessID,
			WarehouseID:     cfg.WarehouseID,
			PackageStrategy: cfg.PackageStrategy,
			Boxes:           datatypes.JSON(boxes),
			Carriers:        datatypes.JSON(carriers),
		}
		if err := r.db.Conn(ctx).Create(&row).Error; err != nil {
			return err
		}
		cfg.ID = row.ID
		cfg.CreatedAt = row.CreatedAt
		cfg.UpdatedAt = row.UpdatedAt
		return nil
	}

	updates := map[string]interface{}{
		"package_strategy": cfg.PackageStrategy,
		"boxes":            datatypes.JSON(boxes),
		"carriers":         datatypes.JSON(carriers),
	}
	if err := r.db.Conn(ctx).Model(&models.ShippingConfig{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
		return err
	}
	cfg.ID = existing.ID
	cfg.CreatedAt = existing.CreatedAt
	return nil
}

func (r *Repository) DeleteConfig(ctx context.Context, businessID uint, warehouseID uint) error {
	return r.db.Conn(ctx).
		Where("business_id = ? AND warehouse_id = ?", businessID, warehouseID).
		Delete(&models.ShippingConfig{}).Error
}

func (r *Repository) ListWarehouseOrigins(ctx context.Context, businessID uint) ([]domain.WarehouseOrigin, error) {
	var rows []struct {
		ID        uint
		Name      string
		Address   string
		City      string
		State     string
		Phone     string
		IsDefault bool
		IsActive  bool
	}

	err := r.db.Conn(ctx).
		Table("warehouses").
		Select("id, name, address, city, state, phone, is_default, is_active").
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Order("is_default DESC, id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]domain.WarehouseOrigin, 0, len(rows))
	for _, w := range rows {
		out = append(out, domain.WarehouseOrigin{
			ID:        w.ID,
			Name:      w.Name,
			Address:   w.Address,
			City:      w.City,
			State:     w.State,
			Phone:     w.Phone,
			IsDefault: w.IsDefault,
			IsActive:  w.IsActive,
		})
	}
	return out, nil
}

func (r *Repository) SetDefaultWarehouse(ctx context.Context, businessID uint, warehouseID uint) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("warehouses").
			Where("business_id = ? AND deleted_at IS NULL", businessID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Table("warehouses").
			Where("business_id = ? AND id = ? AND deleted_at IS NULL", businessID, warehouseID).
			Update("is_default", true).Error
	})
}

func toDomain(row *models.ShippingConfig) domain.ShippingConfig {
	cfg := domain.ShippingConfig{
		ID:              row.ID,
		BusinessID:      row.BusinessID,
		WarehouseID:     row.WarehouseID,
		PackageStrategy: row.PackageStrategy,
		Boxes:           []domain.Box{},
		Carriers:        []domain.CarrierSetting{},
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
	if len(row.Boxes) > 0 {
		_ = json.Unmarshal(row.Boxes, &cfg.Boxes)
	}
	if len(row.Carriers) > 0 {
		_ = json.Unmarshal(row.Carriers, &cfg.Carriers)
	}
	return cfg
}
