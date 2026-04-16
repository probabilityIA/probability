import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/routes/domain/entities.dart';

void main() {
  group('RouteInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'driver_id': 10,
        'driver_name': 'Carlos',
        'vehicle_id': 20,
        'vehicle_plate': 'ABC-123',
        'status': 'in_progress',
        'date': '2026-03-01',
        'start_time': '08:00',
        'end_time': '17:00',
        'origin_address': 'Calle 100 #10-30',
        'total_stops': 5,
        'completed_stops': 3,
        'failed_stops': 1,
        'notes': 'Priority route',
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
      };

      final route = RouteInfo.fromJson(json);

      expect(route.id, 1);
      expect(route.businessId, 5);
      expect(route.driverId, 10);
      expect(route.driverName, 'Carlos');
      expect(route.vehicleId, 20);
      expect(route.vehiclePlate, 'ABC-123');
      expect(route.status, 'in_progress');
      expect(route.date, '2026-03-01');
      expect(route.startTime, '08:00');
      expect(route.endTime, '17:00');
      expect(route.originAddress, 'Calle 100 #10-30');
      expect(route.totalStops, 5);
      expect(route.completedStops, 3);
      expect(route.failedStops, 1);
      expect(route.notes, 'Priority route');
      expect(route.createdAt, '2026-01-01T00:00:00Z');
      expect(route.updatedAt, '2026-01-02T00:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final route = RouteInfo.fromJson(json);

      expect(route.id, 0);
      expect(route.businessId, 0);
      expect(route.status, '');
      expect(route.date, '');
      expect(route.totalStops, 0);
      expect(route.completedStops, 0);
      expect(route.failedStops, 0);
      expect(route.createdAt, '');
      expect(route.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'status': 'pending',
        'date': '2026-03-01',
        'total_stops': 0,
        'completed_stops': 0,
        'failed_stops': 0,
        'created_at': 'c',
        'updated_at': 'u',
      };

      final route = RouteInfo.fromJson(json);

      expect(route.driverId, isNull);
      expect(route.driverName, isNull);
      expect(route.vehicleId, isNull);
      expect(route.vehiclePlate, isNull);
      expect(route.startTime, isNull);
      expect(route.endTime, isNull);
      expect(route.originAddress, isNull);
      expect(route.notes, isNull);
    });
  });

  group('RouteStopInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'route_id': 10,
        'order_id': 'ord-123',
        'sequence': 3,
        'status': 'delivered',
        'address': 'Calle 50 #20-10',
        'city': 'Bogota',
        'lat': 4.624335,
        'lng': -74.063644,
        'customer_name': 'Maria',
        'customer_phone': '+57123456789',
        'estimated_arrival': '10:00',
        'actual_arrival': '10:15',
        'actual_departure': '10:30',
        'delivery_notes': 'Leave at door',
        'failure_reason': null,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
      };

      final stop = RouteStopInfo.fromJson(json);

      expect(stop.id, 1);
      expect(stop.routeId, 10);
      expect(stop.orderId, 'ord-123');
      expect(stop.sequence, 3);
      expect(stop.status, 'delivered');
      expect(stop.address, 'Calle 50 #20-10');
      expect(stop.city, 'Bogota');
      expect(stop.lat, 4.624335);
      expect(stop.lng, -74.063644);
      expect(stop.customerName, 'Maria');
      expect(stop.customerPhone, '+57123456789');
      expect(stop.estimatedArrival, '10:00');
      expect(stop.actualArrival, '10:15');
      expect(stop.actualDeparture, '10:30');
      expect(stop.deliveryNotes, 'Leave at door');
      expect(stop.failureReason, isNull);
      expect(stop.createdAt, '2026-01-01T00:00:00Z');
      expect(stop.updatedAt, '2026-01-02T00:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final stop = RouteStopInfo.fromJson(json);

      expect(stop.id, 0);
      expect(stop.routeId, 0);
      expect(stop.sequence, 0);
      expect(stop.status, '');
      expect(stop.address, '');
      expect(stop.customerName, '');
      expect(stop.createdAt, '');
      expect(stop.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'route_id': 1,
        'sequence': 1,
        'status': 'pending',
        'address': 'addr',
        'customer_name': 'name',
        'created_at': 'c',
        'updated_at': 'u',
      };

      final stop = RouteStopInfo.fromJson(json);

      expect(stop.orderId, isNull);
      expect(stop.city, isNull);
      expect(stop.lat, isNull);
      expect(stop.lng, isNull);
      expect(stop.customerPhone, isNull);
      expect(stop.estimatedArrival, isNull);
      expect(stop.actualArrival, isNull);
      expect(stop.actualDeparture, isNull);
      expect(stop.deliveryNotes, isNull);
      expect(stop.failureReason, isNull);
    });
  });

  group('RouteDetail', () {
    test('fromJson parses all fields including stops', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'driver_id': 10,
        'status': 'in_progress',
        'date': '2026-03-01',
        'total_stops': 2,
        'completed_stops': 1,
        'failed_stops': 0,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
        'actual_start_time': '08:05',
        'actual_end_time': null,
        'origin_warehouse_id': 3,
        'origin_lat': 4.6,
        'origin_lng': -74.0,
        'total_distance_km': 25.5,
        'total_duration_min': 45.0,
        'stops': [
          {
            'id': 1,
            'route_id': 1,
            'sequence': 1,
            'status': 'delivered',
            'address': 'Addr 1',
            'customer_name': 'C1',
            'created_at': 'c',
            'updated_at': 'u',
          },
          {
            'id': 2,
            'route_id': 1,
            'sequence': 2,
            'status': 'pending',
            'address': 'Addr 2',
            'customer_name': 'C2',
            'created_at': 'c',
            'updated_at': 'u',
          },
        ],
      };

      final detail = RouteDetail.fromJson(json);

      expect(detail.id, 1);
      expect(detail.businessId, 5);
      expect(detail.actualStartTime, '08:05');
      expect(detail.actualEndTime, isNull);
      expect(detail.originWarehouseId, 3);
      expect(detail.originLat, 4.6);
      expect(detail.originLng, -74.0);
      expect(detail.totalDistanceKm, 25.5);
      expect(detail.totalDurationMin, 45.0);
      expect(detail.stops.length, 2);
      expect(detail.stops[0].customerName, 'C1');
      expect(detail.stops[1].status, 'pending');
    });

    test('fromJson handles null stops list', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'status': 's',
        'date': 'd',
        'total_stops': 0,
        'completed_stops': 0,
        'failed_stops': 0,
        'created_at': 'c',
        'updated_at': 'u',
      };

      final detail = RouteDetail.fromJson(json);

      expect(detail.stops, isEmpty);
      expect(detail.actualStartTime, isNull);
      expect(detail.originWarehouseId, isNull);
      expect(detail.originLat, isNull);
      expect(detail.originLng, isNull);
      expect(detail.totalDistanceKm, isNull);
      expect(detail.totalDurationMin, isNull);
    });
  });

  group('CreateRouteStopDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateRouteStopDTO(
        address: 'Calle 100',
        customerName: 'Maria',
      );

      final json = dto.toJson();

      expect(json['address'], 'Calle 100');
      expect(json['customer_name'], 'Maria');
    });

    test('toJson includes all non-null optional fields', () {
      final dto = CreateRouteStopDTO(
        orderId: 'ord-1',
        address: 'Calle 100',
        city: 'Bogota',
        lat: 4.6,
        lng: -74.0,
        customerName: 'Maria',
        customerPhone: '+57123',
        deliveryNotes: 'Ring bell',
      );

      final json = dto.toJson();

      expect(json['order_id'], 'ord-1');
      expect(json['city'], 'Bogota');
      expect(json['lat'], 4.6);
      expect(json['lng'], -74.0);
      expect(json['customer_phone'], '+57123');
      expect(json['delivery_notes'], 'Ring bell');
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateRouteStopDTO(
        address: 'Addr',
        customerName: 'Name',
      );

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json.containsKey('order_id'), false);
      expect(json.containsKey('city'), false);
      expect(json.containsKey('lat'), false);
    });
  });

  group('CreateRouteDTO', () {
    test('toJson includes required date', () {
      final dto = CreateRouteDTO(date: '2026-03-01');

      final json = dto.toJson();

      expect(json['date'], '2026-03-01');
    });

    test('toJson includes all non-null optional fields', () {
      final dto = CreateRouteDTO(
        date: '2026-03-01',
        driverId: 10,
        vehicleId: 20,
        originAddress: 'Origin',
        originLat: 4.6,
        originLng: -74.0,
        notes: 'Note',
        stops: [
          CreateRouteStopDTO(address: 'Stop 1', customerName: 'C1'),
        ],
      );

      final json = dto.toJson();

      expect(json['driver_id'], 10);
      expect(json['vehicle_id'], 20);
      expect(json['origin_address'], 'Origin');
      expect(json['origin_lat'], 4.6);
      expect(json['origin_lng'], -74.0);
      expect(json['notes'], 'Note');
      expect(json['stops'], isA<List>());
      expect((json['stops'] as List).length, 1);
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateRouteDTO(date: '2026-03-01');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json.containsKey('driver_id'), false);
      expect(json.containsKey('stops'), false);
    });
  });

  group('UpdateRouteDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateRouteDTO(
        driverId: 10,
        vehicleId: 20,
        date: '2026-04-01',
        originAddress: 'New Origin',
        originLat: 4.7,
        originLng: -74.1,
        notes: 'Updated note',
      );

      final json = dto.toJson();

      expect(json['driver_id'], 10);
      expect(json['vehicle_id'], 20);
      expect(json['date'], '2026-04-01');
      expect(json['origin_address'], 'New Origin');
      expect(json['origin_lat'], 4.7);
      expect(json['origin_lng'], -74.1);
      expect(json['notes'], 'Updated note');
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateRouteDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });
  });

  group('AddStopDTO', () {
    test('toJson includes required fields', () {
      final dto = AddStopDTO(address: 'Addr', customerName: 'Name');

      final json = dto.toJson();

      expect(json['address'], 'Addr');
      expect(json['customer_name'], 'Name');
    });

    test('toJson includes all non-null optional fields', () {
      final dto = AddStopDTO(
        orderId: 'o1',
        address: 'Addr',
        city: 'City',
        lat: 1.0,
        lng: 2.0,
        customerName: 'Name',
        customerPhone: '123',
        deliveryNotes: 'Notes',
      );

      final json = dto.toJson();

      expect(json['order_id'], 'o1');
      expect(json['city'], 'City');
      expect(json['lat'], 1.0);
      expect(json['lng'], 2.0);
      expect(json['customer_phone'], '123');
      expect(json['delivery_notes'], 'Notes');
    });

    test('toJson excludes null optional fields', () {
      final dto = AddStopDTO(address: 'A', customerName: 'N');

      final json = dto.toJson();

      expect(json.length, 2);
    });
  });

  group('UpdateStopDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateStopDTO(
        address: 'NewAddr',
        city: 'NewCity',
        lat: 3.0,
        lng: 4.0,
        customerName: 'NewName',
        customerPhone: '456',
        deliveryNotes: 'NewNotes',
      );

      final json = dto.toJson();

      expect(json['address'], 'NewAddr');
      expect(json['city'], 'NewCity');
      expect(json['lat'], 3.0);
      expect(json['lng'], 4.0);
      expect(json['customer_name'], 'NewName');
      expect(json['customer_phone'], '456');
      expect(json['delivery_notes'], 'NewNotes');
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateStopDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });
  });

  group('UpdateStopStatusDTO', () {
    test('toJson includes status', () {
      final dto = UpdateStopStatusDTO(status: 'delivered');

      final json = dto.toJson();

      expect(json['status'], 'delivered');
    });

    test('toJson includes failureReason when present', () {
      final dto = UpdateStopStatusDTO(
        status: 'failed',
        failureReason: 'Customer not home',
      );

      final json = dto.toJson();

      expect(json['status'], 'failed');
      expect(json['failure_reason'], 'Customer not home');
    });

    test('toJson excludes null failureReason', () {
      final dto = UpdateStopStatusDTO(status: 'delivered');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json.containsKey('failure_reason'), false);
    });
  });

  group('ReorderStopsDTO', () {
    test('toJson produces correct structure', () {
      final dto = ReorderStopsDTO(stopIds: [3, 1, 2]);

      final json = dto.toJson();

      expect(json['stop_ids'], [3, 1, 2]);
    });

    test('toJson handles empty list', () {
      final dto = ReorderStopsDTO(stopIds: []);

      final json = dto.toJson();

      expect(json['stop_ids'], isEmpty);
    });
  });

  group('GetRoutesParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetRoutesParams(
        page: 1,
        pageSize: 20,
        status: 'in_progress',
        driverId: 10,
        dateFrom: '2026-01-01',
        dateTo: '2026-12-31',
        search: 'route',
        businessId: 5,
      );

      final qp = params.toQueryParams();

      expect(qp['page'], 1);
      expect(qp['page_size'], 20);
      expect(qp['status'], 'in_progress');
      expect(qp['driver_id'], 10);
      expect(qp['date_from'], '2026-01-01');
      expect(qp['date_to'], '2026-12-31');
      expect(qp['search'], 'route');
      expect(qp['business_id'], 5);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetRoutesParams(page: 2);

      final qp = params.toQueryParams();

      expect(qp.length, 1);
      expect(qp.containsKey('page'), true);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetRoutesParams();

      final qp = params.toQueryParams();

      expect(qp, isEmpty);
    });
  });

  group('DriverOption', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'first_name': 'Carlos',
        'last_name': 'Garcia',
        'phone': '+57123',
        'identification': '1234567',
        'status': 'active',
        'license_type': 'B2',
      };

      final driver = DriverOption.fromJson(json);

      expect(driver.id, 1);
      expect(driver.firstName, 'Carlos');
      expect(driver.lastName, 'Garcia');
      expect(driver.phone, '+57123');
      expect(driver.identification, '1234567');
      expect(driver.status, 'active');
      expect(driver.licenseType, 'B2');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};
      final driver = DriverOption.fromJson(json);

      expect(driver.id, 0);
      expect(driver.firstName, '');
      expect(driver.lastName, '');
      expect(driver.phone, '');
      expect(driver.identification, '');
      expect(driver.status, '');
      expect(driver.licenseType, '');
    });
  });

  group('VehicleOption', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'type': 'van',
        'license_plate': 'ABC-123',
        'brand': 'Toyota',
        'vehicle_model': 'HiAce',
        'status': 'available',
      };

      final vehicle = VehicleOption.fromJson(json);

      expect(vehicle.id, 1);
      expect(vehicle.type, 'van');
      expect(vehicle.licensePlate, 'ABC-123');
      expect(vehicle.brand, 'Toyota');
      expect(vehicle.vehicleModel, 'HiAce');
      expect(vehicle.status, 'available');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};
      final vehicle = VehicleOption.fromJson(json);

      expect(vehicle.id, 0);
      expect(vehicle.type, '');
      expect(vehicle.licensePlate, '');
      expect(vehicle.brand, '');
      expect(vehicle.vehicleModel, '');
      expect(vehicle.status, '');
    });
  });

  group('AssignableOrder', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'order_number': 'ORD-001',
        'customer_name': 'Maria',
        'customer_phone': '+57123',
        'address': 'Calle 50',
        'city': 'Bogota',
        'lat': 4.6,
        'lng': -74.0,
        'total_amount': 150000.0,
        'item_count': 3,
        'created_at': '2026-01-01',
      };

      final order = AssignableOrder.fromJson(json);

      expect(order.id, '42');
      expect(order.orderNumber, 'ORD-001');
      expect(order.customerName, 'Maria');
      expect(order.customerPhone, '+57123');
      expect(order.address, 'Calle 50');
      expect(order.city, 'Bogota');
      expect(order.lat, 4.6);
      expect(order.lng, -74.0);
      expect(order.totalAmount, 150000.0);
      expect(order.itemCount, 3);
      expect(order.createdAt, '2026-01-01');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};
      final order = AssignableOrder.fromJson(json);

      expect(order.id, '');
      expect(order.orderNumber, '');
      expect(order.customerName, '');
      expect(order.customerPhone, '');
      expect(order.address, '');
      expect(order.city, '');
      expect(order.totalAmount, 0.0);
      expect(order.itemCount, 0);
      expect(order.createdAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': '1',
        'order_number': 'O',
        'customer_name': 'N',
        'customer_phone': 'P',
        'address': 'A',
        'city': 'C',
        'total_amount': 100,
        'item_count': 1,
        'created_at': 'c',
      };

      final order = AssignableOrder.fromJson(json);

      expect(order.lat, isNull);
      expect(order.lng, isNull);
    });
  });
}
