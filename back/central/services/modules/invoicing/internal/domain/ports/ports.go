package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

type IRepository interface {
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

	CreateInvoiceItem(ctx context.Context, item *entities.InvoiceItem) error
	GetInvoiceItemsByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceItem, error)
	UpdateInvoiceItemsBatch(ctx context.Context, items []*entities.InvoiceItem) error

	CreateInvoicingProvider(ctx context.Context, provider *entities.InvoicingProvider) error
	GetInvoicingProviderByID(ctx context.Context, id uint) (*entities.InvoicingProvider, error)
	GetProviderByBusinessAndType(ctx context.Context, businessID uint, providerTypeCode string) (*entities.InvoicingProvider, error)
	GetDefaultProviderByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingProvider, error)
	ListInvoicingProviders(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error)
	UpdateInvoicingProvider(ctx context.Context, provider *entities.InvoicingProvider) error
	DeleteInvoicingProvider(ctx context.Context, id uint) error

	GetProviderTypeByCode(ctx context.Context, code string) (*entities.InvoicingProviderType, error)
	ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error)
	GetActiveProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error)

	CreateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error
	GetInvoicingConfigByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error)
	GetConfigByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error)
	ListInvoicingConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*entities.InvoicingConfig, error)
	UpdateInvoicingConfig(ctx context.Context, config *entities.InvoicingConfig) error
	DeleteInvoicingConfig(ctx context.Context, id uint) error
	ConfigExistsForIntegration(ctx context.Context, integrationID uint) (bool, error)
	GetConfigByInvoicingIntegration(ctx context.Context, businessID uint, invoicingIntegrationID uint) (*entities.InvoicingConfig, error)
	GetEnabledConfigByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingConfig, error)
	GetAnyConfigByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingConfig, error)

	CreateInvoiceSyncLog(ctx context.Context, log *entities.InvoiceSyncLog) error
	GetSyncLogsByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error)
	GetPendingSyncLogRetries(ctx context.Context, limit int) ([]*entities.InvoiceSyncLog, error)
	UpdateInvoiceSyncLog(ctx context.Context, log *entities.InvoiceSyncLog) error

	CreateCreditNote(ctx context.Context, note *entities.CreditNote) error
	GetCreditNoteByID(ctx context.Context, id uint) (*entities.CreditNote, error)
	GetCreditNotesByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.CreditNote, error)
	ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error)
	UpdateCreditNote(ctx context.Context, note *entities.CreditNote) error

	CreateJob(ctx context.Context, job *entities.BulkInvoiceJob) error
	CreateJobItems(ctx context.Context, items []*entities.BulkInvoiceJobItem) error
	GetJobByID(ctx context.Context, jobID string) (*entities.BulkInvoiceJob, error)
	GetJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error)
	GetJobItemByJobAndOrder(ctx context.Context, jobID, orderID string) (*entities.BulkInvoiceJobItem, error)
	GetJobItemByInvoiceID(ctx context.Context, invoiceID uint) (*entities.BulkInvoiceJobItem, error)
	UpdateJob(ctx context.Context, job *entities.BulkInvoiceJob) error
	UpdateJobItem(ctx context.Context, item *entities.BulkInvoiceJobItem) error
	ListJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error)
	IncrementJobCounters(ctx context.Context, jobID string, processed, successful, failed int) error

	GetOrderByID(ctx context.Context, orderID string) (*dtos.OrderData, error)
	CountOrdersOutsideBusiness(ctx context.Context, orderIDs []string, businessID uint) (int64, error)
	UpdateOrderInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error
	GetInvoiceableOrders(ctx context.Context, filter dtos.InvoiceableOrdersFilter) ([]*dtos.OrderData, int64, error)
	CancelInvoiceAndReleaseOrder(ctx context.Context, invoiceID uint) (bool, string, error)

	GetIntegrationTypeByIntegrationID(ctx context.Context, integrationID uint) (int, error)

	GetIntegrationBusinessID(ctx context.Context, integrationID uint) (uint, error)

	GetIssuedInvoicesByDateRange(ctx context.Context, businessID uint, dateFrom, dateTo string) ([]*entities.Invoice, error)

	GetOrderCreatedAtsByIDs(ctx context.Context, orderIDs []string) (map[string]*time.Time, error)

	ListProductsByBusinessID(ctx context.Context, businessID uint) ([]dtos.SystemProduct, error)
}

type IInvoicingProviderClient interface {
	Authenticate(ctx context.Context, credentials map[string]interface{}) (string, error)

	CreateInvoice(ctx context.Context, token string, request *dtos.InvoiceRequest) (*dtos.InvoiceResponse, error)

	CancelInvoice(ctx context.Context, token string, externalID string, reason string) error

	CreateCreditNote(ctx context.Context, token string, request *dtos.CreditNoteRequest) (*dtos.CreditNoteResponse, error)

	GetInvoiceStatus(ctx context.Context, token string, externalID string) (string, error)

	ValidateCredentials(ctx context.Context, credentials map[string]interface{}) error
}

type IEncryptionService interface {
	Encrypt(data map[string]interface{}) (map[string]interface{}, error)
	Decrypt(data map[string]interface{}) (map[string]interface{}, error)
	EncryptString(text string) (string, error)
	DecryptString(encrypted string) (string, error)
}

type IEventPublisher interface {
	PublishInvoiceCreated(ctx context.Context, invoice *entities.Invoice) error
	PublishInvoiceCancelled(ctx context.Context, invoice *entities.Invoice) error
	PublishInvoiceFailed(ctx context.Context, invoice *entities.Invoice, errorMsg string) error
	PublishCreditNoteCreated(ctx context.Context, creditNote *entities.CreditNote) error
	PublishBulkInvoiceJob(ctx context.Context, message *dtos.BulkInvoiceJobMessage) error
}

type IInvoiceSSEPublisher interface {
	PublishInvoiceCreated(ctx context.Context, invoice *entities.Invoice) error
	PublishInvoiceFailed(ctx context.Context, invoice *entities.Invoice, errorMsg string) error
	PublishInvoiceCancelled(ctx context.Context, invoice *entities.Invoice) error
	PublishCreditNoteCreated(ctx context.Context, creditNote *entities.CreditNote) error
	PublishBulkJobProgress(ctx context.Context, job *entities.BulkInvoiceJob) error
	PublishBulkJobCompleted(ctx context.Context, job *entities.BulkInvoiceJob) error
	PublishCompareReady(ctx context.Context, data *dtos.CompareResponseData) error
	PublishListItemsReady(ctx context.Context, data *dtos.ItemCompareResponseData) error
}

type IInvoiceRequestPublisher interface {
	PublishInvoiceRequest(ctx context.Context, request *dtos.InvoiceRequestMessage) error
}

type ICompareCache interface {
	StoreCompareResult(ctx context.Context, correlationID string, data *dtos.CompareResponseData) error
	GetCompareResult(ctx context.Context, correlationID string) (*dtos.CompareResponseData, error)
	StoreItemCompareResult(ctx context.Context, correlationID string, data *dtos.ItemCompareResponseData) error
	GetItemCompareResult(ctx context.Context, correlationID string) (*dtos.ItemCompareResponseData, error)
	StoreBankAccountsResult(ctx context.Context, correlationID string, data *dtos.BankAccountsResponseData) error
	GetBankAccountsResult(ctx context.Context, correlationID string) (*dtos.BankAccountsResponseData, error)
}

type IConfigCache interface {
	Get(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error)
	Set(ctx context.Context, integrationID uint, config *entities.InvoicingConfig) error
	Invalidate(ctx context.Context, integrationID uint) error

	GetByBusinessID(ctx context.Context, businessID uint) (*entities.InvoicingConfig, error)
	SetByBusinessID(ctx context.Context, businessID uint, config *entities.InvoicingConfig) error
	InvalidateByBusinessID(ctx context.Context, businessID uint) error
}

type IOrderRepository interface {
	GetByID(ctx context.Context, orderID string) (*dtos.OrderData, error)
	UpdateInvoiceInfo(ctx context.Context, orderID string, invoiceID string, invoiceURL string) error
	GetInvoiceableOrders(ctx context.Context, filter dtos.InvoiceableOrdersFilter) ([]*dtos.OrderData, int64, error)
}

type IUseCase interface {
	WarmConfigCache(ctx context.Context) error

	CreateInvoice(ctx context.Context, dto *dtos.CreateInvoiceDTO) (*entities.Invoice, error)
	CreateJournal(ctx context.Context, dto *dtos.CreateJournalDTO) (*entities.Invoice, error)
	RegisterManualInvoice(ctx context.Context, dto *dtos.RegisterManualInvoiceDTO) (*entities.Invoice, error)
	CancelInvoice(ctx context.Context, dto *dtos.CancelInvoiceDTO) error
	MarkInvoiceAsCancelled(ctx context.Context, dto *dtos.MarkInvoiceAsCancelledDTO) error
	RetryInvoice(ctx context.Context, invoiceID uint) error
	CheckPendingInvoice(ctx context.Context, invoiceID uint) error
	CancelRetry(ctx context.Context, invoiceID uint) error
	EnableRetry(ctx context.Context, invoiceID uint) error
	GenerateCashReceipt(ctx context.Context, invoiceID uint) error
	DeletePendingInvoice(ctx context.Context, invoiceID uint) error
	GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error)
	ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, int64, error)
	GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error)
	GetInvoiceSyncLogs(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error)

	CreateProvider(ctx context.Context, dto *dtos.CreateProviderDTO) (*entities.InvoicingProvider, error)
	UpdateProvider(ctx context.Context, id uint, dto *dtos.UpdateProviderDTO) error
	GetProvider(ctx context.Context, id uint) (*entities.InvoicingProvider, error)
	ListProviders(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error)
	TestProviderConnection(ctx context.Context, id uint) error

	CreateConfig(ctx context.Context, dto *dtos.CreateConfigDTO) (*entities.InvoicingConfig, error)
	UpdateConfig(ctx context.Context, id uint, dto *dtos.UpdateConfigDTO) (*entities.InvoicingConfig, error)
	GetConfig(ctx context.Context, id uint) (*entities.InvoicingConfig, error)
	ListConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error)
	DeleteConfig(ctx context.Context, id uint) error

	CreateCreditNote(ctx context.Context, dto *dtos.CreateCreditNoteDTO) (*entities.CreditNote, error)
	GetCreditNote(ctx context.Context, id uint) (*entities.CreditNote, error)
	ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error)

	ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error)

	GetSummary(ctx context.Context, businessID uint, period string) (*entities.InvoiceSummary, error)
	GetDetailedStats(ctx context.Context, businessID uint, filters map[string]interface{}) (*entities.DetailedStats, error)
	GetTrends(ctx context.Context, businessID uint, startDate, endDate, granularity, metric string) (*entities.TrendData, error)

	BulkCreateInvoices(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (*dtos.BulkCreateResult, error)

	BulkCreateInvoicesAsync(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (string, error)

	RetryFailedBulk(ctx context.Context, businessID uint) (int, error)
	GetBulkJobStatus(ctx context.Context, jobID string) (*entities.BulkInvoiceJob, error)
	GetBulkJobItems(ctx context.Context, jobID string) ([]*entities.BulkInvoiceJobItem, error)
	ListBulkJobs(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.BulkInvoiceJob, int64, error)

	RequestComparison(ctx context.Context, dto *dtos.CompareRequestDTO) (string, error)

	RequestInventorySync(ctx context.Context, businessID, integrationID uint) (string, error)

	RequestListSiigoWarehouses(ctx context.Context, businessID, integrationID uint) (string, error)

	SyncCancellations(ctx context.Context, dto *dtos.CompareRequestDTO) (string, error)

	GetCompareResult(ctx context.Context, correlationID string) (*dtos.CompareResponseData, error)

	RequestListItems(ctx context.Context, dto *dtos.ListItemsRequestDTO) (string, error)

	GetListItemsResult(ctx context.Context, correlationID string) (*dtos.ItemCompareResponseData, error)

	RequestListBankAccounts(ctx context.Context, dto *dtos.ListBankAccountsRequestDTO) (string, error)

	GetListBankAccountsResult(ctx context.Context, correlationID string) (*dtos.BankAccountsResponseData, error)
}
