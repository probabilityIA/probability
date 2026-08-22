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
  String _orderNumberFilter = '';
  String _statusFilter = '';
  int? _integrationIdFilter;

  List<Order> get orders => list.loadedItems;
  bool get isLoading => list.isLoading;
  bool get isLoadingMore => list.isLoadingMore;
  bool get hasMore => list.hasMore;
  String? get error => _error ?? list.error;

  OrderUseCases get _useCases =>
      _injectedUseCases ?? OrderUseCases(OrderApiRepository(_apiClient));

  Future<PaginatedResponse<Order>> _fetchPage(int page, int pageSize) {
    return _useCases.getOrders(GetOrdersParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      orderNumber: _orderNumberFilter.isNotEmpty ? _orderNumberFilter : null,
      status: _statusFilter.isNotEmpty ? _statusFilter : null,
      integrationId: _integrationIdFilter,
    ));
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

  void setFilters({
    String? orderNumber,
    String? status,
    int? integrationId,
  }) {
    _orderNumberFilter = orderNumber ?? _orderNumberFilter;
    _statusFilter = status ?? _statusFilter;
    _integrationIdFilter = integrationId ?? _integrationIdFilter;
  }

  void resetFilters() {
    _orderNumberFilter = '';
    _statusFilter = '';
    _integrationIdFilter = null;
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
