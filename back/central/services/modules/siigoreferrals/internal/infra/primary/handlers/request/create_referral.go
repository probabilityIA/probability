package request

import (
	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/dtos"
)

type CreateReferralRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	OrderRange string `json:"order_range"`
}

func (r *CreateReferralRequest) ToDTO() dtos.CreateReferralDTO {
	return dtos.CreateReferralDTO{
		Name:       r.Name,
		Email:      r.Email,
		Phone:      r.Phone,
		OrderRange: r.OrderRange,
	}
}
