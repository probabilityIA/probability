import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/dashboard/domain/entities.dart';

void main() {
  group('OrderCountByIntegrationType', () {
    test('fromJson parses all fields correctly', () {
      final json = {'integration_type': 'shopify', 'count': 150};

      final entity = OrderCountByIntegrationType.fromJson(json);

      expect(entity.integrationType, 'shopify');
      expect(entity.count, 150);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = OrderCountByIntegrationType.fromJson(json);

      expect(entity.integrationType, '');
      expect(entity.count, 0);
    });
  });

  group('TopCustomer', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'customer_name': 'John Doe',
        'customer_email': 'john@example.com',
        'order_count': 25,
      };

      final entity = TopCustomer.fromJson(json);

      expect(entity.customerName, 'John Doe');
      expect(entity.customerEmail, 'john@example.com');
      expect(entity.orderCount, 25);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = TopCustomer.fromJson(json);

      expect(entity.customerName, '');
      expect(entity.customerEmail, '');
      expect(entity.orderCount, 0);
    });
  });

  group('OrderCountByLocation', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'city': 'Bogota',
        'state': 'Cundinamarca',
        'order_count': 300,
      };

      final entity = OrderCountByLocation.fromJson(json);

      expect(entity.city, 'Bogota');
      expect(entity.state, 'Cundinamarca');
      expect(entity.orderCount, 300);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = OrderCountByLocation.fromJson(json);

      expect(entity.city, '');
      expect(entity.state, '');
      expect(entity.orderCount, 0);
    });
  });

  group('TopDriver', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'driver_name': 'Carlos',
        'driver_id': 5,
        'order_count': 100,
      };

      final entity = TopDriver.fromJson(json);

      expect(entity.driverName, 'Carlos');
      expect(entity.driverId, 5);
      expect(entity.orderCount, 100);
    });

    test('fromJson handles null driver_id', () {
      final json = {
        'driver_name': 'Carlos',
        'order_count': 50,
      };

      final entity = TopDriver.fromJson(json);

      expect(entity.driverId, isNull);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = TopDriver.fromJson(json);

      expect(entity.driverName, '');
      expect(entity.driverId, isNull);
      expect(entity.orderCount, 0);
    });
  });

  group('DriverByLocation', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'driver_name': 'Carlos',
        'city': 'Medellin',
        'state': 'Antioquia',
        'order_count': 75,
      };

      final entity = DriverByLocation.fromJson(json);

      expect(entity.driverName, 'Carlos');
      expect(entity.city, 'Medellin');
      expect(entity.state, 'Antioquia');
      expect(entity.orderCount, 75);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = DriverByLocation.fromJson(json);

      expect(entity.driverName, '');
      expect(entity.city, '');
      expect(entity.state, '');
      expect(entity.orderCount, 0);
    });
  });

  group('TopProduct', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'product_name': 'Widget',
        'product_id': 'P001',
        'sku': 'SKU-001',
        'order_count': 200,
        'total_sold': 500,
      };

      final entity = TopProduct.fromJson(json);

      expect(entity.productName, 'Widget');
      expect(entity.productId, 'P001');
      expect(entity.sku, 'SKU-001');
      expect(entity.orderCount, 200);
      expect(entity.totalSold, 500);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = TopProduct.fromJson(json);

      expect(entity.productName, '');
      expect(entity.productId, '');
      expect(entity.sku, '');
      expect(entity.orderCount, 0);
      expect(entity.totalSold, 0);
    });
  });

  group('ProductByCategory', () {
    test('fromJson parses all fields correctly', () {
      final json = {'category': 'Electronics', 'count': 42};

      final entity = ProductByCategory.fromJson(json);

      expect(entity.category, 'Electronics');
      expect(entity.count, 42);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = ProductByCategory.fromJson(json);

      expect(entity.category, '');
      expect(entity.count, 0);
    });
  });

  group('ProductByBrand', () {
    test('fromJson parses all fields correctly', () {
      final json = {'brand': 'Samsung', 'count': 30};

      final entity = ProductByBrand.fromJson(json);

      expect(entity.brand, 'Samsung');
      expect(entity.count, 30);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = ProductByBrand.fromJson(json);

      expect(entity.brand, '');
      expect(entity.count, 0);
    });
  });

  group('ShipmentsByStatus', () {
    test('fromJson parses all fields correctly', () {
      final json = {'status': 'delivered', 'count': 120};

      final entity = ShipmentsByStatus.fromJson(json);

      expect(entity.status, 'delivered');
      expect(entity.count, 120);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = ShipmentsByStatus.fromJson(json);

      expect(entity.status, '');
      expect(entity.count, 0);
    });
  });

  group('ShipmentsByCarrier', () {
    test('fromJson parses all fields correctly', () {
      final json = {'carrier': 'Servientrega', 'count': 80};

      final entity = ShipmentsByCarrier.fromJson(json);

      expect(entity.carrier, 'Servientrega');
      expect(entity.count, 80);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = ShipmentsByCarrier.fromJson(json);

      expect(entity.carrier, '');
      expect(entity.count, 0);
    });
  });

  group('ShipmentsByWarehouse', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'warehouse_name': 'Main Warehouse',
        'warehouse_id': 3,
        'count': 200,
      };

      final entity = ShipmentsByWarehouse.fromJson(json);

      expect(entity.warehouseName, 'Main Warehouse');
      expect(entity.warehouseId, 3);
      expect(entity.count, 200);
    });

    test('fromJson handles null warehouse_id', () {
      final json = {
        'warehouse_name': 'Temp',
        'count': 10,
      };

      final entity = ShipmentsByWarehouse.fromJson(json);

      expect(entity.warehouseId, isNull);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = ShipmentsByWarehouse.fromJson(json);

      expect(entity.warehouseName, '');
      expect(entity.warehouseId, isNull);
      expect(entity.count, 0);
    });
  });

  group('OrdersByBusiness', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'business_id': 7,
        'business_name': 'Acme Corp',
        'order_count': 500,
      };

      final entity = OrdersByBusiness.fromJson(json);

      expect(entity.businessId, 7);
      expect(entity.businessName, 'Acme Corp');
      expect(entity.orderCount, 500);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final entity = OrdersByBusiness.fromJson(json);

      expect(entity.businessId, 0);
      expect(entity.businessName, '');
      expect(entity.orderCount, 0);
    });
  });

  group('DashboardStats', () {
    test('fromJson parses all fields with nested lists', () {
      final json = {
        'total_orders': 1000,
        'orders_by_integration_type': [
          {'integration_type': 'shopify', 'count': 500},
          {'integration_type': 'amazon', 'count': 300},
        ],
        'top_customers': [
          {'customer_name': 'John', 'customer_email': 'john@test.com', 'order_count': 50},
        ],
        'orders_by_location': [
          {'city': 'Bogota', 'state': 'Cundinamarca', 'order_count': 400},
        ],
        'top_drivers': [
          {'driver_name': 'Carlos', 'driver_id': 1, 'order_count': 100},
        ],
        'drivers_by_location': [
          {'driver_name': 'Carlos', 'city': 'Bogota', 'state': 'Cundinamarca', 'order_count': 80},
        ],
        'top_products': [
          {'product_name': 'Widget', 'product_id': 'P1', 'sku': 'S1', 'order_count': 200, 'total_sold': 400},
        ],
        'products_by_category': [
          {'category': 'Electronics', 'count': 50},
        ],
        'products_by_brand': [
          {'brand': 'Samsung', 'count': 30},
        ],
        'shipments_by_status': [
          {'status': 'delivered', 'count': 600},
        ],
        'shipments_by_carrier': [
          {'carrier': 'FedEx', 'count': 300},
        ],
        'shipments_by_warehouse': [
          {'warehouse_name': 'Main', 'warehouse_id': 1, 'count': 500},
        ],
        'orders_by_business': [
          {'business_id': 1, 'business_name': 'Acme', 'order_count': 1000},
        ],
      };

      final stats = DashboardStats.fromJson(json);

      expect(stats.totalOrders, 1000);
      expect(stats.ordersByIntegrationType.length, 2);
      expect(stats.ordersByIntegrationType[0].integrationType, 'shopify');
      expect(stats.topCustomers.length, 1);
      expect(stats.topCustomers[0].customerName, 'John');
      expect(stats.ordersByLocation.length, 1);
      expect(stats.topDrivers.length, 1);
      expect(stats.driversByLocation.length, 1);
      expect(stats.topProducts.length, 1);
      expect(stats.productsByCategory.length, 1);
      expect(stats.productsByBrand.length, 1);
      expect(stats.shipmentsByStatus.length, 1);
      expect(stats.shipmentsByCarrier.length, 1);
      expect(stats.shipmentsByWarehouse.length, 1);
      expect(stats.ordersByBusiness, isNotNull);
      expect(stats.ordersByBusiness!.length, 1);
    });

    test('fromJson uses empty lists for missing/null list fields', () {
      final json = <String, dynamic>{};

      final stats = DashboardStats.fromJson(json);

      expect(stats.totalOrders, 0);
      expect(stats.ordersByIntegrationType, isEmpty);
      expect(stats.topCustomers, isEmpty);
      expect(stats.ordersByLocation, isEmpty);
      expect(stats.topDrivers, isEmpty);
      expect(stats.driversByLocation, isEmpty);
      expect(stats.topProducts, isEmpty);
      expect(stats.productsByCategory, isEmpty);
      expect(stats.productsByBrand, isEmpty);
      expect(stats.shipmentsByStatus, isEmpty);
      expect(stats.shipmentsByCarrier, isEmpty);
      expect(stats.shipmentsByWarehouse, isEmpty);
      expect(stats.ordersByBusiness, isNull);
    });
  });

  group('DashboardStatsResponse', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'success': true,
        'message': 'Stats retrieved',
        'data': {
          'total_orders': 100,
        },
      };

      final response = DashboardStatsResponse.fromJson(json);

      expect(response.success, true);
      expect(response.message, 'Stats retrieved');
      expect(response.data.totalOrders, 100);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = DashboardStatsResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, '');
      expect(response.data.totalOrders, 0);
    });

    test('fromJson handles null data as empty map', () {
      final json = {
        'success': true,
        'message': 'OK',
        'data': null,
      };

      final response = DashboardStatsResponse.fromJson(json);

      expect(response.data.totalOrders, 0);
      expect(response.data.ordersByIntegrationType, isEmpty);
    });
  });
}
