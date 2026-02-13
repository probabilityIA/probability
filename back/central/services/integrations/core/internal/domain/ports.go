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

	// Integration Categories
	GetIntegrationCategoryByID(ctx context.Context, id uint) (*IntegrationCategory, error)
	ListIntegrationCategories(ctx context.Context) ([]*IntegrationCategory, error)
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
	SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error
	SyncOrdersByBusiness(ctx context.Context, businessID uint) error
	GetWebhookURL(ctx context.Context, integrationID uint) (*WebhookInfo, error)
	ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
	VerifyWebhooksByURL(ctx context.Context, integrationID string) ([]interface{}, error)
	CreateWebhook(ctx context.Context, integrationID string) (interface{}, error)
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

// ============================================
// INTEGRATION CACHE
// ============================================

// CachedIntegration representa los datos de integración en cache
type CachedIntegration struct {
	ID                  uint                   `json:"id"`
	Name                string                 `json:"name"`
	Code                string                 `json:"code"`
	Category            string                 `json:"category"`
	IntegrationTypeID   uint                   `json:"integration_type_id"`
	IntegrationTypeCode string                 `json:"integration_type_code"`
	BusinessID          *uint                  `json:"business_id"`
	StoreID             string                 `json:"store_id"`
	IsActive            bool                   `json:"is_active"`
	IsDefault           bool                   `json:"is_default"`
	Config              map[string]interface{} `json:"config"`
	Description         string                 `json:"description"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// CachedCredentials representa credenciales desencriptadas en cache
type CachedCredentials struct {
	IntegrationID uint                   `json:"integration_id"`
	Credentials   map[string]interface{} `json:"credentials"` // DESENCRIPTADAS
	CachedAt      time.Time              `json:"cached_at"`
}

// IIntegrationCache define operaciones de cache para integraciones
type IIntegrationCache interface {
	// Metadata (TTL: 24h)
	SetIntegration(ctx context.Context, integration *CachedIntegration) error
	GetIntegration(ctx context.Context, integrationID uint) (*CachedIntegration, error)

	// Credentials (TTL: 1h - sensibles)
	SetCredentials(ctx context.Context, creds *CachedCredentials) error
	GetCredentials(ctx context.Context, integrationID uint) (*CachedCredentials, error)
	GetCredentialField(ctx context.Context, integrationID uint, field string) (string, error)

	// Invalidación
	InvalidateIntegration(ctx context.Context, integrationID uint) error

	// Búsquedas indexadas
	GetByCode(ctx context.Context, code string) (*CachedIntegration, error)
	GetByBusinessAndType(ctx context.Context, businessID, integrationTypeID uint) (*CachedIntegration, error)
}
