package request

import "mime/multipart"

// BusinessRoleAssignmentRequest representa una asignación de rol a un negocio en el request
type BusinessRoleAssignmentRequest struct {
	BusinessID uint `json:"business_id" form:"business_id" binding:"required"`
	RoleID     uint `json:"role_id" form:"role_id" binding:"required"`
}

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Name       string                `form:"name" json:"name" binding:"required,min=2,max=100"`
	Email      string                `form:"email" json:"email" binding:"required,email"`
	Phone      string                `form:"phone" json:"phone" binding:"omitempty,len=10"`
	AvatarURL  string                `form:"avatar_url" json:"avatar_url" binding:"omitempty,url"`
	AvatarFile *multipart.FileHeader `form:"avatarFile"`
	IsActive   bool                  `form:"is_active" json:"is_active"`
	// Relación de businesses
	BusinessIDsRaw string `form:"business_ids"`
	BusinessIDs    []uint `json:"business_ids"`
}
