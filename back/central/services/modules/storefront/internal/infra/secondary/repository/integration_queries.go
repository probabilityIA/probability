package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var defaultCatalogLayout = entities.StorefrontCatalogLayout{Columns: 4, Rows: 3}

func (r *Repository) IsIntegrationActiveOrMissing(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error) {
	var result struct {
		IsActive bool
	}

	err := r.db.Conn(ctx).
		Table("integrations").
		Select("is_active").
		Where("business_id = ? AND integration_type_id = ? AND deleted_at IS NULL", businessID, integrationTypeID).
		Limit(1).
		First(&result).Error

	if err == gorm.ErrRecordNotFound {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	return result.IsActive, nil
}

func (r *Repository) getIntegrationConfigMap(ctx context.Context, businessID, integrationTypeID uint) (map[string]interface{}, error) {
	var row struct {
		Config datatypes.JSON
	}

	err := r.db.Conn(ctx).
		Table("integrations").
		Select("config").
		Where("business_id = ? AND integration_type_id = ? AND deleted_at IS NULL", businessID, integrationTypeID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	configMap := map[string]interface{}{}
	if len(row.Config) > 0 {
		if err := json.Unmarshal(row.Config, &configMap); err != nil {
			return map[string]interface{}{}, nil
		}
	}
	return configMap, nil
}

func (r *Repository) upsertIntegrationConfigKey(ctx context.Context, businessID, integrationTypeID, requesterUserID uint, key string, value interface{}) error {
	var row struct {
		ID     uint
		Config datatypes.JSON
	}

	err := r.db.Conn(ctx).
		Table("integrations").
		Select("id, config").
		Where("business_id = ? AND integration_type_id = ? AND deleted_at IS NULL", businessID, integrationTypeID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return err
	}

	if row.ID == 0 {
		return r.createCatalogIntegration(ctx, businessID, integrationTypeID, requesterUserID, map[string]interface{}{key: value})
	}

	configMap := map[string]interface{}{}
	if len(row.Config) > 0 {
		if err := json.Unmarshal(row.Config, &configMap); err != nil {
			configMap = map[string]interface{}{}
		}
	}
	configMap[key] = value

	updated, err := json.Marshal(configMap)
	if err != nil {
		return err
	}

	return r.db.Conn(ctx).
		Table("integrations").
		Where("id = ?", row.ID).
		Update("config", datatypes.JSON(updated)).Error
}

func (r *Repository) createCatalogIntegration(ctx context.Context, businessID, integrationTypeID, requesterUserID uint, configMap map[string]interface{}) error {
	config, err := json.Marshal(configMap)
	if err != nil {
		return err
	}

	integration := models.Integration{
		Name:              "Catalogo",
		Code:              fmt.Sprintf("tienda_%d", businessID),
		Category:          "storefront",
		IntegrationTypeID: integrationTypeID,
		BusinessID:        &businessID,
		IsActive:          true,
		IsDefault:         true,
		Config:            datatypes.JSON(config),
		CreatedByID:       requesterUserID,
	}
	return r.db.Conn(ctx).Create(&integration).Error
}

func (r *Repository) GetCatalogLayout(ctx context.Context, businessID, integrationTypeID uint) (entities.StorefrontCatalogLayout, error) {
	configMap, err := r.getIntegrationConfigMap(ctx, businessID, integrationTypeID)
	if err != nil {
		return defaultCatalogLayout, err
	}

	raw, ok := configMap["catalog_layout"]
	if !ok {
		return defaultCatalogLayout, nil
	}
	encoded, err := json.Marshal(raw)
	if err != nil {
		return defaultCatalogLayout, nil
	}

	var parsed struct {
		Columns int `json:"columns"`
		Rows    int `json:"rows"`
	}
	if err := json.Unmarshal(encoded, &parsed); err != nil {
		return defaultCatalogLayout, nil
	}

	layout := entities.StorefrontCatalogLayout{Columns: parsed.Columns, Rows: parsed.Rows}
	if layout.Columns < 1 {
		layout.Columns = defaultCatalogLayout.Columns
	}
	if layout.Rows < 1 {
		layout.Rows = defaultCatalogLayout.Rows
	}
	return layout, nil
}

func (r *Repository) UpdateCatalogLayout(ctx context.Context, businessID, integrationTypeID, requesterUserID uint, layout entities.StorefrontCatalogLayout) error {
	return r.upsertIntegrationConfigKey(ctx, businessID, integrationTypeID, requesterUserID, "catalog_layout", map[string]interface{}{
		"columns": layout.Columns,
		"rows":    layout.Rows,
	})
}

func (r *Repository) GetCatalogBanner(ctx context.Context, businessID, integrationTypeID uint) (entities.StorefrontBanner, error) {
	configMap, err := r.getIntegrationConfigMap(ctx, businessID, integrationTypeID)
	if err != nil {
		return entities.StorefrontBanner{}, err
	}

	raw, ok := configMap["banner"]
	if !ok {
		return entities.StorefrontBanner{}, nil
	}
	encoded, err := json.Marshal(raw)
	if err != nil {
		return entities.StorefrontBanner{}, nil
	}

	var parsed struct {
		Enabled  bool   `json:"enabled"`
		ImageURL string `json:"image_url"`
	}
	if err := json.Unmarshal(encoded, &parsed); err != nil {
		return entities.StorefrontBanner{}, nil
	}
	return entities.StorefrontBanner{Enabled: parsed.Enabled, ImageURL: parsed.ImageURL}, nil
}

func (r *Repository) UpdateCatalogBanner(ctx context.Context, businessID, integrationTypeID, requesterUserID uint, banner entities.StorefrontBanner) error {
	return r.upsertIntegrationConfigKey(ctx, businessID, integrationTypeID, requesterUserID, "banner", map[string]interface{}{
		"enabled":   banner.Enabled,
		"image_url": banner.ImageURL,
	})
}
