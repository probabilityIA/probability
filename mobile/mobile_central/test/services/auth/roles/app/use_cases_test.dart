import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/roles/app/use_cases.dart';
import 'package:mobile_central/services/auth/roles/domain/entities.dart';
import 'package:mobile_central/services/auth/roles/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockRoleRepository implements IRoleRepository {
  // Track calls
  final List<String> calls = [];

  // Configurable return values
  PaginatedResponse<Role>? getRolesResult;
  Role? getRoleByIdResult;
  List<Role>? getRolesByScopeResult;
  List<Role>? getRolesByLevelResult;
  List<Role>? getSystemRolesResult;
  Role? createRoleResult;
  Role? updateRoleResult;
  RolePermissionsResponse? getRolePermissionsResult;

  // Configurable errors
  Exception? errorToThrow;

  // Captured arguments
  GetRolesParams? capturedGetRolesParams;
  int? capturedId;
  int? capturedScopeId;
  int? capturedLevel;
  CreateRoleDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateRoleDTO? capturedUpdateData;
  int? capturedDeleteId;
  int? capturedAssignRoleId;
  AssignPermissionsDTO? capturedAssignData;
  int? capturedPermissionsRoleId;
  int? capturedRemoveRoleId;
  int? capturedRemovePermissionId;

  @override
  Future<PaginatedResponse<Role>> getRoles(GetRolesParams? params) async {
    calls.add('getRoles');
    capturedGetRolesParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getRolesResult!;
  }

  @override
  Future<Role> getRoleById(int id) async {
    calls.add('getRoleById');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getRoleByIdResult!;
  }

  @override
  Future<List<Role>> getRolesByScope(int scopeId) async {
    calls.add('getRolesByScope');
    capturedScopeId = scopeId;
    if (errorToThrow != null) throw errorToThrow!;
    return getRolesByScopeResult!;
  }

  @override
  Future<List<Role>> getRolesByLevel(int level) async {
    calls.add('getRolesByLevel');
    capturedLevel = level;
    if (errorToThrow != null) throw errorToThrow!;
    return getRolesByLevelResult!;
  }

  @override
  Future<List<Role>> getSystemRoles() async {
    calls.add('getSystemRoles');
    if (errorToThrow != null) throw errorToThrow!;
    return getSystemRolesResult!;
  }

  @override
  Future<Role> createRole(CreateRoleDTO data) async {
    calls.add('createRole');
    capturedCreateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createRoleResult!;
  }

  @override
  Future<Role> updateRole(int id, UpdateRoleDTO data) async {
    calls.add('updateRole');
    capturedUpdateId = id;
    capturedUpdateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateRoleResult!;
  }

  @override
  Future<void> deleteRole(int id) async {
    calls.add('deleteRole');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> assignPermissions(int roleId, AssignPermissionsDTO data) async {
    calls.add('assignPermissions');
    capturedAssignRoleId = roleId;
    capturedAssignData = data;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<RolePermissionsResponse> getRolePermissions(int roleId) async {
    calls.add('getRolePermissions');
    capturedPermissionsRoleId = roleId;
    if (errorToThrow != null) throw errorToThrow!;
    return getRolePermissionsResult!;
  }

  @override
  Future<void> removePermissionFromRole(int roleId, int permissionId) async {
    calls.add('removePermissionFromRole');
    capturedRemoveRoleId = roleId;
    capturedRemovePermissionId = permissionId;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Helpers ---

Role _makeRole({int id = 1, String name = 'TestRole'}) {
  return Role(id: id, name: name, isSystem: false);
}

Pagination _makePagination() {
  return Pagination(
    currentPage: 1,
    perPage: 10,
    total: 1,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

// --- Tests ---

void main() {
  late MockRoleRepository mockRepo;
  late RoleUseCases useCases;

  setUp(() {
    mockRepo = MockRoleRepository();
    useCases = RoleUseCases(mockRepo);
  });

  group('getRoles', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Role>(
        data: [_makeRole()],
        pagination: _makePagination(),
      );
      mockRepo.getRolesResult = expected;
      final params = GetRolesParams(page: 1, pageSize: 10);

      final result = await useCases.getRoles(params);

      expect(result.data.length, 1);
      expect(result.data[0].name, 'TestRole');
      expect(mockRepo.calls, ['getRoles']);
      expect(mockRepo.capturedGetRolesParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getRolesResult = PaginatedResponse<Role>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getRoles(null);

      expect(mockRepo.capturedGetRolesParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getRoles(null), throwsException);
    });
  });

  group('getRoleById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getRoleByIdResult = _makeRole(id: 42, name: 'Found');

      final result = await useCases.getRoleById(42);

      expect(result.id, 42);
      expect(result.name, 'Found');
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getRoleById']);
    });
  });

  group('getRolesByScope', () {
    test('delegates to repository with correct scopeId', () async {
      mockRepo.getRolesByScopeResult = [_makeRole(), _makeRole(id: 2)];

      final result = await useCases.getRolesByScope(5);

      expect(result.length, 2);
      expect(mockRepo.capturedScopeId, 5);
      expect(mockRepo.calls, ['getRolesByScope']);
    });
  });

  group('getRolesByLevel', () {
    test('delegates to repository with correct level', () async {
      mockRepo.getRolesByLevelResult = [_makeRole()];

      final result = await useCases.getRolesByLevel(10);

      expect(result.length, 1);
      expect(mockRepo.capturedLevel, 10);
      expect(mockRepo.calls, ['getRolesByLevel']);
    });
  });

  group('getSystemRoles', () {
    test('delegates to repository', () async {
      mockRepo.getSystemRolesResult = [
        Role(id: 1, name: 'SuperAdmin', isSystem: true),
      ];

      final result = await useCases.getSystemRoles();

      expect(result.length, 1);
      expect(result[0].isSystem, true);
      expect(mockRepo.calls, ['getSystemRoles']);
    });
  });

  group('createRole', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateRoleDTO(name: 'NewRole', code: 'NR');
      mockRepo.createRoleResult = _makeRole(id: 99, name: 'NewRole');

      final result = await useCases.createRole(dto);

      expect(result.id, 99);
      expect(result.name, 'NewRole');
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createRole']);
    });
  });

  group('updateRole', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateRoleDTO(name: 'Updated');
      mockRepo.updateRoleResult = _makeRole(id: 5, name: 'Updated');

      final result = await useCases.updateRole(5, dto);

      expect(result.id, 5);
      expect(result.name, 'Updated');
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateRole']);
    });
  });

  group('deleteRole', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteRole(7);

      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deleteRole']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.deleteRole(7), throwsException);
    });
  });

  group('assignPermissions', () {
    test('delegates to repository with correct roleId and data', () async {
      final dto = AssignPermissionsDTO(permissionIds: [1, 2, 3]);

      await useCases.assignPermissions(10, dto);

      expect(mockRepo.capturedAssignRoleId, 10);
      expect(mockRepo.capturedAssignData, dto);
      expect(mockRepo.calls, ['assignPermissions']);
    });
  });

  group('getRolePermissions', () {
    test('delegates to repository with correct roleId', () async {
      mockRepo.getRolePermissionsResult = RolePermissionsResponse(
        roleId: 10,
        roleName: 'Admin',
        permissions: [
          RolePermission(id: 1, resource: 'orders', action: 'read'),
        ],
        count: 1,
      );

      final result = await useCases.getRolePermissions(10);

      expect(result.roleId, 10);
      expect(result.permissions.length, 1);
      expect(mockRepo.capturedPermissionsRoleId, 10);
      expect(mockRepo.calls, ['getRolePermissions']);
    });
  });

  group('removePermissionFromRole', () {
    test('delegates to repository with correct roleId and permissionId',
        () async {
      await useCases.removePermissionFromRole(10, 5);

      expect(mockRepo.capturedRemoveRoleId, 10);
      expect(mockRepo.capturedRemovePermissionId, 5);
      expect(mockRepo.calls, ['removePermissionFromRole']);
    });
  });
}
