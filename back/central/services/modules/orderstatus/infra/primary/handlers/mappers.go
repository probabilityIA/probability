package handlers

import (
	"math"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/domain"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/primary/handlers/response"
)

func toDomainCreate(req *request.CreateOrderStatusMappingRequest) *domain.OrderStatusMapping {
	return &domain.OrderStatusMapping{
		IntegrationTypeID: req.IntegrationTypeID,
		OriginalStatus:    req.OriginalStatus,
		OrderStatusID:     req.OrderStatusID,
		Priority:          req.Priority,
		Description:       req.Description,
		IsActive:          true,
	}
}

func toDomainUpdate(req *request.UpdateOrderStatusMappingRequest) *domain.OrderStatusMapping {
	return &domain.OrderStatusMapping{
		OriginalStatus: req.OriginalStatus,
		OrderStatusID:  req.OrderStatusID,
		Priority:       req.Priority,
		Description:    req.Description,
	}
}

func toResponse(m *domain.OrderStatusMapping, imageURLBase string) *response.OrderStatusMappingResponse {
	resp := &response.OrderStatusMappingResponse{
		ID:                m.ID,
		IntegrationTypeID: m.IntegrationTypeID,
		OriginalStatus:    m.OriginalStatus,
		OrderStatusID:     m.OrderStatusID,
		IsActive:          m.IsActive,
		Priority:          m.Priority,
		Description:       m.Description,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	// Incluir información del IntegrationType si está disponible
	if m.IntegrationType != nil {
		imageURL := m.IntegrationType.ImageURL
		// Construir URL completa si es path relativo
		if imageURL != "" && imageURLBase != "" && !strings.HasPrefix(imageURL, "http") {
			imageURL = strings.TrimRight(imageURLBase, "/") + "/" + strings.TrimLeft(imageURL, "/")
		}
		resp.IntegrationType = &response.IntegrationTypeInfo{
			ID:       m.IntegrationType.ID,
			Code:     m.IntegrationType.Code,
			Name:     m.IntegrationType.Name,
			ImageURL: imageURL,
		}
	}

	// Incluir información del OrderStatus si está disponible
	if m.OrderStatus != nil {
		resp.OrderStatus = &response.OrderStatusInfo{
			ID:          m.OrderStatus.ID,
			Code:        m.OrderStatus.Code,
			Name:        m.OrderStatus.Name,
			Description: m.OrderStatus.Description,
			Category:    m.OrderStatus.Category,
			Color:       m.OrderStatus.Color,
		}
	}

	return resp
}

func toListResponse(mappings []domain.OrderStatusMapping, total int64, page, pageSize int, imageURLBase string) *response.OrderStatusMappingsListResponse {
	var data []response.OrderStatusMappingResponse
	for _, m := range mappings {
		data = append(data, *toResponse(&m, imageURLBase))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	return &response.OrderStatusMappingsListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
