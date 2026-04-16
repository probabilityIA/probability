import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/storefront_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class StorefrontProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<StorefrontProduct> _catalogProducts = [];
  List<StorefrontOrder> _orders = [];
  StorefrontProduct? _selectedProduct;
  StorefrontOrder? _selectedOrder;
  Pagination? _catalogPagination;
  Pagination? _ordersPagination;
  bool _isLoading = false;
  String? _error;
  int _catalogPage = 1;
  int _ordersPage = 1;
  final int _pageSize = 20;
  String _searchFilter = '';
  String _categoryFilter = '';

  StorefrontProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<StorefrontProduct> get catalogProducts => _catalogProducts;
  List<StorefrontOrder> get orders => _orders;
  StorefrontProduct? get selectedProduct => _selectedProduct;
  StorefrontOrder? get selectedOrder => _selectedOrder;
  Pagination? get catalogPagination => _catalogPagination;
  Pagination? get ordersPagination => _ordersPagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get catalogPage => _catalogPage;
  int get ordersPage => _ordersPage;

  StorefrontUseCases get _useCases =>
      StorefrontUseCases(StorefrontApiRepository(_apiClient));

  Future<void> fetchCatalog({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetCatalogParams(
        page: _catalogPage,
        pageSize: _pageSize,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
        category: _categoryFilter.isNotEmpty ? _categoryFilter : null,
        businessId: businessId,
      );
      final response = await _useCases.getCatalog(params);
      _catalogProducts = response.data;
      _catalogPagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchProduct(String id, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _selectedProduct = await _useCases.getProduct(id, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Map<String, dynamic>?> createOrder(
    CreateStorefrontOrderDTO data, {
    int? businessId,
  }) async {
    try {
      return await _useCases.createOrder(data, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<void> fetchOrders({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetOrdersParams(
        page: _ordersPage,
        pageSize: _pageSize,
        businessId: businessId,
      );
      final response = await _useCases.getOrders(params);
      _orders = response.data;
      _ordersPagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchOrder(String id, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _selectedOrder = await _useCases.getOrder(id, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Map<String, dynamic>?> register(RegisterDTO data) async {
    try {
      return await _useCases.register(data);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  void setCatalogPage(int page) {
    _catalogPage = page;
  }

  void setOrdersPage(int page) {
    _ordersPage = page;
  }

  void setFilters({String? search, String? category}) {
    _searchFilter = search ?? _searchFilter;
    _categoryFilter = category ?? _categoryFilter;
    _catalogPage = 1;
  }

  void resetFilters() {
    _searchFilter = '';
    _categoryFilter = '';
    _catalogPage = 1;
  }
}
