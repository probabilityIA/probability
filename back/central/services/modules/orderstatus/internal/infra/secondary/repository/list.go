package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) List(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error) {
	var modelsList []models.OrderStatusMapping
	var total int64

	query := r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Preload("IntegrationType").
		Preload("OrderStatus")

	// Aplicar filtros
	if integrationTypeID, ok := filters["integration_type_id"].(uint); ok && integrationTypeID > 0 {
		query = query.Where("integration_type_id = ?", integrationTypeID)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Contar total antes de aplicar paginación
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	if page, ok := filters["page"].(int); ok && page > 0 {
		if pageSize, ok := filters["page_size"].(int); ok && pageSize > 0 {
			offset := (page - 1) * pageSize
			query = query.Offset(offset).Limit(pageSize)
		}
	}

	// Obtener resultados
	if err := query.Order("integration_type_id ASC, priority DESC, created_at DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	// Convertir a dominio
	domainList := make([]entities.OrderStatusMapping, len(modelsList))
	for i, m := range modelsList {
		domainList[i] = m.ToDomain()
	}

	return domainList, total, nil
}
