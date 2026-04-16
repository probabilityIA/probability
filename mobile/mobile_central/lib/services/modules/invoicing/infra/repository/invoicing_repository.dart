import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class InvoicingApiRepository implements IInvoicingRepository {
  final ApiClient _client;

  InvoicingApiRepository(this._client);

  @override
  Future<PaginatedResponse<Invoice>> getInvoices(InvoiceFilters filters) async {
    final response = await _client.get('/invoices', queryParameters: filters.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => Invoice.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<Invoice> getInvoiceById(int id) async {
    final response = await _client.get('/invoices/$id');
    return Invoice.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Invoice> createInvoice(CreateInvoiceDTO data) async {
    final response = await _client.post('/invoices', data: data.toJson());
    return Invoice.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Invoice> cancelInvoice(int id) async {
    final response = await _client.post('/invoices/$id/cancel');
    return Invoice.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Invoice> retryInvoice(int id) async {
    final response = await _client.post('/invoices/$id/retry');
    return Invoice.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<CreditNote> createCreditNote(CreateCreditNoteDTO data) async {
    final response = await _client.post('/invoices/credit-notes', data: data.toJson());
    return CreditNote.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<PaginatedResponse<InvoicingConfig>> getConfigs(ConfigFilters filters) async {
    final response = await _client.get('/invoicing/configs', queryParameters: filters.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => InvoicingConfig.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<InvoicingConfig> getConfigById(int id) async {
    final response = await _client.get('/invoicing/configs/$id');
    return InvoicingConfig.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<InvoicingConfig> createConfig(CreateConfigDTO data) async {
    final response = await _client.post('/invoicing/configs', data: data.toJson());
    return InvoicingConfig.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<InvoicingConfig> updateConfig(int id, UpdateConfigDTO data) async {
    final response = await _client.put('/invoicing/configs/$id', data: data.toJson());
    return InvoicingConfig.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteConfig(int id) async {
    await _client.delete('/invoicing/configs/$id');
  }

  @override
  Future<void> bulkCreateInvoices(BulkCreateInvoicesDTO data) async {
    await _client.post('/invoices/bulk', data: data.toJson());
  }

  @override
  Future<List<SyncLog>> getSyncLogs(int invoiceId) async {
    final response = await _client.get('/invoices/$invoiceId/sync-logs');
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data.map((e) => SyncLog.fromJson(e)).toList();
  }
}
