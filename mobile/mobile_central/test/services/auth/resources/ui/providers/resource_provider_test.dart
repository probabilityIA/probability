import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/resources/app/use_cases.dart';
import 'package:mobile_central/services/auth/resources/domain/entities.dart';
import 'package:mobile_central/services/auth/resources/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockResourceRepository implements IResourceRepository {
  PaginatedResponse<Resource>? getResourcesResult;
  Resource? createResourceResult;
  Resource? updateResourceResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<Resource>> getResources(
      GetResourcesParams? params) async {
    calls.add('getResources');
    if (errorToThrow != null) throw errorToThrow!;
    return getResourcesResult!;
  }

  @override
  Future<Resource> getResourceById(int id) async {
    calls.add('getResourceById');
    if (errorToThrow != null) throw errorToThrow!;
    return Resource(id: id, name: 'test');
  }

  @override
  Future<Resource> createResource(CreateResourceDTO data) async {
    calls.add('createResource');
    if (errorToThrow != null) throw errorToThrow!;
    return createResourceResult!;
  }

  @override
  Future<Resource> updateResource(int id, UpdateResourceDTO data) async {
    calls.add('updateResource');
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

// --- Testable Provider ---

class TestableResourceProvider {
  final ResourceUseCases _useCases;

  List<Resource> _resources = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableResourceProvider(this._useCases);

  List<Resource> get resources => _resources;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchResources({GetResourcesParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.getResources(params);
      _resources = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Resource?> createResource(CreateResourceDTO data) async {
    try {
      return await _useCases.createResource(data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateResource(int id, UpdateResourceDTO data) async {
    try {
      await _useCases.updateResource(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteResource(int id) async {
    try {
      await _useCases.deleteResource(id);
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

Resource _makeResource({int id = 1, String name = 'orders'}) {
  return Resource(id: id, name: name);
}

// --- Tests ---

void main() {
  late MockResourceRepository mockRepo;
  late ResourceUseCases useCases;
  late TestableResourceProvider provider;

  setUp(() {
    mockRepo = MockResourceRepository();
    useCases = ResourceUseCases(mockRepo);
    provider = TestableResourceProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty resources list', () {
      expect(provider.resources, isEmpty);
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

  group('fetchResources', () {
    test('notifies listeners twice (loading start and end)', () async {
      mockRepo.getResourcesResult = PaginatedResponse<Resource>(
        data: [_makeResource()],
        pagination: _makePagination(),
      );

      await provider.fetchResources();

      expect(provider.notifications.length, 2);
    });

    test('populates resources and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getResourcesResult = PaginatedResponse<Resource>(
        data: [
          _makeResource(id: 1, name: 'orders'),
          _makeResource(id: 2, name: 'products'),
        ],
        pagination: pagination,
      );

      await provider.fetchResources();

      expect(provider.resources.length, 2);
      expect(provider.resources[0].name, 'orders');
      expect(provider.resources[1].name, 'products');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchResources();

      expect(provider.error, contains('Server error'));
      expect(provider.resources, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchResources();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getResourcesResult = PaginatedResponse<Resource>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchResources();

      expect(provider.error, isNull);
    });

    test('passes custom params to use cases', () async {
      mockRepo.getResourcesResult = PaginatedResponse<Resource>(
        data: [],
        pagination: _makePagination(),
      );
      final params = GetResourcesParams(page: 2, name: 'products');

      await provider.fetchResources(params: params);

      expect(mockRepo.calls, contains('getResources'));
    });
  });

  group('createResource', () {
    test('returns created resource on success', () async {
      final dto = CreateResourceDTO(name: 'shipments');
      mockRepo.createResourceResult =
          _makeResource(id: 10, name: 'shipments');

      final result = await provider.createResource(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.name, 'shipments');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateResourceDTO(name: 'fail');

      final result = await provider.createResource(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateResource', () {
    test('returns true on success', () async {
      final dto = UpdateResourceDTO(name: 'updated');
      mockRepo.updateResourceResult = _makeResource(id: 5, name: 'updated');

      final result = await provider.updateResource(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateResourceDTO(name: 'fail');

      final result = await provider.updateResource(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteResource', () {
    test('returns true on success', () async {
      final result = await provider.deleteResource(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteResource(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });
}
