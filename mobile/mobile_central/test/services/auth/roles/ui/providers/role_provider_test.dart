import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/roles/app/use_cases.dart';
import 'package:mobile_central/services/auth/roles/domain/entities.dart';
import 'package:mobile_central/services/auth/roles/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockRoleRepository implements IRoleRepository {
  PaginatedResponse<Role>? getRolesResult;
  Role? createRoleResult;
  Role? updateRoleResult;
  RolePermissionsResponse? getRolePermissionsResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;
  int? capturedAssignRoleId;
  AssignPermissionsDTO? capturedAssignData;

  @override
  Future<PaginatedResponse<Role>> getRoles(GetRolesParams? params) async {
    calls.add('getRoles');
    if (errorToThrow != null) throw errorToThrow!;
    return getRolesResult!;
  }

  @override
  Future<Role> getRoleById(int id) async {
    calls.add('getRoleById');
    if (errorToThrow != null) throw errorToThrow!;
    return Role(id: id, name: 'Test', isSystem: false);
  }

  @override
  Future<List<Role>> getRolesByScope(int scopeId) async {
    calls.add('getRolesByScope');
    return [];
  }

  @override
  Future<List<Role>> getRolesByLevel(int level) async {
    calls.add('getRolesByLevel');
    return [];
  }

  @override
  Future<List<Role>> getSystemRoles() async {
    calls.add('getSystemRoles');
    return [];
  }

  @override
  Future<Role> createRole(CreateRoleDTO data) async {
    calls.add('createRole');
    if (errorToThrow != null) throw errorToThrow!;
    return createRoleResult!;
  }

  @override
  Future<Role> updateRole(int id, UpdateRoleDTO data) async {
    calls.add('updateRole');
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
    if (errorToThrow != null) throw errorToThrow!;
    return getRolePermissionsResult!;
  }

  @override
  Future<void> removePermissionFromRole(int roleId, int permissionId) async {
    calls.add('removePermissionFromRole');
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---
// The real RoleProvider creates use cases internally via a getter that depends
// on ApiClient. We create a testable version that accepts use cases directly.

class TestableRoleProvider {
  final RoleUseCases _useCases;

  List<Role> _roles = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableRoleProvider(this._useCases);

  List<Role> get roles => _roles;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchRoles({GetRolesParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final queryParams =
          params ?? GetRolesParams(page: 1, pageSize: 20);
      final response = await _useCases.getRoles(queryParams);
      _roles = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Role?> createRole(CreateRoleDTO data) async {
    try {
      return await _useCases.createRole(data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateRole(int id, UpdateRoleDTO data) async {
    try {
      await _useCases.updateRole(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteRole(int id) async {
    try {
      await _useCases.deleteRole(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<RolePermissionsResponse?> getRolePermissions(int roleId) async {
    try {
      return await _useCases.getRolePermissions(roleId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> assignPermissions(
      int roleId, AssignPermissionsDTO data) async {
    try {
      await _useCases.assignPermissions(roleId, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }
}

// --- Helpers ---

Pagination _makePagination({
  int currentPage = 1,
  int total = 5,
  int lastPage = 1,
}) {
  return Pagination(
    currentPage: currentPage,
    perPage: 20,
    total: total,
    lastPage: lastPage,
    hasNext: currentPage < lastPage,
    hasPrev: currentPage > 1,
  );
}

Role _makeRole({int id = 1, String name = 'TestRole'}) {
  return Role(id: id, name: name, isSystem: false);
}

// --- Tests ---

void main() {
  late MockRoleRepository mockRepo;
  late RoleUseCases useCases;
  late TestableRoleProvider provider;

  setUp(() {
    mockRepo = MockRoleRepository();
    useCases = RoleUseCases(mockRepo);
    provider = TestableRoleProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty roles list', () {
      expect(provider.roles, isEmpty);
    });

    test('starts with null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchRoles', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getRolesResult = PaginatedResponse<Role>(
        data: [_makeRole()],
        pagination: _makePagination(),
      );

      await provider.fetchRoles();

      // Should have notified twice: once for loading start, once for loading end
      expect(provider.notifications.length, 2);
    });

    test('populates roles and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getRolesResult = PaginatedResponse<Role>(
        data: [_makeRole(id: 1), _makeRole(id: 2, name: 'Admin')],
        pagination: pagination,
      );

      await provider.fetchRoles();

      expect(provider.roles.length, 2);
      expect(provider.roles[0].id, 1);
      expect(provider.roles[1].name, 'Admin');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchRoles();

      expect(provider.error, contains('Server error'));
      expect(provider.roles, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      // First call fails
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchRoles();
      expect(provider.error, isNotNull);

      // Second call succeeds
      mockRepo.errorToThrow = null;
      mockRepo.getRolesResult = PaginatedResponse<Role>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchRoles();

      expect(provider.error, isNull);
    });

    test('uses default params when none provided', () async {
      mockRepo.getRolesResult = PaginatedResponse<Role>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchRoles();

      expect(mockRepo.calls, contains('getRoles'));
    });

    test('uses custom params when provided', () async {
      mockRepo.getRolesResult = PaginatedResponse<Role>(
        data: [],
        pagination: _makePagination(),
      );
      final params = GetRolesParams(page: 3, pageSize: 50, name: 'Admin');

      await provider.fetchRoles(params: params);

      expect(mockRepo.calls, contains('getRoles'));
    });
  });

  group('createRole', () {
    test('returns created role on success', () async {
      final dto = CreateRoleDTO(name: 'NewRole');
      mockRepo.createRoleResult = _makeRole(id: 10, name: 'NewRole');

      final result = await provider.createRole(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.name, 'NewRole');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateRoleDTO(name: 'Fail');

      final result = await provider.createRole(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateRole', () {
    test('returns true on success', () async {
      final dto = UpdateRoleDTO(name: 'Updated');
      mockRepo.updateRoleResult = _makeRole(id: 5, name: 'Updated');

      final result = await provider.updateRole(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateRoleDTO(name: 'Fail');

      final result = await provider.updateRole(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteRole', () {
    test('returns true on success', () async {
      final result = await provider.deleteRole(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteRole(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });

  group('getRolePermissions', () {
    test('returns permissions response on success', () async {
      mockRepo.getRolePermissionsResult = RolePermissionsResponse(
        roleId: 1,
        roleName: 'Admin',
        permissions: [
          RolePermission(id: 10, resource: 'orders', action: 'read'),
        ],
        count: 1,
      );

      final result = await provider.getRolePermissions(1);

      expect(result, isNotNull);
      expect(result!.roleId, 1);
      expect(result.permissions.length, 1);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Permission fetch failed');

      final result = await provider.getRolePermissions(1);

      expect(result, isNull);
      expect(provider.error, contains('Permission fetch failed'));
    });
  });

  group('assignPermissions', () {
    test('returns true on success', () async {
      final dto = AssignPermissionsDTO(permissionIds: [1, 2]);

      final result = await provider.assignPermissions(10, dto);

      expect(result, true);
      expect(mockRepo.capturedAssignRoleId, 10);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Assign failed');
      final dto = AssignPermissionsDTO(permissionIds: [1]);

      final result = await provider.assignPermissions(10, dto);

      expect(result, false);
      expect(provider.error, contains('Assign failed'));
    });
  });
}
