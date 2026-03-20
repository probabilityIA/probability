import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/drivers/domain/entities.dart';

void main() {
  group('DriverInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 10,
        'first_name': 'Carlos',
        'last_name': 'Garcia',
        'email': 'carlos@example.com',
        'phone': '+573001234567',
        'identification': 'CC-1234567890',
        'status': 'active',
        'photo_url': 'https://example.com/photo.jpg',
        'license_type': 'B1',
        'license_expiry': '2027-12-31',
        'warehouse_id': 3,
        'notes': 'Experienced driver',
        'created_at': '2026-01-15T10:30:00Z',
        'updated_at': '2026-01-16T12:00:00Z',
      };

      final driver = DriverInfo.fromJson(json);

      expect(driver.id, 1);
      expect(driver.businessId, 10);
      expect(driver.firstName, 'Carlos');
      expect(driver.lastName, 'Garcia');
      expect(driver.email, 'carlos@example.com');
      expect(driver.phone, '+573001234567');
      expect(driver.identification, 'CC-1234567890');
      expect(driver.status, 'active');
      expect(driver.photoUrl, 'https://example.com/photo.jpg');
      expect(driver.licenseType, 'B1');
      expect(driver.licenseExpiry, '2027-12-31');
      expect(driver.warehouseId, 3);
      expect(driver.notes, 'Experienced driver');
      expect(driver.createdAt, '2026-01-15T10:30:00Z');
      expect(driver.updatedAt, '2026-01-16T12:00:00Z');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final driver = DriverInfo.fromJson(json);

      expect(driver.id, 0);
      expect(driver.businessId, 0);
      expect(driver.firstName, '');
      expect(driver.lastName, '');
      expect(driver.email, '');
      expect(driver.phone, '');
      expect(driver.identification, '');
      expect(driver.status, '');
      expect(driver.photoUrl, '');
      expect(driver.licenseType, '');
      expect(driver.createdAt, '');
      expect(driver.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'business_id': 2,
        'first_name': 'Test',
        'last_name': 'Driver',
        'email': 'test@test.com',
        'phone': '555',
        'identification': '123',
        'status': 'active',
        'photo_url': '',
        'license_type': 'B1',
        'created_at': '',
        'updated_at': '',
      };

      final driver = DriverInfo.fromJson(json);

      expect(driver.licenseExpiry, isNull);
      expect(driver.warehouseId, isNull);
      expect(driver.notes, isNull);
    });
  });

  group('CreateDriverDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = CreateDriverDTO(
        firstName: 'Carlos',
        lastName: 'Garcia',
        email: 'carlos@example.com',
        phone: '+573001234567',
        identification: 'CC-123',
        licenseType: 'B1',
        licenseExpiry: '2027-12-31',
        warehouseId: 3,
        notes: 'Good driver',
      );

      final json = dto.toJson();

      expect(json['first_name'], 'Carlos');
      expect(json['last_name'], 'Garcia');
      expect(json['email'], 'carlos@example.com');
      expect(json['phone'], '+573001234567');
      expect(json['identification'], 'CC-123');
      expect(json['license_type'], 'B1');
      expect(json['license_expiry'], '2027-12-31');
      expect(json['warehouse_id'], 3);
      expect(json['notes'], 'Good driver');
    });

    test('toJson includes only required fields when optional fields are null', () {
      final dto = CreateDriverDTO(
        firstName: 'Carlos',
        lastName: 'Garcia',
        phone: '+57300',
        identification: 'CC-123',
      );

      final json = dto.toJson();

      expect(json.length, 4);
      expect(json['first_name'], 'Carlos');
      expect(json['last_name'], 'Garcia');
      expect(json['phone'], '+57300');
      expect(json['identification'], 'CC-123');
      expect(json.containsKey('email'), false);
      expect(json.containsKey('license_type'), false);
      expect(json.containsKey('license_expiry'), false);
      expect(json.containsKey('warehouse_id'), false);
      expect(json.containsKey('notes'), false);
    });
  });

  group('UpdateDriverDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateDriverDTO(
        firstName: 'Updated',
        lastName: 'Name',
        email: 'updated@test.com',
        phone: '555-9999',
        identification: 'CC-NEW',
        status: 'inactive',
        licenseType: 'C1',
        licenseExpiry: '2028-06-30',
        warehouseId: 5,
        notes: 'Updated notes',
      );

      final json = dto.toJson();

      expect(json['first_name'], 'Updated');
      expect(json['last_name'], 'Name');
      expect(json['email'], 'updated@test.com');
      expect(json['phone'], '555-9999');
      expect(json['identification'], 'CC-NEW');
      expect(json['status'], 'inactive');
      expect(json['license_type'], 'C1');
      expect(json['license_expiry'], '2028-06-30');
      expect(json['warehouse_id'], 5);
      expect(json['notes'], 'Updated notes');
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateDriverDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateDriverDTO(firstName: 'NewFirst', status: 'active');

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json['first_name'], 'NewFirst');
      expect(json['status'], 'active');
    });
  });

  group('GetDriversParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetDriversParams(
        page: 2,
        pageSize: 15,
        search: 'Carlos',
        status: 'active',
        businessId: 5,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 15);
      expect(queryParams['search'], 'Carlos');
      expect(queryParams['status'], 'active');
      expect(queryParams['business_id'], 5);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetDriversParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams.containsKey('page'), true);
      expect(queryParams.containsKey('page_size'), false);
      expect(queryParams.containsKey('search'), false);
      expect(queryParams.containsKey('status'), false);
      expect(queryParams.containsKey('business_id'), false);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetDriversParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });
}
