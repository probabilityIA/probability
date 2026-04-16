import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/permissions/app/use_cases.dart';
import 'package:mobile_central/services/auth/permissions/domain/entities.dart';
import 'package:mobile_central/services/auth/permissions/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockPermissionRepository implements IPermissionRepository {
  final List<String> calls = [];

  PaginatedResponse<Permission>? getPermissionsResult;
  Permission? getPermissionByIdResult;
  List<Permission>? getPermissionsByScopeResult;
  List<Permission>? getPermissionsByResourceResult;
  Permission? createPermissionResult;
  Permission? updatePermissionResult;

  Exception? errorToThrow;

  GetPermissionsParams? capturedGetParams;
  int? capturedId;
  int? capturedScopeId;
  String? capturedResource;
  CreatePermissionDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdatePermissionDTO? capturedUpdateData;
  int? capturedDeleteId;
  BulkCreatePermissionsDTO? capturedBulkData;

  @override
  Future<PaginatedResponse<Permission>> getPermissions(
      GetPermissionsParams? params) async {
    calls.add('getPermissions');
    capturedGetParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getPermissionsResult!;
  }

  @override
  Future<Permission> getPermissionById(int id) async {
    calls.add('getPermissionById');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getPermissionByIdResult!;
  }

  @override
  Future<List<Permission>> getPermissionsByScope(int scopeId) async {
    calls.add('getPermissionsByScope');
    capturedScopeId = scopeId;
    if (errorToThrow != null) throw errorToThrow!;
    return getPermissionsByScopeResult!;
  }

  @override
  Future<List<Permission>> getPermissionsByResource(String resource) async {
    calls.add('getPermissionsByResource');
    capturedResource = resource;
    if (errorToThrow != null) throw errorToThrow!;
    return getPermissionsByResourceResult!;
  }

  @override
  Future<Permission> createPermission(CreatePermissionDTO data) async {
    calls.add('createPermission');
    capturedCreateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createPermissionResult!;
  }

  @override
  Future<Permission> updatePermission(int id, UpdatePermissionDTO data) async {
    calls.add('updatePermission');
    capturedUpdateId = id;
    capturedUpdateData = data;
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
    capturedBulkData = data;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Helpers ---

Permission _makePermission({int id = 1, String resource = 'orders'}) {
  return Permission(id: id, resource: resource, action: 'read');
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
  late MockPermissionRepository mockRepo;
  late PermissionUseCases useCases;

  setUp(() {
    mockRepo = MockPermissionRepository();
    useCases = PermissionUseCases(mockRepo);
  });

  group('getPermissions', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Permission>(
        data: [_makePermission()],
        pagination: _makePagination(),
      );
      mockRepo.getPermissionsResult = expected;
      final params = GetPermissionsParams(page: 1, pageSize: 10);

      final result = await useCases.getPermissions(params);

      expect(result.data.length, 1);
      expect(result.data[0].resource, 'orders');
      expect(mockRepo.calls, ['getPermissions']);
      expect(mockRepo.capturedGetParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getPermissionsResult = PaginatedResponse<Permission>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getPermissions(null);

      expect(mockRepo.capturedGetParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getPermissions(null), throwsException);
    });
  });

  group('getPermissionById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getPermissionByIdResult =
          _makePermission(id: 42, resource: 'products');

      final result = await useCases.getPermissionById(42);

      expect(result.id, 42);
      expect(result.resource, 'products');
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getPermissionById']);
    });
  });

  group('getPermissionsByScope', () {
    test('delegates to repository with correct scopeId', () async {
      mockRepo.getPermissionsByScopeResult = [
        _makePermission(id: 1),
        _makePermission(id: 2),
      ];

      final result = await useCases.getPermissionsByScope(5);

      expect(result.length, 2);
      expect(mockRepo.capturedScopeId, 5);
      expect(mockRepo.calls, ['getPermissionsByScope']);
    });
  });

  group('getPermissionsByResource', () {
    test('delegates to repository with correct resource', () async {
      mockRepo.getPermissionsByResourceResult = [
        _makePermission(id: 1, resource: 'orders'),
      ];

      final result = await useCases.getPermissionsByResource('orders');

      expect(result.length, 1);
      expect(result[0].resource, 'orders');
      expect(mockRepo.capturedResource, 'orders');
      expect(mockRepo.calls, ['getPermissionsByResource']);
    });
  });

  group('createPermission', () {
    test('delegates to repository with correct data', () async {
      final dto = CreatePermissionDTO(resource: 'orders', action: 'write');
      mockRepo.createPermissionResult =
          Permission(id: 99, resource: 'orders', action: 'write');

      final result = await useCases.createPermission(dto);

      expect(result.id, 99);
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createPermission']);
    });
  });

  group('updatePermission', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdatePermissionDTO(action: 'delete');
      mockRepo.updatePermissionResult =
          Permission(id: 5, resource: 'orders', action: 'delete');

      final result = await useCases.updatePermission(5, dto);

      expect(result.id, 5);
      expect(result.action, 'delete');
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updatePermission']);
    });
  });

  group('deletePermission', () {
    test('delegates to repository with correct id', () async {
      await useCases.deletePermission(7);

      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deletePermission']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.deletePermission(7), throwsException);
    });
  });

  group('createPermissionsBulk', () {
    test('delegates to repository with correct data', () async {
      final dto = BulkCreatePermissionsDTO(
        permissions: [
          CreatePermissionDTO(resource: 'orders', action: 'read'),
          CreatePermissionDTO(resource: 'orders', action: 'write'),
        ],
      );

      await useCases.createPermissionsBulk(dto);

      expect(mockRepo.capturedBulkData, dto);
      expect(mockRepo.calls, ['createPermissionsBulk']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Bulk create failed');
      final dto = BulkCreatePermissionsDTO(permissions: []);

      expect(() => useCases.createPermissionsBulk(dto), throwsException);
    });
  });
}
