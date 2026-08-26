package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) domain.IRepository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) CreateShipment(ctx context.Context, shipment *domain.Shipment) error {
	dbShipment := mappers.ToDBShipment(shipment)
	if err := r.db.Conn(ctx).Create(dbShipment).Error; err != nil {
		return err
	}
	shipment.ID = dbShipment.ID
	return nil
}

func (r *Repository) GetShipmentByID(ctx context.Context, id uint) (*domain.Shipment, error) {
	var shipment models.Shipment
	err := r.db.Conn(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("id = ?", id).
		First(&shipment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShipmentNotFound
		}
		return nil, err
	}

	return mappers.ToDomainShipment(&shipment), nil
}

func (r *Repository) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	var shipment models.Shipment
	err := r.db.Conn(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("tracking_number = ?", trackingNumber).
		First(&shipment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrShipmentNotFound
		}
		return nil, err
	}

	return mappers.ToDomainShipment(&shipment), nil
}

func (r *Repository) GetShipmentsByOrderID(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	var shipments []models.Shipment
	err := r.db.Conn(ctx).
		Preload("Order").
		Preload("ShippingAddress").
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&shipments).Error

	if err != nil {
		return nil, err
	}

	domainShipments := make([]domain.Shipment, len(shipments))
	for i, shipment := range shipments {
		domainShipments[i] = *mappers.ToDomainShipment(&shipment)
	}

	return domainShipments, nil
}

func (r *Repository) ListShipments(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error) {
	var shipments []models.Shipment
	var total int64

	query := r.db.Conn(ctx).Model(&models.Shipment{})

	if orderID, ok := filters["order_id"].(string); ok && orderID != "" {
		query = query.Where("shipments.order_id = ?", orderID)
	}

	if orderIDs, ok := filters["order_ids"].([]string); ok && len(orderIDs) > 0 {
		query = query.Where("shipments.order_id IN ?", orderIDs)
	}

	if trackingNumber, ok := filters["tracking_number"].(string); ok && trackingNumber != "" {
		query = query.Where("shipments.tracking_number ILIKE ?", "%"+trackingNumber+"%")
	}

	if trackingNumbers, ok := filters["tracking_numbers"].([]string); ok && len(trackingNumbers) > 0 {
		query = query.Where("shipments.tracking_number IN ?", trackingNumbers)
	}

	if carrier, ok := filters["carrier"].(string); ok && carrier != "" {
		query = query.Where("shipments.carrier ILIKE ?", "%"+carrier+"%")
	}

	if carrierCode, ok := filters["carrier_code"].(string); ok && carrierCode != "" {
		query = query.Where("shipments.carrier_code = ?", carrierCode)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("shipments.status = ?", status)
	}

	if statuses, ok := filters["statuses"].([]string); ok && len(statuses) > 0 {
		query = query.Where("shipments.status IN ?", statuses)
	}

	if guideID, ok := filters["guide_id"].(string); ok && guideID != "" {
		query = query.Where("shipments.guide_id = ?", guideID)
	}

	if customerName, ok := filters["customer_name"].(string); ok && customerName != "" {
		query = query.Where("EXISTS (SELECT 1 FROM orders WHERE orders.id = shipments.order_id AND orders.customer_name ILIKE ?)", "%"+customerName+"%")
	}

	if orderNumber, ok := filters["order_number"].(string); ok && orderNumber != "" {
		query = query.Where("EXISTS (SELECT 1 FROM orders WHERE orders.id = shipments.order_id AND orders.order_number ILIKE ?)", "%"+orderNumber+"%")
	}

	if warehouseID, ok := filters["warehouse_id"].(uint); ok && warehouseID > 0 {
		query = query.Where("shipments.warehouse_id = ?", warehouseID)
	}

	if driverID, ok := filters["driver_id"].(uint); ok && driverID > 0 {
		query = query.Where("shipments.driver_id = ?", driverID)
	}

	if isLastMile, ok := filters["is_last_mile"].(bool); ok {
		query = query.Where("shipments.is_last_mile = ?", isLastMile)
	}

	if isTest, ok := filters["is_test"].(bool); ok {
		query = query.Where("shipments.is_test = ?", isTest)
	}

	if shippedAfter, ok := filters["shipped_after"].(string); ok && shippedAfter != "" {
		query = query.Where("shipments.shipped_at >= ?", shippedAfter)
	}

	if shippedBefore, ok := filters["shipped_before"].(string); ok && shippedBefore != "" {
		query = query.Where("shipments.shipped_at <= ?", shippedBefore)
	}

	if deliveredAfter, ok := filters["delivered_after"].(string); ok && deliveredAfter != "" {
		query = query.Where("shipments.delivered_at >= ?", deliveredAfter)
	}

	if deliveredBefore, ok := filters["delivered_before"].(string); ok && deliveredBefore != "" {
		query = query.Where("shipments.delivered_at <= ?", deliveredBefore)
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("shipments.created_at >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("shipments.created_at <= ?", endDate)
	}

	if createdAfter, ok := filters["created_after"].(string); ok && createdAfter != "" {
		query = query.Where("shipments.created_at >= ?", createdAfter)
	}

	if createdBefore, ok := filters["created_before"].(string); ok && createdBefore != "" {
		query = query.Where("shipments.created_at <= ?", createdBefore)
	}

	if updatedAfter, ok := filters["updated_after"].(string); ok && updatedAfter != "" {
		query = query.Where("shipments.updated_at >= ?", updatedAfter)
	}

	if updatedBefore, ok := filters["updated_before"].(string); ok && updatedBefore != "" {
		query = query.Where("shipments.updated_at <= ?", updatedBefore)
	}

	if businessID, ok := filters["business_id"].(uint); ok && businessID > 0 {
		query = query.Where("EXISTS (SELECT 1 FROM orders WHERE orders.id = shipments.order_id AND orders.business_id = ?)", businessID)
	}

	if integrationID, ok := filters["integration_id"].(uint); ok && integrationID > 0 {
		query = query.Where("EXISTS (SELECT 1 FROM orders WHERE orders.id = shipments.order_id AND orders.integration_id = ?)", integrationID)
	}

	if integrationType, ok := filters["integration_type"].(string); ok && integrationType != "" {
		query = query.Where("EXISTS (SELECT 1 FROM orders WHERE orders.id = shipments.order_id AND orders.integration_type = ?)", integrationType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "shipments.created_at"
	if sort, ok := filters["sort_by"].(string); ok && sort != "" {
		sortFieldMap := map[string]string{
			"id":              "shipments.id",
			"order_id":        "shipments.order_id",
			"tracking_number": "shipments.tracking_number",
			"status":          "shipments.status",
			"carrier":         "shipments.carrier",
			"shipped_at":      "shipments.shipped_at",
			"delivered_at":    "shipments.delivered_at",
			"created_at":      "shipments.created_at",
			"updated_at":      "shipments.updated_at",
			"warehouse_id":    "shipments.warehouse_id",
			"driver_id":       "shipments.driver_id",
		}
		if mappedField, exists := sortFieldMap[sort]; exists {
			sortBy = mappedField
		}
	}

	sortOrder := "desc"
	if order, ok := filters["sort_order"].(string); ok && order != "" {
		sortOrder = order
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	query = query.Preload("Order").Preload("ShippingAddress")

	if err := query.Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	domainShipments := make([]domain.Shipment, len(shipments))
	businessIDs := map[uint]bool{}
	for i, shipment := range shipments {
		domainShipments[i] = *mappers.ToDomainShipment(&shipment)
		if shipment.Order != nil && shipment.Order.BusinessID != nil && *shipment.Order.BusinessID > 0 {
			businessIDs[*shipment.Order.BusinessID] = true
		}
	}

	if len(businessIDs) > 0 {
		ids := make([]uint, 0, len(businessIDs))
		for id := range businessIDs {
			ids = append(ids, id)
		}
		var origins []models.OriginAddress
		if err := r.db.Conn(ctx).
			Where("business_id IN ? AND is_default = true AND deleted_at IS NULL", ids).
			Find(&origins).Error; err == nil {
			defaults := map[uint]models.OriginAddress{}
			for _, o := range origins {
				defaults[o.BusinessID] = o
			}
			for i := range domainShipments {
				bid := uint(0)
				if shipments[i].Order != nil && shipments[i].Order.BusinessID != nil {
					bid = *shipments[i].Order.BusinessID
				}
				if origin, ok := defaults[bid]; ok {
					domainShipments[i].OriginAddress = origin.Street
					domainShipments[i].OriginCity = origin.City
					domainShipments[i].OriginState = origin.State
				}
			}
		}
	}

	return domainShipments, total, nil
}

func (r *Repository) UpdateShipment(ctx context.Context, shipment *domain.Shipment) error {
	dbShipment := mappers.ToDBShipment(shipment)
	return r.db.Conn(ctx).Save(dbShipment).Error
}

func (r *Repository) WithOrderGuideLock(ctx context.Context, orderID string, fn func() error) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext(?)::bigint)", orderID).Error; err != nil {
			return err
		}
		return fn()
	})
}

func (r *Repository) MarkShipmentGenerating(ctx context.Context, id uint, staleBefore time.Time) (bool, error) {
	res := r.db.Conn(ctx).Model(&models.Shipment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Where("status <> ? OR updated_at < ?", domain.ShipmentStatusGenerating, staleBefore).
		Updates(map[string]interface{}{
			"status":     domain.ShipmentStatusGenerating,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *Repository) ReleaseShipmentGenerating(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Model(&models.Shipment{}).
		Where("id = ? AND status = ?", id, domain.ShipmentStatusGenerating).
		Updates(map[string]interface{}{
			"status":     domain.ShipmentStatusPending,
			"updated_at": time.Now(),
		}).Error
}

func (r *Repository) DeleteShipment(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Where("id = ?", id).Delete(&models.Shipment{}).Error
}

func (r *Repository) ShipmentExists(ctx context.Context, orderID string, trackingNumber string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Model(&models.Shipment{}).
		Where("order_id = ? AND tracking_number = ?", orderID, trackingNumber).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) CreateOriginAddress(ctx context.Context, address *domain.OriginAddress) error {
	dbAddress := &models.OriginAddress{
		BusinessID:   address.BusinessID,
		Alias:        address.Alias,
		Company:      address.Company,
		FirstName:    address.FirstName,
		LastName:     address.LastName,
		Email:        address.Email,
		Phone:        address.Phone,
		Street:       address.Street,
		Suburb:       address.Suburb,
		CityDaneCode: address.CityDaneCode,
		City:         address.City,
		State:        address.State,
		PostalCode:   address.PostalCode,
		IsDefault:    address.IsDefault,
	}

	if err := r.db.Conn(ctx).Create(dbAddress).Error; err != nil {
		return err
	}
	address.ID = dbAddress.ID
	return nil
}

func (r *Repository) GetOriginAddressByID(ctx context.Context, id uint) (*domain.OriginAddress, error) {
	var dbAddress models.OriginAddress
	if err := r.db.Conn(ctx).First(&dbAddress, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dirección de origen no encontrada")
		}
		return nil, err
	}

	return &domain.OriginAddress{
		ID:           dbAddress.ID,
		CreatedAt:    dbAddress.CreatedAt,
		UpdatedAt:    dbAddress.UpdatedAt,
		DeletedAt:    dbAddress.DeletedAt,
		BusinessID:   dbAddress.BusinessID,
		Alias:        dbAddress.Alias,
		Company:      dbAddress.Company,
		FirstName:    dbAddress.FirstName,
		LastName:     dbAddress.LastName,
		Email:        dbAddress.Email,
		Phone:        dbAddress.Phone,
		Street:       dbAddress.Street,
		Suburb:       dbAddress.Suburb,
		CityDaneCode: dbAddress.CityDaneCode,
		City:         dbAddress.City,
		State:        dbAddress.State,
		PostalCode:   dbAddress.PostalCode,
		IsDefault:    dbAddress.IsDefault,
	}, nil
}

func (r *Repository) ListOriginAddressesByBusiness(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
	var dbAddresses []models.OriginAddress
	if err := r.db.Conn(ctx).Where("business_id = ?", businessID).Find(&dbAddresses).Error; err != nil {
		return nil, err
	}

	addresses := make([]domain.OriginAddress, len(dbAddresses))
	for i, dbAddress := range dbAddresses {
		addresses[i] = domain.OriginAddress{
			ID:           dbAddress.ID,
			CreatedAt:    dbAddress.CreatedAt,
			UpdatedAt:    dbAddress.UpdatedAt,
			DeletedAt:    dbAddress.DeletedAt,
			BusinessID:   dbAddress.BusinessID,
			Alias:        dbAddress.Alias,
			Company:      dbAddress.Company,
			FirstName:    dbAddress.FirstName,
			LastName:     dbAddress.LastName,
			Email:        dbAddress.Email,
			Phone:        dbAddress.Phone,
			Street:       dbAddress.Street,
			Suburb:       dbAddress.Suburb,
			CityDaneCode: dbAddress.CityDaneCode,
			City:         dbAddress.City,
			State:        dbAddress.State,
			PostalCode:   dbAddress.PostalCode,
			IsDefault:    dbAddress.IsDefault,
		}
	}
	return addresses, nil
}

func (r *Repository) GetDefaultOriginAddress(ctx context.Context, businessID uint) (*domain.OriginAddress, error) {
	var dbAddress models.OriginAddress
	if err := r.db.Conn(ctx).Where("business_id = ? AND is_default = ?", businessID, true).First(&dbAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.OriginAddress{
		ID:           dbAddress.ID,
		CreatedAt:    dbAddress.CreatedAt,
		UpdatedAt:    dbAddress.UpdatedAt,
		DeletedAt:    dbAddress.DeletedAt,
		BusinessID:   dbAddress.BusinessID,
		Alias:        dbAddress.Alias,
		Company:      dbAddress.Company,
		FirstName:    dbAddress.FirstName,
		LastName:     dbAddress.LastName,
		Email:        dbAddress.Email,
		Phone:        dbAddress.Phone,
		Street:       dbAddress.Street,
		Suburb:       dbAddress.Suburb,
		CityDaneCode: dbAddress.CityDaneCode,
		City:         dbAddress.City,
		State:        dbAddress.State,
		PostalCode:   dbAddress.PostalCode,
		IsDefault:    dbAddress.IsDefault,
	}, nil
}

func (r *Repository) UpdateOriginAddress(ctx context.Context, address *domain.OriginAddress) error {
	dbAddress := &models.OriginAddress{
		BusinessID:   address.BusinessID,
		Alias:        address.Alias,
		Company:      address.Company,
		FirstName:    address.FirstName,
		LastName:     address.LastName,
		Email:        address.Email,
		Phone:        address.Phone,
		Street:       address.Street,
		Suburb:       address.Suburb,
		CityDaneCode: address.CityDaneCode,
		City:         address.City,
		State:        address.State,
		PostalCode:   address.PostalCode,
		IsDefault:    address.IsDefault,
	}
	dbAddress.ID = address.ID
	dbAddress.CreatedAt = address.CreatedAt
	dbAddress.UpdatedAt = address.UpdatedAt
	dbAddress.DeletedAt = address.DeletedAt

	return r.db.Conn(ctx).Save(dbAddress).Error
}

func (r *Repository) DeleteOriginAddress(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.OriginAddress{}, id).Error
}

func (r *Repository) SetDefaultOriginAddress(ctx context.Context, businessID, addressID uint) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.OriginAddress{}).Where("business_id = ?", businessID).Update("is_default", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.OriginAddress{}).Where("id = ? AND business_id = ?", addressID, businessID).Update("is_default", true).Error; err != nil {
			return err
		}
		return nil
	})
}
