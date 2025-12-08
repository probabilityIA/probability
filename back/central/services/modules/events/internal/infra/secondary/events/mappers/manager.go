package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/events/requests"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/events/responses"
)

// ===== REQUEST → DOMAIN =====

// ToDomainPublishEvent convierte request a domain DTO
func ToDomainPublishEvent(req requests.PublishEventRequest) domain.Event {
	return req.Event
}

// ToDomainAddConnection convierte request a domain DTO
func ToDomainAddConnection(req requests.AddConnectionRequest) domain.ConnectionRequest {
	return domain.ConnectionRequest{
		IntegrationID: req.IntegrationID,
		Connection:    req.Connection,
	}
}

// ToDomainRemoveConnection convierte request a domain DTO
func ToDomainRemoveConnection(req requests.RemoveConnectionRequest) domain.ConnectionRequest {
	return domain.ConnectionRequest{
		IntegrationID: req.IntegrationID,
	}
}

// ToDomainGetConnectionInfo convierte request a domain DTO
func ToDomainGetConnectionInfo(req requests.GetConnectionInfoRequest) domain.ConnectionRequest {
	return domain.ConnectionRequest{
		IntegrationID: req.IntegrationID,
	}
}

// ===== DOMAIN → RESPONSE =====

// FromDomainPublishEvent convierte domain DTO a response
func FromDomainPublishEvent(event domain.Event, result domain.PublishEventResult) responses.PublishEventResponse {
	return responses.PublishEventResponse{
		Status:      "published",
		EventID:     event.ID,
		PublishedAt: event.Timestamp,
		Recipients:  result.Recipients,
		Errors:      result.Errors,
		Metadata:    event.Metadata,
	}
}

// FromDomainConnection convierte domain DTO a response
func FromDomainConnection(connection domain.ConnectionResult) responses.ConnectionResponse {
	status := "success"
	if len(connection.Errors) > 0 {
		status = "partial_success"
	}

	return responses.ConnectionResponse{
		Status:        status,
		IntegrationID: connection.IntegrationID,
		Operation:     connection.Operation,
		Timestamp:     time.Now(),
		Message:       connection.Message,
		Errors:        connection.Errors,
	}
}

// FromDomainConnectionInfo convierte domain DTO a response
func FromDomainConnectionInfo(info domain.ConnectionInfo) responses.ConnectionInfoResponse {
	if info.IntegrationID == 0 {
		return responses.ConnectionInfoResponse{
			Status:            "no_connections",
			ActiveConnections: 0,
			TotalConnections:  0,
			ConnectedAt:       time.Time{},
			LastActivity:      time.Time{},
		}
	}

	return responses.ConnectionInfoResponse{
		IntegrationID:     info.IntegrationID,
		ActiveConnections: info.ActiveConnections,
		TotalConnections:  info.TotalConnections,
		ConnectedAt:       info.ConnectedAt,
		LastActivity:      time.Now(),
		Status:            "active",
	}
}
