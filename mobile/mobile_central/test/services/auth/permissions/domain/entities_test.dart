import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/permissions/domain/entities.dart';

void main() {
  group('Permission', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'resource': 'orders',
        'action': 'read',
        'description': 'Can read orders',
        'scope_id': 2,
        'business_type_id': 3,
        'business_type_name': 'E-commerce',
      };

      final permission = Permission.fromJson(json);

      expect(permission.id, 1);
      expect(permission.resource, 'orders');
      expect(permission.action, 'read');
      expect(permission.description, 'Can read orders');
      expect(permission.scopeId, 2);
      expect(permission.businessTypeId, 3);
      expect(permission.businessTypeName, 'E-commerce');
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final permission = Permission.fromJson(json);

      expect(permission.id, 0);
    });

    test('fromJson handles null optional fields', () {
      final json = {'id': 5};

      final permission = Permission.fromJson(json);

      expect(permission.resource, isNull);
      expect(permission.action, isNull);
      expect(permission.description, isNull);
      expect(permission.scopeId, isNull);
      expect(permission.businessTypeId, isNull);
      expect(permission.businessTypeName, isNull);
    });
  });

  group('GetPermissionsParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetPermissionsParams(
        page: 1,
        pageSize: 20,
        businessTypeId: 3,
        name: 'orders',
        scopeId: 2,
        resource: 'products',
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 1);
      expect(queryParams['page_size'], 20);
      expect(queryParams['business_type_id'], 3);
      expect(queryParams['name'], 'orders');
      expect(queryParams['scope_id'], 2);
      expect(queryParams['resource'], 'products');
    });

    test('toQueryParams excludes null fields', () {
      final params = GetPermissionsParams(page: 2);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams['page'], 2);
      expect(queryParams.containsKey('page_size'), false);
    });

    test('toQueryParams excludes empty name', () {
      final params = GetPermissionsParams(name: '');

      final queryParams = params.toQueryParams();

      expect(queryParams.containsKey('name'), false);
    });

    test('toQueryParams excludes empty resource', () {
      final params = GetPermissionsParams(resource: '');

      final queryParams = params.toQueryParams();

      expect(queryParams.containsKey('resource'), false);
    });

    test('toQueryParams returns empty map when all null', () {
      final params = GetPermissionsParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });

  group('CreatePermissionDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = CreatePermissionDTO(
        resource: 'orders',
        action: 'write',
        description: 'Can write orders',
        scopeId: 1,
        businessTypeId: 2,
      );

      final json = dto.toJson();

      expect(json['resource'], 'orders');
      expect(json['action'], 'write');
      expect(json['description'], 'Can write orders');
      expect(json['scope_id'], 1);
      expect(json['business_type_id'], 2);
    });

    test('toJson includes only required fields when optionals are null', () {
      final dto = CreatePermissionDTO(
        resource: 'products',
        action: 'read',
      );

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json['resource'], 'products');
      expect(json['action'], 'read');
      expect(json.containsKey('description'), false);
      expect(json.containsKey('scope_id'), false);
      expect(json.containsKey('business_type_id'), false);
    });
  });

  group('UpdatePermissionDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdatePermissionDTO(
        resource: 'orders',
        action: 'delete',
        description: 'Updated desc',
        scopeId: 3,
        businessTypeId: 4,
      );

      final json = dto.toJson();

      expect(json['resource'], 'orders');
      expect(json['action'], 'delete');
      expect(json['description'], 'Updated desc');
      expect(json['scope_id'], 3);
      expect(json['business_type_id'], 4);
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdatePermissionDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdatePermissionDTO(action: 'write');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['action'], 'write');
    });
  });

  group('BulkCreatePermissionsDTO', () {
    test('toJson produces correct structure with multiple permissions', () {
      final dto = BulkCreatePermissionsDTO(
        permissions: [
          CreatePermissionDTO(resource: 'orders', action: 'read'),
          CreatePermissionDTO(
            resource: 'products',
            action: 'write',
            description: 'Write products',
          ),
        ],
      );

      final json = dto.toJson();

      expect(json['permissions'], isList);
      final perms = json['permissions'] as List;
      expect(perms.length, 2);
      expect(perms[0]['resource'], 'orders');
      expect(perms[0]['action'], 'read');
      expect(perms[1]['resource'], 'products');
      expect(perms[1]['description'], 'Write products');
    });

    test('toJson handles empty permissions list', () {
      final dto = BulkCreatePermissionsDTO(permissions: []);

      final json = dto.toJson();

      expect(json['permissions'], isEmpty);
    });
  });
}
