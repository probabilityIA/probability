import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/order_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class OrderProvider extends ChangeNotifier {
  OrderProvider({required ApiClient apiClient, OrderUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<Order>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final OrderUseCases? _injectedUseCases;
  late final PagedListController<Order> list;

  int? _businessId;
  String? _error;
  String _searchField = 'order_number';
  String _searchTerm = '';
  String _statusFilter = '';
  String _platformFilter = '';
  bool? _isPaidFilter;
  bool? _isCodFilter;

  List<OrderStatusOption> _statusOptions = const [];
  bool _loadingStatuses = false;

  int _unfilteredTotal = 0;

  List<OrderStatusOption> get statusOptions => _statusOptions;
  String get searchField => _searchField;
  int get unfilteredTotal => _unfilteredTotal;

  bool get hasFilters =>
      _searchTerm.isNotEmpty ||
      _statusFilter.isNotEmpty ||
      _platformFilter.isNotEmpty ||
      _isPaidFilter != null ||
      _isCodFilter != null;

  List<Order> get orders => list.loadedItems;
  bool get isLoading => list.isLoading;
  bool get isLoadingMore => list.isLoadingMore;
  bool get hasMore => list.hasMore;
  String? get error => _error ?? list.error;

  OrderUseCases get _useCases =>
      _injectedUseCases ?? OrderUseCases(OrderApiRepository(_apiClient));

  String? _termFor(String field) =>
      _searchField == field && _searchTerm.isNotEmpty ? _searchTerm : null;

  Future<PaginatedResponse<Order>> _fetchPage(int page, int pageSize) async {
    final response = await _useCases.getOrders(GetOrdersParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      orderNumber: _termFor('order_number'),
      internalNumber: _termFor('internal_number'),
      customerEmail: _termFor('customer_email'),
      customerPhone: _termFor('customer_phone'),
      status: _statusFilter.isNotEmpty ? _statusFilter : null,
      platform: _platformFilter.isNotEmpty ? _platformFilter : null,
      isPaid: _isPaidFilter,
      isCod: _isCodFilter,
    ));

    if (!hasFilters) _unfilteredTotal = response.pagination.total;
    return response;
  }

  Future<void> loadStatusOptions({int? businessId}) async {
    if (_loadingStatuses || _statusOptions.isNotEmpty) return;
    _loadingStatuses = true;
    try {
      _statusOptions =
          await _useCases.getOrderStatuses(businessId: businessId ?? _businessId);
      notifyListeners();
    } catch (_) {
      _statusOptions = const [];
    }
    _loadingStatuses = false;
  }

  void setSearch({String? field, String? term}) {
    if (field != null) _searchField = field;
    if (term != null) _searchTerm = term.trim();
  }

  void applyFilters({
    String? status,
    String? platform,
    bool? isPaid,
    bool? isCod,
  }) {
    _statusFilter = status ?? '';
    _platformFilter = platform ?? '';
    _isPaidFilter = isPaid;
    _isCodFilter = isCod;
  }

  Future<void> fetchOrders({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
  }

  Future<void> loadMore({int? businessId}) {
    if (businessId != null) _businessId = businessId;
    return list.loadMore();
  }

  Future<List<Order>> ordersByCustomerEmail(String email, {int? businessId}) async {
    try {
      final response = await _useCases.getOrders(GetOrdersParams(
        page: 1,
        pageSize: 20,
        businessId: businessId,
        customerEmail: email,
      ));
      return response.data;
    } catch (_) {
      return const [];
    }
  }

  Future<Order?> getOrderById(String id) async {
    try {
      return await _useCases.getOrderById(id);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<Order?> createOrder(CreateOrderDTO data) async {
    try {
      final order = await _useCases.createOrder(data);
      return order;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateOrder(String id, UpdateOrderDTO data) async {
    try {
      await _useCases.updateOrder(id, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteOrder(String id) async {
    try {
      await _useCases.deleteOrder(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  void resetFilters() {
    _searchField = 'order_number';
    _searchTerm = '';
    _statusFilter = '';
    _platformFilter = '';
    _isPaidFilter = null;
    _isCodFilter = null;
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
