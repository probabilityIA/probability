package domain

import "context"

// IRepository define la interfaz para el almacenamiento de configuraciones
type IRepository interface {
	Create(ctx context.Context, config *NotificationConfig) error
	GetByID(ctx context.Context, id uint) (*NotificationConfig, error)
	Update(ctx context.Context, id uint, config *NotificationConfig) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter ConfigFilter) ([]*NotificationConfig, error)
	GetByBusinessAndEventType(ctx context.Context, businessID uint, eventType string) (*NotificationConfig, error)
}

// IUseCase define la interfaz para la lógica de negocio
type IUseCase interface {
	CreateConfig(ctx context.Context, dto CreateConfigDTO) (*NotificationConfig, error)
	GetConfig(ctx context.Context, id uint) (*NotificationConfig, error)
	UpdateConfig(ctx context.Context, id uint, dto UpdateConfigDTO) (*NotificationConfig, error)
	DeleteConfig(ctx context.Context, id uint) error
	ListConfigs(ctx context.Context, filter ConfigFilter) ([]*NotificationConfig, error)
}
