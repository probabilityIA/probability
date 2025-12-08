package mappers

import (
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/primary/handlers/requests"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/primary/handlers/responses"
)

// ===== REQUEST → DOMAIN =====

// ToDomainSSEConnection convierte request a domain DTO
func ToDomainSSEConnection(req requests.SSEConnectionRequest) domain.SSEConnectionRequest {
	return domain.SSEConnectionRequest{
		IntegrationID: req.IntegrationID,
		BusinessID:    req.BusinessID,
	}
}

// ===== DOMAIN → RESPONSE =====

// FromDomainSSEConnection convierte domain DTO a response
func FromDomainSSEConnection(connection domain.SSEConnection) responses.SSEConnectionResponse {
	return responses.SSEConnectionResponse{
		Status:        "connected",
		IntegrationID: 0, // Ya no se usa integration_id como identificador principal
		BusinessID:    fmt.Sprintf("%d", connection.BusinessID),
		ConnectedAt:   time.Now(), // Usar tiempo actual ya que no se guarda ConnectedAt
		ConnectionID:  connection.ConnectionID,
	}
}

// FromDomainSSEStatus convierte domain DTO a response
func FromDomainSSEStatus(status domain.SSEStatus) responses.SSEStatusResponse {
	return responses.SSEStatusResponse{
		IntegrationID:     status.IntegrationID,
		ActiveConnections: status.ActiveConnections,
		LastActivity:      status.LastActivity,
		Status:            status.Status,
	}
}
