import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/storefront/app/use_cases.dart';
import 'package:mobile_central/services/modules/storefront/domain/entities.dart';
import 'package:mobile_central/services/modules/storefront/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockStorefrontRepository implements IStorefrontRepository {
  final List<String> calls = [];

  PaginatedResponse<StorefrontProduct>? getCatalogResult;
  StorefrontProduct? getProductResult;
  Map<String, dynamic>? createOrderResult;
  PaginatedResponse<StorefrontOrder>? getOrdersResult;
  StorefrontOrder? getOrderResult;
  Map<String, dynamic>? registerResult;

  Exception? errorToThrow;

  GetCatalogParams? capturedCatalogParams;
  String? capturedProductId;
  int? capturedProductBusinessId;
  CreateStorefrontOrderDTO? capturedCreateOrderData;
  int? capturedCreateOrderBusinessId;
  GetOrdersParams? capturedOrdersParams;
  String? capturedOrderId;
  int? capturedOrderBusinessId;
  RegisterDTO? capturedRegisterData;

  @override
  Future<PaginatedResponse<StorefrontProduct>> getCatalog(GetCatalogParams? params) async {
    calls.add('getCatalog');
    capturedCatalogParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getCatalogResult!;
  }

  @override
  Future<StorefrontProduct> getProduct(String id, {int? businessId}) async {
    calls.add('getProduct');
    capturedProductId = id;
    capturedProductBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getProductResult!;
  }

  @override
  Future<Map<String, dynamic>> createOrder(CreateStorefrontOrderDTO data, {int? businessId}) async {
    calls.add('createOrder');
    capturedCreateOrderData = data;
    capturedCreateOrderBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createOrderResult!;
  }

  @override
  Future<PaginatedResponse<StorefrontOrder>> getOrders(GetOrdersParams? params) async {
    calls.add('getOrders');
    capturedOrdersParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getOrdersResult!;
  }

  @override
  Future<StorefrontOrder> getOrder(String id, {int? businessId}) async {
    calls.add('getOrder');
    capturedOrderId = id;
    capturedOrderBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getOrderResult!;
  }

  @override
  Future<Map<String, dynamic>> register(RegisterDTO data) async {
    calls.add('register');
    capturedRegisterData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return registerResult!;
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

Pagination _makePagination() {
  return Pagination(
    currentPage: 1,
    perPage: 10,
    total: 1,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

// --- Tests ---

void main() {
  late MockStorefrontRepository mockRepo;
  late StorefrontUseCases useCases;

  setUp(() {
    mockRepo = MockStorefrontRepository();
    useCases = StorefrontUseCases(mockRepo);
  });

  group('getCatalog', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<StorefrontProduct>(
        data: [_makeProduct()],
        pagination: _makePagination(),
      );
      mockRepo.getCatalogResult = expected;
      final params = GetCatalogParams(page: 1, pageSize: 10);

      final result = await useCases.getCatalog(params);

      expect(result.data.length, 1);
      expect(result.data[0].name, 'TestProduct');
      expect(mockRepo.calls, ['getCatalog']);
      expect(mockRepo.capturedCatalogParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getCatalogResult = PaginatedResponse<StorefrontProduct>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getCatalog(null);

      expect(mockRepo.capturedCatalogParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getCatalog(null), throwsException);
    });
  });

  group('getProduct', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getProductResult = _makeProduct(id: 'p42', name: 'Found');

      final result = await useCases.getProduct('p42');

      expect(result.id, 'p42');
      expect(result.name, 'Found');
      expect(mockRepo.capturedProductId, 'p42');
      expect(mockRepo.calls, ['getProduct']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getProductResult = _makeProduct();

      await useCases.getProduct('p1', businessId: 5);

      expect(mockRepo.capturedProductBusinessId, 5);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.getProduct('p1'), throwsException);
    });
  });

  group('createOrder', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateStorefrontOrderDTO(
        items: [CreateStorefrontOrderItemDTO(productId: 'p1', quantity: 2)],
      );
      mockRepo.createOrderResult = {'id': 'ord-1', 'status': 'created'};

      final result = await useCases.createOrder(dto);

      expect(result['id'], 'ord-1');
      expect(mockRepo.capturedCreateOrderData, dto);
      expect(mockRepo.calls, ['createOrder']);
    });

    test('passes businessId to repository', () async {
      final dto = CreateStorefrontOrderDTO(items: []);
      mockRepo.createOrderResult = {};

      await useCases.createOrder(dto, businessId: 7);

      expect(mockRepo.capturedCreateOrderBusinessId, 7);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateStorefrontOrderDTO(items: []);

      expect(() => useCases.createOrder(dto), throwsException);
    });
  });

  group('getOrders', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<StorefrontOrder>(
        data: [_makeOrder()],
        pagination: _makePagination(),
      );
      mockRepo.getOrdersResult = expected;
      final params = GetOrdersParams(page: 1);

      final result = await useCases.getOrders(params);

      expect(result.data.length, 1);
      expect(mockRepo.calls, ['getOrders']);
      expect(mockRepo.capturedOrdersParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getOrdersResult = PaginatedResponse<StorefrontOrder>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getOrders(null);

      expect(mockRepo.capturedOrdersParams, isNull);
    });
  });

  group('getOrder', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getOrderResult = _makeOrder(id: 'ord-99');

      final result = await useCases.getOrder('ord-99');

      expect(result.id, 'ord-99');
      expect(mockRepo.capturedOrderId, 'ord-99');
      expect(mockRepo.calls, ['getOrder']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getOrderResult = _makeOrder();

      await useCases.getOrder('ord-1', businessId: 3);

      expect(mockRepo.capturedOrderBusinessId, 3);
    });
  });

  group('register', () {
    test('delegates to repository with correct data', () async {
      final dto = RegisterDTO(
        name: 'Test',
        email: 'test@test.com',
        password: 'pass',
        businessCode: 'BC1',
      );
      mockRepo.registerResult = {'token': 'jwt-token'};

      final result = await useCases.register(dto);

      expect(result['token'], 'jwt-token');
      expect(mockRepo.capturedRegisterData, dto);
      expect(mockRepo.calls, ['register']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Registration failed');
      final dto = RegisterDTO(
        name: 'T',
        email: 'e',
        password: 'p',
        businessCode: 'b',
      );

      expect(() => useCases.register(dto), throwsException);
    });
  });
}
