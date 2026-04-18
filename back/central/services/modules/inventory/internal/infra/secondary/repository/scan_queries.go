package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) RecordScanEvent(ctx context.Context, event *entities.ScanEvent) (*entities.ScanEvent, error) {
	m := &models.ScanEvent{
		BusinessID:  event.BusinessID,
		UserID:      event.UserID,
		DeviceID:    event.DeviceID,
		ScannedCode: event.ScannedCode,
		CodeType:    event.CodeType,
		Action:      event.Action,
		ScannedAt:   event.ScannedAt,
	}
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.ScanEventModelToEntity(m), nil
}

func (r *Repository) ResolveScanCode(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error) {
	db := r.db.Conn(ctx)

	var lpn models.LicensePlate
	if err := db.Where("business_id = ? AND code = ? AND status = 'active'", businessID, code).First(&lpn).Error; err == nil {
		id := lpn.ID
		return &entities.ScanResolution{
			Code:       code,
			CodeType:   "lpn",
			MatchedID:  &id,
			LpnID:      &id,
			LocationID: lpn.CurrentLocationID,
			Suggested:  "view_lpn",
		}, nil
	}

	var location models.WarehouseLocation
	if err := db.Joins("INNER JOIN warehouses wh ON wh.id = warehouse_locations.warehouse_id AND wh.business_id = ?", businessID).
		Where("warehouse_locations.code = ? AND warehouse_locations.deleted_at IS NULL", code).
		First(&location).Error; err == nil {
		id := location.ID
		return &entities.ScanResolution{
			Code:       code,
			CodeType:   "location",
			MatchedID:  &id,
			LocationID: &id,
			Suggested:  "view_location",
		}, nil
	}

	var serial models.InventorySerial
	if err := db.Where("business_id = ? AND serial_number = ?", businessID, code).First(&serial).Error; err == nil {
		id := serial.ID
		return &entities.ScanResolution{
			Code:       code,
			CodeType:   "serial",
			MatchedID:  &id,
			SerialID:   &id,
			LotID:      serial.LotID,
			LocationID: serial.CurrentLocationID,
			ProductID:  serial.ProductID,
			Suggested:  "view_serial",
		}, nil
	}

	var lot models.InventoryLot
	if err := db.Where("business_id = ? AND lot_code = ?", businessID, code).First(&lot).Error; err == nil {
		id := lot.ID
		return &entities.ScanResolution{
			Code:      code,
			CodeType:  "lot",
			MatchedID: &id,
			LotID:     &id,
			ProductID: lot.ProductID,
			Suggested: "view_lot",
		}, nil
	}

	var pu models.ProductUoM
	if err := db.Where("business_id = ? AND barcode = ?", businessID, code).First(&pu).Error; err == nil {
		return &entities.ScanResolution{
			Code:      code,
			CodeType:  "sku",
			ProductID: pu.ProductID,
			Suggested: "view_product",
		}, nil
	}

	var product struct {
		ID string
	}
	if err := db.Table("products").
		Select("id").
		Where("business_id = ? AND (sku = ? OR id = ?) AND deleted_at IS NULL", businessID, code, code).
		First(&product).Error; err == nil {
		return &entities.ScanResolution{
			Code:      code,
			CodeType:  "sku",
			ProductID: product.ID,
			Suggested: "view_product",
		}, nil
	}

	return nil, nil
}
