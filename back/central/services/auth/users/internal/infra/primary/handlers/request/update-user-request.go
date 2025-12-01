package request

import "mime/multipart"

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Name         string                `form:"name" json:"name" binding:"omitempty,min=2,max=100"`
	Email        string                `form:"email" json:"email" binding:"omitempty,email"`
	Phone        string                `form:"phone" json:"phone" binding:"omitempty,len=10"`
	AvatarURL    string                `form:"avatar_url" json:"avatar_url" binding:"omitempty,url"`
	RemoveAvatar bool                  `form:"remove_avatar" json:"remove_avatar"`
	AvatarFile   *multipart.FileHeader `form:"avatarFile"`
	IsActive     bool                  `form:"is_active" json:"is_active"`
	// Relación de businesses (sustituye relaciones)
	BusinessIDsRaw string `form:"business_ids"`
	BusinessIDs    []uint `json:"business_ids"`
}
