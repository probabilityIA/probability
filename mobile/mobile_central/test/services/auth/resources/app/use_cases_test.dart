import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/resources/app/use_cases.dart';
import 'package:mobile_central/services/auth/resources/domain/entities.dart';
import 'package:mobile_central/services/auth/resources/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockResourceRepository implements IResourceRepository {
  final List<String> calls = [];

  PaginatedResponse<Resource>? getResourcesResult;
  Resource? getResourceByIdResult;
  Resource? createResourceResult;
  Resource? updateResourceResult;

  Exception? errorToThrow;

  GetResourcesParams? capturedGetParams;
  int? capturedId;
  CreateResourceDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateResourceDTO? capturedUpdateData;
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<Resource>> getResources(
      GetResourcesParams? params) async {
    calls.add('getResources');
    capturedGetParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getResourcesResult!;
  }

  @override
  Future<Resource> getResourceById(int id) async {
    calls.add('getResourceById');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getResourceByIdResult!;
  }

  @override
  Future<Resource> createResource(CreateResourceDTO data) async {
    calls.add('createResource');
    capturedCreateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createResourceResult!;
  }

  @override
  Future<Resource> updateResource(int id, UpdateResourceDTO data) async {
    calls.add('updateResource');
    capturedUpdateId = id;
    capturedUpdateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateResourceResult!;
  }

  @override
  Future<void> deleteResource(int id) async {
    calls.add('deleteResource');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Helpers ---

Resource _makeResource({int id = 1, String name = 'orders'}) {
  return Resource(id: id, name: name);
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
  late MockResourceRepository mockRepo;
  late ResourceUseCases useCases;

  setUp(() {
    mockRepo = MockResourceRepository();
    useCases = ResourceUseCases(mockRepo);
  });

  group('getResources', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Resource>(
        data: [_makeResource()],
        pagination: _makePagination(),
      );
      mockRepo.getResourcesResult = expected;
      final params = GetResourcesParams(page: 1, pageSize: 10);

      final result = await useCases.getResources(params);

      expect(result.data.length, 1);
      expect(result.data[0].name, 'orders');
      expect(mockRepo.calls, ['getResources']);
      expect(mockRepo.capturedGetParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getResourcesResult = PaginatedResponse<Resource>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getResources(null);

      expect(mockRepo.capturedGetParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getResources(null), throwsException);
    });
  });

  group('getResourceById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getResourceByIdResult =
          _makeResource(id: 42, name: 'products');

      final result = await useCases.getResourceById(42);

      expect(result.id, 42);
      expect(result.name, 'products');
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getResourceById']);
    });
  });

  group('createResource', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateResourceDTO(name: 'shipments');
      mockRepo.createResourceResult =
          _makeResource(id: 99, name: 'shipments');

      final result = await useCases.createResource(dto);

      expect(result.id, 99);
      expect(result.name, 'shipments');
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createResource']);
    });
  });

  group('updateResource', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateResourceDTO(name: 'updated');
      mockRepo.updateResourceResult = _makeResource(id: 5, name: 'updated');

      final result = await useCases.updateResource(5, dto);

      expect(result.id, 5);
      expect(result.name, 'updated');
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateResource']);
    });
  });

  group('deleteResource', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteResource(7);

      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deleteResource']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.deleteResource(7), throwsException);
    });
  });
}
