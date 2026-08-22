import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/invoicing_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class InvoicingProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  final InvoicingUseCases? _injectedUseCases;

  List<InvoicingConfig> _configs = [];
  bool _isLoadingConfigs = false;
  String? _error;
  String? _statusFilter;
  int? _businessIdFilter;
  int? _businessId;

  late final PagedListController<Invoice> list;

  InvoicingProvider({required ApiClient apiClient, InvoicingUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<Invoice>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  List<Invoice> get invoices => list.loadedItems;
  List<InvoicingConfig> get configs => _configs;
  bool get isLoading => list.isLoading || _isLoadingConfigs;
  String? get error => _error ?? list.error;

  InvoicingUseCases get _useCases => _injectedUseCases ?? InvoicingUseCases(InvoicingApiRepository(_apiClient));

  Future<PaginatedResponse<Invoice>> _fetchPage(int page, int pageSize) {
    return _useCases.getInvoices(InvoiceFilters(
      page: page,
      pageSize: pageSize,
      businessId: _businessId ?? _businessIdFilter,
      status: _statusFilter,
      invoiceNumber: _invoiceNumberFilter,
    ));
  }

  Future<void> fetchInvoices({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
  }

  Future<void> fetchConfigs({int? businessId}) async {
    _isLoadingConfigs = true;
    _error = null;
    notifyListeners();

    try {
      final filters = ConfigFilters(businessId: businessId);
      final response = await _useCases.getConfigs(filters);
      _configs = response.data;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoadingConfigs = false;
    notifyListeners();
  }

  Future<Invoice?> createInvoice(CreateInvoiceDTO data) async {
    try {
      return await _useCases.createInvoice(data);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> cancelInvoice(int id) async {
    try {
      await _useCases.cancelInvoice(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> retryInvoice(int id) async {
    try {
      await _useCases.retryInvoice(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> bulkCreateInvoices(BulkCreateInvoicesDTO data) async {
    try {
      await _useCases.bulkCreateInvoices(data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  String? _invoiceNumberFilter;

  void setFilters({String? status, int? businessId, String? invoiceNumber}) {
    _statusFilter = status;
    _businessIdFilter = businessId;
    _invoiceNumberFilter = invoiceNumber;
  }

  Invoice? invoiceById(int id) {
    for (final invoice in list.loadedItems) {
      if (invoice.id == id) return invoice;
    }
    return null;
  }

  void resetFilters() {
    _statusFilter = null;
    _businessIdFilter = null;
    _invoiceNumberFilter = null;
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
