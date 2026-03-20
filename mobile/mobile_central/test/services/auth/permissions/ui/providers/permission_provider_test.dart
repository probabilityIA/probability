import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/permissions/app/use_cases.dart';
import 'package:mobile_central/services/auth/permissions/domain/entities.dart';
import 'package:mobile_central/services/auth/permissions/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockPermissionRepository implements IPermissionRepository {
  PaginatedResponse<Permission>? getPermissionsResult;
  Permission? createPermissionResult;
  Permission? updatePermissionResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<Permission>> getPermissions(
      GetPermissionsParams? params) async {
    calls.add('getPermissions');
    if (errorToThrow != null) throw errorToThrow!;
    return getPermissionsResult!;
  }

  @override
  Future<Permission> getPermissionById(int id) async {
    calls.add('getPermissionById');
    if (errorToThrow != null) throw errorToThrow!;
    return Permission(id: id, resource: 'test', action: 'read');
  }

  @override
  Future<List<Permission>> getPermissionsByScope(int scopeId) async {
    calls.add('getPermissionsByScope');
    return [];
  }

  @override
  Future<List<Permission>> getPermissionsByResource(String resource) async {
    calls.add('getPermissionsByResource');
    return [];
  }

  @override
  Future<Permission> createPermission(CreatePermissionDTO data) async {
    calls.add('createPermission');
    if (errorToThrow != null) throw errorToThrow!;
    return createPermissionResult!;
  }

  @override
  Future<Permission> updatePermission(
      int id, UpdatePermissionDTO data) async {
    calls.add('updatePermission');
    if (errorToThrow != null) throw errorToThrow!;
    return updatePermissionResult!;
  }

  @override
  Future<void> deletePermission(int id) async {
    calls.add('deletePermission');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> createPermissionsBulk(BulkCreatePermissionsDTO data) async {
    calls.add('createPermissionsBulk');
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---

class TestablePermissionProvider {
  final PermissionUseCases _useCases;

  List<Permission> _permissions = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestablePermissionProvider(this._useCases);

  List<Permission> get permissions => _permissions;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchPermissions({GetPermissionsParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final queryParams =
          params ?? GetPermissionsParams(page: 1, pageSize: 20);
      final response = await _useCases.getPermissions(queryParams);
      _permissions = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Permission?> createPermission(CreatePermissionDTO data) async {
    try {
      return await _useCases.createPermission(data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updatePermission(int id, UpdatePermissionDTO data) async {
    try {
      await _useCases.updatePermission(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deletePermission(int id) async {
    try {
      await _useCases.deletePermission(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }
}

// --- Helpers ---

Pagination _makePagination({int total = 5}) {
  return Pagination(
    currentPage: 1,
    perPage: 20,
    total: total,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

Permission _makePermission({int id = 1, String resource = 'orders'}) {
  return Permission(id: id, resource: resource, action: 'read');
}

// --- Tests ---

void main() {
  late MockPermissionRepository mockRepo;
  late PermissionUseCases useCases;
  late TestablePermissionProvider provider;

  setUp(() {
    mockRepo = MockPermissionRepository();
    useCases = PermissionUseCases(mockRepo);
    provider = TestablePermissionProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty permissions list', () {
      expect(provider.permissions, isEmpty);
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

  group('fetchPermissions', () {
    test('notifies listeners twice (loading start and end)', () async {
      mockRepo.getPermissionsResult = PaginatedResponse<Permission>(
        data: [_makePermission()],
        pagination: _makePagination(),
      );

      await provider.fetchPermissions();

      expect(provider.notifications.length, 2);
    });

    test('populates permissions and pagination on success', () async {
      final pagination = _makePagination(total: 3);
      mockRepo.getPermissionsResult = PaginatedResponse<Permission>(
        data: [
          _makePermission(id: 1, resource: 'orders'),
          _makePermission(id: 2, resource: 'products'),
          _makePermission(id: 3, resource: 'shipments'),
        ],
        pagination: pagination,
      );

      await provider.fetchPermissions();

      expect(provider.permissions.length, 3);
      expect(provider.permissions[0].resource, 'orders');
      expect(provider.permissions[2].resource, 'shipments');
      expect(provider.pagination?.total, 3);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchPermissions();

      expect(provider.error, contains('Server error'));
      expect(provider.permissions, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchPermissions();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getPermissionsResult = PaginatedResponse<Permission>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchPermissions();

      expect(provider.error, isNull);
    });
  });

  group('createPermission', () {
    test('returns created permission on success', () async {
      final dto = CreatePermissionDTO(resource: 'orders', action: 'write');
      mockRepo.createPermissionResult =
          Permission(id: 10, resource: 'orders', action: 'write');

      final result = await provider.createPermission(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.action, 'write');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreatePermissionDTO(resource: 'fail', action: 'fail');

      final result = await provider.createPermission(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updatePermission', () {
    test('returns true on success', () async {
      final dto = UpdatePermissionDTO(action: 'delete');
      mockRepo.updatePermissionResult =
          Permission(id: 5, resource: 'orders', action: 'delete');

      final result = await provider.updatePermission(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdatePermissionDTO(action: 'fail');

      final result = await provider.updatePermission(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deletePermission', () {
    test('returns true on success', () async {
      final result = await provider.deletePermission(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deletePermission(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });
}
