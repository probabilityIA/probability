package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/entities"
)

// IAlertPublisher define el contrato para publicar alertas en la cola
type IAlertPublisher interface {
	Publish(ctx context.Context, event entities.AlertEvent) error
}

// IUseCase define el contrato del caso de uso de monitoreo
type IUseCase interface {
	ProcessGrafanaAlert(ctx context.Context, dto dtos.GrafanaWebhookDTO) error
}
