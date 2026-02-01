package mappers

import (
	"math"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

// DomainToResponse convierte una entidad de dominio a response DTO
func DomainToResponse(m *entities.OrderStatusMapping, imageURLBase string) *response.OrderStatusMappingResponse {
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

// DomainListToResponse convierte una lista de entidades a respuesta paginada
func DomainListToResponse(mappings []entities.OrderStatusMapping, total int64, page, pageSize int, imageURLBase string) *response.OrderStatusMappingsListResponse {
	var data []response.OrderStatusMappingResponse
	for _, m := range mappings {
		data = append(data, *DomainToResponse(&m, imageURLBase))
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

// StatusInfoToSimpleResponse convierte OrderStatusInfo a formato simple
func StatusInfoToSimpleResponse(status entities.OrderStatusInfo) response.OrderStatusSimpleResponse {
	return response.OrderStatusSimpleResponse{
		ID:       status.ID,
		Name:     status.Name,
		Code:     status.Code,
		IsActive: true, // Por defecto activo si viene de la lista
	}
}

// StatusInfoListToSimpleResponse convierte lista de OrderStatusInfo a formato simple
func StatusInfoListToSimpleResponse(statuses []entities.OrderStatusInfo) *response.OrderStatusesSimpleResponse {
	data := make([]response.OrderStatusSimpleResponse, len(statuses))
	for i, status := range statuses {
		data[i] = StatusInfoToSimpleResponse(status)
	}

	return &response.OrderStatusesSimpleResponse{
		Success: true,
		Message: "Order statuses retrieved successfully",
		Data:    data,
	}
}
