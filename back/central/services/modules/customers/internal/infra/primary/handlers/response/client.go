package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

// ClientResponse respuesta básica de cliente (para listado)
type ClientResponse struct {
	ID         uint      `json:"id"`
	BusinessID uint      `json:"business_id"`
	Name       string    `json:"name"`
	Email      *string   `json:"email"`
	Phone      string    `json:"phone"`
	Dni        *string   `json:"dni"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ClientDetailResponse respuesta con stats de órdenes
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

// FromEntity convierte una entidad de dominio a ClientResponse
func FromEntity(c *entities.Client) ClientResponse {
	return ClientResponse{
		ID:         c.ID,
		BusinessID: c.BusinessID,
		Name:       c.Name,
		Email:      c.Email,
		Phone:      c.Phone,
		Dni:        c.Dni,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

// DetailFromEntity convierte una entidad con stats a ClientDetailResponse
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

// ClientsListResponse respuesta paginada de clientes
type ClientsListResponse struct {
	Data       []ClientResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}
