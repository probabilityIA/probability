import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/customer_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class CustomerProvider extends ChangeNotifier {
  CustomerProvider({required ApiClient apiClient, CustomerUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<CustomerInfo>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final CustomerUseCases? _injectedUseCases;
  late final PagedListController<CustomerInfo> list;

  int? _businessId;
  String? _error;
  String _searchFilter = '';
  int _unfilteredTotal = 0;

  List<CustomerInfo> get customers => list.loadedItems;
  bool get isLoading => list.isLoading;
  bool get isLoadingMore => list.isLoadingMore;
  bool get hasMore => list.hasMore;
  String? get error => _error ?? list.error;
  int get unfilteredTotal => _unfilteredTotal;
  bool get hasFilters => _searchFilter.isNotEmpty;

  CustomerUseCases get _useCases =>
      _injectedUseCases ?? CustomerUseCases(CustomerApiRepository(_apiClient));

  Future<PaginatedResponse<CustomerInfo>> _fetchPage(
      int page, int pageSize) async {
    final response = await _useCases.getCustomers(GetCustomersParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      search: _searchFilter.isNotEmpty ? _searchFilter : null,
    ));

    if (!hasFilters) _unfilteredTotal = response.pagination.total;
    return response;
  }

  Future<void> fetchCustomers({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
  }

  Future<void> loadMore({int? businessId}) {
    if (businessId != null) _businessId = businessId;
    return list.loadMore();
  }

  CustomerInfo? customerById(int id) {
    for (final customer in list.loadedItems) {
      if (customer.id == id) return customer;
    }
    return null;
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

  void setSearch(String search) {
    _searchFilter = search;
  }

  void resetFilters() {
    _searchFilter = '';
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
