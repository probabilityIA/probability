import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/roles/domain/entities.dart';

void main() {
  group('Role', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'Admin',
        'code': 'ADMIN',
        'description': 'Administrator role',
        'level': 10,
        'is_system': true,
        'scope_id': 2,
        'business_type_id': 3,
      };

      final role = Role.fromJson(json);

      expect(role.id, 1);
      expect(role.name, 'Admin');
      expect(role.code, 'ADMIN');
      expect(role.description, 'Administrator role');
      expect(role.level, 10);
      expect(role.isSystem, true);
      expect(role.scopeId, 2);
      expect(role.businessTypeId, 3);
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final role = Role.fromJson(json);

      expect(role.id, 0);
      expect(role.name, '');
      expect(role.isSystem, false);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 5,
        'name': 'Test',
        'is_system': false,
      };

      final role = Role.fromJson(json);

      expect(role.code, isNull);
      expect(role.description, isNull);
      expect(role.level, isNull);
      expect(role.scopeId, isNull);
      expect(role.businessTypeId, isNull);
    });
  });

  group('RolePermission', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 10,
        'resource': 'orders',
        'action': 'read',
        'description': 'Can read orders',
        'scope_id': 1,
      };

      final perm = RolePermission.fromJson(json);

      expect(perm.id, 10);
      expect(perm.resource, 'orders');
      expect(perm.action, 'read');
      expect(perm.description, 'Can read orders');
      expect(perm.scopeId, 1);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final perm = RolePermission.fromJson(json);

      expect(perm.id, 0);
      expect(perm.resource, isNull);
      expect(perm.action, isNull);
      expect(perm.description, isNull);
      expect(perm.scopeId, isNull);
    });
  });

  group('RolePermissionsResponse', () {
    test('fromJson parses all fields including permissions list', () {
      final json = {
        'role_id': 1,
        'role_name': 'Admin',
        'permissions': [
          {'id': 10, 'resource': 'orders', 'action': 'read'},
          {'id': 11, 'resource': 'products', 'action': 'write'},
        ],
        'count': 2,
      };

      final response = RolePermissionsResponse.fromJson(json);

      expect(response.roleId, 1);
      expect(response.roleName, 'Admin');
      expect(response.permissions.length, 2);
      expect(response.permissions[0].id, 10);
      expect(response.permissions[0].resource, 'orders');
      expect(response.permissions[1].id, 11);
      expect(response.permissions[1].resource, 'products');
      expect(response.count, 2);
    });

    test('fromJson handles null permissions list', () {
      final json = {
        'role_id': 1,
        'role_name': 'Empty Role',
      };

      final response = RolePermissionsResponse.fromJson(json);

      expect(response.permissions, isEmpty);
      expect(response.count, 0);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = RolePermissionsResponse.fromJson(json);

      expect(response.roleId, 0);
      expect(response.roleName, '');
      expect(response.permissions, isEmpty);
      expect(response.count, 0);
    });
  });

  group('GetRolesParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetRolesParams(
        page: 1,
        pageSize: 20,
        businessTypeId: 3,
        scopeId: 2,
        isSystem: true,
        name: 'Admin',
        level: 10,
      );

      final queryParams = params.toQueryParams();

      expect(queryParams['page'], 1);
      expect(queryParams['page_size'], 20);
      expect(queryParams['business_type_id'], 3);
      expect(queryParams['scope_id'], 2);
      expect(queryParams['is_system'], true);
      expect(queryParams['name'], 'Admin');
      expect(queryParams['level'], 10);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetRolesParams(page: 1);

      final queryParams = params.toQueryParams();

      expect(queryParams.length, 1);
      expect(queryParams.containsKey('page'), true);
      expect(queryParams.containsKey('page_size'), false);
      expect(queryParams.containsKey('name'), false);
    });

    test('toQueryParams excludes empty name', () {
      final params = GetRolesParams(name: '');

      final queryParams = params.toQueryParams();

      expect(queryParams.containsKey('name'), false);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetRolesParams();

      final queryParams = params.toQueryParams();

      expect(queryParams, isEmpty);
    });
  });

  group('CreateRoleDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = CreateRoleDTO(
        name: 'Manager',
        code: 'MGR',
        description: 'Manager role',
        level: 5,
        scopeId: 1,
        businessTypeId: 2,
      );

      final json = dto.toJson();

      expect(json['name'], 'Manager');
      expect(json['code'], 'MGR');
      expect(json['description'], 'Manager role');
      expect(json['level'], 5);
      expect(json['scope_id'], 1);
      expect(json['business_type_id'], 2);
    });

    test('toJson includes only name when optional fields are null', () {
      final dto = CreateRoleDTO(name: 'Basic');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['name'], 'Basic');
      expect(json.containsKey('code'), false);
      expect(json.containsKey('description'), false);
    });
  });

  group('UpdateRoleDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateRoleDTO(
        name: 'Updated',
        code: 'UPD',
        description: 'Updated description',
        level: 3,
        scopeId: 4,
        businessTypeId: 5,
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated');
      expect(json['code'], 'UPD');
      expect(json['description'], 'Updated description');
      expect(json['level'], 3);
      expect(json['scope_id'], 4);
      expect(json['business_type_id'], 5);
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateRoleDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateRoleDTO(name: 'NewName');

      final json = dto.toJson();

      expect(json.length, 1);
      expect(json['name'], 'NewName');
    });
  });

  group('AssignPermissionsDTO', () {
    test('toJson produces correct structure', () {
      final dto = AssignPermissionsDTO(permissionIds: [1, 2, 3]);

      final json = dto.toJson();

      expect(json['permission_ids'], [1, 2, 3]);
    });

    test('toJson handles empty list', () {
      final dto = AssignPermissionsDTO(permissionIds: []);

      final json = dto.toJson();

      expect(json['permission_ids'], isEmpty);
    });
  });
}
