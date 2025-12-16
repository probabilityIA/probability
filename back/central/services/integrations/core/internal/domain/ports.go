package domain

import (
	"context"
	"io"
	"mime/multipart"
	"time"
)

type IRepository interface {
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

	CreateIntegrationType(ctx context.Context, integrationType *IntegrationType) error
	UpdateIntegrationType(ctx context.Context, id uint, integrationType *IntegrationType) error
	GetIntegrationTypeByID(ctx context.Context, id uint) (*IntegrationType, error)
	GetIntegrationTypeByCode(ctx context.Context, code string) (*IntegrationType, error)
	GetIntegrationTypeByName(ctx context.Context, name string) (*IntegrationType, error)
	DeleteIntegrationType(ctx context.Context, id uint) error
	ListIntegrationTypes(ctx context.Context) ([]*IntegrationType, error)
	ListActiveIntegrationTypes(ctx context.Context) ([]*IntegrationType, error)
}

type IEncryptionService interface {
	EncryptCredentials(ctx context.Context, credentials map[string]interface{}) ([]byte, error)
	DecryptCredentials(ctx context.Context, encryptedData []byte) (map[string]interface{}, error)
	EncryptValue(ctx context.Context, value string) (string, error)
	DecryptValue(ctx context.Context, encryptedValue string) (string, error)
}

// WebhookInfo contiene la información del webhook para una integración
type WebhookInfo struct {
	URL         string   `json:"url"`
	Method      string   `json:"method"`
	Description string   `json:"description"`
	Events      []string `json:"events,omitempty"`
}

type IOrderSyncService interface {
	SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error
	SyncOrdersByBusiness(ctx context.Context, businessID uint) error
	GetWebhookURL(ctx context.Context, integrationID uint) (*WebhookInfo, error)
	ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
}

type IS3Service interface {
	GetImageURL(filename string) string
	DeleteImage(ctx context.Context, filename string) error
	ImageExists(ctx context.Context, filename string) (bool, error)
	UploadFile(ctx context.Context, file io.ReadSeeker, filename string) (string, error)
	DownloadFile(ctx context.Context, filename string) (io.ReadSeeker, error)
	FileExists(ctx context.Context, filename string) (bool, error)
	GetFileURL(ctx context.Context, filename string) (string, error)
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
}
