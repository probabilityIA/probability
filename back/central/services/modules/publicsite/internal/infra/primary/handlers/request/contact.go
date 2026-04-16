package request

import "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"

type ContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message" binding:"required"`
}

func (r *ContactRequest) ToDTO() *dtos.ContactFormDTO {
	return &dtos.ContactFormDTO{
		Name:    r.Name,
		Email:   r.Email,
		Phone:   r.Phone,
		Message: r.Message,
	}
}
