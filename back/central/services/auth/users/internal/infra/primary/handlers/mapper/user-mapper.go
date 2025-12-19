package mapper

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"
)

// ToUserFilters convierte GetUsersRequest a UserFilters del dominio
func ToUserFilters(req request.GetUsersRequest) domain.UserFilters {
	// Parsear user_ids desde string separado por comas
	var userIDs []uint
	if req.UserIDs != "" {
		userIDsStr := strings.Split(req.UserIDs, ",")
		for _, idStr := range userIDsStr {
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
					userIDs = append(userIDs, uint(id))
				}
			}
		}
	}

	return domain.UserFilters{
		Page:           req.Page,
		PageSize:       req.PageSize,
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		UserIDs:        userIDs,
		IsActive:       req.IsActive,
		RoleID:         req.RoleID,
		BusinessID:     req.BusinessID,
		CreatedAt:      req.CreatedAt,
		SortBy:         req.SortBy,
		SortOrder:      req.SortOrder,
		IncludeDeleted: req.IncludeDeleted,
	}
}

// ToCreateUserDTO convierte CreateUserRequest a CreateUserDTO del dominio
func ToCreateUserDTO(req request.CreateUserRequest) domain.CreateUserDTO {
	businessIDs := make([]uint, 0)
	if len(req.BusinessIDs) > 0 {
		businessIDs = append(businessIDs, req.BusinessIDs...)
	} else if req.BusinessIDsRaw != "" {
		var arr []uint
		if err := json.Unmarshal([]byte(req.BusinessIDsRaw), &arr); err == nil {
			businessIDs = append(businessIDs, arr...)
		} else {
			parts := strings.Split(req.BusinessIDsRaw, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				if id, err := strconv.ParseUint(p, 10, 32); err == nil {
					businessIDs = append(businessIDs, uint(id))
				}
			}
		}
	}

	return domain.CreateUserDTO{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		AvatarURL:   req.AvatarURL,
		AvatarFile:  req.AvatarFile,
		IsActive:    req.IsActive,
		ScopeID:     req.ScopeID,
		BusinessIDs: businessIDs,
	}
}

// ToUpdateUserDTO convierte UpdateUserRequest a UpdateUserDTO del dominio
func ToUpdateUserDTO(req request.UpdateUserRequest) domain.UpdateUserDTO {
	businessIDs := make([]uint, 0)
	if len(req.BusinessIDs) > 0 {
		businessIDs = append(businessIDs, req.BusinessIDs...)
	} else if req.BusinessIDsRaw != "" {
		var arr []uint
		if err := json.Unmarshal([]byte(req.BusinessIDsRaw), &arr); err == nil {
			businessIDs = append(businessIDs, arr...)
		} else {
			parts := strings.Split(req.BusinessIDsRaw, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				if id, err := strconv.ParseUint(p, 10, 32); err == nil {
					businessIDs = append(businessIDs, uint(id))
				}
			}
		}
	}

	return domain.UpdateUserDTO{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		AvatarURL:    req.AvatarURL,
		AvatarFile:   req.AvatarFile,
		RemoveAvatar: req.RemoveAvatar,
		IsActive:     req.IsActive,
		BusinessIDs:  businessIDs,
	}
}

// ToUserResponse convierte UserDTO a UserResponse
func ToUserResponse(dto domain.UserDTO) response.UserResponse {
	// Convertir BusinessRoleAssignments
	businessRoleAssignments := make([]response.BusinessRoleAssignmentResponse, len(dto.BusinessRoleAssignments))
	for i, assignment := range dto.BusinessRoleAssignments {
		businessRoleAssignments[i] = response.BusinessRoleAssignmentResponse{
			BusinessID:   assignment.BusinessID,
			BusinessName: assignment.BusinessName,
			RoleID:       assignment.RoleID,
			RoleName:     assignment.RoleName,
		}
	}

	return response.UserResponse{
		ID:                      dto.ID,
		Name:                    dto.Name,
		Email:                   dto.Email,
		Phone:                   dto.Phone,
		AvatarURL:               dto.AvatarURL,
		IsActive:                dto.IsActive,
		IsSuperUser:             dto.IsSuperUser,
		ScopeID:                 dto.ScopeID,
		ScopeCode:               dto.ScopeCode,
		ScopeName:               dto.ScopeName,
		LastLoginAt:             dto.LastLoginAt,
		BusinessRoleAssignments: businessRoleAssignments,
		CreatedAt:               dto.CreatedAt,
		UpdatedAt:               dto.UpdatedAt,
	}
}

// ToUserListResponse convierte un UserListDTO a UserListResponse con paginación
func ToUserListResponse(userListDTO *domain.UserListDTO) response.UserListResponse {
	userResponses := make([]response.UserResponse, len(userListDTO.Users))
	for i, user := range userListDTO.Users {
		userResponses[i] = ToUserResponse(user)
	}

	// Calcular has_next y has_prev
	hasNext := userListDTO.Page < userListDTO.TotalPages
	hasPrev := userListDTO.Page > 1

	return response.UserListResponse{
		Success: true,
		Data:    userResponses,
		Pagination: response.PaginationInfo{
			CurrentPage: userListDTO.Page,
			PerPage:     userListDTO.PageSize,
			Total:       userListDTO.Total,
			LastPage:    userListDTO.TotalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}
}
