import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/dashboard/domain/entities.dart';

void main() {
  group('numeros del dashboard', () {
    test('total_sold decimal no rompe el parseo', () {
      final p = TopProduct.fromJson(const {
        'product_name': 'Citrato de Magnesio',
        'product_id': 'PRD_1',
        'sku': 'SKU-1',
        'order_count': 25579,
        'total_sold': 2285154437.19,
      });

      expect(p.totalSold, 2285154437.19);
      expect(p.orderCount, 25579);
    });

    test('total_sold entero sigue funcionando', () {
      final p = TopProduct.fromJson(const {
        'product_name': 'X',
        'product_id': 'PRD_2',
        'sku': 'SKU-2',
        'order_count': 25,
        'total_sold': 2750000,
      });

      expect(p.totalSold, 2750000.0);
    });

    test('un contador decimal se redondea en vez de explotar', () {
      final s = ShipmentsByStatus.fromJson(const {
        'status': 'pending',
        'count': 2.0,
      });

      expect(s.count, 2);
    });

    test('campos ausentes o nulos quedan en cero', () {
      final s = ShipmentsByStatus.fromJson(const {'status': 'pending'});
      expect(s.count, 0);

      final p = TopProduct.fromJson(const {'product_name': 'Y'});
      expect(p.totalSold, 0);
      expect(p.orderCount, 0);
    });

    test('un id nulo se conserva nulo', () {
      final d = TopDriver.fromJson(const {
        'driver_name': 'Sebastian',
        'order_count': 2,
      });

      expect(d.driverId, isNull);
    });

    test('el sobre completo de produccion se parsea', () {
      final stats = DashboardStats.fromJson(const {
        'total_orders': 114,
        'orders_today': 0,
        'orders_by_integration_type': [
          {'integration_type': 'shopify', 'count': 105},
        ],
        'top_products': [
          {
            'product_name': 'Citrato',
            'product_id': 'PRD_1',
            'sku': 'S1',
            'order_count': 12,
            'total_sold': 1985995.2,
          },
        ],
        'shipments_by_status': [
          {'status': 'pending', 'count': 2},
        ],
        'top_customers': [
          {
            'customer_name': 'Luis Ramos',
            'customer_email': 'luis@example.com',
            'order_count': 3,
          },
        ],
      });

      expect(stats.totalOrders, 114);
      expect(stats.topProducts.single.totalSold, 1985995.2);
      expect(stats.shipmentsByStatus.single.count, 2);
    });
  });

  _periodTests();
}

void _periodTests() {
  group('periodo del dashboard', () {
    final now = DateTime(2026, 8, 22, 15, 30);

    test('todo no manda rango', () {
      expect(DashboardPeriod.all.rangeFrom(now), isNull);
    });

    test('hoy es un solo dia', () {
      final r = DashboardPeriod.today.rangeFrom(now)!;
      expect(r.start, '2026-08-22');
      expect(r.end, '2026-08-22');
    });

    test('7 dias incluye hoy', () {
      final r = DashboardPeriod.week.rangeFrom(now)!;
      expect(r.start, '2026-08-16');
      expect(r.end, '2026-08-22');
    });

    test('30 dias incluye hoy', () {
      final r = DashboardPeriod.month.rangeFrom(now)!;
      expect(r.start, '2026-07-24');
      expect(r.end, '2026-08-22');
    });

    test('el rango cruza el cambio de mes', () {
      final r = DashboardPeriod.week.rangeFrom(DateTime(2026, 3, 2))!;
      expect(r.start, '2026-02-24');
      expect(r.end, '2026-03-02');
    });

    test('las etiquetas son las esperadas', () {
      expect(DashboardPeriod.values.map((p) => p.label).toList(),
          ['Hoy', '7 dias', '30 dias', 'Todo']);
    });
  });
}
