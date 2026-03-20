import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/publicsite/app/use_cases.dart';
import 'package:mobile_central/services/modules/publicsite/domain/entities.dart';
import 'package:mobile_central/services/modules/publicsite/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockPublicSiteRepository implements IPublicSiteRepository {
  PublicBusiness? getBusinessPageResult;
  PaginatedResponse<PublicProduct>? getCatalogResult;
  PublicProduct? getProductResult;
  Map<String, dynamic>? submitContactResult;
  Exception? errorToThrow;

  final List<String> calls = [];

  @override
  Future<PublicBusiness> getBusinessPage(String slug) async {
    calls.add('getBusinessPage');
    if (errorToThrow != null) throw errorToThrow!;
    return getBusinessPageResult!;
  }

  @override
  Future<PaginatedResponse<PublicProduct>> getCatalog(String slug, GetPublicCatalogParams? params) async {
    calls.add('getCatalog');
    if (errorToThrow != null) throw errorToThrow!;
    return getCatalogResult!;
  }

  @override
  Future<PublicProduct> getProduct(String slug, String productId) async {
    calls.add('getProduct');
    if (errorToThrow != null) throw errorToThrow!;
    return getProductResult!;
  }

  @override
  Future<Map<String, dynamic>> submitContact(String slug, ContactFormDTO data) async {
    calls.add('submitContact');
    if (errorToThrow != null) throw errorToThrow!;
    return submitContactResult!;
  }
}

// --- Testable Provider ---

class TestablePublicSiteProvider {
  final PublicSiteUseCases _useCases;

  PublicBusiness? _business;
  List<PublicProduct> _catalogProducts = [];
  PublicProduct? _selectedProduct;
  Pagination? _catalogPagination;
  bool _isLoading = false;
  String? _error;
  int _catalogPage = 1;
  final int _pageSize = 20;
  String _searchFilter = '';
  String _categoryFilter = '';

  final List<String> notifications = [];

  TestablePublicSiteProvider(this._useCases);

  PublicBusiness? get business => _business;
  List<PublicProduct> get catalogProducts => _catalogProducts;
  PublicProduct? get selectedProduct => _selectedProduct;
  Pagination? get catalogPagination => _catalogPagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get catalogPage => _catalogPage;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchBusinessPage(String slug) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _business = await _useCases.getBusinessPage(slug);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchCatalog(String slug) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final params = GetPublicCatalogParams(
        page: _catalogPage,
        pageSize: _pageSize,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
        category: _categoryFilter.isNotEmpty ? _categoryFilter : null,
      );
      final response = await _useCases.getCatalog(slug, params);
      _catalogProducts = response.data;
      _catalogPagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchProduct(String slug, String productId) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _selectedProduct = await _useCases.getProduct(slug, productId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Map<String, dynamic>?> submitContact(String slug, ContactFormDTO data) async {
    try {
      return await _useCases.submitContact(slug, data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  void setCatalogPage(int page) {
    _catalogPage = page;
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

// --- Helpers ---

PublicBusiness _makeBusiness({int id = 1, String name = 'TestStore'}) {
  return PublicBusiness(
    id: id, name: name, code: 'teststore', description: 'A store',
    logoUrl: '', primaryColor: '', secondaryColor: '', tertiaryColor: '',
    quaternaryColor: '', navbarImageUrl: '', featuredProducts: [],
  );
}

PublicProduct _makeProduct({String id = '1', String name = 'TestProduct'}) {
  return PublicProduct(
    id: id, name: name, description: '', shortDescription: '', price: 10.0,
    currency: 'COP', imageUrl: '', sku: 'SKU', stockQuantity: 5,
    category: '', brand: '', isFeatured: false, createdAt: '',
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
  late MockPublicSiteRepository mockRepo;
  late PublicSiteUseCases useCases;
  late TestablePublicSiteProvider provider;

  setUp(() {
    mockRepo = MockPublicSiteRepository();
    useCases = PublicSiteUseCases(mockRepo);
    provider = TestablePublicSiteProvider(useCases);
  });

  group('initial state', () {
    test('starts with null business', () {
      expect(provider.business, isNull);
    });

    test('starts with empty catalog products', () {
      expect(provider.catalogProducts, isEmpty);
    });

    test('starts with null selected product', () {
      expect(provider.selectedProduct, isNull);
    });

    test('starts with null pagination', () {
      expect(provider.catalogPagination, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });

    test('starts on catalog page 1', () {
      expect(provider.catalogPage, 1);
    });
  });

  group('fetchBusinessPage', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getBusinessPageResult = _makeBusiness();

      await provider.fetchBusinessPage('mystore');

      expect(provider.notifications.length, 2);
    });

    test('populates business on success', () async {
      mockRepo.getBusinessPageResult = _makeBusiness(id: 1, name: 'My Store');

      await provider.fetchBusinessPage('mystore');

      expect(provider.business, isNotNull);
      expect(provider.business!.name, 'My Store');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Not found');

      await provider.fetchBusinessPage('invalid');

      expect(provider.error, contains('Not found'));
      expect(provider.business, isNull);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchBusinessPage('invalid');
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getBusinessPageResult = _makeBusiness();
      await provider.fetchBusinessPage('valid');

      expect(provider.error, isNull);
    });
  });

  group('fetchCatalog', () {
    test('populates catalog products and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getCatalogResult = PaginatedResponse<PublicProduct>(
        data: [_makeProduct(id: '1'), _makeProduct(id: '2', name: 'Second')],
        pagination: pagination,
      );

      await provider.fetchCatalog('mystore');

      expect(provider.catalogProducts.length, 2);
      expect(provider.catalogProducts[1].name, 'Second');
      expect(provider.catalogPagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Catalog error');

      await provider.fetchCatalog('mystore');

      expect(provider.error, contains('Catalog error'));
      expect(provider.catalogProducts, isEmpty);
    });
  });

  group('fetchProduct', () {
    test('populates selected product on success', () async {
      mockRepo.getProductResult = _makeProduct(id: '42', name: 'Widget');

      await provider.fetchProduct('mystore', '42');

      expect(provider.selectedProduct, isNotNull);
      expect(provider.selectedProduct!.id, '42');
      expect(provider.selectedProduct!.name, 'Widget');
      expect(provider.isLoading, false);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Product not found');

      await provider.fetchProduct('mystore', '999');

      expect(provider.error, contains('Product not found'));
      expect(provider.selectedProduct, isNull);
    });
  });

  group('submitContact', () {
    test('returns result on success', () async {
      final dto = ContactFormDTO(name: 'John', message: 'Hello');
      mockRepo.submitContactResult = {'success': true};

      final result = await provider.submitContact('mystore', dto);

      expect(result, isNotNull);
      expect(result!['success'], true);
    });

    test('returns null and sets error on failure', () async {
      final dto = ContactFormDTO(name: 'John', message: 'Hello');
      mockRepo.errorToThrow = Exception('Submit failed');

      final result = await provider.submitContact('mystore', dto);

      expect(result, isNull);
      expect(provider.error, contains('Submit failed'));
    });
  });

  group('pagination and filters', () {
    test('setCatalogPage updates page', () {
      provider.setCatalogPage(3);
      expect(provider.catalogPage, 3);
    });

    test('setFilters updates filters and resets page to 1', () {
      provider.setCatalogPage(5);
      provider.setFilters(search: 'widget', category: 'electronics');

      expect(provider.catalogPage, 1);
    });

    test('resetFilters clears filters and resets page to 1', () {
      provider.setFilters(search: 'widget', category: 'electronics');
      provider.setCatalogPage(3);

      provider.resetFilters();

      expect(provider.catalogPage, 1);
    });
  });
}
