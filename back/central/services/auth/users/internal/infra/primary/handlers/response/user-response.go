package response

import "time"

// RoleInfo representa información simplificada de un rol en la respuesta de usuario
type RoleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	IsSystem    bool   `json:"is_system"`
	ScopeID     uint   `json:"scope_id"`
}

// RoleInfoDetailed representa información completa de un rol con información de tipo de business
type RoleInfoDetailed struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Level            int    `json:"level"`
	IsSystem         bool   `json:"is_system"`
	ScopeID          uint   `json:"scope_id"`
	ScopeName        string `json:"scope_name,omitempty"`
	ScopeCode        string `json:"scope_code,omitempty"`
	BusinessTypeID   uint   `json:"business_type_id,omitempty"`
	BusinessTypeName string `json:"business_type_name,omitempty"`
}

// BusinessInfo representa información simplificada de un business en la respuesta de usuario
type BusinessInfo struct {
	ID               uint              `json:"id"`
	Name             string            `json:"name"`
	LogoURL          string            `json:"logo_url"`
	BusinessTypeID   uint              `json:"business_type_id"`
	BusinessTypeName string            `json:"business_type_name"`
	Role             *RoleInfoDetailed `json:"role,omitempty"` // Rol del usuario en este business
}

// BusinessRoleAssignmentResponse representa una asignación business-rol en la respuesta
type BusinessRoleAssignmentResponse struct {
	BusinessID   uint   `json:"business_id"`
	BusinessName string `json:"business_name,omitempty"`
	RoleID       uint   `json:"role_id"`
	RoleName     string `json:"role_name,omitempty"`
}

// UserResponse representa la respuesta de un usuario
type UserResponse struct {
	ID                      uint                             `json:"id"`
	Name                    string                           `json:"name"`
	Email                   string                           `json:"email"`
	Phone                   string                           `json:"phone"`
	AvatarURL               string                           `json:"avatar_url"`
	IsActive                bool                             `json:"is_active"`
	IsSuperUser             bool                             `json:"is_super_user"`
	LastLoginAt             *time.Time                       `json:"last_login_at"`
	BusinessRoleAssignments []BusinessRoleAssignmentResponse `json:"business_role_assignments"`
	CreatedAt               time.Time                        `json:"created_at"`
	UpdatedAt               time.Time                        `json:"updated_at"`
}

// PaginationInfo representa información de paginación
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// UserListResponse representa la respuesta de una lista de usuarios con paginación
type UserListResponse struct {
	Success    bool           `json:"success"`
	Data       []UserResponse `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// UserSuccessResponse representa la respuesta exitosa de un usuario individual
type UserSuccessResponse struct {
	Success bool         `json:"success"`
	Data    UserResponse `json:"data"`
}

// UserMessageResponse representa la respuesta de mensaje para usuarios
type UserMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UserErrorResponse representa la respuesta de error para usuarios
type UserErrorResponse struct {
	Error string `json:"error"`
}

// UserCreatedResponse representa la respuesta al crear un usuario, incluyendo la contraseña generada
// ¡La contraseña solo se muestra una vez en la creación!
type UserCreatedResponse struct {
	Success  bool   `json:"success"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Message  string `json:"message"`
}
