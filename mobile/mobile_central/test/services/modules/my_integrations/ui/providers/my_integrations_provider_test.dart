import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/core/network/api_client.dart';
import 'package:mobile_central/services/modules/my_integrations/app/use_cases.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/entities.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/ports.dart';
import 'package:mobile_central/services/modules/my_integrations/ui/providers/my_integrations_provider.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IMyIntegrationsRepository
// ---------------------------------------------------------------------------
class MockMyIntegrationsRepository implements IMyIntegrationsRepository {
  final List<String> calls = [];

  PaginatedResponse<MyIntegration>? getIntegrationsResult;
  MyIntegration? getIntegrationByIdResult;

  Exception? errorToThrow;

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<PaginatedResponse<MyIntegration>> getIntegrations(
      GetMyIntegrationsParams? params) async {
    _trackCall('getIntegrations');
    return getIntegrationsResult!;
  }

  @override
  Future<MyIntegration> getIntegrationById(int id,
      {int? businessId}) async {
    _trackCall('getIntegrationById');
    return getIntegrationByIdResult!;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
Pagination _makePagination({int total = 50}) {
  return Pagination(
    currentPage: 1,
    perPage: 20,
    total: total,
    lastPage: 3,
    hasNext: true,
    hasPrev: false,
  );
}

MyIntegration _makeIntegration({int id = 1, String name = 'Test Store'}) {
  return MyIntegration(
    id: id,
    createdAt: '',
    updatedAt: '',
    businessId: 1,
    integrationTypeId: 1,
    integrationTypeName: 'Shopify',
    name: name,
    isActive: true,
  );
}

MyIntegrationsProvider _createProvider(MockMyIntegrationsRepository mockRepo) {
  final apiClient = ApiClient();
  final useCases = MyIntegrationsUseCases(mockRepo);
  return MyIntegrationsProvider(apiClient: apiClient, useCases: useCases);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockMyIntegrationsRepository mockRepo;
  late MyIntegrationsProvider provider;

  setUp(() {
    mockRepo = MockMyIntegrationsRepository();
    provider = _createProvider(mockRepo);
  });

  group('Initial state', () {
    test('has empty integrations list', () {
      expect(provider.integrations, isEmpty);
    });

    test('has null selectedIntegration', () {
      expect(provider.selectedIntegration, isNull);
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

    test('has default page 1', () {
      expect(provider.page, 1);
    });
  });

  group('fetchIntegrations', () {
    test('updates integrations and pagination on success', () async {
      mockRepo.getIntegrationsResult = PaginatedResponse<MyIntegration>(
        data: [
          _makeIntegration(id: 1, name: 'Store A'),
          _makeIntegration(id: 2, name: 'Store B'),
        ],
        pagination: _makePagination(total: 2),
      );

      await provider.fetchIntegrations();

      expect(provider.integrations.length, 2);
      expect(provider.integrations[0].name, 'Store A');
      expect(provider.integrations[1].name, 'Store B');
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

      mockRepo.getIntegrationsResult = PaginatedResponse<MyIntegration>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchIntegrations();

      expect(loadingStates, [true, false]);
    });

    test('clears previous error before fetching', () async {
      mockRepo.errorToThrow = Exception('first error');
      await provider.fetchIntegrations();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getIntegrationsResult = PaginatedResponse<MyIntegration>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchIntegrations();

      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('fetch failed');

      await provider.fetchIntegrations();

      expect(provider.error, contains('fetch failed'));
      expect(provider.isLoading, false);
    });
  });

  group('fetchIntegrationById', () {
    test('updates selectedIntegration on success', () async {
      mockRepo.getIntegrationByIdResult =
          _makeIntegration(id: 42, name: 'Amazon Store');

      await provider.fetchIntegrationById(42);

      expect(provider.selectedIntegration, isNotNull);
      expect(provider.selectedIntegration!.id, 42);
      expect(provider.selectedIntegration!.name, 'Amazon Store');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets isLoading during fetch and clears after', () async {
      final loadingStates = <bool>[];
      provider.addListener(() {
        loadingStates.add(provider.isLoading);
      });

      mockRepo.getIntegrationByIdResult = _makeIntegration(id: 1);

      await provider.fetchIntegrationById(1);

      expect(loadingStates, [true, false]);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('not found');

      await provider.fetchIntegrationById(999);

      expect(provider.error, contains('not found'));
      expect(provider.isLoading, false);
    });
  });

  group('setPage', () {
    test('updates page', () {
      provider.setPage(3);
      expect(provider.page, 3);
    });
  });

  group('setFilters', () {
    test('resets page to 1', () {
      provider.setPage(5);
      provider.setFilters(categoryCode: 'ecommerce');
      expect(provider.page, 1);
    });
  });

  group('resetFilters', () {
    test('resets page to 1', () {
      provider.setPage(5);
      provider.resetFilters();
      expect(provider.page, 1);
    });
  });
}
