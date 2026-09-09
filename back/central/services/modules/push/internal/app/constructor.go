package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	RegisterDevice(ctx context.Context, dto dtos.RegisterDeviceDTO) error
	UnregisterDevice(ctx context.Context, userID uint, token string) error
	ListDevices(ctx context.Context, userID uint) ([]entities.DeviceToken, error)
	NotifyBusiness(ctx context.Context, event dtos.PushEventDTO) error
}

type UseCase struct {
	repo   ports.IRepository
	sender ports.IPushSender
	log    log.ILogger
}

func New(repo ports.IRepository, sender ports.IPushSender, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, sender: sender, log: logger}
}
