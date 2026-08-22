import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/product_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class ProductProvider extends ChangeNotifier {
  ProductProvider({required ApiClient apiClient, ProductUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<Product>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final ProductUseCases? _injectedUseCases;
  late final PagedListController<Product> list;

  int? _businessId;
  String? _error;
  String _nameFilter = '';
  String _skuFilter = '';
  List<ProductIntegration> _integrations = [];

  List<Product> get products => list.loadedItems;
  bool get isLoading => list.isLoading;
  bool get isLoadingMore => list.isLoadingMore;
  bool get hasMore => list.hasMore;
  String? get error => _error ?? list.error;
  List<ProductIntegration> get integrations => _integrations;

  ProductUseCases get _useCases =>
      _injectedUseCases ?? ProductUseCases(ProductApiRepository(_apiClient));

  Future<PaginatedResponse<Product>> _fetchPage(int page, int pageSize) {
    return _useCases.getProducts(GetProductsParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      name: _nameFilter.isNotEmpty ? _nameFilter : null,
      sku: _skuFilter.isNotEmpty ? _skuFilter : null,
    ));
  }

  Future<void> fetchProducts({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
  }

  Future<void> loadMore({int? businessId}) {
    if (businessId != null) _businessId = businessId;
    return list.loadMore();
  }

  Future<void> fetchIntegrations(String productId) async {
    try {
      _integrations = await _useCases.getProductIntegrations(productId);
    } catch (_) {
      _integrations = [];
    }
    notifyListeners();
  }

  Product? productById(String id) {
    for (final product in list.loadedItems) {
      if (product.id == id) return product;
    }
    return null;
  }

  Future<Product?> createProduct(CreateProductDTO data) async {
    try {
      return await _useCases.createProduct(data);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateProduct(String id, UpdateProductDTO data) async {
    try {
      await _useCases.updateProduct(id, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteProduct(String id) async {
    try {
      await _useCases.deleteProduct(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  void setFilters({String? name, String? sku}) {
    _nameFilter = name ?? _nameFilter;
    _skuFilter = sku ?? _skuFilter;
  }

  void resetFilters() {
    _nameFilter = '';
    _skuFilter = '';
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
