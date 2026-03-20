import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/invoicing_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class InvoicingProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  final InvoicingUseCases? _injectedUseCases;

  List<Invoice> _invoices = [];
  List<InvoicingConfig> _configs = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String? _statusFilter;
  int? _businessIdFilter;

  InvoicingProvider({required ApiClient apiClient, InvoicingUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases;

  List<Invoice> get invoices => _invoices;
  List<InvoicingConfig> get configs => _configs;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;

  InvoicingUseCases get _useCases => _injectedUseCases ?? InvoicingUseCases(InvoicingApiRepository(_apiClient));

  Future<void> fetchInvoices({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final filters = InvoiceFilters(
        page: _page,
        pageSize: _pageSize,
        businessId: businessId ?? _businessIdFilter,
        status: _statusFilter,
      );
      final response = await _useCases.getInvoices(filters);
      _invoices = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchConfigs({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final filters = ConfigFilters(businessId: businessId);
      final response = await _useCases.getConfigs(filters);
      _configs = response.data;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
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

  void setPage(int page) {
    _page = page;
  }

  void setFilters({String? status, int? businessId}) {
    _statusFilter = status;
    _businessIdFilter = businessId;
    _page = 1;
  }

  void resetFilters() {
    _statusFilter = null;
    _businessIdFilter = null;
    _page = 1;
  }
}
