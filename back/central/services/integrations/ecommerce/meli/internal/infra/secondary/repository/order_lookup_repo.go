package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type OrderLookupRepo struct {
	db  db.IDatabase
	log log.ILogger
}

func NewOrderLookup(database db.IDatabase, logger log.ILogger) domain.IOrderLookupRepository {
	return &OrderLookupRepo{
		db:  database,
		log: logger.WithModule("meli.order_lookup"),
	}
}

func (r *OrderLookupRepo) GetMeliShipmentByOrderID(ctx context.Context, orderID string) (*domain.MeliOrderRef, error) {
	var row struct {
		IntegrationID uint
		GuideID       string
	}
	err := r.db.Conn(ctx).
		Table("orders AS o").
		Select("o.integration_id, s.guide_id").
		Joins("JOIN shipments s ON s.order_id = o.id AND s.deleted_at IS NULL").
		Where("o.id = ? AND o.integration_type = ? AND o.deleted_at IS NULL", orderID, "mercado_libre").
		Order("s.created_at DESC").
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.IntegrationID == 0 || row.GuideID == "" {
		return nil, nil
	}
	shipmentID, perr := strconv.ParseInt(row.GuideID, 10, 64)
	if perr != nil {
		return nil, nil
	}
	return &domain.MeliOrderRef{
		IntegrationID: row.IntegrationID,
		ShipmentID:    shipmentID,
	}, nil
}

func (r *OrderLookupRepo) GetMeliLabelRefByShipmentID(ctx context.Context, shipmentID uint) (*domain.MeliLabelRef, error) {
	var row struct {
		BusinessID     uint
		IntegrationID  uint
		GuideID        string
		TrackingNumber string
	}
	err := r.db.Conn(ctx).
		Table("shipments AS s").
		Select("o.business_id, o.integration_id, COALESCE(s.guide_id,'') AS guide_id, COALESCE(s.tracking_number,'') AS tracking_number").
		Joins("JOIN orders o ON o.id = s.order_id AND o.deleted_at IS NULL").
		Where("s.id = ? AND s.deleted_at IS NULL", shipmentID).
		Where("(o.integration_type = ? OR o.platform = ?)", "mercado_libre", "mercadolibre").
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.BusinessID == 0 {
		return nil, nil
	}

	meliShipmentID := parseMeliShipmentID(row.GuideID)
	if meliShipmentID == 0 {
		meliShipmentID = parseMeliShipmentID(row.TrackingNumber)
	}

	return &domain.MeliLabelRef{
		BusinessID:     row.BusinessID,
		IntegrationID:  row.IntegrationID,
		MeliShipmentID: meliShipmentID,
	}, nil
}

func parseMeliShipmentID(raw string) int64 {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0
	}
	return id
}
