package requests

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
)

// PublishEventRequest representa la solicitud para publicar un evento
type PublishEventRequest struct {
	Event domain.Event `json:"event" validate:"required"`
}

// AddConnectionRequest representa la solicitud para agregar una conexión
type AddConnectionRequest struct {
	IntegrationID int64               `json:"integration_id" validate:"required"`
	Connection    http.ResponseWriter `json:"-"` // No se serializa
}

// RemoveConnectionRequest representa la solicitud para remover una conexión
type RemoveConnectionRequest struct {
	IntegrationID int64 `json:"integration_id" validate:"required"`
}

// GetConnectionInfoRequest representa la solicitud para obtener información de conexiones
type GetConnectionInfoRequest struct {
	IntegrationID int64 `json:"integration_id" validate:"required"`
}
