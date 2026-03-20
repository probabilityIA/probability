import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/warehouses/domain/entities.dart';

void main() {
  group('Warehouse', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'name': 'Main Warehouse',
        'code': 'WH-001',
        'address': 'Calle 100 #15-20',
        'city': 'Bogota',
        'state': 'Cundinamarca',
        'country': 'CO',
        'zip_code': '110111',
        'phone': '+573001234567',
        'contact_name': 'Juan',
        'contact_email': 'juan@example.com',
        'is_active': true,
        'is_default': true,
        'is_fulfillment': false,
        'company': 'Test Corp',
        'first_name': 'Juan',
        'last_name': 'Perez',
        'email': 'warehouse@example.com',
        'suburb': 'Centro',
        'city_dane_code': '11001',
        'postal_code': '110111',
        'street': 'Calle 100',
        'latitude': 4.6097,
        'longitude': -74.0817,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-03-01T00:00:00Z',
      };

      final warehouse = Warehouse.fromJson(json);

      expect(warehouse.id, 1);
      expect(warehouse.businessId, 5);
      expect(warehouse.name, 'Main Warehouse');
      expect(warehouse.code, 'WH-001');
      expect(warehouse.address, 'Calle 100 #15-20');
      expect(warehouse.city, 'Bogota');
      expect(warehouse.state, 'Cundinamarca');
      expect(warehouse.country, 'CO');
      expect(warehouse.zipCode, '110111');
      expect(warehouse.phone, '+573001234567');
      expect(warehouse.contactName, 'Juan');
      expect(warehouse.contactEmail, 'juan@example.com');
      expect(warehouse.isActive, true);
      expect(warehouse.isDefault, true);
      expect(warehouse.isFulfillment, false);
      expect(warehouse.company, 'Test Corp');
      expect(warehouse.firstName, 'Juan');
      expect(warehouse.lastName, 'Perez');
      expect(warehouse.email, 'warehouse@example.com');
      expect(warehouse.suburb, 'Centro');
      expect(warehouse.cityDaneCode, '11001');
      expect(warehouse.postalCode, '110111');
      expect(warehouse.street, 'Calle 100');
      expect(warehouse.latitude, 4.6097);
      expect(warehouse.longitude, -74.0817);
      expect(warehouse.createdAt, '2026-01-01T00:00:00Z');
      expect(warehouse.updatedAt, '2026-03-01T00:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final warehouse = Warehouse.fromJson(json);

      expect(warehouse.id, 0);
      expect(warehouse.businessId, 0);
      expect(warehouse.name, '');
      expect(warehouse.code, '');
      expect(warehouse.address, '');
      expect(warehouse.city, '');
      expect(warehouse.state, '');
      expect(warehouse.country, '');
      expect(warehouse.zipCode, '');
      expect(warehouse.phone, '');
      expect(warehouse.contactName, '');
      expect(warehouse.contactEmail, '');
      expect(warehouse.isActive, false);
      expect(warehouse.isDefault, false);
      expect(warehouse.isFulfillment, false);
      expect(warehouse.company, '');
      expect(warehouse.firstName, '');
      expect(warehouse.lastName, '');
      expect(warehouse.email, '');
      expect(warehouse.suburb, '');
      expect(warehouse.cityDaneCode, '');
      expect(warehouse.postalCode, '');
      expect(warehouse.street, '');
      expect(warehouse.latitude, isNull);
      expect(warehouse.longitude, isNull);
      expect(warehouse.createdAt, '');
      expect(warehouse.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'latitude': null,
        'longitude': null,
      };

      final warehouse = Warehouse.fromJson(json);

      expect(warehouse.latitude, isNull);
      expect(warehouse.longitude, isNull);
    });
  });

  group('WarehouseDetail', () {
    test('fromJson parses all fields including locations', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'name': 'Detail WH',
        'code': 'DWH-001',
        'address': 'Addr',
        'city': 'City',
        'state': 'State',
        'country': 'CO',
        'zip_code': '110111',
        'phone': '123',
        'contact_name': 'Contact',
        'contact_email': 'contact@test.com',
        'is_active': true,
        'is_default': false,
        'is_fulfillment': true,
        'company': 'Co',
        'first_name': 'First',
        'last_name': 'Last',
        'email': 'e@t.com',
        'suburb': 'Sub',
        'city_dane_code': '11001',
        'postal_code': '110111',
        'street': 'St',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
        'locations': [
          {
            'id': 10,
            'warehouse_id': 1,
            'name': 'Shelf A',
            'code': 'SA',
            'type': 'shelf',
            'is_active': true,
            'is_fulfillment': false,
            'capacity': 100,
            'created_at': '2026-01-01',
            'updated_at': '2026-01-01',
          },
        ],
      };

      final detail = WarehouseDetail.fromJson(json);

      expect(detail.id, 1);
      expect(detail.name, 'Detail WH');
      expect(detail.locations.length, 1);
      expect(detail.locations[0].name, 'Shelf A');
      expect(detail.locations[0].capacity, 100);
    });

    test('fromJson handles null locations list', () {
      final json = {
        'id': 1,
        'locations': null,
      };

      final detail = WarehouseDetail.fromJson(json);

      expect(detail.locations, isEmpty);
    });

    test('fromJson handles missing locations key', () {
      final json = {'id': 1};

      final detail = WarehouseDetail.fromJson(json);

      expect(detail.locations, isEmpty);
    });
  });

  group('WarehouseLocation', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 10,
        'warehouse_id': 1,
        'name': 'Rack B',
        'code': 'RB',
        'type': 'rack',
        'is_active': true,
        'is_fulfillment': true,
        'capacity': 500,
        'created_at': '2026-01-01',
        'updated_at': '2026-03-01',
      };

      final location = WarehouseLocation.fromJson(json);

      expect(location.id, 10);
      expect(location.warehouseId, 1);
      expect(location.name, 'Rack B');
      expect(location.code, 'RB');
      expect(location.type, 'rack');
      expect(location.isActive, true);
      expect(location.isFulfillment, true);
      expect(location.capacity, 500);
      expect(location.createdAt, '2026-01-01');
      expect(location.updatedAt, '2026-03-01');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final location = WarehouseLocation.fromJson(json);

      expect(location.id, 0);
      expect(location.warehouseId, 0);
      expect(location.name, '');
      expect(location.code, '');
      expect(location.type, '');
      expect(location.isActive, false);
      expect(location.isFulfillment, false);
      expect(location.capacity, isNull);
      expect(location.createdAt, '');
      expect(location.updatedAt, '');
    });

    test('fromJson handles null capacity', () {
      final json = {'id': 1, 'capacity': null};

      final location = WarehouseLocation.fromJson(json);

      expect(location.capacity, isNull);
    });
  });

  group('CreateWarehouseDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateWarehouseDTO(name: 'New WH', code: 'NWH');

      final json = dto.toJson();

      expect(json['name'], 'New WH');
      expect(json['code'], 'NWH');
      expect(json.length, 2);
    });

    test('toJson includes all optional fields when provided', () {
      final dto = CreateWarehouseDTO(
        name: 'Full WH',
        code: 'FWH',
        address: 'Addr',
        city: 'City',
        state: 'State',
        country: 'CO',
        zipCode: '110111',
        phone: '123',
        contactName: 'CN',
        contactEmail: 'ce@t.com',
        isDefault: true,
        isFulfillment: false,
        company: 'Co',
        firstName: 'First',
        lastName: 'Last',
        email: 'e@t.com',
        suburb: 'Sub',
        cityDaneCode: '11001',
        postalCode: '110111',
        street: 'St',
        latitude: 4.6,
        longitude: -74.0,
      );

      final json = dto.toJson();

      expect(json['name'], 'Full WH');
      expect(json['code'], 'FWH');
      expect(json['address'], 'Addr');
      expect(json['city'], 'City');
      expect(json['state'], 'State');
      expect(json['country'], 'CO');
      expect(json['zip_code'], '110111');
      expect(json['phone'], '123');
      expect(json['contact_name'], 'CN');
      expect(json['contact_email'], 'ce@t.com');
      expect(json['is_default'], true);
      expect(json['is_fulfillment'], false);
      expect(json['company'], 'Co');
      expect(json['first_name'], 'First');
      expect(json['last_name'], 'Last');
      expect(json['email'], 'e@t.com');
      expect(json['suburb'], 'Sub');
      expect(json['city_dane_code'], '11001');
      expect(json['postal_code'], '110111');
      expect(json['street'], 'St');
      expect(json['latitude'], 4.6);
      expect(json['longitude'], -74.0);
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateWarehouseDTO(name: 'Min', code: 'M');

      final json = dto.toJson();

      expect(json.containsKey('address'), false);
      expect(json.containsKey('latitude'), false);
      expect(json.containsKey('longitude'), false);
    });
  });

  group('UpdateWarehouseDTO', () {
    test('toJson includes required and optional fields', () {
      final dto = UpdateWarehouseDTO(
        name: 'Updated',
        code: 'UPD',
        isActive: false,
        latitude: 5.0,
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated');
      expect(json['code'], 'UPD');
      expect(json['is_active'], false);
      expect(json['latitude'], 5.0);
    });

    test('toJson only includes name and code when no optionals', () {
      final dto = UpdateWarehouseDTO(name: 'N', code: 'C');

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json['name'], 'N');
      expect(json['code'], 'C');
    });
  });

  group('GetWarehousesParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetWarehousesParams(
        page: 2,
        pageSize: 25,
        search: 'main',
        isActive: true,
        isFulfillment: false,
        businessId: 3,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 25);
      expect(queryParams['search'], 'main');
      expect(queryParams['is_active'], true);
      expect(queryParams['is_fulfillment'], false);
      expect(queryParams['business_id'], 3);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetWarehousesParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams.containsKey('page'), true);
    });

    test('toQueryParams returns empty map when all null', () {
      final params = GetWarehousesParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });

  group('CreateLocationDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateLocationDTO(name: 'Loc A', code: 'LA');

      final json = dto.toJson();

      expect(json['name'], 'Loc A');
      expect(json['code'], 'LA');
      expect(json.length, 2);
    });

    test('toJson includes optional fields when provided', () {
      final dto = CreateLocationDTO(
        name: 'Loc B',
        code: 'LB',
        type: 'shelf',
        isFulfillment: true,
        capacity: 200,
      );

      final json = dto.toJson();

      expect(json['type'], 'shelf');
      expect(json['is_fulfillment'], true);
      expect(json['capacity'], 200);
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateLocationDTO(name: 'N', code: 'C');

      final json = dto.toJson();

      expect(json.containsKey('type'), false);
      expect(json.containsKey('is_fulfillment'), false);
      expect(json.containsKey('capacity'), false);
    });
  });

  group('UpdateLocationDTO', () {
    test('toJson includes required and optional fields', () {
      final dto = UpdateLocationDTO(
        name: 'Updated Loc',
        code: 'UL',
        type: 'rack',
        isActive: false,
        isFulfillment: true,
        capacity: 300,
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated Loc');
      expect(json['code'], 'UL');
      expect(json['type'], 'rack');
      expect(json['is_active'], false);
      expect(json['is_fulfillment'], true);
      expect(json['capacity'], 300);
    });

    test('toJson only includes name and code when no optionals', () {
      final dto = UpdateLocationDTO(name: 'N', code: 'C');

      final json = dto.toJson();

      expect(json.length, 2);
    });
  });
}
