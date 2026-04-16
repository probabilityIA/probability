package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

type ClientResponse struct {
	ID          uint      `json:"id"`
	BusinessID  uint      `json:"business_id"`
	Name        string    `json:"name"`
	Email       *string   `json:"email"`
	Phone       string    `json:"phone"`
	Dni         *string   `json:"dni"`
	TotalOrders int64     `json:"total_orders"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClientDetailResponse struct {
	ID          uint       `json:"id"`
	BusinessID  uint       `json:"business_id"`
	Name        string     `json:"name"`
	Email       *string    `json:"email"`
	Phone       string     `json:"phone"`
	Dni         *string    `json:"dni"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	OrderCount  int64      `json:"order_count"`
	TotalSpent  float64    `json:"total_spent"`
	LastOrderAt *time.Time `json:"last_order_at"`
}

func FromEntity(c *entities.Client) ClientResponse {
	return ClientResponse{
		ID:          c.ID,
		BusinessID:  c.BusinessID,
		Name:        c.Name,
		Email:       c.Email,
		Phone:       c.Phone,
		Dni:         c.Dni,
		TotalOrders: c.OrderCount,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func DetailFromEntity(c *entities.Client) ClientDetailResponse {
	return ClientDetailResponse{
		ID:          c.ID,
		BusinessID:  c.BusinessID,
		Name:        c.Name,
		Email:       c.Email,
		Phone:       c.Phone,
		Dni:         c.Dni,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		OrderCount:  c.OrderCount,
		TotalSpent:  c.TotalSpent,
		LastOrderAt: c.LastOrderAt,
	}
}

type ClientsListResponse struct {
	Data       []ClientResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}
