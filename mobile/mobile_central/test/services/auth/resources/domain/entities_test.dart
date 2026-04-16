import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/resources/domain/entities.dart';

void main() {
  group('Resource', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'orders',
        'description': 'Order management',
        'business_type_id': 3,
        'business_type_name': 'E-commerce',
      };

      final resource = Resource.fromJson(json);

      expect(resource.id, 1);
      expect(resource.name, 'orders');
      expect(resource.description, 'Order management');
      expect(resource.businessTypeId, 3);
      expect(resource.businessTypeName, 'E-commerce');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final resource = Resource.fromJson(json);

      expect(resource.id, 0);
      expect(resource.name, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'name': 'products',
      };

      final resource = Resource.fromJson(json);

      expect(resource.description, isNull);
      expect(resource.businessTypeId, isNull);
      expect(resource.businessTypeName, isNull);
    });
  });

  group('GetResourcesParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetResourcesParams(
        page: 2,
        pageSize: 15,
        name: 'orders',
        description: 'manage',
        businessTypeId: 3,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 2);
      expect(queryParams['page_size'], 15);
      expect(queryParams['name'], 'orders');
      expect(queryParams['description'], 'manage');
      expect(queryParams['business_type_id'], 3);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetResourcesParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams['page'], 1);
      expect(queryParams.containsKey('page_size'), false);
      expect(queryParams.containsKey('name'), false);
    });

    test('toQueryParams excludes empty name', () {
      final params = GetResourcesParams(name: '');

      final queryParams = params.toQueryParams();

      expect(queryParams.containsKey('name'), false);
    });

    test('toQueryParams excludes empty description', () {
      final params = GetResourcesParams(description: '');

      final queryParams = params.toQueryParams();

      expect(queryParams.containsKey('description'), false);
    });

    test('toQueryParams returns empty map when all null', () {
      final params = GetResourcesParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });

  group('CreateResourceDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = CreateResourceDTO(
        name: 'shipments',
        description: 'Shipment management',
        businessTypeId: 2,
      );

      final json = dto.toJson();

      expect(json['name'], 'shipments');
      expect(json['description'], 'Shipment management');
      expect(json['business_type_id'], 2);
    });

    test('toJson includes only name when optionals are null', () {
      final dto = CreateResourceDTO(name: 'basic');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['name'], 'basic');
      expect(json.containsKey('description'), false);
      expect(json.containsKey('business_type_id'), false);
    });
  });

  group('UpdateResourceDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateResourceDTO(
        name: 'updated',
        description: 'New desc',
        businessTypeId: 4,
      );

      final json = dto.toJson();

      expect(json['name'], 'updated');
      expect(json['description'], 'New desc');
      expect(json['business_type_id'], 4);
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateResourceDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateResourceDTO(description: 'Only desc');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['description'], 'Only desc');
    });
  });
}
