package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/domain"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/primary/handlers/response"
)

func toDomainCreate(req *request.CreateOrderStatusMappingRequest) *domain.OrderStatusMapping {
	return &domain.OrderStatusMapping{
		IntegrationType: req.IntegrationType,
		OriginalStatus:  req.OriginalStatus,
		MappedStatus:    req.MappedStatus,
		Priority:        req.Priority,
		Description:     req.Description,
		IsActive:        true,
	}
}

func toDomainUpdate(req *request.UpdateOrderStatusMappingRequest) *domain.OrderStatusMapping {
	return &domain.OrderStatusMapping{
		OriginalStatus: req.OriginalStatus,
		MappedStatus:   req.MappedStatus,
		Priority:       req.Priority,
		Description:    req.Description,
	}
}

func toResponse(m *domain.OrderStatusMapping) *response.OrderStatusMappingResponse {
	return &response.OrderStatusMappingResponse{
		ID:              m.ID,
		IntegrationType: m.IntegrationType,
		OriginalStatus:  m.OriginalStatus,
		MappedStatus:    m.MappedStatus,
		IsActive:        m.IsActive,
		Priority:        m.Priority,
		Description:     m.Description,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func toListResponse(mappings []domain.OrderStatusMapping, total int64) *response.OrderStatusMappingsListResponse {
	var data []response.OrderStatusMappingResponse
	for _, m := range mappings {
		data = append(data, *toResponse(&m))
	}
	return &response.OrderStatusMappingsListResponse{
		Data:  data,
		Total: total,
	}
}
