package domain

import "context"

// IIntegrationRepository define la interfaz del repositorio de integraciones
type IIntegrationRepository interface {
	// CRUD básico
	Create(ctx context.Context, integration *Integration) error
	Update(ctx context.Context, id uint, integration *Integration) error
	GetByID(ctx context.Context, id uint) (*Integration, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filters IntegrationFilters) ([]*Integration, int64, error)

	// Métodos específicos
	GetByType(ctx context.Context, integrationType string, businessID *uint) (*Integration, error)
	GetActiveByType(ctx context.Context, integrationType string, businessID *uint) (*Integration, error)
	ListByBusiness(ctx context.Context, businessID uint) ([]*Integration, error)
	ListByType(ctx context.Context, integrationType string) ([]*Integration, error)
	SetAsDefault(ctx context.Context, id uint) error
	ExistsByCode(ctx context.Context, code string, businessID *uint) (bool, error)
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

// IIntegrationUseCase define la interfaz del caso de uso de integraciones
type IIntegrationUseCase interface {
	CreateIntegration(ctx context.Context, dto CreateIntegrationDTO) (*Integration, error)
	UpdateIntegration(ctx context.Context, id uint, dto UpdateIntegrationDTO) (*Integration, error)
	GetIntegrationByID(ctx context.Context, id uint) (*Integration, error)
	GetIntegrationByType(ctx context.Context, integrationType string, businessID *uint) (*IntegrationWithCredentials, error)
	DeleteIntegration(ctx context.Context, id uint) error
	ListIntegrations(ctx context.Context, filters IntegrationFilters) ([]*Integration, int64, error)
	TestIntegration(ctx context.Context, id uint) error
	ActivateIntegration(ctx context.Context, id uint) error
	DeactivateIntegration(ctx context.Context, id uint) error
	SetAsDefault(ctx context.Context, id uint) error
}
