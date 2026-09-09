package response

import (
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

type ClientResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Phone string  `json:"phone"`
	Dni   *string `json:"dni"`
}

type CreateClientResponse struct {
	Client       ClientResponse `json:"client"`
	TempPassword string         `json:"temp_password,omitempty"`
}

type ClientListResponse struct {
	Data       []ClientResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

func ClientFromEntity(c *entities.StorefrontClient) ClientResponse {
	return ClientResponse{
		ID:    c.ID,
		Name:  c.Name,
		Email: c.Email,
		Phone: c.Phone,
		Dni:   c.Dni,
	}
}

func ClientsFromEntities(clients []entities.StorefrontClient) []ClientResponse {
	result := make([]ClientResponse, len(clients))
	for i := range clients {
		result[i] = ClientFromEntity(&clients[i])
	}
	return result
}
