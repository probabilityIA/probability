import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/core/network/api_client.dart';
import 'package:mobile_central/services/modules/orders/app/use_cases.dart';
import 'package:mobile_central/services/modules/orders/domain/entities.dart';
import 'package:mobile_central/services/modules/orders/domain/ports.dart';
import 'package:mobile_central/services/modules/orders/ui/providers/order_provider.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

class FakeOrderRepository implements IOrderRepository {
  FakeOrderRepository({required this.total});

  final int total;
  final List<int> requestedPages = [];

  @override
  Future<PaginatedResponse<Order>> getOrders(GetOrdersParams? params) async {
    final page = params?.page ?? 1;
    final size = params?.pageSize ?? 20;
    requestedPages.add(page);

    final start = (page - 1) * size;
    final end = (start + size) > total ? total : (start + size);
    return PaginatedResponse<Order>(
      data: [for (var i = start; i < end; i++) _order(i)],
      pagination: Pagination(
        currentPage: page,
        perPage: size,
        total: total,
        lastPage: (total / size).ceil(),
        hasNext: end < total,
        hasPrev: page > 1,
      ),
    );
  }

  @override
  Future<Order> getOrderById(String id) async => _order(0);

  @override
  Future<Order> createOrder(CreateOrderDTO data) async => _order(0);

  @override
  Future<Order> updateOrder(String id, UpdateOrderDTO data) async => _order(0);

  @override
  Future<void> deleteOrder(String id) async {}

  @override
  Future<Map<String, dynamic>> getOrderRaw(String id) async => {};

  static Order _order(int i) => Order(
        id: '$i',
        createdAt: '2026-01-01T00:00:00Z',
        updatedAt: '2026-01-01T00:00:00Z',
        integrationId: 1,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-$i',
        orderNumber: 'ORD-$i',
        internalNumber: 'INT-$i',
        subtotal: 100.0,
        tax: 19.0,
        discount: 0.0,
        shippingCost: 10.0,
        totalAmount: 129.0,
        currency: 'COP',
        customerName: 'Cliente $i',
        customerEmail: 'c$i@test.com',
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
}

OrderProvider _provider(FakeOrderRepository repo) => OrderProvider(
      apiClient: ApiClient(),
      useCases: OrderUseCases(repo),
    );

void main() {
  group('OrderProvider paginado real', () {
    test('la primera carga trae una pagina y expone el total', () async {
      final repo = FakeOrderRepository(total: 5000);
      final provider = _provider(repo);

      await provider.fetchOrders();

      expect(provider.orders.length, 20);
      expect(provider.list.total, 5000);
      expect(provider.hasMore, isTrue);
      expect(repo.requestedPages, [1]);
    });

    test('loadMore encadena paginas sin repetirlas', () async {
      final repo = FakeOrderRepository(total: 5000);
      final provider = _provider(repo);

      await provider.fetchOrders();
      await provider.loadMore();
      await provider.loadMore();

      expect(repo.requestedPages, [1, 2, 3]);
      expect(provider.list.itemCount, 60);
    });

    test('recorrer 10.000 ordenes no llena la memoria', () async {
      final repo = FakeOrderRepository(total: 10000);
      final provider = _provider(repo);

      await provider.fetchOrders();
      for (var i = 0; i < 120; i++) {
        provider.list.itemAt(provider.list.itemCount - 1);
        await Future<void>.delayed(Duration.zero);
        await provider.loadMore();
      }
      await Future<void>.delayed(Duration.zero);

      expect(provider.list.itemCount, 20 * 121);
      expect(provider.list.liveItemCount,
          lessThanOrEqualTo(provider.list.maxPagesInMemory * 20));
      expect(provider.list.pagesInMemory,
          lessThanOrEqualTo(provider.list.maxPagesInMemory));
    });

    test('una posicion expulsada se vuelve a pedir al mirarla', () async {
      final repo = FakeOrderRepository(total: 10000);
      final provider = _provider(repo);

      await provider.fetchOrders();
      for (var i = 0; i < 20; i++) {
        provider.list.itemAt(provider.list.itemCount - 1);
        await Future<void>.delayed(Duration.zero);
        await provider.loadMore();
      }
      await Future<void>.delayed(Duration.zero);
      expect(provider.list.isHole(0), isTrue);

      provider.list.itemAt(0);
      await Future<void>.delayed(Duration.zero);
      await Future<void>.delayed(Duration.zero);

      expect(provider.list.itemAt(0)?.orderNumber, 'ORD-0');
    });

    test('cambiar de filtro repagina desde cero', () async {
      final repo = FakeOrderRepository(total: 5000);
      final provider = _provider(repo);

      await provider.fetchOrders();
      await provider.loadMore();
      expect(provider.list.itemCount, 40);

      provider.setFilters(status: 'delivered');
      await provider.fetchOrders();

      expect(provider.list.itemCount, 20);
      expect(provider.list.itemAt(0)?.orderNumber, 'ORD-0');
    });
  });
}
