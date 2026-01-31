package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// ═══════════════════════════════════════════════════════════════
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

// ═══════════════════════════════════════════════════════════════
// CLIENTE DE PROVEEDOR (Secondary Port - Driven Adapter)
// ═══════════════════════════════════════════════════════════════

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

// IInvoicingProviderClient define las operaciones que debe implementar un cliente de proveedor
type IInvoicingProviderClient interface {
	// Autenticación
	Authenticate(ctx context.Context, credentials map[string]interface{}) (string, error)

	// Crear factura
	CreateInvoice(ctx context.Context, token string, request *InvoiceRequest) (*InvoiceResponse, error)

	// Cancelar factura
	CancelInvoice(ctx context.Context, token string, externalID string, reason string) error

	// Crear nota de crédito
	CreateCreditNote(ctx context.Context, token string, request *CreditNoteRequest) (*CreditNoteResponse, error)

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
}

// ═══════════════════════════════════════════════════════════════
// REPOSITORIO DE ÓRDENES (Secondary Port - Dependencia externa)
// ═══════════════════════════════════════════════════════════════

// OrderData representa los datos mínimos necesarios de una orden para facturación
type OrderData struct {
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
}

// IOrderRepository define las operaciones para obtener datos de órdenes
type IOrderRepository interface {
	GetByID(ctx context.Context, orderID string) (*OrderData, error)
	UpdateInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error
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
	GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, error)
	GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error)

	// Proveedores
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

	// Tipos de proveedores
	ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error)
}
