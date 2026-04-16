import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IInvoicingRepository {
  // Invoices
  Future<PaginatedResponse<Invoice>> getInvoices(InvoiceFilters filters);
  Future<Invoice> getInvoiceById(int id);
  Future<Invoice> createInvoice(CreateInvoiceDTO data);
  Future<Invoice> cancelInvoice(int id);
  Future<Invoice> retryInvoice(int id);
  Future<CreditNote> createCreditNote(CreateCreditNoteDTO data);

  // Configs
  Future<PaginatedResponse<InvoicingConfig>> getConfigs(ConfigFilters filters);
  Future<InvoicingConfig> getConfigById(int id);
  Future<InvoicingConfig> createConfig(CreateConfigDTO data);
  Future<InvoicingConfig> updateConfig(int id, UpdateConfigDTO data);
  Future<void> deleteConfig(int id);

  // Bulk
  Future<void> bulkCreateInvoices(BulkCreateInvoicesDTO data);

  // Sync logs
  Future<List<SyncLog>> getSyncLogs(int invoiceId);
}
