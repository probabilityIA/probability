import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/product_repository.dart';

class ProductProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<Product> _products = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _nameFilter = '';
  String _skuFilter = '';

  ProductProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<Product> get products => _products;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;

  ProductUseCases get _useCases => ProductUseCases(ProductApiRepository(_apiClient));

  Future<void> fetchProducts({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetProductsParams(
        page: _page,
        pageSize: _pageSize,
        businessId: businessId,
        name: _nameFilter.isNotEmpty ? _nameFilter : null,
        sku: _skuFilter.isNotEmpty ? _skuFilter : null,
      );
      final response = await _useCases.getProducts(params);
      _products = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Product?> createProduct(CreateProductDTO data) async {
    try {
      return await _useCases.createProduct(data);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateProduct(String id, UpdateProductDTO data) async {
    try {
      await _useCases.updateProduct(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteProduct(String id) async {
    try {
      await _useCases.deleteProduct(id);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  void setPage(int page) {
    _page = page;
  }

  void setFilters({String? name, String? sku}) {
    _nameFilter = name ?? _nameFilter;
    _skuFilter = sku ?? _skuFilter;
    _page = 1;
  }

  void resetFilters() {
    _nameFilter = '';
    _skuFilter = '';
    _page = 1;
  }
}
