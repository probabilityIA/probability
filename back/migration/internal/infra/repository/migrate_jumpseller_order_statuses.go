package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	probStatusPending    = 1
	probStatusProcessing = 2
	probStatusDelivered  = 4
	probStatusCancelled  = 6
	probStatusInTransit  = 17
	probStatusFailed     = 19
)

var jumpsellerChannelStatuses = []models.IntegrationChannelStatus{
	{Code: "Pending Payment", Name: "Pago pendiente", Description: "La orden se creo pero aun no se ha pagado", DisplayOrder: 1, IsActive: true},
	{Code: "Paid", Name: "Pagada", Description: "La orden fue pagada y esta lista para preparacion", DisplayOrder: 2, IsActive: true},
	{Code: "Canceled", Name: "Cancelada", Description: "La orden fue cancelada en Jumpseller", DisplayOrder: 3, IsActive: true},
	{Code: "Abandoned", Name: "Abandonada", Description: "El comprador no completo el checkout", DisplayOrder: 4, IsActive: true},
	{Code: "requested", Name: "Envio solicitado", Description: "Shipment status: el envio fue solicitado", DisplayOrder: 5, IsActive: true},
	{Code: "in_transit", Name: "En transito", Description: "Shipment status: el envio va en camino", DisplayOrder: 6, IsActive: true},
	{Code: "delivered", Name: "Entregado", Description: "Shipment status: el envio fue entregado", DisplayOrder: 7, IsActive: true},
	{Code: "failed", Name: "Envio fallido", Description: "Shipment status: el envio no pudo entregarse", DisplayOrder: 8, IsActive: true},
}

var jumpsellerStatusMappings = []models.OrderStatusMapping{
	{OriginalStatus: "Pending Payment", OrderStatusID: probStatusPending, Description: "Jumpseller: orden sin pagar", IsActive: true},
	{OriginalStatus: "Paid", OrderStatusID: probStatusProcessing, Description: "Jumpseller: orden pagada, entra a preparacion", IsActive: true},
	{OriginalStatus: "Canceled", OrderStatusID: probStatusCancelled, Description: "Jumpseller: orden cancelada", IsActive: true},
	{OriginalStatus: "Abandoned", OrderStatusID: probStatusCancelled, Description: "Jumpseller: checkout abandonado", IsActive: true},

	{OriginalStatus: "pending", OrderStatusID: probStatusPending, Description: "Jumpseller canonico: pendiente", IsActive: true},
	{OriginalStatus: "paid", OrderStatusID: probStatusProcessing, Description: "Jumpseller canonico: pagada", IsActive: true},
	{OriginalStatus: "cancelled", OrderStatusID: probStatusCancelled, Description: "Jumpseller canonico: cancelada", IsActive: true},
	{OriginalStatus: "abandoned", OrderStatusID: probStatusCancelled, Description: "Jumpseller canonico: abandonada", IsActive: true},

	{OriginalStatus: "in_transit", OrderStatusID: probStatusInTransit, Description: "Jumpseller shipment: en camino", IsActive: true},
	{OriginalStatus: "delivered", OrderStatusID: probStatusDelivered, Description: "Jumpseller shipment: entregado", IsActive: true},
	{OriginalStatus: "failed", OrderStatusID: probStatusFailed, Description: "Jumpseller shipment: entrega fallida", IsActive: true},
}

func (r *Repository) migrateJumpsellerOrderStatuses(ctx context.Context) error {
	var tipo models.IntegrationType
	if err := r.db.Conn(ctx).Where("id = ?", jumpsellerTypeID).First(&tipo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("migrateJumpsellerOrderStatuses: consultando tipo: %w", err)
	}

	channelStatuses := make([]models.IntegrationChannelStatus, 0, len(jumpsellerChannelStatuses))
	for _, s := range jumpsellerChannelStatuses {
		s.IntegrationTypeID = jumpsellerTypeID
		channelStatuses = append(channelStatuses, s)
	}

	if err := r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "integration_type_id"}, {Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "description", "display_order", "is_active", "updated_at"}),
		}).
		Create(&channelStatuses).Error; err != nil {
		return fmt.Errorf("migrateJumpsellerOrderStatuses: sembrando channel statuses: %w", err)
	}

	mappings := make([]models.OrderStatusMapping, 0, len(jumpsellerStatusMappings))
	for _, m := range jumpsellerStatusMappings {
		m.IntegrationTypeID = jumpsellerTypeID
		mappings = append(mappings, m)
	}

	if err := r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "integration_type_id"}, {Name: "original_status"}},
			DoUpdates: clause.AssignmentColumns([]string{"order_status_id", "description", "is_active", "updated_at"}),
		}).
		Create(&mappings).Error; err != nil {
		return fmt.Errorf("migrateJumpsellerOrderStatuses: sembrando mapeos de estado: %w", err)
	}

	return nil
}
