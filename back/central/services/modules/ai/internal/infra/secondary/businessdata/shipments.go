package businessdata

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

const maxTrackingEvents = 5

type trackingMetadata struct {
	TrackingEvents []struct {
		Date        string `json:"date"`
		Status      string `json:"status"`
		RawStatus   string `json:"raw_status"`
		Description string `json:"description"`
	} `json:"tracking_events"`
}

func (r *Reader) FindShipments(ctx context.Context, businessID uint, trackingNumber string) ([]entities.ShipmentInfo, error) {
	clean := strings.TrimSpace(trackingNumber)
	if clean == "" {
		return nil, nil
	}

	var rows []models.Shipment
	err := r.db.Conn(ctx).Model(&models.Shipment{}).
		Select("shipments.*").
		Joins("JOIN orders ON orders.id = shipments.order_id").
		Where("orders.business_id = ? AND orders.deleted_at IS NULL", businessID).
		Where("(shipments.tracking_number = ? OR shipments.guide_id = ?)", clean, clean).
		Preload("Order", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_number", "internal_number")
		}).
		Order("shipments.created_at DESC").
		Limit(maxMatches).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]entities.ShipmentInfo, 0, len(rows))
	for _, row := range rows {
		number := ""
		if row.Order != nil {
			number = orderNumber(*row.Order)
		}
		out = append(out, toShipmentInfo(row, number))
	}
	return out, nil
}

func toShipmentInfo(row models.Shipment, orderNumber string) entities.ShipmentInfo {
	info := entities.ShipmentInfo{
		TrackingNumber:      deref(row.TrackingNumber),
		Carrier:             deref(row.Carrier),
		Status:              row.Status,
		CarrierStatus:       deref(row.CarrierStatus),
		CarrierStatusDetail: deref(row.CarrierStatusDetail),
		CodCollectAmount:    row.CodCollectAmount,
		TotalCost:           row.TotalCost,
		DestinationCity:     row.DestinationCity,
		OrderNumber:         orderNumber,
		IsTest:              row.IsTest,
		CreatedAt:           row.CreatedAt,
		ShippedAt:           row.ShippedAt,
		DeliveredAt:         row.DeliveredAt,
		EstimatedDelivery:   row.EstimatedDelivery,
	}
	if info.TrackingNumber == "" {
		info.TrackingNumber = deref(row.GuideID)
	}

	if len(row.Metadata) > 0 {
		var meta trackingMetadata
		if err := json.Unmarshal(row.Metadata, &meta); err == nil {
			events := meta.TrackingEvents
			for i := len(events) - 1; i >= 0 && len(info.Events) < maxTrackingEvents; i-- {
				info.Events = append(info.Events, entities.TrackingEvent{
					Date:        events[i].Date,
					Status:      events[i].Status,
					RawStatus:   events[i].RawStatus,
					Description: events[i].Description,
				})
			}
		}
	}
	return info
}
