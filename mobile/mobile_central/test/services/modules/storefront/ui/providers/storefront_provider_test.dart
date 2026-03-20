import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/storefront/app/use_cases.dart';
import 'package:mobile_central/services/modules/storefront/domain/entities.dart';
import 'package:mobile_central/services/modules/storefront/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockStorefrontRepository implements IStorefrontRepository {
  PaginatedResponse<StorefrontProduct>? getCatalogResult;
  StorefrontProduct? getProductResult;
  Map<String, dynamic>? createOrderResult;
  PaginatedResponse<StorefrontOrder>? getOrdersResult;
  StorefrontOrder? getOrderResult;
  Map<String, dynamic>? registerResult;
  Exception? errorToThrow;

  final List<String> calls = [];

  @override
  Future<PaginatedResponse<StorefrontProduct>> getCatalog(GetCatalogParams? params) async {
    calls.add('getCatalog');
    if (errorToThrow != null) throw errorToThrow!;
    return getCatalogResult!;
  }

  @override
  Future<StorefrontProduct> getProduct(String id, {int? businessId}) async {
    calls.add('getProduct');
    if (errorToThrow != null) throw errorToThrow!;
    return getProductResult!;
  }

  @override
  Future<Map<String, dynamic>> createOrder(CreateStorefrontOrderDTO data, {int? businessId}) async {
    calls.add('createOrder');
    if (errorToThrow != null) throw errorToThrow!;
    return createOrderResult!;
  }

  @override
  Future<PaginatedResponse<StorefrontOrder>> getOrders(GetOrdersParams? params) async {
    calls.add('getOrders');
    if (errorToThrow != null) throw errorToThrow!;
    return getOrdersResult!;
  }

  @override
  Future<StorefrontOrder> getOrder(String id, {int? businessId}) async {
    calls.add('getOrder');
    if (errorToThrow != null) throw errorToThrow!;
    return getOrderResult!;
  }

  @override
  Future<Map<String, dynamic>> register(RegisterDTO data) async {
    calls.add('register');
    if (errorToThrow != null) throw errorToThrow!;
    return registerResult!;
  }
}

// --- Testable Provider ---

class TestableStorefrontProvider {
  final StorefrontUseCases _useCases;

  List<StorefrontProduct> _products = [];
  List<StorefrontOrder> _orders = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableStorefrontProvider(this._useCases);

  List<StorefrontProduct> get products => _products;
  List<StorefrontOrder> get orders => _orders;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchCatalog({GetCatalogParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.getCatalog(params);
      _products = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchOrders({GetOrdersParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.getOrders(params);
      _orders = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Map<String, dynamic>?> createOrder(CreateStorefrontOrderDTO data, {int? businessId}) async {
    try {
      return await _useCases.createOrder(data, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Map<String, dynamic>?> register(RegisterDTO data) async {
    try {
      return await _useCases.register(data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }
}

// --- Helpers ---

StorefrontProduct _makeProduct({String id = '1', String name = 'TestProduct'}) {
  return StorefrontProduct(
    id: id,
    name: name,
    description: 'desc',
    shortDescription: 'short',
    price: 10000.0,
    currency: 'COP',
    imageUrl: 'https://example.com/img.png',
    sku: 'SKU-1',
    stockQuantity: 10,
    category: 'Cat',
    brand: 'Brand',
    isFeatured: false,
    createdAt: '2026-01-01',
  );
}

StorefrontOrder _makeOrder({String id = '1'}) {
  return StorefrontOrder(
    id: id,
    orderNumber: 'ON-001',
    status: 'pending',
    totalAmount: 50000.0,
    currency: 'COP',
    createdAt: '2026-01-01',
    items: [],
  );
}

Pagination _makePagination({int total = 1, int currentPage = 1, int lastPage = 1}) {
  return Pagination(
    currentPage: currentPage,
    perPage: 10,
    total: total,
    lastPage: lastPage,
    hasNext: currentPage < lastPage,
    hasPrev: currentPage > 1,
  );
}

// --- Tests ---

void main() {
  late MockStorefrontRepository mockRepo;
  late StorefrontUseCases useCases;
  late TestableStorefrontProvider provider;

  setUp(() {
    mockRepo = MockStorefrontRepository();
    useCases = StorefrontUseCases(mockRepo);
    provider = TestableStorefrontProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty products list', () {
      expect(provider.products, isEmpty);
    });

    test('starts with empty orders list', () {
      expect(provider.orders, isEmpty);
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
  });

  group('fetchCatalog', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getCatalogResult = PaginatedResponse<StorefrontProduct>(
        data: [_makeProduct()],
        pagination: _makePagination(),
      );

      await provider.fetchCatalog();

      expect(provider.notifications.length, 2);
    });

    test('populates products and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getCatalogResult = PaginatedResponse<StorefrontProduct>(
        data: [_makeProduct(id: '1'), _makeProduct(id: '2', name: 'Second')],
        pagination: pagination,
      );

      await provider.fetchCatalog();

      expect(provider.products.length, 2);
      expect(provider.products[0].id, '1');
      expect(provider.products[1].name, 'Second');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchCatalog();

      expect(provider.error, contains('Server error'));
      expect(provider.products, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchCatalog();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getCatalogResult = PaginatedResponse<StorefrontProduct>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchCatalog();

      expect(provider.error, isNull);
    });
  });

  group('fetchOrders', () {
    test('populates orders and pagination on success', () async {
      mockRepo.getOrdersResult = PaginatedResponse<StorefrontOrder>(
        data: [_makeOrder(id: 'ord-1')],
        pagination: _makePagination(),
      );

      await provider.fetchOrders();

      expect(provider.orders.length, 1);
      expect(provider.orders[0].id, 'ord-1');
      expect(provider.isLoading, false);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Orders fetch failed');

      await provider.fetchOrders();

      expect(provider.error, contains('Orders fetch failed'));
      expect(provider.orders, isEmpty);
    });
  });

  group('createOrder', () {
    test('returns result on success', () async {
      final dto = CreateStorefrontOrderDTO(
        items: [CreateStorefrontOrderItemDTO(productId: 'p1', quantity: 1)],
      );
      mockRepo.createOrderResult = {'id': 'new-ord', 'status': 'created'};

      final result = await provider.createOrder(dto);

      expect(result, isNotNull);
      expect(result!['id'], 'new-ord');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Order creation failed');
      final dto = CreateStorefrontOrderDTO(items: []);

      final result = await provider.createOrder(dto);

      expect(result, isNull);
      expect(provider.error, contains('Order creation failed'));
    });
  });

  group('register', () {
    test('returns result on success', () async {
      final dto = RegisterDTO(
        name: 'Test',
        email: 'test@test.com',
        password: 'pass',
        businessCode: 'BC',
      );
      mockRepo.registerResult = {'token': 'abc'};

      final result = await provider.register(dto);

      expect(result, isNotNull);
      expect(result!['token'], 'abc');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Registration failed');
      final dto = RegisterDTO(
        name: 'T',
        email: 'e',
        password: 'p',
        businessCode: 'b',
      );

      final result = await provider.register(dto);

      expect(result, isNull);
      expect(provider.error, contains('Registration failed'));
    });
  });
}
