import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class InvoicingUseCases {
  final IInvoicingRepository _repository;

  InvoicingUseCases(this._repository);

  Future<PaginatedResponse<Invoice>> getInvoices(InvoiceFilters filters) {
    return _repository.getInvoices(filters);
  }

  Future<Invoice> getInvoiceById(int id) {
    return _repository.getInvoiceById(id);
  }

  Future<Invoice> createInvoice(CreateInvoiceDTO data) {
    return _repository.createInvoice(data);
  }

  Future<Invoice> cancelInvoice(int id) {
    return _repository.cancelInvoice(id);
  }

  Future<Invoice> retryInvoice(int id) {
    return _repository.retryInvoice(id);
  }

  Future<CreditNote> createCreditNote(CreateCreditNoteDTO data) {
    return _repository.createCreditNote(data);
  }

  Future<PaginatedResponse<InvoicingConfig>> getConfigs(ConfigFilters filters) {
    return _repository.getConfigs(filters);
  }

  Future<InvoicingConfig> getConfigById(int id) {
    return _repository.getConfigById(id);
  }

  Future<InvoicingConfig> createConfig(CreateConfigDTO data) {
    return _repository.createConfig(data);
  }

  Future<InvoicingConfig> updateConfig(int id, UpdateConfigDTO data) {
    return _repository.updateConfig(id, data);
  }

  Future<void> deleteConfig(int id) {
    return _repository.deleteConfig(id);
  }

  Future<void> bulkCreateInvoices(BulkCreateInvoicesDTO data) {
    return _repository.bulkCreateInvoices(data);
  }

  Future<List<SyncLog>> getSyncLogs(int invoiceId) {
    return _repository.getSyncLogs(invoiceId);
  }
}
