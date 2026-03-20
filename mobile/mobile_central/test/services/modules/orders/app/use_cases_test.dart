import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/orders/app/use_cases.dart';
import 'package:mobile_central/services/modules/orders/domain/entities.dart';
import 'package:mobile_central/services/modules/orders/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IOrderRepository
// ---------------------------------------------------------------------------
class MockOrderRepository implements IOrderRepository {
  // Captured arguments
  GetOrdersParams? lastGetOrdersParams;
  String? lastGetOrderByIdArg;
  CreateOrderDTO? lastCreateOrderDTO;
  String? lastUpdateOrderId;
  UpdateOrderDTO? lastUpdateOrderDTO;
  String? lastDeleteOrderId;
  String? lastGetOrderRawId;

  // Configurable return values / errors
  PaginatedResponse<Order>? getOrdersResult;
  Order? getOrderByIdResult;
  Order? createOrderResult;
  Order? updateOrderResult;
  Map<String, dynamic>? getOrderRawResult;
  Exception? errorToThrow;

  @override
  Future<PaginatedResponse<Order>> getOrders(GetOrdersParams? params) async {
    lastGetOrdersParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getOrdersResult ??
        PaginatedResponse(
          data: [],
          pagination: _emptyPagination(),
        );
  }

  @override
  Future<Order> getOrderById(String id) async {
    lastGetOrderByIdArg = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getOrderByIdResult ?? _defaultOrder();
  }

  @override
  Future<Order> createOrder(CreateOrderDTO data) async {
    lastCreateOrderDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createOrderResult ?? _defaultOrder();
  }

  @override
  Future<Order> updateOrder(String id, UpdateOrderDTO data) async {
    lastUpdateOrderId = id;
    lastUpdateOrderDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateOrderResult ?? _defaultOrder();
  }

  @override
  Future<void> deleteOrder(String id) async {
    lastDeleteOrderId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<Map<String, dynamic>> getOrderRaw(String id) async {
    lastGetOrderRawId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getOrderRawResult ?? {};
  }

  // Helpers
  static Order _defaultOrder() => Order(
        id: '1',
        createdAt: '2026-01-01T00:00:00Z',
        updatedAt: '2026-01-01T00:00:00Z',
        integrationId: 1,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        orderNumber: 'ORD-001',
        internalNumber: 'INT-001',
        subtotal: 100.0,
        tax: 19.0,
        discount: 0.0,
        shippingCost: 10.0,
        totalAmount: 129.0,
        currency: 'COP',
        customerName: 'Test Customer',
        customerEmail: 'test@test.com',
        customerPhone: '+573001234567',
        customerDni: '123456789',
        shippingStreet: 'Calle 10',
        shippingCity: 'Bogota',
        shippingState: 'Cundinamarca',
        shippingCountry: 'CO',
        shippingPostalCode: '110111',
        paymentMethodId: 1,
        isPaid: true,
        warehouseName: 'Main',
        driverName: 'Carlos',
        isLastMile: false,
        orderTypeName: 'Standard',
        status: 'confirmed',
        originalStatus: 'paid',
        userName: 'admin',
        invoiceable: true,
        occurredAt: '2026-01-01T00:00:00Z',
        importedAt: '2026-01-01T01:00:00Z',
      );

  static Pagination _emptyPagination() => Pagination(
        currentPage: 1,
        perPage: 10,
        total: 0,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
}

void main() {
  late MockOrderRepository mockRepo;
  late OrderUseCases useCases;

  setUp(() {
    mockRepo = MockOrderRepository();
    useCases = OrderUseCases(mockRepo);
  });

  group('getOrders', () {
    test('delegates to repository with params', () async {
      final params = GetOrdersParams(page: 2, pageSize: 15, status: 'pending');

      await useCases.getOrders(params);

      expect(mockRepo.lastGetOrdersParams, same(params));
    });

    test('delegates to repository with null params', () async {
      await useCases.getOrders(null);

      expect(mockRepo.lastGetOrdersParams, isNull);
    });

    test('returns the response from repository', () async {
      final expectedOrders = [MockOrderRepository._defaultOrder()];
      final expectedPagination = Pagination(
        currentPage: 1,
        perPage: 10,
        total: 1,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
      mockRepo.getOrdersResult = PaginatedResponse(
        data: expectedOrders,
        pagination: expectedPagination,
      );

      final result = await useCases.getOrders(null);

      expect(result.data, hasLength(1));
      expect(result.data[0].id, '1');
      expect(result.pagination.total, 1);
    });
  });

  group('getOrderById', () {
    test('delegates to repository with correct id', () async {
      await useCases.getOrderById('42');

      expect(mockRepo.lastGetOrderByIdArg, '42');
    });

    test('returns order from repository', () async {
      final result = await useCases.getOrderById('1');

      expect(result.id, '1');
      expect(result.orderNumber, 'ORD-001');
    });
  });

  group('createOrder', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      await useCases.createOrder(dto);

      expect(mockRepo.lastCreateOrderDTO, same(dto));
    });

    test('returns created order from repository', () async {
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      final order = await useCases.createOrder(dto);

      expect(order.id, '1');
    });
  });

  group('updateOrder', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = UpdateOrderDTO(status: 'shipped');

      await useCases.updateOrder('10', dto);

      expect(mockRepo.lastUpdateOrderId, '10');
      expect(mockRepo.lastUpdateOrderDTO, same(dto));
    });

    test('returns updated order from repository', () async {
      final dto = UpdateOrderDTO(status: 'shipped');

      final order = await useCases.updateOrder('10', dto);

      expect(order.id, '1');
    });
  });

  group('deleteOrder', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteOrder('77');

      expect(mockRepo.lastDeleteOrderId, '77');
    });
  });

  group('getOrderRaw', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getOrderRawResult = {'id': '123', 'raw': true};

      await useCases.getOrderRaw('123');

      expect(mockRepo.lastGetOrderRawId, '123');
    });

    test('returns raw data from repository', () async {
      mockRepo.getOrderRawResult = {'id': '123', 'raw': true};

      final result = await useCases.getOrderRaw('123');

      expect(result['id'], '123');
      expect(result['raw'], true);
    });
  });

  group('error propagation', () {
    test('getOrders propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.getOrders(null),
        throwsA(isA<Exception>()),
      );
    });

    test('getOrderById propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('not found');

      expect(
        () => useCases.getOrderById('1'),
        throwsA(isA<Exception>()),
      );
    });

    test('createOrder propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('validation error');

      expect(
        () => useCases.createOrder(
          CreateOrderDTO(
            integrationId: 10,
            integrationType: 'shopify',
            platform: 'shopify',
            externalId: 'ext-1',
            subtotal: 100.0,
            totalAmount: 119.0,
            paymentMethodId: 1,
          ),
        ),
        throwsA(isA<Exception>()),
      );
    });

    test('updateOrder propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('server error');

      expect(
        () => useCases.updateOrder('1', UpdateOrderDTO(status: 'X')),
        throwsA(isA<Exception>()),
      );
    });

    test('deleteOrder propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('forbidden');

      expect(
        () => useCases.deleteOrder('1'),
        throwsA(isA<Exception>()),
      );
    });

    test('getOrderRaw propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('raw error');

      expect(
        () => useCases.getOrderRaw('1'),
        throwsA(isA<Exception>()),
      );
    });
  });
}
