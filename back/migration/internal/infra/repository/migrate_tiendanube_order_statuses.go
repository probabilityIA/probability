package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	tiendanubeTypeID    = 17
	probStatusCompleted = 5
	probStatusRefunded  = 7
)

var tiendanubeChannelStatuses = []models.IntegrationChannelStatus{
	{Code: "open", Name: "Abierta", Description: "La orden esta abierta en Tiendanube", DisplayOrder: 1, IsActive: true},
	{Code: "closed", Name: "Cerrada", Description: "La orden fue cerrada en Tiendanube", DisplayOrder: 2, IsActive: true},
	{Code: "cancelled", Name: "Cancelada", Description: "La orden fue cancelada en Tiendanube", DisplayOrder: 3, IsActive: true},
	{Code: "pending", Name: "Pago pendiente", Description: "Payment status: la orden aun no se pago", DisplayOrder: 4, IsActive: true},
	{Code: "paid", Name: "Pagada", Description: "Payment status: la orden fue pagada", DisplayOrder: 5, IsActive: true},
	{Code: "refunded", Name: "Reembolsada", Description: "Payment status: la orden fue reembolsada", DisplayOrder: 6, IsActive: true},
	{Code: "voided", Name: "Anulada", Description: "Payment status: el pago fue anulado", DisplayOrder: 7, IsActive: true},
	{Code: "abandoned", Name: "Abandonada", Description: "Payment status: el comprador no completo el checkout", DisplayOrder: 8, IsActive: true},
	{Code: "unpacked", Name: "Sin empacar", Description: "Shipping status: la orden no se ha empacado", DisplayOrder: 9, IsActive: true},
	{Code: "shipped", Name: "Enviada", Description: "Shipping status: la orden fue despachada", DisplayOrder: 10, IsActive: true},
	{Code: "delivered", Name: "Entregada", Description: "Shipping status: la orden fue entregada", DisplayOrder: 11, IsActive: true},
}

var tiendanubeStatusMappings = []models.OrderStatusMapping{
	{OriginalStatus: "open", OrderStatusID: probStatusProcessing, Description: "Tiendanube: orden abierta, entra a preparacion", IsActive: true},
	{OriginalStatus: "closed", OrderStatusID: probStatusCompleted, Description: "Tiendanube: orden cerrada", IsActive: true},
	{OriginalStatus: "cancelled", OrderStatusID: probStatusCancelled, Description: "Tiendanube: orden cancelada", IsActive: true},

	{OriginalStatus: "pending", OrderStatusID: probStatusPending, Description: "Tiendanube: pago pendiente", IsActive: true},
	{OriginalStatus: "paid", OrderStatusID: probStatusProcessing, Description: "Tiendanube: pagada, entra a preparacion", IsActive: true},
	{OriginalStatus: "completed", OrderStatusID: probStatusCompleted, Description: "Tiendanube: pagada y entregada", IsActive: true},
	{OriginalStatus: "refunded", OrderStatusID: probStatusRefunded, Description: "Tiendanube: reembolsada", IsActive: true},
	{OriginalStatus: "voided", OrderStatusID: probStatusCancelled, Description: "Tiendanube: pago anulado", IsActive: true},
	{OriginalStatus: "abandoned", OrderStatusID: probStatusCancelled, Description: "Tiendanube: checkout abandonado", IsActive: true},

	{OriginalStatus: "shipped", OrderStatusID: probStatusInTransit, Description: "Tiendanube shipping: despachada", IsActive: true},
	{OriginalStatus: "in_transit", OrderStatusID: probStatusInTransit, Description: "Tiendanube shipping: en camino", IsActive: true},
	{OriginalStatus: "delivered", OrderStatusID: probStatusDelivered, Description: "Tiendanube shipping: entregada", IsActive: true},
}

func (r *Repository) migrateTiendanubeOrderStatuses(ctx context.Context) error {
	var tipo models.IntegrationType
	if err := r.db.Conn(ctx).Where("id = ?", tiendanubeTypeID).First(&tipo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("migrateTiendanubeOrderStatuses: consultando tipo: %w", err)
	}

	channelStatuses := make([]models.IntegrationChannelStatus, 0, len(tiendanubeChannelStatuses))
	for _, s := range tiendanubeChannelStatuses {
		s.IntegrationTypeID = tiendanubeTypeID
		channelStatuses = append(channelStatuses, s)
	}

	if err := r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "integration_type_id"}, {Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "description", "display_order", "is_active", "updated_at"}),
		}).
		Create(&channelStatuses).Error; err != nil {
		return fmt.Errorf("migrateTiendanubeOrderStatuses: sembrando channel statuses: %w", err)
	}

	mappings := make([]models.OrderStatusMapping, 0, len(tiendanubeStatusMappings))
	for _, m := range tiendanubeStatusMappings {
		m.IntegrationTypeID = tiendanubeTypeID
		mappings = append(mappings, m)
	}

	if err := r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "integration_type_id"}, {Name: "original_status"}},
			DoUpdates: clause.AssignmentColumns([]string{"order_status_id", "description", "is_active", "updated_at"}),
		}).
		Create(&mappings).Error; err != nil {
		return fmt.Errorf("migrateTiendanubeOrderStatuses: sembrando mapeos de estado: %w", err)
	}

	return nil
}
