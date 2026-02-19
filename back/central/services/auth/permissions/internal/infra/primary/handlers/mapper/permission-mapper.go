package mapper

import (
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/response"
)

// ToCreatePermissionDTO convierte CreatePermissionRequest a CreatePermissionDTO
func ToCreatePermissionDTO(req request.CreatePermissionRequest) domain.CreatePermissionDTO {
	return domain.CreatePermissionDTO{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		ResourceID:     req.ResourceID,
		ActionID:       req.ActionID,
		ScopeID:        req.ScopeID,
		BusinessTypeID: req.BusinessTypeID,
	}
}

// ToUpdatePermissionDTO convierte UpdatePermissionRequest a UpdatePermissionDTO
func ToUpdatePermissionDTO(req request.UpdatePermissionRequest) domain.UpdatePermissionDTO {
	return domain.UpdatePermissionDTO{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		ResourceID:     req.ResourceID,
		ActionID:       req.ActionID,
		ScopeID:        req.ScopeID,
		BusinessTypeID: req.BusinessTypeID,
	}
}

// ToBulkItemToCreateDTO convierte BulkCreatePermissionItem a CreatePermissionDTO
func ToBulkItemToCreateDTO(item request.BulkCreatePermissionItem) domain.CreatePermissionDTO {
	return domain.CreatePermissionDTO{
		Name:           item.Name,
		ResourceID:     item.ResourceID,
		ActionID:       item.ActionID,
		ScopeID:        item.ScopeID,
		BusinessTypeID: item.BusinessTypeID,
	}
}

// ToPermissionResponse convierte PermissionDTO a PermissionResponse
func ToPermissionResponse(dto domain.PermissionDTO) response.PermissionResponse {
	return response.PermissionResponse{
		ID:               dto.ID,
		Name:             dto.Name,
		Code:             dto.Code,
		Description:      dto.Description,
		Resource:         dto.Resource,
		Action:           dto.Action,
		ResourceID:       dto.ResourceID,
		ActionID:         dto.ActionID,
		ScopeID:          dto.ScopeID,
		ScopeName:        dto.ScopeName,
		ScopeCode:        dto.ScopeCode,
		BusinessTypeID:   dto.BusinessTypeID,
		BusinessTypeName: dto.BusinessTypeName,
	}
}

// ToPermissionListResponse convierte []PermissionDTO a PermissionListResponse
func ToPermissionListResponse(dtos []domain.PermissionDTO) response.PermissionListResponse {
	permissions := make([]response.PermissionResponse, len(dtos))
	for i, dto := range dtos {
		permissions[i] = ToPermissionResponse(dto)
	}

	return response.PermissionListResponse{
		Success: true,
		Data:    permissions,
		Total:   len(permissions),
	}
}
