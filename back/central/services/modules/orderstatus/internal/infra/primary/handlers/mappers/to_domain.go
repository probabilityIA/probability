package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/request"
)

// CreateRequestToDomain convierte un CreateRequest a entidad de dominio
func CreateRequestToDomain(req *request.CreateOrderStatusMappingRequest) *entities.OrderStatusMapping {
	return &entities.OrderStatusMapping{
		IntegrationTypeID: req.IntegrationTypeID,
		OriginalStatus:    req.OriginalStatus,
		OrderStatusID:     req.OrderStatusID,
		Description:       req.Description,
		IsActive:          true,
	}
}

// UpdateRequestToDomain convierte un UpdateRequest a entidad de dominio
func UpdateRequestToDomain(req *request.UpdateOrderStatusMappingRequest) *entities.OrderStatusMapping {
	return &entities.OrderStatusMapping{
		OriginalStatus: req.OriginalStatus,
		OrderStatusID:  req.OrderStatusID,
		Description:    req.Description,
	}
}
