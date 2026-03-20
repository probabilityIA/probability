import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/customers/domain/entities.dart';

void main() {
  group('CustomerInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 10,
        'name': 'John Doe',
        'email': 'john@example.com',
        'phone': '+573001234567',
        'dni': '1234567890',
        'created_at': '2026-01-15T10:30:00Z',
        'updated_at': '2026-01-16T12:00:00Z',
      };

      final customer = CustomerInfo.fromJson(json);

      expect(customer.id, 1);
      expect(customer.businessId, 10);
      expect(customer.name, 'John Doe');
      expect(customer.email, 'john@example.com');
      expect(customer.phone, '+573001234567');
      expect(customer.dni, '1234567890');
      expect(customer.createdAt, '2026-01-15T10:30:00Z');
      expect(customer.updatedAt, '2026-01-16T12:00:00Z');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final customer = CustomerInfo.fromJson(json);

      expect(customer.id, 0);
      expect(customer.businessId, 0);
      expect(customer.name, '');
      expect(customer.phone, '');
      expect(customer.createdAt, '');
      expect(customer.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'business_id': 2,
        'name': 'Jane',
        'phone': '555-1234',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
      };

      final customer = CustomerInfo.fromJson(json);

      expect(customer.email, isNull);
      expect(customer.dni, isNull);
    });
  });

  group('CustomerDetail', () {
    test('fromJson parses all fields including inherited and extra', () {
      final json = {
        'id': 1,
        'business_id': 10,
        'name': 'John Doe',
        'email': 'john@example.com',
        'phone': '+573001234567',
        'dni': '1234567890',
        'created_at': '2026-01-15T10:30:00Z',
        'updated_at': '2026-01-16T12:00:00Z',
        'order_count': 25,
        'total_spent': 150000.50,
        'last_order_at': '2026-03-10T08:00:00Z',
      };

      final detail = CustomerDetail.fromJson(json);

      expect(detail.id, 1);
      expect(detail.businessId, 10);
      expect(detail.name, 'John Doe');
      expect(detail.email, 'john@example.com');
      expect(detail.phone, '+573001234567');
      expect(detail.dni, '1234567890');
      expect(detail.createdAt, '2026-01-15T10:30:00Z');
      expect(detail.updatedAt, '2026-01-16T12:00:00Z');
      expect(detail.orderCount, 25);
      expect(detail.totalSpent, 150000.50);
      expect(detail.lastOrderAt, '2026-03-10T08:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final detail = CustomerDetail.fromJson(json);

      expect(detail.id, 0);
      expect(detail.businessId, 0);
      expect(detail.name, '');
      expect(detail.phone, '');
      expect(detail.orderCount, 0);
      expect(detail.totalSpent, 0.0);
      expect(detail.lastOrderAt, isNull);
    });

    test('fromJson converts total_spent integer to double', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'name': 'Test',
        'phone': '123',
        'created_at': '',
        'updated_at': '',
        'order_count': 0,
        'total_spent': 5000,
      };

      final detail = CustomerDetail.fromJson(json);

      expect(detail.totalSpent, 5000.0);
      expect(detail.totalSpent, isA<double>());
    });

    test('fromJson handles null lastOrderAt', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'name': 'Test',
        'phone': '123',
        'created_at': '',
        'updated_at': '',
        'order_count': 3,
        'total_spent': 100.0,
        'last_order_at': null,
      };

      final detail = CustomerDetail.fromJson(json);

      expect(detail.lastOrderAt, isNull);
    });
  });

  group('GetCustomersParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetCustomersParams(
        page: 2,
        pageSize: 15,
        search: 'John',
        businessId: 5,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 15);
      expect(queryParams['search'], 'John');
      expect(queryParams['business_id'], 5);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetCustomersParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams.containsKey('page'), true);
      expect(queryParams.containsKey('page_size'), false);
      expect(queryParams.containsKey('search'), false);
      expect(queryParams.containsKey('business_id'), false);
    });

    test('toQueryParams excludes empty search string', () {
      final params = GetCustomersParams(search: '');

      final queryParams = params.toQueryParams();

      expect(queryParams.containsKey('search'), false);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetCustomersParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });

  group('CreateCustomerDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = CreateCustomerDTO(
        name: 'John Doe',
        email: 'john@example.com',
        phone: '+573001234567',
        dni: '1234567890',
      );

      final json = dto.toJson();

      expect(json['name'], 'John Doe');
      expect(json['email'], 'john@example.com');
      expect(json['phone'], '+573001234567');
      expect(json['dni'], '1234567890');
    });

    test('toJson includes only name when optional fields are null', () {
      final dto = CreateCustomerDTO(name: 'Jane');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['name'], 'Jane');
      expect(json.containsKey('email'), false);
      expect(json.containsKey('phone'), false);
      expect(json.containsKey('dni'), false);
    });
  });

  group('UpdateCustomerDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateCustomerDTO(
        name: 'Updated Name',
        email: 'updated@example.com',
        phone: '+573009876543',
        dni: '9876543210',
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated Name');
      expect(json['email'], 'updated@example.com');
      expect(json['phone'], '+573009876543');
      expect(json['dni'], '9876543210');
    });

    test('toJson includes only name when optional fields are null', () {
      final dto = UpdateCustomerDTO(name: 'Only Name');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['name'], 'Only Name');
      expect(json.containsKey('email'), false);
      expect(json.containsKey('phone'), false);
      expect(json.containsKey('dni'), false);
    });
  });
}
