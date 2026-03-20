import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/vehicles/domain/entities.dart';

void main() {
  group('VehicleInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'type': 'truck',
        'license_plate': 'ABC-123',
        'brand': 'Toyota',
        'model': 'Hilux',
        'year': 2024,
        'color': 'White',
        'status': 'active',
        'weight_capacity_kg': 1500.5,
        'volume_capacity_m3': 8.0,
        'photo_url': 'https://example.com/photo.png',
        'insurance_expiry': '2027-01-01',
        'registration_expiry': '2026-12-31',
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-03-01T00:00:00Z',
      };

      final vehicle = VehicleInfo.fromJson(json);

      expect(vehicle.id, 1);
      expect(vehicle.businessId, 5);
      expect(vehicle.type, 'truck');
      expect(vehicle.licensePlate, 'ABC-123');
      expect(vehicle.brand, 'Toyota');
      expect(vehicle.model, 'Hilux');
      expect(vehicle.year, 2024);
      expect(vehicle.color, 'White');
      expect(vehicle.status, 'active');
      expect(vehicle.weightCapacityKg, 1500.5);
      expect(vehicle.volumeCapacityM3, 8.0);
      expect(vehicle.photoUrl, 'https://example.com/photo.png');
      expect(vehicle.insuranceExpiry, '2027-01-01');
      expect(vehicle.registrationExpiry, '2026-12-31');
      expect(vehicle.createdAt, '2026-01-01T00:00:00Z');
      expect(vehicle.updatedAt, '2026-03-01T00:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final vehicle = VehicleInfo.fromJson(json);

      expect(vehicle.id, 0);
      expect(vehicle.businessId, 0);
      expect(vehicle.type, '');
      expect(vehicle.licensePlate, '');
      expect(vehicle.brand, '');
      expect(vehicle.model, '');
      expect(vehicle.year, isNull);
      expect(vehicle.color, '');
      expect(vehicle.status, '');
      expect(vehicle.weightCapacityKg, isNull);
      expect(vehicle.volumeCapacityM3, isNull);
      expect(vehicle.photoUrl, '');
      expect(vehicle.insuranceExpiry, isNull);
      expect(vehicle.registrationExpiry, isNull);
      expect(vehicle.createdAt, '');
      expect(vehicle.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 10,
        'year': null,
        'weight_capacity_kg': null,
        'volume_capacity_m3': null,
        'insurance_expiry': null,
        'registration_expiry': null,
      };

      final vehicle = VehicleInfo.fromJson(json);

      expect(vehicle.id, 10);
      expect(vehicle.year, isNull);
      expect(vehicle.weightCapacityKg, isNull);
      expect(vehicle.volumeCapacityM3, isNull);
      expect(vehicle.insuranceExpiry, isNull);
      expect(vehicle.registrationExpiry, isNull);
    });
  });

  group('CreateVehicleDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateVehicleDTO(
        type: 'van',
        licensePlate: 'XYZ-789',
      );

      final json = dto.toJson();

      expect(json['type'], 'van');
      expect(json['license_plate'], 'XYZ-789');
      expect(json.length, 2);
    });

    test('toJson includes all optional fields when provided', () {
      final dto = CreateVehicleDTO(
        type: 'truck',
        licensePlate: 'ABC-123',
        brand: 'Chevrolet',
        model: 'NHR',
        year: 2025,
        color: 'Blue',
        weightCapacityKg: 2000.0,
        volumeCapacityM3: 12.5,
        insuranceExpiry: '2027-06-01',
        registrationExpiry: '2026-12-31',
      );

      final json = dto.toJson();

      expect(json['type'], 'truck');
      expect(json['license_plate'], 'ABC-123');
      expect(json['brand'], 'Chevrolet');
      expect(json['model'], 'NHR');
      expect(json['year'], 2025);
      expect(json['color'], 'Blue');
      expect(json['weight_capacity_kg'], 2000.0);
      expect(json['volume_capacity_m3'], 12.5);
      expect(json['insurance_expiry'], '2027-06-01');
      expect(json['registration_expiry'], '2026-12-31');
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateVehicleDTO(
        type: 'bike',
        licensePlate: 'B-001',
      );

      final json = dto.toJson();

      expect(json.containsKey('brand'), false);
      expect(json.containsKey('model'), false);
      expect(json.containsKey('year'), false);
      expect(json.containsKey('color'), false);
      expect(json.containsKey('weight_capacity_kg'), false);
      expect(json.containsKey('volume_capacity_m3'), false);
      expect(json.containsKey('insurance_expiry'), false);
      expect(json.containsKey('registration_expiry'), false);
    });
  });

  group('UpdateVehicleDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateVehicleDTO(
        type: 'sedan',
        licensePlate: 'UPD-001',
        brand: 'Honda',
        model: 'Civic',
        year: 2023,
        color: 'Red',
        status: 'maintenance',
        weightCapacityKg: 500.0,
        volumeCapacityM3: 3.0,
        insuranceExpiry: '2027-01-01',
        registrationExpiry: '2026-06-30',
      );

      final json = dto.toJson();

      expect(json['type'], 'sedan');
      expect(json['license_plate'], 'UPD-001');
      expect(json['brand'], 'Honda');
      expect(json['model'], 'Civic');
      expect(json['year'], 2023);
      expect(json['color'], 'Red');
      expect(json['status'], 'maintenance');
      expect(json['weight_capacity_kg'], 500.0);
      expect(json['volume_capacity_m3'], 3.0);
      expect(json['insurance_expiry'], '2027-01-01');
      expect(json['registration_expiry'], '2026-06-30');
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateVehicleDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateVehicleDTO(status: 'inactive');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['status'], 'inactive');
    });
  });

  group('GetVehiclesParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetVehiclesParams(
        page: 2,
        pageSize: 25,
        search: 'Toyota',
        type: 'truck',
        status: 'active',
        businessId: 3,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 25);
      expect(queryParams['search'], 'Toyota');
      expect(queryParams['type'], 'truck');
      expect(queryParams['status'], 'active');
      expect(queryParams['business_id'], 3);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetVehiclesParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams.containsKey('page'), true);
      expect(queryParams.containsKey('search'), false);
    });

    test('toQueryParams returns empty map when all null', () {
      final params = GetVehiclesParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });
}
