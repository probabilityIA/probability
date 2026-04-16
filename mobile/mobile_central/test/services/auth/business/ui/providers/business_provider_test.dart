import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/core/network/api_client.dart';
import 'package:mobile_central/services/auth/business/app/use_cases.dart';
import 'package:mobile_central/services/auth/business/domain/entities.dart';
import 'package:mobile_central/services/auth/business/domain/ports.dart';
import 'package:mobile_central/services/auth/business/ui/providers/business_provider.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IBusinessRepository
// ---------------------------------------------------------------------------
class MockBusinessRepository implements IBusinessRepository {
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

  // Call tracking
  final List<String> calls = [];

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<PaginatedResponse<Business>> getBusinesses(
      GetBusinessesParams? params) async {
    _trackCall('getBusinesses');
    return getBusinessesResult!;
  }

  @override
  Future<Business> getBusinessById(int id) async {
    _trackCall('getBusinessById');
    return getBusinessByIdResult!;
  }

  @override
  Future<Business> createBusiness(CreateBusinessDTO data) async {
    _trackCall('createBusiness');
    return createBusinessResult!;
  }

  @override
  Future<Business> updateBusiness(int id, UpdateBusinessDTO data) async {
    _trackCall('updateBusiness');
    return updateBusinessResult!;
  }

  @override
  Future<void> deleteBusiness(int id) async {
    _trackCall('deleteBusiness');
  }

  @override
  Future<void> activateBusiness(int id) async {
    _trackCall('activateBusiness');
  }

  @override
  Future<void> deactivateBusiness(int id) async {
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
    _trackCall('getConfiguredResources');
    return getConfiguredResourcesResult!;
  }

  @override
  Future<void> activateConfiguredResource(int resourceId) async {
    _trackCall('activateConfiguredResource');
  }

  @override
  Future<void> deactivateConfiguredResource(int resourceId) async {
    _trackCall('deactivateConfiguredResource');
  }

  @override
  Future<List<BusinessType>> getBusinessTypes() async {
    _trackCall('getBusinessTypes');
    return getBusinessTypesResult!;
  }

  @override
  Future<BusinessType> createBusinessType(Map<String, dynamic> data) async {
    _trackCall('createBusinessType');
    return createBusinessTypeResult!;
  }

  @override
  Future<BusinessType> updateBusinessType(
      int id, Map<String, dynamic> data) async {
    _trackCall('updateBusinessType');
    return updateBusinessTypeResult!;
  }

  @override
  Future<void> deleteBusinessType(int id) async {
    _trackCall('deleteBusinessType');
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
Business _makeBusiness({int id = 1, String name = 'Test', bool isActive = true}) {
  return Business(id: id, name: name, isActive: isActive);
}

Pagination _makePagination({
  int currentPage = 1,
  int perPage = 10,
  int total = 50,
  int lastPage = 5,
  bool hasNext = true,
  bool hasPrev = false,
}) {
  return Pagination(
    currentPage: currentPage,
    perPage: perPage,
    total: total,
    lastPage: lastPage,
    hasNext: hasNext,
    hasPrev: hasPrev,
  );
}

BusinessProvider _createProvider(MockBusinessRepository mockRepo) {
  final apiClient = ApiClient();
  final useCases = BusinessUseCases(mockRepo);
  return BusinessProvider(apiClient: apiClient, useCases: useCases);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockBusinessRepository mockRepo;
  late BusinessProvider provider;

  setUp(() {
    mockRepo = MockBusinessRepository();
    provider = _createProvider(mockRepo);
  });

  group('Initial state', () {
    test('has empty businesses list', () {
      expect(provider.businesses, isEmpty);
    });

    test('has empty businessesSimple list', () {
      expect(provider.businessesSimple, isEmpty);
    });

    test('has empty businessTypes list', () {
      expect(provider.businessTypes, isEmpty);
    });

    test('has null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('is not loading', () {
      expect(provider.isLoading, false);
    });

    test('has no error', () {
      expect(provider.error, isNull);
    });

    test('has no selectedBusinessId', () {
      expect(provider.selectedBusinessId, isNull);
    });
  });

  group('setSelectedBusinessId', () {
    test('updates selectedBusinessId and notifies listeners', () {
      var notified = false;
      provider.addListener(() => notified = true);

      provider.setSelectedBusinessId(42);

      expect(provider.selectedBusinessId, 42);
      expect(notified, true);
    });

    test('can set to null', () {
      provider.setSelectedBusinessId(10);
      provider.setSelectedBusinessId(null);

      expect(provider.selectedBusinessId, isNull);
    });
  });

  group('fetchBusinesses', () {
    test('updates businesses and pagination on success', () async {
      mockRepo.getBusinessesResult = PaginatedResponse<Business>(
        data: [
          _makeBusiness(id: 1, name: 'Biz A'),
          _makeBusiness(id: 2, name: 'Biz B'),
        ],
        pagination: _makePagination(total: 2, lastPage: 1, hasNext: false),
      );

      await provider.fetchBusinesses();

      expect(provider.businesses.length, 2);
      expect(provider.businesses[0].name, 'Biz A');
      expect(provider.businesses[1].name, 'Biz B');
      expect(provider.pagination, isNotNull);
      expect(provider.pagination!.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets isLoading during fetch and clears after', () async {
      final loadingStates = <bool>[];
      provider.addListener(() {
        loadingStates.add(provider.isLoading);
      });

      mockRepo.getBusinessesResult = PaginatedResponse<Business>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchBusinesses();

      // First notification: isLoading = true, second: isLoading = false
      expect(loadingStates, [true, false]);
    });

    test('clears previous error before fetching', () async {
      // First set an error state
      mockRepo.errorToThrow = Exception('first error');
      await provider.fetchBusinesses();
      expect(provider.error, isNotNull);

      // Now succeed -- error should be cleared
      mockRepo.errorToThrow = null;
      mockRepo.getBusinessesResult = PaginatedResponse<Business>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchBusinesses();

      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('fetch failed');

      await provider.fetchBusinesses();

      expect(provider.error, contains('fetch failed'));
      expect(provider.isLoading, false);
    });

    test('passes params to use case', () async {
      mockRepo.getBusinessesResult = PaginatedResponse<Business>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchBusinesses(
        params: GetBusinessesParams(page: 2, pageSize: 20),
      );

      expect(mockRepo.calls, ['getBusinesses']);
    });
  });

  group('fetchBusinessesSimple', () {
    test('updates businessesSimple list on success', () async {
      mockRepo.getBusinessesSimpleResult = [
        BusinessSimple(id: 1, name: 'Simple A'),
        BusinessSimple(id: 2, name: 'Simple B'),
      ];

      await provider.fetchBusinessesSimple();

      expect(provider.businessesSimple.length, 2);
      expect(provider.businessesSimple[0].name, 'Simple A');
      expect(provider.businessesSimple[1].name, 'Simple B');
    });

    test('notifies listeners on success', () async {
      var notifyCount = 0;
      provider.addListener(() => notifyCount++);

      mockRepo.getBusinessesSimpleResult = [
        BusinessSimple(id: 1, name: 'S'),
      ];

      await provider.fetchBusinessesSimple();

      expect(notifyCount, 1);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('simple fetch failed');

      await provider.fetchBusinessesSimple();

      expect(provider.error, contains('simple fetch failed'));
    });
  });

  group('fetchBusinessTypes', () {
    test('updates businessTypes list on success', () async {
      mockRepo.getBusinessTypesResult = [
        BusinessType(id: 1, name: 'Retail'),
        BusinessType(id: 2, name: 'Food'),
      ];

      await provider.fetchBusinessTypes();

      expect(provider.businessTypes.length, 2);
      expect(provider.businessTypes[0].name, 'Retail');
      expect(provider.businessTypes[1].name, 'Food');
    });

    test('notifies listeners on success', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.getBusinessTypesResult = [];

      await provider.fetchBusinessTypes();

      expect(notified, true);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('types error');

      await provider.fetchBusinessTypes();

      expect(provider.error, contains('types error'));
    });
  });

  group('createBusiness', () {
    test('returns created Business on success', () async {
      mockRepo.createBusinessResult = _makeBusiness(id: 99, name: 'Created');

      final result = await provider.createBusiness(
        CreateBusinessDTO(name: 'Created'),
      );

      expect(result, isNotNull);
      expect(result!.id, 99);
      expect(result.name, 'Created');
      expect(provider.error, isNull);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('create failed');

      final result = await provider.createBusiness(
        CreateBusinessDTO(name: 'Fail'),
      );

      expect(result, isNull);
      expect(provider.error, contains('create failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('create error');

      await provider.createBusiness(CreateBusinessDTO(name: 'Fail'));

      expect(notified, true);
    });
  });

  group('updateBusiness', () {
    test('returns true on success', () async {
      mockRepo.updateBusinessResult = _makeBusiness(id: 1, name: 'Updated');

      final result = await provider.updateBusiness(
        1,
        UpdateBusinessDTO(name: 'Updated'),
      );

      expect(result, true);
      expect(provider.error, isNull);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('update failed');

      final result = await provider.updateBusiness(
        1,
        UpdateBusinessDTO(name: 'Fail'),
      );

      expect(result, false);
      expect(provider.error, contains('update failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('update error');

      await provider.updateBusiness(1, UpdateBusinessDTO(name: 'Fail'));

      expect(notified, true);
    });
  });

  group('deleteBusiness', () {
    test('returns true on success', () async {
      final result = await provider.deleteBusiness(5);

      expect(result, true);
      expect(mockRepo.calls, ['deleteBusiness']);
      expect(provider.error, isNull);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('delete failed');

      final result = await provider.deleteBusiness(5);

      expect(result, false);
      expect(provider.error, contains('delete failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('delete error');

      await provider.deleteBusiness(1);

      expect(notified, true);
    });
  });

  group('activateBusiness', () {
    test('returns true on success', () async {
      final result = await provider.activateBusiness(3);

      expect(result, true);
      expect(mockRepo.calls, ['activateBusiness']);
      expect(provider.error, isNull);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('activate failed');

      final result = await provider.activateBusiness(3);

      expect(result, false);
      expect(provider.error, contains('activate failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('activate error');

      await provider.activateBusiness(1);

      expect(notified, true);
    });
  });

  group('deactivateBusiness', () {
    test('returns true on success', () async {
      final result = await provider.deactivateBusiness(4);

      expect(result, true);
      expect(mockRepo.calls, ['deactivateBusiness']);
      expect(provider.error, isNull);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('deactivate failed');

      final result = await provider.deactivateBusiness(4);

      expect(result, false);
      expect(provider.error, contains('deactivate failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('deactivate error');

      await provider.deactivateBusiness(1);

      expect(notified, true);
    });
  });
}
