package ports

import (
	"context"
	"io"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

type IDockerService interface {
	ListContainers(ctx context.Context) ([]entities.Container, error)
	GetContainer(ctx context.Context, id string) (*entities.Container, error)
	GetContainerStats(ctx context.Context, id string) (*entities.ContainerStats, error)
	GetContainerLogs(ctx context.Context, id string, tail int) ([]entities.LogEntry, error)
	StreamContainerLogs(ctx context.Context, id string) (io.ReadCloser, error)
	RestartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	StartContainer(ctx context.Context, id string) error
	GetComposeServices(ctx context.Context) ([]entities.ComposeService, error)
	GetSystemStats(ctx context.Context) (*entities.SystemStats, error)
}

type IUserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.MonitoringUser, string, error) // returns user, passwordHash, error
}

type IUseCase interface {
	Login(ctx context.Context, email, password string) (*entities.MonitoringUser, error)
	GenerateToken(user *entities.MonitoringUser) (string, error)
	ListContainers(ctx context.Context) ([]entities.Container, error)
	GetContainer(ctx context.Context, id string) (*entities.Container, error)
	GetContainerStats(ctx context.Context, id string) (*entities.ContainerStats, error)
	GetContainerLogs(ctx context.Context, id string, tail int) ([]entities.LogEntry, error)
	StreamContainerLogs(ctx context.Context, id string) (io.ReadCloser, error)
	ContainerAction(ctx context.Context, id string, action string) error
	GetComposeServices(ctx context.Context) ([]entities.ComposeService, error)
	GetSystemStats(ctx context.Context) (*entities.SystemStats, error)
}
