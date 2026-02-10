package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// ═══════════════════════════════════════════════════════════════
<<<<<<< HEAD
// REPOSITORIOS (Secondary Ports - Driven Adapters)
// ═══════════════════════════════════════════════════════════════

// IInvoiceRepository define las operaciones de persistencia para facturas
type IInvoiceRepository interface {
	Create(ctx context.Context, invoice *entities.Invoice) error
	GetByID(ctx context.Context, id uint) (*entities.Invoice, error)
	GetByOrderID(ctx context.Context, orderID string) (*entities.Invoice, error)
	GetByOrderAndProvider(ctx context.Context, orderID string, providerID uint) (*entities.Invoice, error)
	List(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, error)
	Update(ctx context.Context, invoice *entities.Invoice) error
	Delete(ctx context.Context, id uint) error
	ExistsForOrder(ctx context.Context, orderID string, providerID uint) (bool, error)

	// Estadísticas y resúmenes
	GetSummary(ctx context.Context, businessID uint, start, end time.Time) (*entities.InvoiceSummary, error)
	GetDetailedStats(ctx context.Context, businessID uint, filters map[string]interface{}) (*entities.DetailedStats, error)
	GetTrends(ctx context.Context, businessID uint, start, end time.Time, granularity, metric string) (*entities.TrendData, error)
}

// IInvoiceItemRepository define las operaciones de persistencia para items de factura
type IInvoiceItemRepository interface {
	Create(ctx context.Context, item *entities.InvoiceItem) error
	GetByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceItem, error)
	UpdateBatch(ctx context.Context, items []*entities.InvoiceItem) error
}

// IInvoicingProviderRepository define las operaciones de persistencia para proveedores
type IInvoicingProviderRepository interface {
	Create(ctx context.Context, provider *entities.InvoicingProvider) error
	GetByID(ctx context.Context, id uint) (*entities.InvoicingProvider, error)
	GetByBusinessAndType(ctx context.Context, businessID uint, providerTypeCode string) (*entities.InvoicingProvider, error)
	GetDefaultByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingProvider, error)
	List(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error)
	Update(ctx context.Context, provider *entities.InvoicingProvider) error
	Delete(ctx context.Context, id uint) error
}

// IInvoicingProviderTypeRepository define las operaciones de persistencia para tipos de proveedor
type IInvoicingProviderTypeRepository interface {
	GetByCode(ctx context.Context, code string) (*entities.InvoicingProviderType, error)
	List(ctx context.Context) ([]*entities.InvoicingProviderType, error)
	GetActive(ctx context.Context) ([]*entities.InvoicingProviderType, error)
}

// IInvoicingConfigRepository define las operaciones de persistencia para configuraciones
type IInvoicingConfigRepository interface {
	Create(ctx context.Context, config *entities.InvoicingConfig) error
	GetByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error)
	GetByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error)
	List(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error)
	Update(ctx context.Context, config *entities.InvoicingConfig) error
	Delete(ctx context.Context, id uint) error
	ExistsForIntegration(ctx context.Context, integrationID uint) (bool, error)
}

// IInvoiceSyncLogRepository define las operaciones de persistencia para logs de sincronización
type IInvoiceSyncLogRepository interface {
	Create(ctx context.Context, log *entities.InvoiceSyncLog) error
	GetByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error)
	GetPendingRetries(ctx context.Context, limit int) ([]*entities.InvoiceSyncLog, error)
	Update(ctx context.Context, log *entities.InvoiceSyncLog) error
}

// ICreditNoteRepository define las operaciones de persistencia para notas de crédito
type ICreditNoteRepository interface {
	Create(ctx context.Context, note *entities.CreditNote) error
	GetByID(ctx context.Context, id uint) (*entities.CreditNote, error)
	GetByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.CreditNote, error)
	List(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error)
	Update(ctx context.Context, note *entities.CreditNote) error
}
=======
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

>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

// ═══════════════════════════════════════════════════════════════
// CLIENTE DE PROVEEDOR (Secondary Port - Driven Adapter)
// ═══════════════════════════════════════════════════════════════

<<<<<<< HEAD
// InvoiceRequest representa los datos necesarios para crear una factura en el proveedor
type InvoiceRequest struct {
	Invoice      *entities.Invoice
	InvoiceItems []*entities.InvoiceItem
	Provider     *entities.InvoicingProvider
	Config       map[string]interface{}
}

// InvoiceResponse representa la respuesta del proveedor al crear una factura
type InvoiceResponse struct {
	InvoiceNumber string
	ExternalID    string
	InvoiceURL    *string
	PDFURL        *string
	XMLURL        *string
	CUFE          *string
	IssuedAt      string
	RawResponse   map[string]interface{}
}

// CreditNoteRequest representa los datos necesarios para crear una nota de crédito
type CreditNoteRequest struct {
	Invoice     *entities.Invoice
	CreditNote  *entities.CreditNote
	Provider    *entities.InvoicingProvider
}

// CreditNoteResponse representa la respuesta del proveedor al crear una nota de crédito
type CreditNoteResponse struct {
	CreditNoteNumber string
	ExternalID       string
	NoteURL          *string
	PDFURL           *string
	XMLURL           *string
	CUFE             *string
	IssuedAt         string
	RawResponse      map[string]interface{}
}

=======
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
// IInvoicingProviderClient define las operaciones que debe implementar un cliente de proveedor
type IInvoicingProviderClient interface {
	// Autenticación
	Authenticate(ctx context.Context, credentials map[string]interface{}) (string, error)

	// Crear factura
<<<<<<< HEAD
	CreateInvoice(ctx context.Context, token string, request *InvoiceRequest) (*InvoiceResponse, error)
=======
	CreateInvoice(ctx context.Context, token string, request *dtos.InvoiceRequest) (*dtos.InvoiceResponse, error)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

	// Cancelar factura
	CancelInvoice(ctx context.Context, token string, externalID string, reason string) error

	// Crear nota de crédito
<<<<<<< HEAD
	CreateCreditNote(ctx context.Context, token string, request *CreditNoteRequest) (*CreditNoteResponse, error)
=======
	CreateCreditNote(ctx context.Context, token string, request *dtos.CreditNoteRequest) (*dtos.CreditNoteResponse, error)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

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
<<<<<<< HEAD
=======
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
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
}

// ═══════════════════════════════════════════════════════════════
// REPOSITORIO DE ÓRDENES (Secondary Port - Dependencia externa)
// ═══════════════════════════════════════════════════════════════

<<<<<<< HEAD
// OrderData representa los datos mínimos necesarios de una orden para facturación
type OrderData struct {
	// Campos existentes
	ID               string
	BusinessID       uint
	IntegrationID    uint
	OrderNumber      string
	TotalAmount      float64
	Subtotal         float64
	Tax              float64
	Discount         float64
	ShippingCost     float64
	Currency         string
	CustomerName     string
	CustomerEmail    string
	CustomerPhone    string
	CustomerDNI      string
	IsPaid           bool
	PaymentMethodID  uint
	Invoiceable      bool
	Items            []OrderItemData

	// Campos nuevos - Necesarios para filtros avanzados
	Status          string     // Estado de la orden (pending, confirmed, etc.)
	OrderTypeID     uint       // Tipo de orden (delivery, pickup, etc.)
	OrderTypeName   string     // Nombre del tipo de orden
	CustomerID      *string    // ID del cliente (para exclusiones)
	CustomerType    *string    // Tipo de cliente (natural, juridica)
	ShippingCity    *string    // Ciudad de envío
	ShippingState   *string    // Departamento/Estado de envío
	ShippingCountry *string    // País de envío
	CreatedAt       time.Time  // Fecha de creación
}

// OrderItemData representa un item de orden
type OrderItemData struct {
	ProductID   *string
	SKU         string
	Name        string
	Description *string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
	Tax         float64
	TaxRate     *float64
	Discount    float64

	// Campos nuevos - Para filtros de productos
	CategoryID   *uint   // Categoría del producto
	CategoryName *string // Nombre de la categoría
}

// IOrderRepository define las operaciones para obtener datos de órdenes
type IOrderRepository interface {
	GetByID(ctx context.Context, orderID string) (*OrderData, error)
	UpdateInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error
=======
// IOrderRepository define las operaciones para obtener datos de órdenes
type IOrderRepository interface {
	GetByID(ctx context.Context, orderID string) (*dtos.OrderData, error)
	UpdateInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error
	GetInvoiceableOrders(ctx context.Context, businessID uint, page, pageSize int) ([]*dtos.OrderData, int64, error)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
}

// ═══════════════════════════════════════════════════════════════
// CASOS DE USO (Primary Port - Driver)
// ═══════════════════════════════════════════════════════════════

// IUseCase define todos los casos de uso del módulo de facturación
type IUseCase interface {
	// Facturas
	CreateInvoice(ctx context.Context, dto *dtos.CreateInvoiceDTO) (*entities.Invoice, error)
	CancelInvoice(ctx context.Context, dto *dtos.CancelInvoiceDTO) error
	RetryInvoice(ctx context.Context, invoiceID uint) error
<<<<<<< HEAD
	GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, error)
	GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error)
=======
	CancelRetry(ctx context.Context, invoiceID uint) error
	EnableRetry(ctx context.Context, invoiceID uint) error
	GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, int64, error)
	GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error)
	GetInvoiceSyncLogs(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

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
<<<<<<< HEAD
=======

	// Creación masiva de facturas (DEPRECADO - Síncrono)
	// DEPRECATED: Usar BulkCreateInvoicesAsync para procesamiento asíncrono
	BulkCreateInvoices(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (*dtos.BulkCreateResult, error)

	// Creación masiva de facturas (Asíncrono con RabbitMQ)
	BulkCreateInvoicesAsync(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (string, error)
	GetBulkJobStatus(ctx context.Context, jobID string) (*entities.BulkInvoiceJob, error)
	GetBulkJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error)
	ListBulkJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
}
