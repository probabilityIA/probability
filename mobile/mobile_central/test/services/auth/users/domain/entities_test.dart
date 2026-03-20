import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/users/domain/entities.dart';

void main() {
  group('BusinessRoleAssignment', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'business_id': 5,
        'business_name': 'Acme Corp',
        'role_id': 3,
        'role_name': 'Admin',
      };

      final assignment = BusinessRoleAssignment.fromJson(json);

      expect(assignment.businessId, 5);
      expect(assignment.businessName, 'Acme Corp');
      expect(assignment.roleId, 3);
      expect(assignment.roleName, 'Admin');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'business_id': 1,
        'role_id': 2,
      };

      final assignment = BusinessRoleAssignment.fromJson(json);

      expect(assignment.businessId, 1);
      expect(assignment.businessName, isNull);
      expect(assignment.roleId, 2);
      expect(assignment.roleName, isNull);
    });

    test('fromJson defaults ids to 0 when missing', () {
      final json = <String, dynamic>{};

      final assignment = BusinessRoleAssignment.fromJson(json);

      expect(assignment.businessId, 0);
      expect(assignment.roleId, 0);
    });
  });

  group('User', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'name': 'John Doe',
        'email': 'john@example.com',
        'phone': '+573001234567',
        'avatar_url': 'https://example.com/avatar.png',
        'is_active': true,
        'is_super_user': false,
        'scope_id': 10,
        'business_role_assignments': [
          {
            'business_id': 1,
            'business_name': 'Biz A',
            'role_id': 2,
            'role_name': 'Editor',
          },
          {
            'business_id': 3,
            'business_name': 'Biz B',
            'role_id': 4,
            'role_name': 'Viewer',
          },
        ],
      };

      final user = User.fromJson(json);

      expect(user.id, 42);
      expect(user.name, 'John Doe');
      expect(user.email, 'john@example.com');
      expect(user.phone, '+573001234567');
      expect(user.avatarUrl, 'https://example.com/avatar.png');
      expect(user.isActive, true);
      expect(user.isSuperUser, false);
      expect(user.scopeId, 10);
      expect(user.businessRoleAssignments, hasLength(2));
      expect(user.businessRoleAssignments[0].businessId, 1);
      expect(user.businessRoleAssignments[0].roleName, 'Editor');
      expect(user.businessRoleAssignments[1].businessId, 3);
      expect(user.businessRoleAssignments[1].roleName, 'Viewer');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'Jane',
        'email': 'jane@example.com',
        'is_active': true,
        'is_super_user': true,
      };

      final user = User.fromJson(json);

      expect(user.id, 1);
      expect(user.name, 'Jane');
      expect(user.email, 'jane@example.com');
      expect(user.phone, isNull);
      expect(user.avatarUrl, isNull);
      expect(user.isActive, true);
      expect(user.isSuperUser, true);
      expect(user.scopeId, isNull);
      expect(user.businessRoleAssignments, isEmpty);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final user = User.fromJson(json);

      expect(user.id, 0);
      expect(user.name, '');
      expect(user.email, '');
      expect(user.isActive, false);
      expect(user.isSuperUser, false);
      expect(user.businessRoleAssignments, isEmpty);
    });

    test('fromJson handles null business_role_assignments', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'email': 'test@test.com',
        'is_active': true,
        'is_super_user': false,
        'business_role_assignments': null,
      };

      final user = User.fromJson(json);

      expect(user.businessRoleAssignments, isEmpty);
    });

    test('fromJson handles empty business_role_assignments list', () {
      final json = {
        'id': 1,
        'name': 'Test',
        'email': 'test@test.com',
        'is_active': true,
        'is_super_user': false,
        'business_role_assignments': [],
      };

      final user = User.fromJson(json);

      expect(user.businessRoleAssignments, isEmpty);
    });
  });

  group('GetUsersParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetUsersParams(
        page: 2,
        pageSize: 25,
        name: 'John',
        email: 'john@example.com',
        phone: '+57300',
        isActive: true,
        roleId: 3,
        businessId: 5,
      );

      final query = params.toQueryParams();

      expect(query['page'], 2);
      expect(query['page_size'], 25);
      expect(query['name'], 'John');
      expect(query['email'], 'john@example.com');
      expect(query['phone'], '+57300');
      expect(query['is_active'], true);
      expect(query['role_id'], 3);
      expect(query['business_id'], 5);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetUsersParams();

      final query = params.toQueryParams();

      expect(query, isEmpty);
    });

    test('toQueryParams excludes empty strings for name, email, phone', () {
      final params = GetUsersParams(
        page: 1,
        name: '',
        email: '',
        phone: '',
      );

      final query = params.toQueryParams();

      expect(query.containsKey('page'), true);
      expect(query.containsKey('name'), false);
      expect(query.containsKey('email'), false);
      expect(query.containsKey('phone'), false);
    });

    test('toQueryParams includes non-empty string filters', () {
      final params = GetUsersParams(
        name: 'A',
        email: 'b@c.com',
        phone: '1',
      );

      final query = params.toQueryParams();

      expect(query['name'], 'A');
      expect(query['email'], 'b@c.com');
      expect(query['phone'], '1');
    });

    test('toQueryParams includes isActive when false', () {
      final params = GetUsersParams(isActive: false);

      final query = params.toQueryParams();

      expect(query['is_active'], false);
    });
  });

  group('CreateUserDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateUserDTO(
        name: 'New User',
        email: 'new@example.com',
      );

      final json = dto.toJson();

      expect(json['name'], 'New User');
      expect(json['email'], 'new@example.com');
      expect(json['is_active'], true); // default
      expect(json.containsKey('phone'), false);
      expect(json.containsKey('scope_id'), false);
      expect(json.containsKey('business_ids'), false);
    });

    test('toJson includes all optional fields when provided', () {
      final dto = CreateUserDTO(
        name: 'Full User',
        email: 'full@example.com',
        phone: '+573009999999',
        isActive: false,
        scopeId: 7,
        businessIds: [1, 2, 3],
      );

      final json = dto.toJson();

      expect(json['name'], 'Full User');
      expect(json['email'], 'full@example.com');
      expect(json['phone'], '+573009999999');
      expect(json['is_active'], false);
      expect(json['scope_id'], 7);
      expect(json['business_ids'], [1, 2, 3]);
    });

    test('toJson excludes null phone, scopeId, and businessIds', () {
      final dto = CreateUserDTO(
        name: 'Minimal',
        email: 'min@example.com',
        isActive: true,
      );

      final json = dto.toJson();

      expect(json.containsKey('phone'), false);
      expect(json.containsKey('scope_id'), false);
      expect(json.containsKey('business_ids'), false);
    });
  });

  group('UpdateUserDTO', () {
    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateUserDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateUserDTO(
        name: 'Updated Name',
        isActive: false,
      );

      final json = dto.toJson();

      expect(json['name'], 'Updated Name');
      expect(json['is_active'], false);
      expect(json.containsKey('email'), false);
      expect(json.containsKey('phone'), false);
      expect(json.containsKey('scope_id'), false);
      expect(json.containsKey('remove_avatar'), false);
    });

    test('toJson includes all fields when all provided', () {
      final dto = UpdateUserDTO(
        name: 'New Name',
        email: 'new@email.com',
        phone: '+573001111111',
        isActive: true,
        scopeId: 5,
        removeAvatar: true,
      );

      final json = dto.toJson();

      expect(json['name'], 'New Name');
      expect(json['email'], 'new@email.com');
      expect(json['phone'], '+573001111111');
      expect(json['is_active'], true);
      expect(json['scope_id'], 5);
      expect(json['remove_avatar'], true);
    });
  });

  group('RoleAssignment', () {
    test('toJson produces correct map', () {
      final assignment = RoleAssignment(businessId: 10, roleId: 20);

      final json = assignment.toJson();

      expect(json['business_id'], 10);
      expect(json['role_id'], 20);
      expect(json.length, 2);
    });
  });

  group('AssignRolesDTO', () {
    test('toJson serializes list of assignments', () {
      final dto = AssignRolesDTO(
        assignments: [
          RoleAssignment(businessId: 1, roleId: 10),
          RoleAssignment(businessId: 2, roleId: 20),
        ],
      );

      final json = dto.toJson();

      expect(json['assignments'], isList);
      final assignments = json['assignments'] as List;
      expect(assignments, hasLength(2));
      expect(assignments[0]['business_id'], 1);
      expect(assignments[0]['role_id'], 10);
      expect(assignments[1]['business_id'], 2);
      expect(assignments[1]['role_id'], 20);
    });

    test('toJson handles empty assignments list', () {
      final dto = AssignRolesDTO(assignments: []);

      final json = dto.toJson();

      expect(json['assignments'], isList);
      expect((json['assignments'] as List), isEmpty);
    });
  });
}
