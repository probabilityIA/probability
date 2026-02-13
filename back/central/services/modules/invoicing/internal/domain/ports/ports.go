package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// ═══════════════════════════════════════════════════════════════
// REPOSITORIO (Secondary Port - Driven Adapter)
// ═══════════════════════════════════════════════════════════════

// IRepository define TODAS las operaciones de persistencia del módulo de facturación
type IRepository interface {
	// ═══════════════════════════════════════════
	// INVOICES
	// ═══════════════════════════════════════════
	CreateInvoice(ctx context.Context, invoice *entities.Invoice) error
	GetInvoiceByID(ctx context.Context, id uint) (*entities.Invoice, error)
	GetInvoiceByOrderID(ctx context.Context, orderID string) (*entities.Invoice, error)
	GetInvoiceByOrderAndProvider(ctx context.Context, orderID string, providerID uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, int64, error)
	UpdateInvoice(ctx context.Context, invoice *entities.Invoice) error
	DeleteInvoice(ctx context.Context, id uint) error
	InvoiceExistsForOrder(ctx context.Context, orderID string, providerID uint) (bool, error)
	GetInvoiceSummary(ctx context.Context, businessID uint, start, end time.Time) (*entities.InvoiceSummary, error)
	GetInvoiceDetailedStats(ctx context.Context, businessID uint, filters map[string]interface{}) (*entities.DetailedStats, error)
	GetInvoiceTrends(ctx context.Context, businessID uint, start, end time.Time, granularity, metric string) (*entities.TrendData, error)

	// ═══════════════════════════════════════════
	// INVOICE ITEMS
	// ═══════════════════════════════════════════
	CreateInvoiceItem(ctx context.Context, item *entities.InvoiceItem) error
	GetInvoiceItemsByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceItem, error)
	UpdateInvoiceItemsBatch(ctx context.Context, items []*entities.InvoiceItem) error

	// ═══════════════════════════════════════════
	// INVOICING PROVIDERS
	// ═══════════════════════════════════════════
	CreateInvoicingProvider(ctx context.Context, provider *entities.InvoicingProvider) error
	GetInvoicingProviderByID(ctx context.Context, id uint) (*entities.InvoicingProvider, error)
	GetProviderByBusinessAndType(ctx context.Context, businessID uint, providerTypeCode string) (*entities.InvoicingProvider, error)
	GetDefaultProviderByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingProvider, error)
	ListInvoicingProviders(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error)
	UpdateInvoicingProvider(ctx context.Context, provider *entities.InvoicingProvider) error
	DeleteInvoicingProvider(ctx context.Context, id uint) error

	// ═══════════════════════════════════════════
	// INVOICING PROVIDER TYPES
	// ═══════════════════════════════════════════
	GetProviderTypeByCode(ctx context.Context, code string) (*entities.InvoicingProviderType, error)
	ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error)
	GetActiveProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error)

	// ═══════════════════════════════════════════
	// INVOICING CONFIGS
	// ═══════════════════════════════════════════
	CreateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error
	GetInvoicingConfigByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error)
	GetConfigByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error)
	ListInvoicingConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*entities.InvoicingConfig, error)
	UpdateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error
	DeleteInvoicingConfig(ctx context.Context, id uint) error
	ConfigExistsForIntegration(ctx context.Context, integrationID uint) (bool, error)

	// ═══════════════════════════════════════════
	// INVOICE SYNC LOGS
	// ═══════════════════════════════════════════
	CreateInvoiceSyncLog(ctx context.Context, log *entities.InvoiceSyncLog) error
	GetSyncLogsByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error)
	GetPendingSyncLogRetries(ctx context.Context, limit int) ([]*entities.InvoiceSyncLog, error)
	UpdateInvoiceSyncLog(ctx context.Context, log *entities.InvoiceSyncLog) error

	// ═══════════════════════════════════════════
	// CREDIT NOTES
	// ═══════════════════════════════════════════
	CreateCreditNote(ctx context.Context, note *entities.CreditNote) error
	GetCreditNoteByID(ctx context.Context, id uint) (*entities.CreditNote, error)
	GetCreditNotesByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.CreditNote, error)
	ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error)
	UpdateCreditNote(ctx context.Context, note *entities.CreditNote) error

	// ═══════════════════════════════════════════
	// BULK INVOICE JOBS
	// ═══════════════════════════════════════════
	CreateJob(ctx context.Context, job *entities.BulkInvoiceJob) error
	CreateJobItems(ctx context.Context, items []*entities.BulkInvoiceJobItem) error
	GetJobByID(ctx context.Context, jobID string) (*entities.BulkInvoiceJob, error)
	GetJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error)
	UpdateJob(ctx context.Context, job *entities.BulkInvoiceJob) error
	UpdateJobItem(ctx context.Context, item *entities.BulkInvoiceJobItem) error
	ListJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error)
	IncrementJobCounters(ctx context.Context, jobID string, processed, successful, failed int) error

	// ═══════════════════════════════════════════
	// ORDERS
	// ═══════════════════════════════════════════
	GetOrderByID(ctx context.Context, orderID string) (*dtos.OrderData, error)
	UpdateOrderInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error
	GetInvoiceableOrders(ctx context.Context, businessID uint, page, pageSize int) ([]*dtos.OrderData, int64, error)
}

// ═══════════════════════════════════════════════════════════════
// CLIENTE DE PROVEEDOR (Secondary Port - Driven Adapter)
// ═══════════════════════════════════════════════════════════════

// IInvoicingProviderClient define las operaciones que debe implementar un cliente de proveedor
type IInvoicingProviderClient interface {
	// Autenticación
	Authenticate(ctx context.Context, credentials map[string]interface{}) (string, error)

	// Crear factura
	CreateInvoice(ctx context.Context, token string, request *dtos.InvoiceRequest) (*dtos.InvoiceResponse, error)

	// Cancelar factura
	CancelInvoice(ctx context.Context, token string, externalID string, reason string) error

	// Crear nota de crédito
	CreateCreditNote(ctx context.Context, token string, request *dtos.CreditNoteRequest) (*dtos.CreditNoteResponse, error)

	// Consultar estado de factura
	GetInvoiceStatus(ctx context.Context, token string, externalID string) (string, error)

	// Validar credenciales (test de conexión)
	ValidateCredentials(ctx context.Context, credentials map[string]interface{}) error
}

// ═══════════════════════════════════════════════════════════════
// SERVICIOS EXTERNOS (Secondary Ports - Driven Adapters)
// ═══════════════════════════════════════════════════════════════

// IEncryptionService define las operaciones de encriptación/desencriptación
type IEncryptionService interface {
	Encrypt(data map[string]interface{}) (map[string]interface{}, error)
	Decrypt(data map[string]interface{}) (map[string]interface{}, error)
	EncryptString(text string) (string, error)
	DecryptString(encrypted string) (string, error)
}

// IEventPublisher define las operaciones para publicar eventos
type IEventPublisher interface {
	PublishInvoiceCreated(ctx context.Context, invoice *entities.Invoice) error
	PublishInvoiceCancelled(ctx context.Context, invoice *entities.Invoice) error
	PublishInvoiceFailed(ctx context.Context, invoice *entities.Invoice, errorMsg string) error
	PublishCreditNoteCreated(ctx context.Context, creditNote *entities.CreditNote) error
	PublishBulkInvoiceJob(ctx context.Context, message *dtos.BulkInvoiceJobMessage) error
}

// IInvoiceSSEPublisher publica eventos a Redis Pub/Sub para SSE en tiempo real
type IInvoiceSSEPublisher interface {
	PublishInvoiceCreated(ctx context.Context, invoice *entities.Invoice) error
	PublishInvoiceFailed(ctx context.Context, invoice *entities.Invoice, errorMsg string) error
	PublishInvoiceCancelled(ctx context.Context, invoice *entities.Invoice) error
	PublishCreditNoteCreated(ctx context.Context, creditNote *entities.CreditNote) error
	PublishBulkJobProgress(ctx context.Context, job *entities.BulkInvoiceJob) error
	PublishBulkJobCompleted(ctx context.Context, job *entities.BulkInvoiceJob) error
}

// IInvoiceRequestPublisher publica solicitudes de facturación a colas específicas de proveedores
// Los proveedores (Softpymes, Siigo, Factus) consumen estas solicitudes y publican respuestas
type IInvoiceRequestPublisher interface {
	PublishInvoiceRequest(ctx context.Context, request *dtos.InvoiceRequestMessage) error
}

// ═══════════════════════════════════════════════════════════════
// REPOSITORIO DE ÓRDENES (Secondary Port - Dependencia externa)
// ═══════════════════════════════════════════════════════════════

// IOrderRepository define las operaciones para obtener datos de órdenes
type IOrderRepository interface {
	GetByID(ctx context.Context, orderID string) (*dtos.OrderData, error)
	UpdateInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error
	GetInvoiceableOrders(ctx context.Context, businessID uint, page, pageSize int) ([]*dtos.OrderData, int64, error)
}

// ═══════════════════════════════════════════════════════════════
// CASOS DE USO (Primary Port - Driver)
// ═══════════════════════════════════════════════════════════════

// IUseCase define todos los casos de uso del módulo de facturación
type IUseCase interface {
	// Cache warming
	WarmConfigCache(ctx context.Context) error

	// Facturas
	CreateInvoice(ctx context.Context, dto *dtos.CreateInvoiceDTO) (*entities.Invoice, error)
	CancelInvoice(ctx context.Context, dto *dtos.CancelInvoiceDTO) error
	RetryInvoice(ctx context.Context, invoiceID uint) error
	CancelRetry(ctx context.Context, invoiceID uint) error
	EnableRetry(ctx context.Context, invoiceID uint) error
	GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, int64, error)
	GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error)
	GetInvoiceSyncLogs(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error)

	// Proveedores (DEPRECADOS - Migrados a integrations/core)
	// NOTA: Estos métodos están deprecados y serán eliminados en una futura versión
	// Usar integrations/core para gestión de proveedores de facturación
	CreateProvider(ctx context.Context, dto *dtos.CreateProviderDTO) (*entities.InvoicingProvider, error)
	UpdateProvider(ctx context.Context, id uint, dto *dtos.UpdateProviderDTO) error
	GetProvider(ctx context.Context, id uint) (*entities.InvoicingProvider, error)
	ListProviders(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error)
	TestProviderConnection(ctx context.Context, id uint) error

	// Configuraciones
	CreateConfig(ctx context.Context, dto *dtos.CreateConfigDTO) (*entities.InvoicingConfig, error)
	UpdateConfig(ctx context.Context, id uint, dto *dtos.UpdateConfigDTO) (*entities.InvoicingConfig, error)
	GetConfig(ctx context.Context, id uint) (*entities.InvoicingConfig, error)
	ListConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error)
	DeleteConfig(ctx context.Context, id uint) error

	// Notas de crédito
	CreateCreditNote(ctx context.Context, dto *dtos.CreateCreditNoteDTO) (*entities.CreditNote, error)
	GetCreditNote(ctx context.Context, id uint) (*entities.CreditNote, error)
	ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error)

	// Tipos de proveedores (DEPRECADO - Migrado a integrations/core)
	// NOTA: Este método está deprecado y será eliminado en una futura versión
	// Usar integrations/core para listar tipos de integraciones de facturación
	ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error)

	// Estadísticas y resúmenes
	GetSummary(ctx context.Context, businessID uint, period string) (*entities.InvoiceSummary, error)
	GetDetailedStats(ctx context.Context, businessID uint, filters map[string]interface{}) (*entities.DetailedStats, error)
	GetTrends(ctx context.Context, businessID uint, startDate, endDate, granularity, metric string) (*entities.TrendData, error)

	// Creación masiva de facturas (DEPRECADO - Síncrono)
	// DEPRECATED: Usar BulkCreateInvoicesAsync para procesamiento asíncrono
	BulkCreateInvoices(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (*dtos.BulkCreateResult, error)

	// Creación masiva de facturas (Asíncrono con RabbitMQ)
	BulkCreateInvoicesAsync(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (string, error)
	GetBulkJobStatus(ctx context.Context, jobID string) (*entities.BulkInvoiceJob, error)
	GetBulkJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error)
	ListBulkJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error)
}
