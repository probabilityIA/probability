import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/customer_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class CustomerProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<CustomerInfo> _customers = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _searchFilter = '';

  CustomerProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<CustomerInfo> get customers => _customers;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;

  CustomerUseCases get _useCases => CustomerUseCases(CustomerApiRepository(_apiClient));

  Future<void> fetchCustomers({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetCustomersParams(
        page: _page,
        pageSize: _pageSize,
        businessId: businessId,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
      );
      final response = await _useCases.getCustomers(params);
      _customers = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<CustomerInfo?> createCustomer(CreateCustomerDTO data) async {
    try {
      return await _useCases.createCustomer(data);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateCustomer(int id, UpdateCustomerDTO data) async {
    try {
      await _useCases.updateCustomer(id, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteCustomer(int id) async {
    try {
      await _useCases.deleteCustomer(id);
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

  void setSearch(String search) {
    _searchFilter = search;
    _page = 1;
  }

  void resetFilters() {
    _searchFilter = '';
    _page = 1;
  }
}
