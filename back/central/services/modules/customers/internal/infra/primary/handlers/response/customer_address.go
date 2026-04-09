package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

type CustomerAddressResponse struct {
	ID         uint      `json:"id"`
	CustomerID uint      `json:"customer_id"`
	BusinessID uint      `json:"business_id"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	TimesUsed  int       `json:"times_used"`
	LastUsedAt time.Time `json:"last_used_at"`
}

func AddressFromEntity(a *entities.CustomerAddress) CustomerAddressResponse {
	return CustomerAddressResponse{
		ID:         a.ID,
		CustomerID: a.CustomerID,
		BusinessID: a.BusinessID,
		Street:     a.Street,
		City:       a.City,
		State:      a.State,
		Country:    a.Country,
		PostalCode: a.PostalCode,
		TimesUsed:  a.TimesUsed,
		LastUsedAt: a.LastUsedAt,
	}
}

type CustomerAddressListResponse struct {
	Data       []CustomerAddressResponse `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
	TotalPages int                       `json:"total_pages"`
}
