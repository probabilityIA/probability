package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/entities"
)

type ReferralResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	OrderRange string    `json:"order_range"`
	CreatedAt  time.Time `json:"created_at"`
}

func FromEntity(r *entities.SiigoReferral) ReferralResponse {
	return ReferralResponse{
		ID:         r.ID,
		Name:       r.Name,
		Email:      r.Email,
		Phone:      r.Phone,
		OrderRange: r.OrderRange,
		CreatedAt:  r.CreatedAt,
	}
}

func FromEntities(referrals []entities.SiigoReferral) []ReferralResponse {
	out := make([]ReferralResponse, 0, len(referrals))
	for i := range referrals {
		out = append(out, FromEntity(&referrals[i]))
	}
	return out
}
