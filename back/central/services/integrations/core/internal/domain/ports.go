package domain

import (
	"context"
	"time"
)

// IRepository define la interfaz unificada del repositorio de integraciones y tipos de integración
type IRepository interface {
	// Métodos de Integrations
	CreateIntegration(ctx context.Context, integration *Integration) error
	UpdateIntegration(ctx context.Context, id uint, integration *Integration) error
	GetIntegrationByID(ctx context.Context, id uint) (*Integration, error)
	DeleteIntegration(ctx context.Context, id uint) error
	ListIntegrations(ctx context.Context, filters IntegrationFilters) ([]*Integration, int64, error)
	GetIntegrationByIntegrationTypeID(ctx context.Context, integrationTypeID uint, businessID *uint) (*Integration, error)
	GetActiveIntegrationByIntegrationTypeID(ctx context.Context, integrationTypeID uint, businessID *uint) (*Integration, error)
	ListIntegrationsByBusiness(ctx context.Context, businessID uint) ([]*Integration, error)
	ListIntegrationsByIntegrationTypeID(ctx context.Context, integrationTypeID uint) ([]*Integration, error)
	SetIntegrationAsDefault(ctx context.Context, id uint) error
	ExistsIntegrationByCode(ctx context.Context, code string, businessID *uint) (bool, error)
	UpdateLastSync(ctx context.Context, id uint, lastSync time.Time) error

	// Métodos de IntegrationTypes
	CreateIntegrationType(ctx context.Context, integrationType *IntegrationType) error
	UpdateIntegrationType(ctx context.Context, id uint, integrationType *IntegrationType) error
	GetIntegrationTypeByID(ctx context.Context, id uint) (*IntegrationType, error)
	GetIntegrationTypeByCode(ctx context.Context, code string) (*IntegrationType, error)
	GetIntegrationTypeByName(ctx context.Context, name string) (*IntegrationType, error)
	DeleteIntegrationType(ctx context.Context, id uint) error
	ListIntegrationTypes(ctx context.Context) ([]*IntegrationType, error)
	ListActiveIntegrationTypes(ctx context.Context) ([]*IntegrationType, error)
}

// IEncryptionService define la interfaz del servicio de encriptación
type IEncryptionService interface {
	// Encriptar credenciales antes de guardar
	EncryptCredentials(ctx context.Context, credentials map[string]interface{}) ([]byte, error)

	// Desencriptar credenciales para usar
	DecryptCredentials(ctx context.Context, encryptedData []byte) (map[string]interface{}, error)

	// Encriptar un valor individual
	EncryptValue(ctx context.Context, value string) (string, error)

	// Desencriptar un valor individual
	DecryptValue(ctx context.Context, encryptedValue string) (string, error)
}

// IIntegrationTypeUseCase define la interfaz del caso de uso de tipos de integración

// IIntegrationUseCase define la interfaz del caso de uso de integraciones
