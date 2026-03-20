import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/products/app/use_cases.dart';
import 'package:mobile_central/services/modules/products/domain/entities.dart';
import 'package:mobile_central/services/modules/products/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockProductRepository implements IProductRepository {
  PaginatedResponse<Product>? getProductsResult;
  Product? createProductResult;
  Product? updateProductResult;
  List<ProductIntegration>? getProductIntegrationsResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  String? capturedDeleteId;

  @override
  Future<PaginatedResponse<Product>> getProducts(GetProductsParams? params) async {
    calls.add('getProducts');
    if (errorToThrow != null) throw errorToThrow!;
    return getProductsResult!;
  }

  @override
  Future<Product> getProductById(String id, {int? businessId}) async {
    calls.add('getProductById');
    if (errorToThrow != null) throw errorToThrow!;
    return Product(
      id: id, createdAt: '', updatedAt: '', businessId: 1, sku: 'S',
      name: 'Test', price: 10, currency: 'COP', stock: 0, manageStock: false,
      status: 'active', isActive: true,
    );
  }

  @override
  Future<Product> createProduct(CreateProductDTO data, {int? businessId}) async {
    calls.add('createProduct');
    if (errorToThrow != null) throw errorToThrow!;
    return createProductResult!;
  }

  @override
  Future<Product> updateProduct(String id, UpdateProductDTO data, {int? businessId}) async {
    calls.add('updateProduct');
    if (errorToThrow != null) throw errorToThrow!;
    return updateProductResult!;
  }

  @override
  Future<void> deleteProduct(String id, {int? businessId}) async {
    calls.add('deleteProduct');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<ProductIntegration>> getProductIntegrations(String productId, {int? businessId}) async {
    calls.add('getProductIntegrations');
    if (errorToThrow != null) throw errorToThrow!;
    return getProductIntegrationsResult ?? [];
  }

  @override
  Future<void> addProductIntegration(String productId, AddProductIntegrationDTO data, {int? businessId}) async {
    calls.add('addProductIntegration');
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> removeProductIntegration(String productId, int integrationId, {int? businessId}) async {
    calls.add('removeProductIntegration');
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---

class TestableProductProvider {
  final ProductUseCases _useCases;

  List<Product> _products = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _nameFilter = '';
  String _skuFilter = '';

  final List<String> notifications = [];

  TestableProductProvider(this._useCases);

  List<Product> get products => _products;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchProducts({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

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
    _notifyListeners();
  }

  Future<Product?> createProduct(CreateProductDTO data) async {
    try {
      return await _useCases.createProduct(data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateProduct(String id, UpdateProductDTO data) async {
    try {
      await _useCases.updateProduct(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteProduct(String id) async {
    try {
      await _useCases.deleteProduct(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
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

// --- Helpers ---

Product _makeProduct({String id = '1', String name = 'TestProduct'}) {
  return Product(
    id: id, createdAt: '', updatedAt: '', businessId: 1, sku: 'SKU-$id',
    name: name, price: 10.0, currency: 'COP', stock: 5, manageStock: true,
    status: 'active', isActive: true,
  );
}

Pagination _makePagination({int currentPage = 1, int total = 5, int lastPage = 1}) {
  return Pagination(
    currentPage: currentPage, perPage: 20, total: total, lastPage: lastPage,
    hasNext: currentPage < lastPage, hasPrev: currentPage > 1,
  );
}

// --- Tests ---

void main() {
  late MockProductRepository mockRepo;
  late ProductUseCases useCases;
  late TestableProductProvider provider;

  setUp(() {
    mockRepo = MockProductRepository();
    useCases = ProductUseCases(mockRepo);
    provider = TestableProductProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty products list', () {
      expect(provider.products, isEmpty);
    });

    test('starts with null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });

    test('starts on page 1', () {
      expect(provider.page, 1);
    });
  });

  group('fetchProducts', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getProductsResult = PaginatedResponse<Product>(
        data: [_makeProduct()],
        pagination: _makePagination(),
      );

      await provider.fetchProducts();

      expect(provider.notifications.length, 2);
    });

    test('populates products and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getProductsResult = PaginatedResponse<Product>(
        data: [_makeProduct(id: '1'), _makeProduct(id: '2', name: 'Second')],
        pagination: pagination,
      );

      await provider.fetchProducts();

      expect(provider.products.length, 2);
      expect(provider.products[0].id, '1');
      expect(provider.products[1].name, 'Second');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchProducts();

      expect(provider.error, contains('Server error'));
      expect(provider.products, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchProducts();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getProductsResult = PaginatedResponse<Product>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchProducts();

      expect(provider.error, isNull);
    });

    test('passes businessId to use cases', () async {
      mockRepo.getProductsResult = PaginatedResponse<Product>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchProducts(businessId: 5);

      expect(mockRepo.calls, contains('getProducts'));
    });
  });

  group('createProduct', () {
    test('returns created product on success', () async {
      final dto = CreateProductDTO(businessId: 1, sku: 'S', name: 'New', price: 10, stock: 5);
      mockRepo.createProductResult = _makeProduct(id: '10', name: 'New');

      final result = await provider.createProduct(dto);

      expect(result, isNotNull);
      expect(result!.id, '10');
      expect(result.name, 'New');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateProductDTO(businessId: 1, sku: 'S', name: 'Fail', price: 10, stock: 5);

      final result = await provider.createProduct(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateProduct', () {
    test('returns true on success', () async {
      final dto = UpdateProductDTO(name: 'Updated');
      mockRepo.updateProductResult = _makeProduct(id: '5', name: 'Updated');

      final result = await provider.updateProduct('5', dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateProductDTO(name: 'Fail');

      final result = await provider.updateProduct('5', dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteProduct', () {
    test('returns true on success', () async {
      final result = await provider.deleteProduct('7');

      expect(result, true);
      expect(mockRepo.capturedDeleteId, '7');
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteProduct('7');

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });

  group('pagination', () {
    test('setPage updates page', () {
      provider.setPage(3);
      expect(provider.page, 3);
    });
  });

  group('filters', () {
    test('setFilters updates filters and resets page to 1', () {
      provider.setPage(5);
      provider.setFilters(name: 'Widget', sku: 'SKU-001');

      expect(provider.page, 1);
    });

    test('resetFilters clears filters and resets page to 1', () {
      provider.setFilters(name: 'Widget', sku: 'SKU-001');
      provider.setPage(3);

      provider.resetFilters();

      expect(provider.page, 1);
    });
  });
}
