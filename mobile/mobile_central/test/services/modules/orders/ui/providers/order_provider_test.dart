import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/orders/app/use_cases.dart';
import 'package:mobile_central/services/modules/orders/domain/entities.dart';
import 'package:mobile_central/services/modules/orders/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IOrderRepository
// ---------------------------------------------------------------------------
class MockOrderRepository implements IOrderRepository {
  PaginatedResponse<Order>? getOrdersResult;
  Order? getOrderByIdResult;
  Order? createOrderResult;
  Order? updateOrderResult;
  Map<String, dynamic>? getOrderRawResult;
  Exception? errorToThrow;

  GetOrdersParams? lastGetOrdersParams;
  String? lastDeleteId;

  @override
  Future<PaginatedResponse<Order>> getOrders(GetOrdersParams? params) async {
    lastGetOrdersParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getOrdersResult ??
        PaginatedResponse(
          data: [],
          pagination: _defaultPagination(),
        );
  }

  @override
  Future<Order> getOrderById(String id) async {
    if (errorToThrow != null) throw errorToThrow!;
    return getOrderByIdResult ?? _defaultOrder();
  }

  @override
  Future<Order> createOrder(CreateOrderDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return createOrderResult ?? _defaultOrder();
  }

  @override
  Future<Order> updateOrder(String id, UpdateOrderDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return updateOrderResult ?? _defaultOrder();
  }

  @override
  Future<void> deleteOrder(String id) async {
    lastDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<Map<String, dynamic>> getOrderRaw(String id) async {
    if (errorToThrow != null) throw errorToThrow!;
    return getOrderRawResult ?? {};
  }

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

  static Pagination _defaultPagination() => Pagination(
        currentPage: 1,
        perPage: 20,
        total: 0,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
}

// ---------------------------------------------------------------------------
// Testable provider that mirrors OrderProvider logic but accepts a repository
// ---------------------------------------------------------------------------
class TestableOrderProvider {
  final IOrderRepository _repository;

  List<Order> _orders = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _orderNumberFilter = '';
  String _statusFilter = '';
  int? _integrationIdFilter;

  int _notifyCount = 0;

  TestableOrderProvider(this._repository);

  List<Order> get orders => _orders;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;
  int get notifyCount => _notifyCount;
  String get orderNumberFilter => _orderNumberFilter;
  String get statusFilter => _statusFilter;
  int? get integrationIdFilter => _integrationIdFilter;

  OrderUseCases get _useCases => OrderUseCases(_repository);

  void _notifyListeners() {
    _notifyCount++;
  }

  Future<void> fetchOrders({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final params = GetOrdersParams(
        page: _page,
        pageSize: _pageSize,
        businessId: businessId,
        orderNumber:
            _orderNumberFilter.isNotEmpty ? _orderNumberFilter : null,
        status: _statusFilter.isNotEmpty ? _statusFilter : null,
        integrationId: _integrationIdFilter,
      );
      final response = await _useCases.getOrders(params);
      _orders = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Order?> getOrderById(String id) async {
    try {
      return await _useCases.getOrderById(id);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Order?> createOrder(CreateOrderDTO data) async {
    try {
      final order = await _useCases.createOrder(data);
      return order;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateOrder(String id, UpdateOrderDTO data) async {
    try {
      await _useCases.updateOrder(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteOrder(String id) async {
    try {
      await _useCases.deleteOrder(id);
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

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockOrderRepository mockRepo;
  late TestableOrderProvider provider;

  setUp(() {
    mockRepo = MockOrderRepository();
    provider = TestableOrderProvider(mockRepo);
  });

  group('initial state', () {
    test('has empty orders list', () {
      expect(provider.orders, isEmpty);
    });

    test('pagination is null', () {
      expect(provider.pagination, isNull);
    });

    test('isLoading is false', () {
      expect(provider.isLoading, false);
    });

    test('error is null', () {
      expect(provider.error, isNull);
    });

    test('page defaults to 1', () {
      expect(provider.page, 1);
    });

    test('pageSize defaults to 20', () {
      expect(provider.pageSize, 20);
    });
  });

  group('fetchOrders', () {
    test('updates orders list and pagination on success', () async {
      final testOrders = [
        MockOrderRepository._defaultOrder(),
      ];
      final testPagination = Pagination(
        currentPage: 1,
        perPage: 20,
        total: 1,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
      mockRepo.getOrdersResult = PaginatedResponse(
        data: testOrders,
        pagination: testPagination,
      );

      await provider.fetchOrders();

      expect(provider.orders, hasLength(1));
      expect(provider.orders[0].id, '1');
      expect(provider.pagination, isNotNull);
      expect(provider.pagination!.total, 1);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('server down');

      await provider.fetchOrders();

      expect(provider.error, contains('server down'));
      expect(provider.orders, isEmpty);
      expect(provider.isLoading, false);
    });

    test('notifies listeners twice (loading start and end)', () async {
      await provider.fetchOrders();

      expect(provider.notifyCount, 2);
    });

    test('passes correct params including page and pageSize', () async {
      await provider.fetchOrders();

      expect(mockRepo.lastGetOrdersParams, isNotNull);
      expect(mockRepo.lastGetOrdersParams!.page, 1);
      expect(mockRepo.lastGetOrdersParams!.pageSize, 20);
    });

    test('passes filters when set', () async {
      provider.setFilters(
        orderNumber: 'ORD-001',
        status: 'pending',
        integrationId: 5,
      );

      await provider.fetchOrders();

      expect(mockRepo.lastGetOrdersParams!.orderNumber, 'ORD-001');
      expect(mockRepo.lastGetOrdersParams!.status, 'pending');
      expect(mockRepo.lastGetOrdersParams!.integrationId, 5);
    });

    test('does not pass empty filters', () async {
      await provider.fetchOrders();

      expect(mockRepo.lastGetOrdersParams!.orderNumber, isNull);
      expect(mockRepo.lastGetOrdersParams!.status, isNull);
      expect(mockRepo.lastGetOrdersParams!.integrationId, isNull);
    });

    test('passes businessId when provided', () async {
      await provider.fetchOrders(businessId: 7);

      expect(mockRepo.lastGetOrdersParams!.businessId, 7);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('fail');
      await provider.fetchOrders();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      await provider.fetchOrders();

      expect(provider.error, isNull);
    });
  });

  group('getOrderById', () {
    test('returns order on success', () async {
      final result = await provider.getOrderById('1');

      expect(result, isNotNull);
      expect(result!.id, '1');
    });

    test('returns null on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('not found');

      final result = await provider.getOrderById('999');

      expect(result, isNull);
      expect(provider.error, contains('not found'));
    });
  });

  group('createOrder', () {
    test('returns Order on success', () async {
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      final result = await provider.createOrder(dto);

      expect(result, isNotNull);
      expect(result!.id, '1');
    });

    test('returns null on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('create failed');
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      final result = await provider.createOrder(dto);

      expect(result, isNull);
      expect(provider.error, contains('create failed'));
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      await provider.createOrder(dto);

      expect(provider.notifyCount, 1);
    });

    test('does not notify listeners on success', () async {
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      await provider.createOrder(dto);

      expect(provider.notifyCount, 0);
    });
  });

  group('updateOrder', () {
    test('returns true on success', () async {
      final result = await provider.updateOrder(
        '1',
        UpdateOrderDTO(status: 'shipped'),
      );

      expect(result, true);
    });

    test('returns false on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('update failed');

      final result = await provider.updateOrder(
        '1',
        UpdateOrderDTO(status: 'shipped'),
      );

      expect(result, false);
      expect(provider.error, contains('update failed'));
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');

      await provider.updateOrder('1', UpdateOrderDTO(status: 'X'));

      expect(provider.notifyCount, 1);
    });

    test('does not notify listeners on success', () async {
      await provider.updateOrder('1', UpdateOrderDTO(status: 'X'));

      expect(provider.notifyCount, 0);
    });
  });

  group('deleteOrder', () {
    test('returns true on success', () async {
      final result = await provider.deleteOrder('1');

      expect(result, true);
      expect(mockRepo.lastDeleteId, '1');
    });

    test('returns false on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('delete failed');

      final result = await provider.deleteOrder('1');

      expect(result, false);
      expect(provider.error, contains('delete failed'));
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');

      await provider.deleteOrder('1');

      expect(provider.notifyCount, 1);
    });

    test('does not notify listeners on success', () async {
      await provider.deleteOrder('1');

      expect(provider.notifyCount, 0);
    });
  });

  group('setPage', () {
    test('updates internal page value', () {
      provider.setPage(3);

      expect(provider.page, 3);
    });

    test('does not notify listeners', () {
      provider.setPage(5);

      expect(provider.notifyCount, 0);
    });

    test('fetchOrders uses the updated page', () async {
      provider.setPage(4);

      await provider.fetchOrders();

      expect(mockRepo.lastGetOrdersParams!.page, 4);
    });
  });

  group('setFilters', () {
    test('updates orderNumber filter', () {
      provider.setFilters(orderNumber: 'ORD-123');

      expect(provider.orderNumberFilter, 'ORD-123');
    });

    test('updates status filter', () {
      provider.setFilters(status: 'pending');

      expect(provider.statusFilter, 'pending');
    });

    test('updates integrationId filter', () {
      provider.setFilters(integrationId: 5);

      expect(provider.integrationIdFilter, 5);
    });

    test('resets page to 1', () {
      provider.setPage(5);
      provider.setFilters(orderNumber: 'X');

      expect(provider.page, 1);
    });

    test('preserves existing filters when not overridden', () {
      provider.setFilters(orderNumber: 'A');
      provider.setFilters(status: 'B');

      expect(provider.orderNumberFilter, 'A');
      expect(provider.statusFilter, 'B');
    });
  });

  group('resetFilters', () {
    test('clears all filters', () {
      provider.setFilters(
        orderNumber: 'X',
        status: 'Y',
        integrationId: 5,
      );
      provider.resetFilters();

      expect(provider.orderNumberFilter, '');
      expect(provider.statusFilter, '');
      expect(provider.integrationIdFilter, isNull);
    });

    test('resets page to 1', () {
      provider.setPage(10);
      provider.resetFilters();

      expect(provider.page, 1);
    });
  });

  group('loading states', () {
    test('isLoading is false before fetch', () {
      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch completes', () async {
      await provider.fetchOrders();

      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch fails', () async {
      mockRepo.errorToThrow = Exception('fail');

      await provider.fetchOrders();

      expect(provider.isLoading, false);
    });
  });
}
