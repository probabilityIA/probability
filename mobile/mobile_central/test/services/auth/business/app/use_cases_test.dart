import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/business/app/use_cases.dart';
import 'package:mobile_central/services/auth/business/domain/entities.dart';
import 'package:mobile_central/services/auth/business/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IBusinessRepository
// ---------------------------------------------------------------------------
class MockBusinessRepository implements IBusinessRepository {
  // Tracking calls
  final List<String> calls = [];

  // Configurable return values
  PaginatedResponse<Business>? getBusinessesResult;
  Business? getBusinessByIdResult;
  Business? createBusinessResult;
  Business? updateBusinessResult;
  List<BusinessSimple>? getBusinessesSimpleResult;
  List<ConfiguredResource>? getConfiguredResourcesResult;
  List<BusinessType>? getBusinessTypesResult;
  BusinessType? createBusinessTypeResult;
  BusinessType? updateBusinessTypeResult;

  // Configurable errors
  Exception? errorToThrow;

  // Captured arguments
  GetBusinessesParams? lastGetBusinessesParams;
  int? lastId;
  CreateBusinessDTO? lastCreateDTO;
  UpdateBusinessDTO? lastUpdateDTO;
  Map<String, dynamic>? lastBusinessTypeData;

  void _trackCall(String method) {
    calls.add(method);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<PaginatedResponse<Business>> getBusinesses(
      GetBusinessesParams? params) async {
    lastGetBusinessesParams = params;
    _trackCall('getBusinesses');
    return getBusinessesResult!;
  }

  @override
  Future<Business> getBusinessById(int id) async {
    lastId = id;
    _trackCall('getBusinessById');
    return getBusinessByIdResult!;
  }

  @override
  Future<Business> createBusiness(CreateBusinessDTO data) async {
    lastCreateDTO = data;
    _trackCall('createBusiness');
    return createBusinessResult!;
  }

  @override
  Future<Business> updateBusiness(int id, UpdateBusinessDTO data) async {
    lastId = id;
    lastUpdateDTO = data;
    _trackCall('updateBusiness');
    return updateBusinessResult!;
  }

  @override
  Future<void> deleteBusiness(int id) async {
    lastId = id;
    _trackCall('deleteBusiness');
  }

  @override
  Future<void> activateBusiness(int id) async {
    lastId = id;
    _trackCall('activateBusiness');
  }

  @override
  Future<void> deactivateBusiness(int id) async {
    lastId = id;
    _trackCall('deactivateBusiness');
  }

  @override
  Future<List<BusinessSimple>> getBusinessesSimple() async {
    _trackCall('getBusinessesSimple');
    return getBusinessesSimpleResult!;
  }

  @override
  Future<List<ConfiguredResource>> getConfiguredResources(
      int businessId) async {
    lastId = businessId;
    _trackCall('getConfiguredResources');
    return getConfiguredResourcesResult!;
  }

  @override
  Future<void> activateConfiguredResource(int resourceId) async {
    lastId = resourceId;
    _trackCall('activateConfiguredResource');
  }

  @override
  Future<void> deactivateConfiguredResource(int resourceId) async {
    lastId = resourceId;
    _trackCall('deactivateConfiguredResource');
  }

  @override
  Future<List<BusinessType>> getBusinessTypes() async {
    _trackCall('getBusinessTypes');
    return getBusinessTypesResult!;
  }

  @override
  Future<BusinessType> createBusinessType(Map<String, dynamic> data) async {
    lastBusinessTypeData = data;
    _trackCall('createBusinessType');
    return createBusinessTypeResult!;
  }

  @override
  Future<BusinessType> updateBusinessType(
      int id, Map<String, dynamic> data) async {
    lastId = id;
    lastBusinessTypeData = data;
    _trackCall('updateBusinessType');
    return updateBusinessTypeResult!;
  }

  @override
  Future<void> deleteBusinessType(int id) async {
    lastId = id;
    _trackCall('deleteBusinessType');
  }
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------
Business _makeBusiness({int id = 1, String name = 'Test'}) {
  return Business(id: id, name: name, isActive: true);
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

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockBusinessRepository mockRepo;
  late BusinessUseCases useCases;

  setUp(() {
    mockRepo = MockBusinessRepository();
    useCases = BusinessUseCases(mockRepo);
  });

  group('getBusinesses', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Business>(
        data: [_makeBusiness()],
        pagination: _makePagination(),
      );
      mockRepo.getBusinessesResult = expected;

      final params = GetBusinessesParams(page: 1, pageSize: 10);
      final result = await useCases.getBusinesses(params);

      expect(mockRepo.calls, ['getBusinesses']);
      expect(mockRepo.lastGetBusinessesParams, params);
      expect(result.data.length, 1);
      expect(result.data.first.name, 'Test');
    });

    test('delegates with null params', () async {
      mockRepo.getBusinessesResult = PaginatedResponse<Business>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getBusinesses(null);

      expect(mockRepo.calls, ['getBusinesses']);
      expect(mockRepo.lastGetBusinessesParams, isNull);
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.getBusinesses(null),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getBusinessById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getBusinessByIdResult = _makeBusiness(id: 42, name: 'Found');

      final result = await useCases.getBusinessById(42);

      expect(mockRepo.calls, ['getBusinessById']);
      expect(mockRepo.lastId, 42);
      expect(result.id, 42);
      expect(result.name, 'Found');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('not found');

      expect(
        () => useCases.getBusinessById(999),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('createBusiness', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateBusinessDTO(name: 'New Biz', domain: 'new.com');
      mockRepo.createBusinessResult = _makeBusiness(id: 10, name: 'New Biz');

      final result = await useCases.createBusiness(dto);

      expect(mockRepo.calls, ['createBusiness']);
      expect(mockRepo.lastCreateDTO, dto);
      expect(result.id, 10);
      expect(result.name, 'New Biz');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('validation error');

      expect(
        () => useCases.createBusiness(CreateBusinessDTO(name: 'Fail')),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('updateBusiness', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = UpdateBusinessDTO(name: 'Updated');
      mockRepo.updateBusinessResult = _makeBusiness(id: 5, name: 'Updated');

      final result = await useCases.updateBusiness(5, dto);

      expect(mockRepo.calls, ['updateBusiness']);
      expect(mockRepo.lastId, 5);
      expect(mockRepo.lastUpdateDTO, dto);
      expect(result.name, 'Updated');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('update failed');

      expect(
        () => useCases.updateBusiness(1, UpdateBusinessDTO(name: 'X')),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('deleteBusiness', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteBusiness(7);

      expect(mockRepo.calls, ['deleteBusiness']);
      expect(mockRepo.lastId, 7);
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('delete failed');

      expect(
        () => useCases.deleteBusiness(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('activateBusiness', () {
    test('delegates to repository with correct id', () async {
      await useCases.activateBusiness(3);

      expect(mockRepo.calls, ['activateBusiness']);
      expect(mockRepo.lastId, 3);
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('activate failed');

      expect(
        () => useCases.activateBusiness(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('deactivateBusiness', () {
    test('delegates to repository with correct id', () async {
      await useCases.deactivateBusiness(4);

      expect(mockRepo.calls, ['deactivateBusiness']);
      expect(mockRepo.lastId, 4);
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('deactivate failed');

      expect(
        () => useCases.deactivateBusiness(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getBusinessesSimple', () {
    test('delegates to repository and returns list', () async {
      mockRepo.getBusinessesSimpleResult = [
        BusinessSimple(id: 1, name: 'A'),
        BusinessSimple(id: 2, name: 'B'),
      ];

      final result = await useCases.getBusinessesSimple();

      expect(mockRepo.calls, ['getBusinessesSimple']);
      expect(result.length, 2);
      expect(result[0].name, 'A');
      expect(result[1].name, 'B');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('simple fetch failed');

      expect(
        () => useCases.getBusinessesSimple(),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getConfiguredResources', () {
    test('delegates to repository with correct businessId', () async {
      mockRepo.getConfiguredResourcesResult = [
        ConfiguredResource(id: 1, name: 'Inventory', isActive: true),
      ];

      final result = await useCases.getConfiguredResources(99);

      expect(mockRepo.calls, ['getConfiguredResources']);
      expect(mockRepo.lastId, 99);
      expect(result.length, 1);
      expect(result.first.name, 'Inventory');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('resources error');

      expect(
        () => useCases.getConfiguredResources(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('activateConfiguredResource', () {
    test('delegates to repository with correct resourceId', () async {
      await useCases.activateConfiguredResource(15);

      expect(mockRepo.calls, ['activateConfiguredResource']);
      expect(mockRepo.lastId, 15);
    });
  });

  group('deactivateConfiguredResource', () {
    test('delegates to repository with correct resourceId', () async {
      await useCases.deactivateConfiguredResource(16);

      expect(mockRepo.calls, ['deactivateConfiguredResource']);
      expect(mockRepo.lastId, 16);
    });
  });

  group('getBusinessTypes', () {
    test('delegates to repository and returns list', () async {
      mockRepo.getBusinessTypesResult = [
        BusinessType(id: 1, name: 'Retail'),
        BusinessType(id: 2, name: 'Food'),
      ];

      final result = await useCases.getBusinessTypes();

      expect(mockRepo.calls, ['getBusinessTypes']);
      expect(result.length, 2);
      expect(result[0].name, 'Retail');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('types error');

      expect(
        () => useCases.getBusinessTypes(),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('createBusinessType', () {
    test('delegates to repository with correct data', () async {
      final data = {'name': 'New Type', 'code': 'NT'};
      mockRepo.createBusinessTypeResult =
          BusinessType(id: 5, name: 'New Type', code: 'NT');

      final result = await useCases.createBusinessType(data);

      expect(mockRepo.calls, ['createBusinessType']);
      expect(mockRepo.lastBusinessTypeData, data);
      expect(result.id, 5);
      expect(result.name, 'New Type');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('create type failed');

      expect(
        () => useCases.createBusinessType({'name': 'Fail'}),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('updateBusinessType', () {
    test('delegates to repository with correct id and data', () async {
      final data = {'name': 'Updated Type'};
      mockRepo.updateBusinessTypeResult =
          BusinessType(id: 3, name: 'Updated Type');

      final result = await useCases.updateBusinessType(3, data);

      expect(mockRepo.calls, ['updateBusinessType']);
      expect(mockRepo.lastId, 3);
      expect(mockRepo.lastBusinessTypeData, data);
      expect(result.name, 'Updated Type');
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('update type failed');

      expect(
        () => useCases.updateBusinessType(1, {'name': 'X'}),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('deleteBusinessType', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteBusinessType(8);

      expect(mockRepo.calls, ['deleteBusinessType']);
      expect(mockRepo.lastId, 8);
    });

    test('propagates error from repository', () async {
      mockRepo.errorToThrow = Exception('delete type failed');

      expect(
        () => useCases.deleteBusinessType(1),
        throwsA(isA<Exception>()),
      );
    });
  });
}
