import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/order_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class OrderProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<Order> _orders = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _orderNumberFilter = '';
  String _statusFilter = '';
  int? _integrationIdFilter;

  OrderProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<Order> get orders => _orders;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;

  OrderUseCases get _useCases =>
      OrderUseCases(OrderApiRepository(_apiClient));

  Future<void> fetchOrders({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetOrdersParams(
        page: _page,
        pageSize: _pageSize,
        businessId: businessId,
        orderNumber: _orderNumberFilter.isNotEmpty ? _orderNumberFilter : null,
        status: _statusFilter.isNotEmpty ? _statusFilter : null,
        integrationId: _integrationIdFilter,
      );
      final response = await _useCases.getOrders(params);
      _orders = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
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

  void setPage(int page) {
    _page = page;
  }

  void setFilters({
    String? orderNumber,
    String? status,
    int? integrationId,
  }) {
    _orderNumberFilter = orderNumber ?? _orderNumberFilter;
    _statusFilter = status ?? _statusFilter;
    _integrationIdFilter = integrationId ?? _integrationIdFilter;
    _page = 1;
  }

  void resetFilters() {
    _orderNumberFilter = '';
    _statusFilter = '';
    _integrationIdFilter = null;
    _page = 1;
  }
}
