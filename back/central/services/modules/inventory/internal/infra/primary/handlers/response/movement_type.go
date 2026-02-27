package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

// MovementTypeResponse respuesta de tipo de movimiento
type MovementTypeResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	Direction   string    `json:"direction"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MovementTypeListResponse respuesta paginada de tipos de movimiento
type MovementTypeListResponse struct {
	Data       []MovementTypeResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

// MovementTypeFromEntity convierte entidad a response
func MovementTypeFromEntity(e *entities.StockMovementType) MovementTypeResponse {
	return MovementTypeResponse{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		IsActive:    e.IsActive,
		Direction:   e.Direction,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
